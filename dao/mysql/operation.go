package mysql

import (
	"errors"
	"log"
	"reflect"

	"github.com/go-xorm/xorm"

	"github.com/chinx/coupon/dao"
)

var (
	NoRecords      = errors.New("no eligible records in the table")
	insertFailed   = "insert to table %s failed, error: %s\n"
	updateFailed   = "update from table %s failed, error: %s\n"
	findFailed     = "find in table %s failed, error: %s\n"
	deleteFailed   = "delete from table %s failed, error: %s\n"
	scanRowsFailed = "scan rows from table %s result failed, error: %s\n"
)

func init() {
	register(&dao.OfficialAccount{}, &dao.Admin{}, &dao.User{},
		&dao.Activity{}, &dao.Task{}, &dao.Bargain{})
}

type tabler interface {
	TableName() string
}

type DB interface {
	Insert(...interface{}) (int64, error)
	Where(interface{}, ...interface{}) *xorm.Session
	Update(interface{}, ...interface{}) (int64, error)
	In(string, ...interface{}) *xorm.Session
	Count(...interface{}) (int64, error)
	Desc(...string) *xorm.Session
}

func insert(db DB, tb tabler) (err error) {
	n, err := db.Insert(tb)
	if err == nil {
		if n == 0 {
			err = NoRecords
		}
	}
	if err != nil {
		log.Printf(insertFailed, tb.TableName(), err)
	}
	return
}

func update(db DB, tb tabler, condition string, args ...interface{}) (err error) {
	n, err := db.Where(condition, args...).Update(tb)
	if err == nil {
		if n == 0 {
			err = NoRecords
		}
	}
	if err != nil {
		log.Printf(updateFailed, tb.TableName(), err)
	}
	return
}

func find(db DB, tb tabler, condition string, args ...interface{}) (list []interface{}, err error) {
	rows, err := db.Where(condition, args...).Rows(tb)
	if err != nil {
		log.Printf(findFailed, tb.TableName(), err)
	}
	defer rows.Close()
	return scanRows(rows, tb), err
}

func in(db DB, tb tabler, field string, conditions []interface{}) []interface{} {
	rows, err := db.In(field, conditions...).Rows(tb)
	if err != nil {
		log.Printf(findFailed, tb.TableName(), err)
		return []interface{}{}
	}
	defer rows.Close()
	return scanRows(rows, tb)
}

func list(db DB, tb tabler, offset, count int, condition string, args ...interface{}) (int64, []interface{}) {
	n, _ := db.Count(tb)
	if n == 0 || offset >= int(n) {
		return 0, []interface{}{}
	}

	rows, err := db.Desc("id").Limit(count, offset).Rows(tb)
	if err != nil {
		log.Printf(findFailed, tb.TableName(), err)
		return 0, []interface{}{}
	}
	defer rows.Close()
	return n, scanRows(rows, tb)
}

func getLast(db DB, tb tabler, order, condition string, args ...interface{}) (err error) {
	existed, err := db.Where(condition, args...).OrderBy(order).Get(tb)
	if err == nil {
		if !existed {
			err = NoRecords
		}
	}
	if err != nil {
		log.Printf(findFailed, tb.TableName(), err)
	}
	return
}

func get(db DB, tb tabler, condition string, args ...interface{}) (err error) {
	existed, err := db.Where(condition, args...).Get(tb)
	if err == nil {
		if !existed {
			err = NoRecords
		}
	}
	if err != nil {
		log.Printf(findFailed, tb.TableName(), err)
	}
	return
}

func exist(db DB, tb tabler, condition string, args ...interface{}) bool {
	existed, err := db.Where(condition, args...).Exist(tb)
	if err != nil || !existed {
		return false
	}
	return true
}

func delete(db DB, tb tabler, condition string, args ...interface{}) (err error) {
	n, err := db.Where(condition, args).Delete(tb)
	if err == nil {
		if n == 0 {
			err = NoRecords
		}
	}
	if err != nil {
		log.Printf(deleteFailed, tb.TableName(), err)
	}
	return
}

func scanRows(rows *xorm.Rows, tb tabler) []interface{} {
	items := make([]interface{}, 0, 10)
	for rows.Next() {
		item := newInterface(tb)
		err := rows.Scan(item)
		if err != nil {
			log.Fatalf(scanRowsFailed, tb.TableName(), err)
			break
		}
		items = append(items, item)
	}
	return items
}

func newInterface(tb interface{}) interface{} {
	val := reflect.Indirect(reflect.ValueOf(tb))
	t := val.Type()
	return reflect.New(t).Interface()
}

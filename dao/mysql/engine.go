package mysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

var tables = make([]xorm.TableName, 0, 16)
var engine *xorm.Engine

func InitORM(driver, source string) error {
	eng, err := xorm.NewEngine(driver, source)
	if err != nil {
		return err
	}

	if err = eng.Ping(); err != nil {
		return err
	}

	eng.SetMapper(core.GonicMapper{})
	engine = eng
	engine.ShowSQL(true)

	syncTables()
	return nil
}

func syncTables() {
	for _, value := range tables {
		engine.Sync2(value)
	}
}

func register(tabs ...xorm.TableName) {
	for _, tab := range tabs {
		if isExist(tab.TableName()) {
			fmt.Printf("table %s is already exists", tab.TableName())
		}
		tables = append(tables, tab)
	}
}

func isExist(name string) bool {
	for _, value := range tables {
		if value.TableName() == name {
			return true
		}
	}
	return false
}

func Insert(tb tabler) (err error) {
	return insert(engine, tb)
}

func Update(tb tabler, condition string, args ...interface{}) (err error) {
	return update(engine, tb, condition, args...)
}

func Find(tb tabler, condition string, args ...interface{}) (list []interface{}, err error) {
	return find(engine, tb, condition, args...)
}

func In(tb tabler, field string, conditions []interface{}) []interface{} {
	return in(engine, tb, field, conditions)
}

func List(tb tabler, offset, count int, condition string, args ...interface{}) (int64, []interface{}) {
	return list(engine, tb, offset, count, condition, args...)
}

func GetLast(tb tabler, order, condition string, args ...interface{}) (err error) {
	return getLast(engine, tb, order, condition, args...)
}

func Get(tb tabler, condition string, args ...interface{}) (err error) {
	return get(engine, tb, condition, args...)
}

func Exist(tb tabler, condition string, args ...interface{}) bool {
	return exist(engine, tb, condition, args...)
}

func Delete(tb tabler, condition string, args ...interface{}) (err error) {
	return delete(engine, tb, condition, args...)
}

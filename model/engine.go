package model

import (
	"fmt"
	"reflect"

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

func NewSession() *xorm.Session {
	return engine.NewSession()
}

func Insert(tb interface{}) bool {
	n, err := engine.Insert(tb)
	if err != nil || n == 0 {
		return false
	}
	return true
}

func Update(tb interface{}) bool {
	n, err := engine.Update(tb)
	if err != nil || n == 0 {
		return false
	}
	return true
}

func Find(tb interface{}) bool {
	err := engine.Find(tb)
	if err != nil {
		return false
	}
	return true
}

func Get(tb interface{}) bool {
	ok, err := engine.Get(tb)
	if err != nil || !ok {
		return false
	}
	return true
}

func Delete(tb interface{}) bool {
	n, err := engine.Delete(tb)
	if err != nil || n == 0 {
		return false
	}
	return true
}

func In(tb interface{}, field string, conditions []interface{}) []interface{} {
	rows, err := engine.In(field, conditions...).Desc(field).Rows(tb)
	if err != nil {
		return nil
	}
	items := make([]interface{}, 0, len(conditions))
	defer rows.Close()
	for rows.Next() {
		item := newInterface(tb)
		err = rows.Scan(item)
		if err != nil {
			break
		}
		items = append(items, item)
	}
	if err != nil {
		return nil
	}
	return items
}

func newInterface(tb interface{}) interface{} {
	t := reflect.ValueOf(tb).Type()
	return reflect.New(t).Interface()
}

package orm

import (
	"fmt"
	"reflect"
)

const (
	Exact OpType = iota // =
	Gt                  // >
	Gte                 // >=
	Lt                  // <
	Lte                 // <=
	Inn                 // in
)

type OpType int

type UzOrmI interface {
	Insert(md interface{}) (int64, error)

	Update(md interface{}, cols ...string) (int64, error)

	Delete(md interface{}) (int64, error)

	Select(md interface{}) error

	QueryTable(md interface{}) QueryPrepareI

	RegisterDataBase(host, user, password, dbName string)
}

type QueryPrepareI interface {
	QueryEnd()

	ExecEnd()

	Insert(mds ...interface{}) (int64, error)

	Update()

	Delete()

	Select()

	Limit()

	Count()

	OrderBy()

	GroupBy()

	And(col string, op OpType, value interface{}) QueryPrepareI

	Or(col string, op OpType, value interface{}) QueryPrepareI
}

func getInd(md interface{}) (ind reflect.Value) {
	val := reflect.ValueOf(md)
	ind = reflect.Indirect(val)
	typ := ind.Type()
	if val.Kind() != reflect.Ptr {
		panic(fmt.Errorf("[getInd] cannot use non-ptr model struct `%s`", typ))
	}
	if typ.Kind() == reflect.Ptr {
		panic(fmt.Errorf("[getInd] only allow ptr model struct, it looks you use two reference to the struct `%s`", typ))
	}
	return
}

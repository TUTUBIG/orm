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
	QueryEnd(mds interface{}, sql string, values []interface{}) error

	ExecEnd(sql string, values ...interface{}) (int64, error)

	Insert(mds ...interface{}) (int64, error)

	Update(cols ...string) (int64, error)

	Delete() (int64, error)

	Select(mds interface{}) error

	Limit(limit ...int) QueryPrepareI

	Count() (int64, error)

	OrderBy(col string) QueryPrepareI

	GroupBy(col string) QueryPrepareI

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

func getOpType(op OpType) string {
	switch op {
	case Exact:
		return "= ?"
	case Gt:
		return "> ?"
	case Gte:
		return ">= ?"
	case Lt:
		return "< ?"
	case Lte:
		return "<= ?"
	default:
		return "="
	}

}

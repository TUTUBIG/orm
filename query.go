package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
)

type FilterS struct {
	Col   string
	Op    OpType
	Value interface{}
	AndOr bool
}

type queryPrepare struct {
	md      interface{}
	db      *sql.DB
	filters []FilterS
}

func NewQueryPrepare(md interface{}, db *sql.DB) QueryPrepareI {
	q := new(queryPrepare)

	q.md = md
	q.db = db

	return q
}

func (q *queryPrepare) QueryEnd() {

}

func (q *queryPrepare) ExecEnd() {

}

func (q *queryPrepare) Insert(mds ...interface{}) (int64, error) {
	ind := getInd(q.md)

	tableName := ind.Type().Name()

	var sqlT, sqlV, sqlI string
	var values []interface{}

	for i := 0; i < ind.NumField(); i++ {
		if i == 0 {
			sqlT += tableName + "." + ind.Type().Field(i).Name
			sqlV += "?"
		} else {
			sqlT += "," + tableName + "." + ind.Type().Field(i).Name
			sqlV += "," + "?"
		}

		values = append(values, ind.Field(i).Interface())
	}

	sqlI = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, sqlT, sqlV)
	fmt.Println("insert model data to database sql: ", sqlI)

	if len(mds) > 0 {
		var sqlV string
		// var values  []interface{}

		sInd := reflect.Indirect(reflect.ValueOf(mds))

		switch sInd.Kind() {
		case reflect.Array, reflect.Slice:
			if sInd.Len() == 0 {
				return 0, errors.New("[queryPrepare.Insert] the models which is going to be inserted is empty ")
			}
		default:
			return 0, errors.New("[queryPrepare.Insert] the models which is going to be inserted is not a array or slice ")
		}

		// TODO

		sqlI = fmt.Sprintf("INSERT INTO %s (%s) %s", tableName, sqlT, sqlV)
	}

	res, err := q.db.Exec(sqlI, values...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()

}

func (q *queryPrepare) Update() {

}

func (q *queryPrepare) Delete() {

}

func (q *queryPrepare) Select() {

}

func (q *queryPrepare) Limit() {

}

func (q *queryPrepare) Count() {

}

func (q *queryPrepare) OrderBy() {

}

func (q *queryPrepare) GroupBy() {

}

func (q *queryPrepare) And(col string, op OpType, value interface{}) QueryPrepareI {
	filter := FilterS{
		Col:   col,
		Op:    op,
		Value: value,
		AndOr: true,
	}

	q.filters = append(q.filters, filter)

	return q
}

func (q *queryPrepare) Or(col string, op OpType, value interface{}) QueryPrepareI {
	filter := FilterS{
		Col:   col,
		Op:    op,
		Value: value,
		AndOr: false,
	}

	q.filters = append(q.filters, filter)

	return q
}

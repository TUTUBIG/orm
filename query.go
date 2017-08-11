package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
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
	orderBy string
	groupBy string
	limit   []int
}

func NewQueryPrepare(md interface{}, db *sql.DB) QueryPrepareI {
	q := new(queryPrepare)

	q.md = md
	q.db = db

	return q
}

/*
	[mds]

	type S struct {
		A int
		B string
	}

	type T struct {
		A []S
	}

	var t *T
*/

func (q *queryPrepare) QueryEnd(mds interface{}, sql string, values []interface{}) error {

	ind := reflect.Indirect(reflect.ValueOf(mds))

	switch ind.Kind() {
	case reflect.Struct:
	default:
		return errors.New("[queryPrepare.Insert] the models which is going to query data from database is not a ptr of struct ")
	}

	if ind.Type().NumField() != 1 {
		return errors.New("[queryPrepare.Insert] the models which is going to query data from database contains fields not just one or none ")
	}

	rvals := reflect.New(ind.Field(0).Type())
	rinds := reflect.Indirect(rvals)

	rows, err := q.db.Query(sql, values...)
	if err != nil {
		return err
	}

	var fields []interface{}

	for rows.Next() {
		fields = make([]interface{}, 0)
		rval := reflect.New(ind.Field(0).Type().Elem())
		rind := reflect.Indirect(rval)

		for i := 0; i < rind.NumField(); i++ {
			fields = append(fields, rind.Field(i).Addr().Interface())
		}

		if err := rows.Scan(fields...); err != nil {
			return err
		}

		rinds = reflect.Append(rinds, rind)

	}

	ind.Field(0).Set(rinds)

	return nil
}

func (q *queryPrepare) ExecEnd(sql string, values ...interface{}) (int64, error) {
	res, err := q.db.Exec(sql, values...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
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
		values = make([]interface{}, 0)
		sqlV = ""

		sInd := reflect.Indirect(reflect.ValueOf(mds))

		switch sInd.Kind() {
		case reflect.Array, reflect.Slice:
			if sInd.Len() == 0 {
				return 0, errors.New("[queryPrepare.Insert] the models which is going to be inserted is empty ")
			}
		default:
			return 0, errors.New("[queryPrepare.Insert] the models which is going to be inserted is not a array or slice ")
		}

		var m interface{}
		m = sInd.Index(0).Interface()

		mInd := reflect.Indirect(reflect.ValueOf(m))

		for i := 0; i < mInd.Len(); i++ {
			if i == 0 {
				sqlV += "("
			} else {
				sqlV += ",("
			}
			for j := 0; j < mInd.Index(i).NumField(); j++ {
				if j == 0 {
					sqlV += "?"
				} else {
					sqlV += ",?"
				}

				values = append(values, mInd.Index(i).Field(j).Interface())
			}
			sqlV += ")"
		}

		sqlI = fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", tableName, sqlT, sqlV)

	}

	fmt.Println("[queryPrepare.Insert] sqlI: ", sqlI)
	fmt.Println("[queryPrepare.Insert] values: ", values)

	return q.ExecEnd(sqlI, values...)
}

func (q *queryPrepare) Update(cols ...string) (int64, error) {
	ind := getInd(q.md)

	tableName := ind.Type().Name()

	var valuesW []interface{}
	pkKey := ind.Type().Field(0).Name
	pkValue := ind.Field(0).Interface()
	valuesW = append(valuesW, pkValue)

	sqlW := tableName + "." + pkKey + "=?"

	if len(q.filters) > 0 {
		sqlW = ""
		valuesW = make([]interface{}, 0)
		for i, f := range q.filters {
			if i == 0 {
				sqlW += tableName + "." + f.Col + getOpType(f.Op)
			} else {
				if f.AndOr {
					sqlW += " AND " + tableName + "." + f.Col + getOpType(f.Op)
				} else {
					sqlW += " OR " + tableName + "." + f.Col + getOpType(f.Op)
				}
			}

			valuesW = append(valuesW, ind.FieldByName(f.Col).Interface())
		}

	}

	var sqlT string
	var valuesC []interface{}
	if len(cols) > 0 {
		for i, c := range cols {
			if i == 0 {
				sqlT += tableName + "." + c + "=?"
			} else {
				sqlT += "," + tableName + "." + c + "=?"
			}

			valuesC = append(valuesC, ind.FieldByName(c).Interface())
		}
	} else {
		for i := 0; i < ind.NumField(); i++ {
			if i == 0 {
				sqlT += tableName + "." + ind.Type().Field(i).Name + "=?"
			} else {
				sqlT += "," + tableName + "." + ind.Type().Field(i).Name + "=?"
			}

			valuesC = append(valuesC, ind.Field(i).Interface())
		}
	}

	sqlD := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, sqlT, sqlW)
	fmt.Println("delete model data to database sql: ", sqlD)

	var values []interface{}

	values = append(values, valuesC...)
	values = append(values, valuesW...)

	return q.ExecEnd(sqlD, values...)
}

func (q *queryPrepare) Delete() (int64, error) {
	ind := getInd(q.md)

	tableName := ind.Type().Name()

	var values []interface{}
	var sqlW string
	pkKey := ind.Type().Field(0).Name
	pkValue := ind.Field(0).Interface()
	values = append(values, pkValue)

	sqlW = tableName + "." + pkKey + "=?"

	if len(q.filters) > 0 {
		sqlW = ""
		values = make([]interface{}, 0)
		for i, f := range q.filters {
			if i == 0 {
				sqlW += tableName + "." + f.Col + getOpType(f.Op)
			} else {
				if f.AndOr {
					sqlW += " AND " + tableName + "." + f.Col + getOpType(f.Op)
				} else {
					sqlW += " OR " + tableName + "." + f.Col + getOpType(f.Op)
				}
			}

			values = append(values, f.Value)
		}
	}

	sqlD := fmt.Sprintf("DELETE FROM %s WHERE %s", tableName, sqlW)
	fmt.Println("delete model data to database sql: ", sqlD)
	fmt.Println("delete model data to database values: ", values)

	return q.ExecEnd(sqlD, values...)
}

func (q *queryPrepare) Select(mds interface{}) error {

	sInd := reflect.Indirect(reflect.ValueOf(mds))

	switch sInd.Kind() {
	case reflect.Struct:
	default:
		return errors.New("[queryPrepare.Select] the models which is going to query data from database is not a ptr of struct ")
	}

	if sInd.NumField() != 1 {
		return errors.New("[queryPrepare.Select] the models which is going to query data from database contains fields not just one or none ")
	}

	elem := sInd.Field(0).Type().Elem()
	tableName := elem.Name()

	var sqlT, sqlQ, sqlL string
	var valuesW []interface{}

	sqlW := "1=?"
	sqlO := q.orderBy
	sqlG := q.groupBy
	valuesW = append(valuesW, 1)

	if len(q.limit) == 2 {
		sqlL = " limit " + strconv.Itoa(q.limit[0]) + "," + strconv.Itoa(q.limit[1])
	}

	if len(q.filters) > 0 {
		sqlW = ""
		valuesW = make([]interface{}, 0)
		for i, f := range q.filters {
			if i == 0 {
				sqlW += tableName + "." + f.Col + getOpType(f.Op)
			} else {
				if f.AndOr {
					sqlW += " AND " + tableName + "." + f.Col + getOpType(f.Op)
				} else {
					sqlW += " OR " + tableName + "." + f.Col + getOpType(f.Op)
				}
			}

			valuesW = append(valuesW, f.Value)
		}

	}

	for i := 0; i < elem.NumField(); i++ {
		if i == 0 {
			sqlT += tableName + "." + elem.Field(i).Name
		} else {
			sqlT += "," + tableName + "." + elem.Field(i).Name
		}
	}

	sqlQ = fmt.Sprintf("SELECT %s FROM %s WHERE %s %s %s %s", sqlT, tableName, sqlW, sqlL, sqlG, sqlO)
	fmt.Println("select model data from database, sql: ", sqlQ)

	return q.QueryEnd(mds, sqlQ, valuesW)
}

func (q *queryPrepare) Limit(limit ...int) QueryPrepareI {
	if len(limit) == 2 {
		q.limit = limit
	}

	return q
}

func (q *queryPrepare) Count() (int64, error) {
	ind := getInd(q.md)
	tableName := ind.Type().Name()
	pkKey := ind.Type().Field(0).Name

	sqlT := tableName + "." + pkKey
	sqlW := "1=?"
	var valuesW []interface{}
	valuesW = append(valuesW, 1)

	if len(q.filters) > 0 {
		sqlW = ""
		valuesW = make([]interface{}, 0)
		for i, f := range q.filters {
			if i == 0 {
				sqlW += tableName + "." + f.Col + getOpType(f.Op)
			} else {
				if f.AndOr {
					sqlW += " AND " + tableName + "." + f.Col + getOpType(f.Op)
				} else {
					sqlW += " OR " + tableName + "." + f.Col + getOpType(f.Op)
				}
			}

			valuesW = append(valuesW, f.Value)
		}

	}

	sqlC := fmt.Sprintf("SELECT COUNT(%s) FROM %s WHERE %s ", sqlT, tableName, sqlW)
	fmt.Println("count data, sql: ", sqlC)

	var num int64
	row := q.db.QueryRow(sqlC, valuesW...)
	if err := row.Scan(&num); err != nil {
		return 0, err
	}

	return num, nil
}

func (q *queryPrepare) OrderBy(col string) QueryPrepareI {
	q.orderBy = col

	return q
}

func (q *queryPrepare) GroupBy(col string) QueryPrepareI {
	q.groupBy = col

	return q
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

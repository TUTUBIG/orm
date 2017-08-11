package orm

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
)

const (
	dbDriver           = "mysql"
	dbCharset          = "utf8"
	maxOpenConnections = 2000
	maxIdleConnections = 1000
)

type uzOrm struct {
	db *sql.DB
}

func NewOrm() UzOrmI {
	o := new(uzOrm)

	return o
}

// RegisterDataBase Setting the database connect params. Use the database driver self dataSource args.
func (o *uzOrm) RegisterDataBase(host, user, password, dbName string) {

	dataSource := user + ":" + password + "@tcp(" + host + ")/" + dbName + "?charset=" + dbCharset

	db, err := sql.Open(dbDriver, dataSource)
	if err != nil {
		panic(fmt.Errorf("[RegisterDataBase] register Db error: %s", err.Error()))
	}

	db.SetMaxOpenConns(maxOpenConnections)
	db.SetMaxIdleConns(maxIdleConnections)
	db.Ping()

	o.db = db

	return
}

// insert model data to database
func (o *uzOrm) Insert(md interface{}) (int64, error) {
	ind := getInd(md)

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

	res, err := o.db.Exec(sqlI, values...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()

}

// update model data to database
func (o *uzOrm) Update(md interface{}, cols ...string) (int64, error) {
	ind := getInd(md)

	tableName := ind.Type().Name()

	var sqlT, sqlU string
	var values []interface{}

	// the first field is primary key by default
	pkKey := ind.Type().Field(0).Name
	pkValue := ind.Field(0).Interface()

	if len(cols) > 0 {
		for i, c := range cols {
			if i == 0 {
				sqlT += tableName + "." + c + "=?"
			} else {
				sqlT += "," + tableName + "." + c + "=?"
			}

			values = append(values, ind.FieldByName(c).Interface())
		}
	} else {
		for i := 0; i < ind.NumField(); i++ {
			if i == 0 {
				sqlT += tableName + "." + ind.Type().Field(i).Name + "=?"
			} else {
				sqlT += "," + tableName + "." + ind.Type().Field(i).Name + "=?"
			}

			values = append(values, ind.Field(i).Interface())
		}
	}

	values = append(values, pkValue)

	sqlU = fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?", tableName, sqlT, pkKey)
	fmt.Println("update model data to database sql: ", sqlU)

	res, err := o.db.Exec(sqlU, values...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()

}

// delete model data which is in the database
func (o *uzOrm) Delete(md interface{}) (int64, error) {
	ind := getInd(md)

	tableName := ind.Type().Name()

	// the first field is primary key by default
	pkKey := ind.Type().Field(0).Name
	pkValue := ind.Field(0).Interface()

	sqlD := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", tableName, pkKey)
	fmt.Println("delete model data to database sql: ", sqlD)

	res, err := o.db.Exec(sqlD, pkValue)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()

}

// query model data from the database
func (o *uzOrm) Select(md interface{}) error {

	defer o.db.Close()

	ind := getInd(md)
	tableName := ind.Type().Name()

	var sqlT, sqlQ string
	var fields []interface{}

	for i := 0; i < ind.NumField(); i++ {
		if i == 0 {
			sqlT += tableName + "." + ind.Type().Field(i).Name
		} else {
			sqlT += "," + tableName + "." + ind.Type().Field(i).Name
		}

		fields = append(fields, ind.Field(i).Addr().Interface())
	}

	pkKey := ind.Type().Field(0).Name
	pkValue := ind.Field(0).Interface()

	sqlQ = fmt.Sprintf("SELECT %s FROM %s WHERE %s = ? ", sqlT, tableName, pkKey)
	fmt.Println("select model data from database, sql: ", sqlQ)

	row := o.db.QueryRow(sqlQ, pkValue)
	if err := row.Scan(fields...); err != nil {
		return err
	}

	return nil
}

// advanced query
func (o *uzOrm) QueryTable(md interface{}) QueryPrepareI {
	val := reflect.ValueOf(md)
	sInd := reflect.Indirect(val)

	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		if sInd.Len() == 0 {
			panic(fmt.Errorf("[uzOrm.QueryTable] the models which is going to be inserted is empty "))
		}

		var temp interface{}
		temp = sInd.Index(0).Addr().Interface()
		md = temp

	case reflect.Ptr:
	default:
		panic(fmt.Errorf("[uzOrm.QueryTable] %s is not support ", sInd.Kind()))
	}

	q := NewQueryPrepare(md, o.db)
	return q
}

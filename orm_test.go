package orm

import (
	"fmt"
	"testing"
)

type Test_model struct {
	Id         int64
	Name       string
	Age        int
	Is_student bool
	Mark       float32
}

func TestUzOrm_Select(t *testing.T) {
	test := Test_model{
		Id: 1,
	}

	o := NewOrm()

	o.RegisterDataBase("10.6.64.68:3308", "dbnsmobilegame", "K8ZZPjwz9ySAOAwu", "test")

	err := o.Select(&test)
	fmt.Println(err)

	fmt.Println(test)

	return
}

func TestUzOrm_Insert(t *testing.T) {
	test := Test_model{
		Id:         2,
		Name:       "kate",
		Age:        23,
		Is_student: true,
		Mark:       89.90,
	}

	o := NewOrm()

	o.RegisterDataBase("10.6.64.68:3308", "dbnsmobilegame", "K8ZZPjwz9ySAOAwu", "test")

	i, err := o.Insert(&test)
	fmt.Println(err)

	fmt.Println(i)

	return
}

func TestUzOrm_Update(t *testing.T) {
	test := Test_model{
		Id:         2,
		Name:       "susan",
		Age:        23,
		Is_student: true,
		Mark:       69.90,
	}

	o := NewOrm()

	o.RegisterDataBase("10.6.64.68:3308", "dbnsmobilegame", "K8ZZPjwz9ySAOAwu", "test")

	i, err := o.Update(&test)
	fmt.Println(err)

	fmt.Println(i)

	return
}

func TestUzOrm_Delete(t *testing.T) {
	test := Test_model{
		Id:         2,
		Name:       "susan",
		Age:        23,
		Is_student: true,
		Mark:       69.90,
	}

	o := NewOrm()

	o.RegisterDataBase("10.6.64.68:3308", "dbnsmobilegame", "K8ZZPjwz9ySAOAwu", "test")

	i, err := o.Delete(&test)
	fmt.Println(err)

	fmt.Println(i)

	return
}

func TestQueryPrepare_Insert(t *testing.T) {
	tests := []Test_model{
		{
			Id:         1,
			Name:       "susan",
			Age:        23,
			Is_student: true,
			Mark:       69.90,
		},
		{
			Id:         2,
			Name:       "katy",
			Age:        23,
			Is_student: true,
			Mark:       70.00,
		},
		{
			Id:         3,
			Name:       "lucy",
			Age:        33,
			Is_student: false,
			Mark:       00.00,
		},
	}

	o := NewOrm()

	o.RegisterDataBase("10.6.64.68:3308", "dbnsmobilegame", "K8ZZPjwz9ySAOAwu", "test")

	i, err := o.QueryTable(tests).Insert(tests)
	fmt.Println(err)

	fmt.Println(i)

	return
}

func TestQueryPrepare_Update(t *testing.T) {
	test := Test_model{
		Id:         2,
		Name:       "alvin",
		Age:        24,
		Is_student: true,
		Mark:       69.90,
	}

	o := NewOrm()

	o.RegisterDataBase("10.6.64.68:3308", "dbnsmobilegame", "K8ZZPjwz9ySAOAwu", "test")

	i, err := o.QueryTable(&test).And("Id", Exact, 1).Update("Name", "Age")
	fmt.Println(err)

	fmt.Println(i)

	return
}

func TestQueryPrepare_Delete(t *testing.T) {
	test := Test_model{}

	o := NewOrm()

	o.RegisterDataBase("10.6.64.68:3308", "dbnsmobilegame", "K8ZZPjwz9ySAOAwu", "test")

	i, err := o.QueryTable(&test).And("Id", Exact, 1).Delete()
	fmt.Println(err)

	fmt.Println(i)

	return
}

func TestQueryPrepare_Select(t *testing.T) {
	type Temp struct {
		Test []Test_model
	}

	var temp Temp

	o := NewOrm()

	o.RegisterDataBase("10.6.64.68:3308", "dbnsmobilegame", "K8ZZPjwz9ySAOAwu", "test")

	if err := o.QueryTable(&temp).And("Age", Exact, 23).Select(&temp); err != nil {
		fmt.Println(err)
	}

	fmt.Println("temp: ", temp)

	return
}

func TestQueryPrepare_Count(t *testing.T) {
	test := Test_model{}

	o := NewOrm()

	o.RegisterDataBase("10.6.64.68:3308", "dbnsmobilegame", "K8ZZPjwz9ySAOAwu", "test")

	i, err := o.QueryTable(&test).And("Age", Exact, 23).Count()

	fmt.Println(i, err)

	return
}

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

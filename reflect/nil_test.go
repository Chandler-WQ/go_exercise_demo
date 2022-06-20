package reflect_test

import (
	"fmt"
	"reflect"
	"testing"
)

type A struct {
	B []string
}

func TestZero(t *testing.T) {
	a := &A{
		B: []string{"1", "2"},
	}
	v := reflect.Indirect(reflect.ValueOf(a))
	v = v.FieldByName("B")
	v.Set(reflect.Zero(v.Type()))
	fmt.Println(a)
	a = &A{
		B: []string{},
	}
	fmt.Println(a.B == nil)
	c := &A{}
	fmt.Println(c.B == nil)
}

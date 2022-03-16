package chores

import "reflect"

var a = 10

func Base() {
	va := reflect.ValueOf(a)
	ta := reflect.TypeOf(a)

}

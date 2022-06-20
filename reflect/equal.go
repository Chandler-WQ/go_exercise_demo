package reflect

import (
	"fmt"
	"go/ast"
	"reflect"
	"time"

	"github.com/Chandler-WQ/go_exercise_demo/pkg/log"
)

func Equal(got, expect interface{}) bool {
	if !reflect.DeepEqual(got, expect) {
		isEqual := func() bool {
			if curTime, ok := got.(time.Time); ok {
				format := "2006-01-02T15:04:05Z07:00"

				if curTime.Round(time.Second).UTC().Format(format) != expect.(time.Time).Round(time.Second).UTC().Format(format) &&
					curTime.Truncate(time.Second).UTC().Format(format) != expect.(time.Time).Truncate(time.Second).UTC().Format(format) {
					log.Errorf("expect: %v, got %v after time round", expect.(time.Time), curTime)

					return false
				}
			} else if fmt.Sprint(got) != fmt.Sprint(expect) {
				log.Errorf("%expect: %#v, got %#v", expect, got)
				return false
			}
			//这里值得商榷
			return true
		}

		if fmt.Sprint(got) == fmt.Sprint(expect) {
			return true
		}

		if reflect.Indirect(reflect.ValueOf(got)).IsValid() != reflect.Indirect(reflect.ValueOf(expect)).IsValid() {
			log.Errorf("expect: %+v, got %+v", expect, got)
			return false
		}

		if got != nil {
			//Interface()获取v的当前值作为接口{}
			got = reflect.Indirect(reflect.ValueOf(got)).Interface()
		}

		if expect != nil {
			expect = reflect.Indirect(reflect.ValueOf(expect)).Interface()
		}

		if reflect.ValueOf(got).IsValid() != reflect.ValueOf(expect).IsValid() {
			log.Errorf(" expect: %+v, got %+v", expect, got)
			return false
		}

		//kind可以为struct，type为struct的名字 数据类型可以有无限种，但是它们可以被分为几种类型。这个就是reflect.Kind.
		if reflect.ValueOf(got).Kind() == reflect.Slice {
			if reflect.ValueOf(expect).Kind() == reflect.Slice {
				if reflect.ValueOf(got).Len() == reflect.ValueOf(expect).Len() {
					for i := 0; i < reflect.ValueOf(got).Len(); i++ {
						// name := fmt.Sprintf(reflect.ValueOf(got).Type().Name()+" #%v", i)
						return Equal(reflect.ValueOf(got).Index(i).Interface(), reflect.ValueOf(expect).Index(i).Interface())
					}
				} else {
					name := reflect.ValueOf(got).Type().Elem().Name()
					log.Errorf("%v expects length: %v, got %v (expects: %+v, got %+v)",
						name, reflect.ValueOf(expect).Len(), reflect.ValueOf(got).Len(), expect, got)
					return false
				}
			}
		}

		if reflect.ValueOf(got).Kind() == reflect.Struct {
			if reflect.ValueOf(expect).Kind() == reflect.Struct {
				if reflect.ValueOf(got).NumField() == reflect.ValueOf(expect).NumField() {
					exported := false
					for i := 0; i < reflect.ValueOf(got).NumField(); i++ {
						if fieldStruct := reflect.ValueOf(got).Type().Field(i); ast.IsExported(fieldStruct.Name) {
							exported = true
							field := reflect.ValueOf(got).Field(i)
							Equal(field.Interface(), reflect.ValueOf(expect).Field(i).Interface())
						}
					}

					if exported {
						return true
					}
				}
			}
		}

		if reflect.ValueOf(got).Type().ConvertibleTo(reflect.ValueOf(expect).Type()) {
			got = reflect.ValueOf(got).Convert(reflect.ValueOf(expect).Type()).Interface()
			isEqual()
		} else if reflect.ValueOf(expect).Type().ConvertibleTo(reflect.ValueOf(got).Type()) {
			expect = reflect.ValueOf(got).Convert(reflect.ValueOf(got).Type()).Interface()
			isEqual()
		} else {
			log.Errorf("expect: %+v, got %+v", expect, got)
			return false
		}
	}
	return false
}

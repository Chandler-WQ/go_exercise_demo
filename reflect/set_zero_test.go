package reflect_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

type User struct {
	Id   int
	Name string
}

func ChangeSlice(s interface{}) {
	sT := reflect.TypeOf(s)
	fmt.Println(sT.Kind())
	if sT.Kind() != reflect.Ptr {
		fmt.Println("参数必须是ptr类型")
		os.Exit(-1)
	}
	sV := reflect.ValueOf(s)
	// 取得数组中元素的类型
	sEE := sT.Elem().Elem()
	// 数组的值
	sVE := sV.Elem()
	//fmt.Println(sEE)
	//fmt.Println(sVE)

	// new一个数组中的元素对象
	sON := reflect.New(sEE)
	// 对象的值
	sONE := sON.Elem()
	// 给对象复制
	sONEId := sONE.FieldByName("Id")
	sONEName := sONE.FieldByName("Name")
	sONEId.SetInt(10)
	sONEName.SetString("李四")

	// 创建一个新数组并把元素的值追加进去
	newArr := make([]reflect.Value, 0)
	newArr = append(newArr, sON.Elem())

	// 把原数组的值和新的数组合并
	resArr := reflect.Append(sVE, newArr...)

	// 最终结果给原数组
	sVE.Set(resArr)
}

func TestEraseList(t *testing.T) {
	type A struct {
		Id int16
	}
	As := []A{{1}, {2}, {3}, {4}}
	v := reflect.ValueOf(&As)
	v1 := reflect.Indirect(v)
	v2 := v1.Slice(0, 1)
	v3 := v1.Slice(2, v1.Len())
	v4 := reflect.AppendSlice(v2, v3)
	v.Elem().Set(v4)
	t.Logf("%v", As)
	// &[{1} {3} {4}]
}

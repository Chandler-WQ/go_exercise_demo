package reflect

import (
	"reflect"
	"testing"
)

type A struct {
	B *B
}
type B struct {
	C *C
}

type C struct {
	Num int
}

func BenchmarkField(b *testing.B) {
	// feilds := [3]int64{1, 1, 1}

	b.Run("dfs", func(b *testing.B) {
		for x := 0; x < b.N; x++ {
			c1 := A{B: &B{C: &C{}}}
			a := reflect.ValueOf(c1)
			a = reflect.Indirect(a)
			b := a.FieldByName("B")
			b = reflect.Indirect(b)
			c := b.FieldByName("C")
			c = reflect.Indirect(c)
			d := c.FieldByName("Num")
			d.Set(reflect.ValueOf(10))
			// fmt.Println(c1.B.C.Num)
		}
	})

	b.Run("index", func(b *testing.B) {
		for x := 0; x < b.N; x++ {
			c1 := A{B: &B{C: &C{}}}
			a := reflect.ValueOf(c1)
			a = reflect.Indirect(a)
			b := a.Field(0)
			b = reflect.Indirect(b)
			c := b.Field(0)
			c = reflect.Indirect(c)
			d := c.Field(0)
			d.Set(reflect.ValueOf(10))
			// fmt.Println(c1.B.C.Num)
		}
	})

}

/*
goos: darwin
goarch: amd64
pkg: github.com/Chandler-WQ/go_exercise_demo/reflect
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkField/dfs-12            4747224               254.1 ns/op            40 B/op          5 allocs/op
BenchmarkField/index-12         19140554                62.91 ns/op           16 B/op          2 allocs/op
PASS
ok      github.com/Chandler-WQ/go_exercise_demo/reflect 2.768s
*/

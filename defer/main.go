package main

import (
	"errors"
	"fmt"
)

func main() {
	ErrTest()
}

func ErrTest() (err error) {
	defer func() {
		fmt.Println(err)
	}()
	a, err := 1, errors.New("test")
	_ = a
	return err
}

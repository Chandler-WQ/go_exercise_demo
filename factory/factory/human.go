package factory

import (
	"fmt"
)

/*
原始工厂模式
这里以女娲造人为例
*/

//简单工厂的模式的话就不需要这个接口了
//定义人类接口
type Human interface {
	Eat()
	Play()
}

//实现黑人y
type yellowHuman struct {
}

func (man *yellowHuman) Eat() {
	fmt.Println("Yellow man eat")
}

func (man *yellowHuman) Play() {
	fmt.Println("Yellow man play")
}

//实现白人
type whiteHuman struct {
}

func (man *whiteHuman) Eat() {
	fmt.Println("White man eat")
}

func (man *whiteHuman) Play() {
	fmt.Println("White man play")
}

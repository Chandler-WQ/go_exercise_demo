package main

import (
	"fmt"

	"github.com/Chandler-WQ/go_exercise_demo/factory/factory"
)

func main() {
	humanFactory := factory.HumanFactory{}
	err, hu := humanFactory.Create("white")
	if err != nil {
		fmt.Println(err)
	} else {
		hu.Eat()
		hu.Play()
	}

	err, hu = humanFactory.Create("yellow")
	if err != nil {
		fmt.Println(err)
	} else {
		hu.Eat()
		hu.Play()
	}

	err, hu = humanFactory.Create("hhh")
	if err != nil {
		fmt.Println(err)
	} else {
		hu.Eat()
		hu.Play()
	}
}

/*
White man eat
White man play
Yellow man eat
Yellow man play
The color is not currently supported
*/

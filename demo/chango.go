package main

import (
	"fmt"
	"time"
)

type hh interface {
	add()
}

type a struct {
	d int
	e int
}

func (a *a) add() {

}

func main() {
	// b := sync.WaitGroup{}
	// b.Add(1)
	hhchan := make(chan hh, 10)
	go pro(hhchan)
	go con(hhchan)
	time.Sleep(time.Minute)

}

func pro(b chan hh) {
	i := 0
	for {
		i++
		time.Sleep(1 * time.Second)
		b <- &a{
			i, 2,
		}
	}
}

func con(b chan hh) {
	hhs := make([]hh, 0, 10)
	for {
		if len(hhs) == 5 {
			go func(ha []hh) {
				time.Sleep(20 * time.Second)
				fmt.Printf("go ha is %v\n", ha[0])
			}(hhs)
			hhs = make([]hh, 0, 10)
		}
		c := <-b
		hhs = append(hhs, c)
	}
}

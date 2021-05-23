package singleton1

/*
带数量限制的单例模式
*/

import (
	"math/rand"
	"sync"
)

type singleton struct {
}

var once sync.Once

var instances [2]*singleton
var numOfnInstances = 2 //数量限制

//获取现有的单例列表中的某一个
func GetInstance() *singleton {
	once.Do(func() {
		for i := 0; i < numOfnInstances; i++ {
			instances[i] = &singleton{}
		}
	})
	i := rand.Intn(numOfnInstances)
	return instances[i]
}

package factory

import (
	"errors"
)

var _ HumanFactoryInf = (*HumanFactory)(nil) //可以检测是否某一个接口被这个实现类所实现

type HumanFactoryInf interface {
	Create(color string) (error, Human)
}

type HumanFactory struct {
}

func (factory *HumanFactory) Create(color string) (error, Human) {
	switch color {
	case "white":
		return nil, new(whiteHuman)
	case "yellow":
		return nil, new(yellowHuman)
	default:
		return errors.New("The color is not currently supported"), nil
	}
}

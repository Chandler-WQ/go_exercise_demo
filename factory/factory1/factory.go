package factory1


//升级为多个工厂，这样的话扩展就
var _ HumanFactoryInf = (*WhiteHumanFactory)(nil) //可以检测是否某一个接口被这个实现类所实现

type HumanFactoryInf interface {
	Create() Human
}

type WhiteHumanFactory struct {
}

func (factory *WhiteHumanFactory) Create() Human {
	return new(whiteHuman)
}

type YellowHumanFactory struct {
}

func (factory *YellowHumanFactory) Create() Human {
	return new(yellowHuman)
}

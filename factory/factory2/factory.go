package factory2

//抽象工厂
var _ HumanFactoryInf = (*HumanFactory)(nil) //可以检测是否某一个接口被这个实现类所实现

//缺点：当需要增加创建black时则代码改动很大
//优点：可以控制创建产品的关系，我们可以将CreatWhite和CreatYellow内置，然后对外提供Create限制CreatWhite和CreatYellow关系
type HumanFactoryInf interface {
	CreateWhite() Human
	CreateYellow() Human
}

type HumanFactory struct {
}

func (factory *HumanFactory) CreateWhite() Human {
	return new(whiteHuman)
}

func (factory *HumanFactory) CreateYellow() Human {
	return new(yellowHuman)
}

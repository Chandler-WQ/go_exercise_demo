package main

import (
	"container/list"
	"errors"
	"fmt"
	"math"
	"strings"
)

//题目： 两点之间过指定路径的最短路径

type Node struct {
	Id      int64           //编号单调递增from 0
	IsDepot bool            //是否是维修点
	EdgeOut map[int64]*Edge //node的出度 ,key为出向nodeid
	Extra   interface{}
}

func (node *Node) String() string {
	return fmt.Sprintf("id :%v,depot %v", node.Id, node.IsDepot)
}

type Graph struct {
	Nodes   map[int64]*Node
	Edges   map[int64]*Edge
	DepotS  map[int64]*Node //停留维修车的地址
	MarkSeg int64
}

func (gra *Graph) String() string {
	str := strings.Builder{}
	str.WriteString("nodes:\n")
	for _, node := range gra.Nodes {
		str.WriteString(fmt.Sprintf("node:%s\n", node))
	}
	str.WriteString("edges:\n")
	for _, edge := range gra.Edges {
		str.WriteString(fmt.Sprintf("edge:%s\n", edge))
	}
	return str.String()
}

type RoadGra struct {
	Graph         *Graph
	FixCars       map[int64]*FixCar    //key为编号
	FixCarsByType map[string][]*FixCar //key为维修车类型，value为carIdlist
}

func (road *RoadGra) String() string {
	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("graph:%s\n", road.Graph))
	str.WriteString(fmt.Sprintf("cars:\n"))
	for _, car := range road.FixCars {
		str.WriteString(fmt.Sprintf("%s\n", car))
	}
	return str.String()
}

type Edge struct {
	Id          int64 //编号单调递增from 0
	ReverseEdge *Edge //相反的Edge
	Caller      *Node
	Callee      *Node
	Length      int64
	BrokenMark  int64
}

func (edge *Edge) String() string {
	return fmt.Sprintf("id:%v,caller:%v,callee:%v,length:%v", edge.Id, edge.Caller.Id, edge.Callee.Id, edge.Length)
}

type FixCar struct {
	Id      int64
	FixType string
	Depot   *Node //初始化停止Id
}

func (car *FixCar) String() string {
	return fmt.Sprintf("id:%v,fix_type:%v,node_id:%v", car.Id, car.FixType, car.Depot.Id)
}

func NewGraph() *Graph {
	return &Graph{
		Nodes:  make(map[int64]*Node),
		Edges:  make(map[int64]*Edge),
		DepotS: make(map[int64]*Node),
	}
}

func NewNode(id int64) *Node {
	return &Node{
		Id:      id,
		EdgeOut: make(map[int64]*Edge, 10),
	}
}

func NewEdge(id int64, Caller, Callee *Node, Length int64) *Edge {
	return &Edge{
		Id:     id,
		Caller: Caller,
		Callee: Callee,
		Length: Length,
	}
}

func NewRoad() *RoadGra {
	return &RoadGra{
		Graph:         NewGraph(),
		FixCars:       make(map[int64]*FixCar),
		FixCarsByType: make(map[string][]*FixCar),
	}
}

func (road *RoadGra) AddNewFixCar(NodeId int64, FixCarType string) error {
	depot, ok := road.Graph.Nodes[NodeId]
	if !ok {
		return errors.New("the nodeid is invalid")
	}
	id := int64(len(road.FixCars))
	depot.IsDepot = true
	car := &FixCar{
		Id:      id,
		Depot:   depot,
		FixType: FixCarType,
	}
	road.AddFixCar(car)
	return nil
}

func (road *RoadGra) AddFixCar(fixCar *FixCar) {
	if road.FixCars == nil {
		road.FixCars = make(map[int64]*FixCar)
	}
	if road.FixCarsByType == nil {
		road.FixCarsByType = make(map[string][]*FixCar)
	}
	road.FixCars[fixCar.Id] = fixCar

	cars := road.FixCarsByType[fixCar.FixType]
	cars = append(cars, fixCar)
	road.FixCarsByType[fixCar.FixType] = cars

	road.Graph.AddNode(fixCar.Depot, true)
}

func (gra *Graph) AddNode(node *Node, IsDepot bool) {
	if gra.Nodes == nil {
		gra.Nodes = make(map[int64]*Node)
	}

	gra.Nodes[node.Id] = node
	if IsDepot {
		if gra.DepotS == nil {
			gra.DepotS = make(map[int64]*Node)
		}
		gra.DepotS[node.Id] = node
	}
}

func (gra *Graph) AddBroken(nodeId1, nodeId2 int64) error {
	node1, ok1 := gra.Nodes[nodeId1]
	node2, ok2 := gra.Nodes[nodeId2]
	if !ok1 || !ok2 {
		return errors.New("the node id is invalid.")
	}
	mark := 1 << gra.MarkSeg
	gra.MarkSeg++
	if _, ok := node1.EdgeOut[nodeId2]; ok {
		node1.EdgeOut[nodeId2].BrokenMark = int64(mark)
	}
	if _, ok := node2.EdgeOut[nodeId1]; ok {
		node2.EdgeOut[nodeId1].BrokenMark = int64(mark)
	}

	return nil
}

func (gra *Graph) AddNewEdge(nodeid1, nodeId2 int64, length int64) error {
	id := int64(len(gra.Edges))
	node1, ok1 := gra.Nodes[nodeid1]
	node2, ok2 := gra.Nodes[nodeId2]
	if !ok1 || !ok2 {
		return errors.New("the node id is invalid.")
	}
	if _, ok := node1.EdgeOut[nodeId2]; ok {
		return nil
	}
	edge1 := NewEdge(id, node1, node2, length)
	edge2 := NewEdge(id+1, node2, node1, length)
	if gra.Edges == nil {
		gra.Edges = make(map[int64]*Edge)
	}
	edge1.ReverseEdge = edge2
	edge2.ReverseEdge = edge1
	gra.Edges[edge1.Id] = edge1
	gra.Edges[edge2.Id] = edge2
	node1.addEdgeOut(edge1)
	node2.addEdgeOut(edge2)
	return nil
}

func (node *Node) addEdgeOut(out *Edge) {
	if node.EdgeOut == nil {
		node.EdgeOut = make(map[int64]*Edge, 10)
	}

	node.EdgeOut[out.Callee.Id] = out
}

type Path struct {
	gra    *Graph
	start  *Node
	end    *Node
	list   []*Node
	length int64
}

func (pa *Path) String() string {
	str := strings.Builder{}
	str.WriteString(fmt.Sprintf("points:%v....->", pa.start.Id))
	str.WriteString(fmt.Sprintf("%v \t", pa.end.Id))
	str.WriteString("details:")
	for i := 0; i < len(pa.list)-1; i++ {
		str.WriteString(fmt.Sprintf("%v->", pa.list[i].Id))
	}
	if len(pa.list) > 0 {
		str.WriteString(fmt.Sprintf("%v", pa.end.Id))
	}
	str.WriteString(fmt.Sprintf("\tlength:%v", pa.length))

	return str.String()
}

type State map[int64]map[int64]*Path //dist 过的点量，路径

func (state State) String() string {
	str := strings.Builder{}
	for k, paths := range state {
		str.WriteString(fmt.Sprintf("dist:%v", k))
		for k, path := range paths {
			str.WriteString(fmt.Sprintf("\t broken:%010b %s\n", k, path))
		}
	}
	return str.String()
}

var states State
var queue *list.List

func initPath(gra *Graph, startId int64) (State, *list.List, map[int64]bool, error) {
	start, ok := gra.Nodes[startId]
	if !ok {
		return nil, nil, nil, errors.New("the startId is invalid")
	}

	inQueue := make(map[int64]bool, len(gra.Nodes))
	states := make(State)
	for _, node := range gra.Nodes {
		path := &Path{
			gra:    gra,
			start:  start,
			end:    node,
			list:   make([]*Node, 0, 10),
			length: math.MaxInt64,
		}
		if node.Id == start.Id {
			path.length = 0
			path.list = append(path.list, start)
		}
		states[node.Id] = make(map[int64]*Path)
		states[node.Id] = map[int64]*Path{
			0: path,
		}
	}

	queue := list.New()
	queue.PushBack(start)
	inQueue[startId] = true
	return states, queue, inQueue, nil
}

func findPath(gra *Graph, startId int64) (State, error) {
	states, queue, inQueue, err := initPath(gra, startId)
	if err != nil {
		fmt.Printf("initpath err %v", err)
		return nil, err
	}

	for queue.Len() != 0 {
		element := queue.Front()
		node := element.Value.(*Node)
		queue.Remove(element)
		inQueue[node.Id] = false
		paths := states[node.Id]
		// fmt.Printf("%s\n\n", states)
		for _, edge := range node.EdgeOut {
			outId := edge.Callee.Id
			calleePaths := states[outId]
			update := false
			for k, path := range paths {
				if path.length != math.MaxInt64 {
					newk := k | edge.BrokenMark
					if oldpath, ok := calleePaths[newk]; ok {
						_ = oldpath
						_ = path
						if oldpath.length > path.length+edge.Length {
							updatePath(gra, node, edge, states, k, newk)
							update = true
						}
					} else {
						addPath(gra, node, edge, states, k, newk)
						update = true
					}
				}
			}

			if update && !inQueue[outId] {
				queue.PushBack(edge.Callee)
				inQueue[outId] = true
			}
		}
	}
	return states, nil
}

func addPath(gra *Graph, caller *Node, edge *Edge, states State, oldmark, mark int64) {
	callerPaths := states[caller.Id]
	callerPath := callerPaths[oldmark]

	calleePaths := states[edge.Callee.Id]

	calleePath := &Path{}
	list := make([]*Node, 0, len(callerPath.list)+1)
	for _, node := range callerPath.list {
		list = append(list, node)
	}
	list = append(list, edge.Callee)
	calleePath.list = list
	calleePath.gra = gra
	calleePath.start = callerPath.start
	calleePath.end = edge.Callee
	calleePath.length = callerPath.length + edge.Length

	calleePaths[mark] = calleePath
	states[edge.Callee.Id] = calleePaths
}

func updatePath(gra *Graph, caller *Node, edge *Edge, states State, oldmark, mark int64) {
	callerPaths := states[caller.Id]
	callerPath := callerPaths[oldmark]

	calleePaths := states[edge.Callee.Id]
	calleePath := calleePaths[mark]

	list := make([]*Node, 0, len(callerPath.list)+1)
	for _, node := range callerPath.list {
		list = append(list, node)
	}
	list = append(list, edge.Callee)
	calleePath.list = list
	calleePath.gra = gra
	calleePath.length = callerPath.length + edge.Length

	calleePaths[mark] = calleePath
	states[edge.Callee.Id] = calleePaths
}

var road = NewRoad()

func case1() {
	for i := 1; i < 14; i++ {
		road.Graph.AddNode(NewNode(int64(i)), false)
	}
	road.Graph.AddNewEdge(1, 2, 15)
	road.Graph.AddNewEdge(1, 3, 24)
	road.Graph.AddNewEdge(2, 5, 29)
	road.Graph.AddNewEdge(2, 6, 40)
	road.Graph.AddNewEdge(3, 4, 21)
	road.Graph.AddNewEdge(3, 7, 22)
	road.Graph.AddNewEdge(4, 7, 19)
	road.Graph.AddNewEdge(4, 5, 7)
	road.Graph.AddNewEdge(5, 8, 15)
	road.Graph.AddNewEdge(6, 9, 15)
	road.Graph.AddNewEdge(7, 10, 20)
	road.Graph.AddNewEdge(7, 11, 19)
	road.Graph.AddNewEdge(8, 9, 14)
	road.Graph.AddNewEdge(8, 11, 20)
	road.Graph.AddNewEdge(8, 13, 34)
	road.Graph.AddNewEdge(8, 13, 34)
	road.Graph.AddNewEdge(9, 13, 29)
	road.Graph.AddNewEdge(10, 11, 18)
	road.Graph.AddNewEdge(11, 12, 22)
	road.Graph.AddNewEdge(12, 13, 23)

	// road.Graph.AddBroken(2, 1)
	// road.Graph.AddBroken(8, 13)
	// road.Graph.AddBroken(7, 10)
	road.Graph.AddBroken(4, 7)
	road.Graph.AddBroken(5, 8)
	road.Graph.AddBroken(10, 11)
}

func main() {
	case1()
	//求点2出发的最短路径
	states, err := findPath(road.Graph, 2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s", states)

}

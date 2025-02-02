package main

import (
	"errors"
	"reflect"
)

type node struct {
	Parent   *node
	Position *Position
	g        int
	h        int
	f        int
}

func newNode(parent *node, position *Position) *node {
	n := node{}
	n.Parent = parent
	n.Position = position
	n.g = 0
	n.h = 0
	n.f = 0

	return &n
}

func (n *node) isEqual(other *node) bool {
	return n.Position.IsEqual(other.Position)
}

type AStar struct{}

func (as AStar) GetPath(level Level, start *Position, end *Position) []Position {
	gd := NewGameData()
	openList := make([]*node, 0)
	closedList := make([]*node, 0)

	startNode := newNode(nil, start)
	startNode.g = 0
	startNode.h = 0
	startNode.f = 0

	endNode := newNode(nil, end)

	openList = append(openList, startNode)

	for {
		if len(openList) == 0 {
			break
		}

		currentNode := openList[0]
		currentIndex := 0

		for index, item := range openList {
			if item.f < currentNode.f {
				currentNode = item
				currentIndex = index
			}
		}

		openList = append(openList[:currentIndex], openList[currentIndex+1:]...)
		closedList = append(closedList, currentNode)

		if currentNode.isEqual(endNode) {
			path := make([]Position, 0)
			current := currentNode
			for current != nil {
				path = append(path, *current.Position)
				current = current.Parent
			}

			reverseSlice(path)
			return path
		}

		edges := make([]*node, 0)

		if currentNode.Position.Y > 0 {
			tile := level.Tiles[level.GetIndexFromXY(currentNode.Position.X, currentNode.Position.Y-1)]
			if tile.TileType != WALL {
				upNodePosition := Position{X: currentNode.Position.X, Y: currentNode.Position.Y - 1}
				upNode := newNode(currentNode, &upNodePosition)
				edges = append(edges, upNode)
			}
		}

		if currentNode.Position.Y < gd.ScreenHeight {
			tile := level.Tiles[level.GetIndexFromXY(currentNode.Position.X, currentNode.Position.Y+1)]
			if tile.TileType != WALL {
				downNodePosition := Position{X: currentNode.Position.X, Y: currentNode.Position.Y + 1}
				downNode := newNode(currentNode, &downNodePosition)
				edges = append(edges, downNode)
			}
		}

		if currentNode.Position.X > 0 {
			tile := level.Tiles[level.GetIndexFromXY(currentNode.Position.X-1, currentNode.Position.Y)]
			if tile.TileType != WALL {
				leftNodePosition := Position{X: currentNode.Position.X - 1, Y: currentNode.Position.Y}
				leftNode := newNode(currentNode, &leftNodePosition)
				edges = append(edges, leftNode)
			}
		}

		if currentNode.Position.X < gd.ScreenWidth {
			tile := level.Tiles[level.GetIndexFromXY(currentNode.Position.X+1, currentNode.Position.Y)]
			if tile.TileType != WALL {
				rightNodePosition := Position{X: currentNode.Position.X + 1, Y: currentNode.Position.Y}
				rightNode := newNode(currentNode, &rightNodePosition)
				edges = append(edges, rightNode)
			}
		}

		for _, edge := range edges {
			if isInSlice(closedList, edge) {
				continue
			}
			edge.g = currentNode.g + 1
			edge.h = edge.Position.GetManhattanDistance(endNode.Position)
			edge.f = edge.g + edge.h

			if isInSlice(openList, edge) {
				isFurther := false
				for _, n := range openList {
					if edge.g > n.g {
						isFurther = true
						break
					}
				}
				if isFurther {
					continue
				}
			}
			openList = append(openList, edge)
		}
	}

	return nil
}

func isInSlice(s []*node, target *node) bool {
	for _, n := range s {
		if n.isEqual(target) {
			return true
		}
	}
	return false
}

func reverseSlice(data interface{}) {
	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Slice {
		panic(errors.New("data must be a slice type"))
	}
	valueLen := value.Len()
	for i := 0; i <= int((valueLen-1)/2); i++ {
		reverseIndex := valueLen - 1 - i
		tmp := value.Index(reverseIndex).Interface()
		value.Index(reverseIndex).Set(value.Index(i))
		value.Index(i).Set(reflect.ValueOf(tmp))
	}
}

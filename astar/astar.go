package astar

import (
	"math"
	"reflect"
)

type Node interface {
	GetNeighbors() []Node
	H(goal Node) int
	D(neighbor Node) int
}

type NodeSet map[Node]bool
type GScore map[Node]int
type FScore map[Node]int
type CameFrom map[Node]*Node

// https://en.wikipedia.org/wiki/A*_search_algorithm
//
// h(n) - Heuristic - estimates the cost to reach goal from node n
// d(n, n) - Is the weight of the edge from current to neighbor
func AStar(start Node, goal Node) ([]*Node, bool) {
	// The set of discovered nodes that may need to be (re-)expanded.
	// Initially, only the start node is known.
	// This is usually implemented as a min-heap or priority queue rather than a hash-set.
	openSet := NodeSet{
		start: true,
	}

	// For node n, cameFrom[n] is the node immediately preceding it on the cheapest path from start.
	cameFrom := CameFrom{}

	// For node n, gScore[n] is the cost of the cheapest path from start to n currently known.
	gScore := GScore{
		start: 0,
	}
	// For node n, fScore[n] := gScore[n] + h(n). fScore[n] represents our current best guess as to
	// how short a path from start to finish can be if it goes through n.
	fScore := FScore{
		start: start.H(goal),
	}

	for len(openSet) > 0 {
		current := lowestFScore(openSet, fScore)

		if reflect.DeepEqual(current, goal) {
			return reconstructPath(cameFrom, &current), true
		}

		delete(openSet, current)

		for _, neighbor := range current.GetNeighbors() {
			// tentative_gScore is the distance from start to the neighbor through current
			tentativeGScore := gScore[current] + current.D(neighbor)
			gScoreNeighbor, gScoreFound := gScore[neighbor]
			if tentativeGScore < gScoreNeighbor || !gScoreFound {
				// This path to neighbor is better than any previous one. Record it!
				cameFrom[neighbor] = &current
				gScore[neighbor] = tentativeGScore
				fScore[neighbor] = gScore[neighbor] + neighbor.H(goal)
				if _, ok := openSet[neighbor]; !ok {
					openSet[neighbor] = true
				}
			}
		}
	}

	return []*Node{}, false
}

func reconstructPath(cameFrom CameFrom, current *Node) []*Node {
	totalPath := []*Node{current}

	var temp *Node
	var ok bool

	temp, ok = cameFrom[*current]
	for ok {
		totalPath = append([]*Node{temp}, totalPath...)
		current = temp
		temp, ok = cameFrom[*current]
	}

	return totalPath
}

func lowestFScore(set NodeSet, fScore FScore) Node {
	lowestScore := math.MaxInt32
	var lowestNode Node
	for node := range set {
		score := fScore[node]
		if score < lowestScore {
			lowestScore = score
			lowestNode = node
		}
	}
	return lowestNode
}

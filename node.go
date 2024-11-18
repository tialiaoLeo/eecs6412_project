// node.go
package main

// Node represents a node in a graph
type Node struct {
	NodeID  string // Unique identifier for the node
	CoreNum int    // Core number of the node
}

func (n *Node) IncrementCoreNum() {
	n.CoreNum++
}

func (n *Node) k_core(graph map[string][]Edge, nodes map[string]*Node) {
	for n.CoreNum > n.num_neighbors(n.CoreNum, graph, nodes) {
		n.CoreNum -= 1
	}

}

func (n *Node) num_neighbors(CoreNum int, graph map[string][]Edge, nodes map[string]*Node) int {
	res := 0
	for _, value := range graph[n.NodeID] {
		if nodes[value.Target].CoreNum >= CoreNum {
			res += 1
		}

	}
	return res

}

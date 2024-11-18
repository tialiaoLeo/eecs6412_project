// node.go
package main

// Node represents a node in a graph
type Node struct {
	NodeID        string   // Unique identifier for the node
	NodeNeighbors []string // Slice of node IDs representing neighbors
	CoreNum       int      // Core number of the node
}

func (n *Node) IncrementCoreNum() {
	n.CoreNum++
}

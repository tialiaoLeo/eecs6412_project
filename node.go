// node.go
package main

import "sync"

// Node represents a node in a graph
type Node struct {
	NodeID        string   // Unique identifier for the node
	NodeNeighbors []string // Slice of node IDs representing neighbors
	CoreNum       int      // Core number of the node
	mutex         sync.Mutex
}

func (n *Node) IncrementCoreNum() {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.CoreNum++
}

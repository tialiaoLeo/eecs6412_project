// node.go
package main

// global variable for message id

// Node represents a node in a graph
type Node struct {
	NodeID    string // Unique identifier for the node
	CoreNum   int    // Core number of the node
	Msg_queue []Message
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

func (n *Node) publish(nodes map[string]*Node, global *Global) {
	for _, target := range nodes {
		n.Msg_queue = append(n.Msg_queue, Message{global.ms_id, n, target})
		global.ms_id = global.ms_id + 1
	}
}

func (n *Node) consume() bool {
	pub_again := len(n.Msg_queue) > 0
	n.Msg_queue = nil
	return pub_again
}

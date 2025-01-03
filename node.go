// node.go
package main

import (
	"fmt"
	"math"
	"time"
)

// global variable for message id

// Node represents a node in a graph
type Node struct {
	NodeID          string           // Unique identifier for the node
	CoreNum         int              // Core number of the node
	Msg_queue       []heart_beat_msg // heart beat msg
	secure_msg      []*secure_msg    // secure msg received
	k_core_msg      map[string]bool  // boolean array for k core
	terminate_times int
}

func (n *Node) k_core(graph map[string][]Edge, nodes map[string]*Node) {
	num_n := n.num_neighbors(graph, nodes)
	for n.CoreNum > num_n {
		n.CoreNum -= 1
		num_n = n.num_neighbors(graph, nodes)
	}

}

func (n *Node) num_neighbors(graph map[string][]Edge, nodes map[string]*Node) int {
	time.Sleep(100 * time.Millisecond)
	for _, value := range graph[n.NodeID] {
		n.send(nodes[value.Target])
	}
	count := 0
	for _, v := range n.k_core_msg {
		if v {
			count++
		}
	}
	if len(n.k_core_msg) < len(graph[n.NodeID]) {
		return math.MaxInt32
	}
	n.k_core_msg = make(map[string]bool)
	return count

}

// heart beat publish
func (n *Node) publish(nodes map[string]*Node, neighbors map[string][]Edge) {
	time.Sleep(300 * time.Millisecond)
	for _, neighbor := range neighbors[n.NodeID] {
		var n_node = nodes[neighbor.Target]
		n_node.Msg_queue = append(n_node.Msg_queue, heart_beat_msg{n, n_node})
	}
}

func (n *Node) terminate(nodes map[string]*Node, graph map[string][]Edge) bool {
	tms := 30 * (int(n.NodeID[0])%26 + 1)
	time.Sleep(time.Duration(tms) * time.Millisecond)
	fmt.Printf("Waiting for: %v \n", tms)
	var allNeighstop = true
	cur_len := len(graph[n.NodeID])
	for _, nn := range nodes {
		if len(nn.Msg_queue) > 0 || (cur_len < len(graph[nn.NodeID]) && nn.terminate_times == 0) {
			allNeighstop = false
			break
		}
	}
	return allNeighstop
}

// heart beat consume
func (n *Node) consume() bool {
	time.Sleep(300 * time.Millisecond)
	pub_again := len(n.Msg_queue) > 0
	n.Msg_queue = nil
	return pub_again
}

func (n *Node) send(to *Node) {
	time.Sleep(300 * time.Millisecond)
	for _, m := range to.secure_msg {
		if m.from.NodeID == n.NodeID {
			return
		}
	}
	cur_msg := Newsecure_msg(n.CoreNum, n, to)
	to.secure_msg = append(to.secure_msg, cur_msg)
	fmt.Printf("node %v send: %v to: %v \n", n.NodeID, cur_msg.from_v, to.NodeID)
}

func (n *Node) reply() {
	time.Sleep(300 * time.Millisecond)
	curMsg := n.secure_msg
	n.secure_msg = nil
	for _, m := range curMsg {
		toV := Encrypt(n.CoreNum, m)
		m.to_v = &toV
		if _, exists := m.from.k_core_msg[n.NodeID]; !exists {
			// Key does not exist, write it
			m.from.k_core_msg[n.NodeID] = m.compare()
		}
		fmt.Printf("node: %v reply with %v back to %v \n", n.NodeID, m.to_v, m.from.NodeID)
	}
}

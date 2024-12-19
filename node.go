// node.go
package main

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

// global variable for message id

// Node represents a node in a graph
type Node struct {
	NodeID      string           // Unique identifier for the node
	CoreNum     int              // Core number of the node
	Msg_queue   []heart_beat_msg // heart beat msg
	secure_msg  []secure_msg     // secure msg received
	k_core_msg  map[string]bool  // boolean array for k core
	private_key string
}

func (n *Node) k_core(graph map[string][]Edge, nodes map[string]*Node, global *Global) {
	num_n := n.num_neighbors(graph, nodes, global)
	for n.CoreNum > num_n {
		n.CoreNum -= 1
		num_n = n.num_neighbors(graph, nodes, global)
	}

}

func (n *Node) num_neighbors(graph map[string][]Edge, nodes map[string]*Node, global *Global) int {
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
	if n.NodeID == "d" {
		fmt.Printf("d node k core: %v core num: %v \n", n.k_core_msg, n.CoreNum)
		fmt.Println("d node k core len: ", len(n.k_core_msg))
	}
	n.k_core_msg = make(map[string]bool)
	return count

}

// heart beat publish
func (n *Node) publish(nodes map[string]*Node, global *Global, neighbors map[string][]Edge) {
	time.Sleep(300 * time.Millisecond)
	for _, neighbor := range neighbors[n.NodeID] {
		var n_node = nodes[neighbor.Target]
		n_node.Msg_queue = append(n_node.Msg_queue, heart_beat_msg{global.hb_msg_id, n, n_node})
		global.hb_msg_id = global.hb_msg_id + 1
	}
}

func (n *Node) terminate(nodes map[string]*Node, neighbors map[string][]Edge) bool {
	var allNeighstop = true
	for _, n := range neighbors[n.NodeID] {
		if len(nodes[n.Target].Msg_queue) > 0 {
			allNeighstop = false
		}
	}
	if len(n.Msg_queue) == 0 && allNeighstop {
		time.Sleep(200 * time.Millisecond)
		return true
	}
	return false
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
	fromV, _ := Encrypt(strconv.Itoa(n.CoreNum), n.private_key)
	for _, m := range to.secure_msg {
		if m.from.NodeID == n.NodeID {
			return
		}
	}
	fmt.Printf("node %v send: %v to: %v \n", n.NodeID, fromV, to.NodeID)

	to.secure_msg = append(to.secure_msg, secure_msg{n, fromV, to, "", n.private_key, ""})
}

func (n *Node) reply() {
	time.Sleep(300 * time.Millisecond)
	curMsg := n.secure_msg
	n.secure_msg = nil
	for _, m := range curMsg {
		m.to_v, _ = Encrypt(strconv.Itoa(n.CoreNum), n.private_key)
		m.to_key = n.private_key
		if _, exists := m.from.k_core_msg[n.NodeID]; !exists {
			// Key does not exist, write it
			m.from.k_core_msg[n.NodeID] = m.compare()
		}
		fmt.Printf("node: %v reply with %v back to %v \n", n.NodeID, m.to_v, m.from.NodeID)
	}
}

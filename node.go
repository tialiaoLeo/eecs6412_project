// node.go
package main

import (
	"crypto/rsa"
	"fmt"
	"math"
	"strconv"
	"time"
)

// global variable for message id

// Node represents a node in a graph

type Node struct {
	NodeID     string           // Unique identifier for the node
	CoreNum    int              // Core number of the node
	Msg_queue  []heart_beat_msg // Heartbeat messages
	secure_msg []secure_msg     // Secure messages received
	k_core_msg map[string]bool  // Boolean array for k-core
	PrivateKey *rsa.PrivateKey  // RSA private key
	PublicKey  *rsa.PublicKey   // RSA public key
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
	n.Msg_queue = []heart_beat_msg{}
	return pub_again
}
func (n *Node) send(to *Node) {
	time.Sleep(300 * time.Millisecond)

	// using the public key of the target node
	encryptedCoreNum, err := EncryptWithPublicKey(strconv.Itoa(n.CoreNum), to.PublicKey)
	if err != nil {
		fmt.Println("Error encrypting:", err)
		return
	}

	for _, m := range to.secure_msg {
		if m.from.NodeID == n.NodeID {
			return
		}
	}

	fmt.Printf("Node %v sends encrypted coreNum to Node %v: %v\n", n.NodeID, to.NodeID, encryptedCoreNum)

	to.secure_msg = append(to.secure_msg, secure_msg{from: n, from_v: encryptedCoreNum, to: to})
}
func (n *Node) reply() {
	time.Sleep(300 * time.Millisecond)
	curMsg := n.secure_msg
	n.secure_msg = nil

	for _, m := range curMsg {
		// encrypt reply message using the public key from the sender
		encryptedReply, err := EncryptWithPublicKey(strconv.Itoa(n.CoreNum), m.from.PublicKey)
		if err != nil {
			fmt.Printf("Node %v failed to encrypt reply: %v\n", n.NodeID, err)
			continue
		}

		// store the encrypted message
		m.to_v = encryptedReply

		// using compare method to compare, which will be firstly decrypted
		if _, exists := m.from.k_core_msg[n.NodeID]; !exists {
			m.from.k_core_msg[n.NodeID] = m.compare()
		}

		fmt.Printf("Node %v replied with encrypted coreNum back to Node %v\n", n.NodeID, m.from.NodeID)
	}
}

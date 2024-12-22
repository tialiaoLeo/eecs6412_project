package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	filePath := "fixed_graph.txt"

	// Call the ReadGraphFile function
	graph, err := ReadGraphFile(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Initialize nodes and generate RSA keys
	nodes := make(map[string]*Node)
	for key, value := range graph {
		privateKey, err := GenerateRSAKeys()
		if err != nil {
			fmt.Printf("Error generating keys for node %s: %v\n", key, err)
			return
		}
		nodes[key] = &Node{
			NodeID:          key,
			CoreNum:         len(value),
			Msg_queue:       []heart_beat_msg{},
			secure_msg:      []secure_msg{},
			k_core_msg:      make(map[string]bool),
			PrivateKey:      privateKey,
			PublicKey:       &privateKey.PublicKey,
			terminate_times: 0,
		}
	}
	// WaitGroup to wait for both Go routines to finish
	var wg sync.WaitGroup
	for _, n := range nodes {
		wg.Add(1)
		go func(n *Node) {
			defer wg.Done()
			n.k_core(graph, nodes)
			fmt.Printf("Node: %v CoreNum: %v \n", n.NodeID, n.CoreNum)
			n.publish(nodes, graph)

			for {
				// Regular work
				time.Sleep(100 * time.Millisecond)
				n.reply()
				n.k_core(graph, nodes)
				if n.terminate(nodes, graph) {
					n.terminate_times = n.terminate_times + 1
				}
				if n.consume() {
					pre := n.CoreNum
					if pre != n.CoreNum {
						n.publish(nodes, graph)
					}
				}
				fmt.Printf("Node: %v CoreNum: %v \n", n.NodeID, n.CoreNum)
				if n.terminate_times > 1 {
					fmt.Println("Termination condition met. Shutting down...")
					fmt.Printf("Final Node: %v CoreNum: %v \n", n.NodeID, n.CoreNum)
					break
				}
			}
		}(n)
	}

	// Wait for both goroutines to complete
	wg.Wait()
}

/*
func terminate(ctx context.Context, cancel context.CancelFunc, nodes map[string]*Node, global *Global) {
	for {
		time.Sleep(500 * time.Millisecond)
		// Check if all nodes have empty Msg_queue
		allIdle := true
		for _, node := range nodes {
			if len(node.Msg_queue) > 0 {
				allIdle = false
				break
			}
		}
		if allIdle && global.hb_msg_id > 0 {
			fmt.Println("Termination condition met. Shutting down...")
			cancel() // Signal cancellation to all goroutines
			return
		}
	}
}


func terminate1(ctx context.Context, cancel context.CancelFunc, nodes map[string]*Node, initiator *Node) {
	fmt.Printf("Node %v initiating termination check...\n", initiator.NodeID)

	for _, node := range nodes {
		go func(node *Node) {
			// Check if Msg_queue is empty
			isIdle := len(node.Msg_queue) == 0

			// Encrypt the state with initiator's public key
			encryptedState, err := EncryptWithPublicKey(strconv.FormatBool(isIdle), initiator.PublicKey)
			if err != nil {
				fmt.Printf("Error encrypting state for Node %v: %v\n", node.NodeID, err)
				return
			}

			// Send encrypted state back to the initiator
			initiator.secure_msg = append(initiator.secure_msg, secure_msg{
				from:   node,
				from_v: encryptedState,
				to:     initiator,
			})
		}(node)
	}

	// Wait for responses
	time.Sleep(200 * time.Second)

	// Process responses
	allIdle := true
	for _, msg := range initiator.secure_msg {
		decryptedState, err := DecryptWithPrivateKey(msg.from_v, initiator.PrivateKey)
		if err != nil {
			fmt.Printf("Error decrypting state from Node %v: %v\n", msg.from.NodeID, err)
			allIdle = false
			break
		}

		isIdle, err := strconv.ParseBool(decryptedState)
		if err != nil {
			fmt.Printf("Error parsing state from Node %v: %v\n", msg.from.NodeID, err)
			allIdle = false
			break
		}

		if !isIdle {
			allIdle = false
			break
		}
	}

	// Clear secure_msg queue
	initiator.secure_msg = nil

	if allIdle {
		fmt.Println("Termination condition met. Shutting down...")
		cancel()
	} else {
		fmt.Println("Termination condition not met. Continuing execution...")
	}
}
*/

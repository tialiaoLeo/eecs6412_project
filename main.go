package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Global struct {
	// global heart beat msg id
	hb_msg_id int
}

func main() {
	filePath := "graph.txt"
	global := &Global{0}
	// Call the ReadGraphFile function
	graph, err := ReadGraphFile(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	nodes := make(map[string]*Node)
	// Print the graph
	for key, value := range graph {
		nodes[key] = &Node{key, len(value), []heart_beat_msg{}, []secure_msg{}, map[string]bool{}, createKey(key)}
	}
	ctx, cancel := context.WithCancel(context.Background())
	// WaitGroup to wait for both Go routines to finish
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		terminate(ctx, cancel, nodes, global)
	}()
	for _, n := range nodes {
		wg.Add(1)
		go func(ctx context.Context, n *Node, global *Global) {
			defer wg.Done()
			n.k_core(graph, nodes, global)
			fmt.Printf("Node: %v CoreNum: %v \n", n.NodeID, n.CoreNum)
			n.publish(nodes, global, graph)

			for {
				select {
				case <-ctx.Done():
					// Termination signal received
					fmt.Printf("Node %v shutting down...\n", n.NodeID)
					return
				default:
					// Regular work
					time.Sleep(100 * time.Millisecond)
					n.reply()
					n.k_core(graph, nodes, global)
					if n.consume() {
						pre := n.CoreNum
						if pre != n.CoreNum {
							n.publish(nodes, global, graph)
						}
					}
					fmt.Printf("Node: %v CoreNum: %v \n", n.NodeID, n.CoreNum)
				}
			}
		}(ctx, n, global)
	}

	// Wait for both goroutines to complete
	wg.Wait()

}

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

func debug(n *Node) {
	if n.NodeID == "a" {
		fmt.Println("a: ", n.k_core_msg)
	}
}

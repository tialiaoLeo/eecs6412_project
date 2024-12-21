package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	filePath := "graph.txt"
	// Call the ReadGraphFile function
	graph, err := ReadGraphFile(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	nodes := make(map[string]*Node)
	// Print the graph
	for key, value := range graph {
		nodes[key] = &Node{key, len(value), []heart_beat_msg{}, []secure_msg{}, map[string]bool{}, createKey(key), 0}
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
				if n.terminate_times > 3 {
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

package main

import (
	"fmt"
	"sync"
	"time"
)

func (n *Node) k_core_exe(graph map[string][]Edge, nodes map[string]*Node) {
	// Example logic for k_core, modify as needed
	n.k_core(graph, nodes)
	fmt.Printf("Node: %v CoreNum: %v \n", n.NodeID, n.CoreNum)
}

var done = false

type Global struct {
	ms_id int
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
		nodes[key] = &Node{key, len(value), []Message{}}
	}

	// WaitGroup to wait for both Go routines to finish
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		terminate(nodes, global)
	}()
	for _, n := range nodes {
		wg.Add(1)
		go func(n *Node, global *Global) {
			defer wg.Done()
			n.k_core(graph, nodes)
			fmt.Printf("Node: %v CoreNum: %v \n", n.NodeID, n.CoreNum)
			n.publish(nodes, global)
			for !done {
				if n.consume() {
					pre := n.CoreNum
					n.k_core(graph, nodes)
					fmt.Printf("Node: %v CoreNum: %v \n", n.NodeID, n.CoreNum)
					if pre != n.CoreNum {
						n.publish(nodes, global)
					}
				}
			}
		}(n, global)
	}

	// Wait for both goroutines to complete
	wg.Wait()
}

func terminate(nodes map[string]*Node, global *Global) {
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
		if allIdle && global.ms_id > 0 {
			fmt.Println("Termination condition met. Shutting down...")
			done = true // Notify all Goroutines to terminate
			return
		}
	}
}

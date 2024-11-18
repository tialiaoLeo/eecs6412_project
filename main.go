package main

import (
	"fmt"
	"sync"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func (n *Node) k_core_exe(graph map[string][]Edge, nodes map[string]Node) {
	// Example logic for k_core, modify as needed
	n.k_core(graph, nodes)
	fmt.Printf("Node: %v CoreNum: %v \n", n.NodeID, n.CoreNum)
}
func main() {
	filePath := "graph.txt"

	// Call the ReadGraphFile function
	graph, err := ReadGraphFile(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	nodes := make(map[string]Node)
	// Print the graph
	for key, value := range graph {
		nodes[key] = Node{key, len(value)}
	}

	// WaitGroup to wait for both Go routines to finish
	var wg sync.WaitGroup

	for _, n := range nodes {
		wg.Add(1)
		go func(n Node) {
			defer wg.Done()
			n.k_core(graph, nodes)
			fmt.Printf("Node: %v CoreNum: %v \n", n.NodeID, n.CoreNum)
		}(n)
	}

	// Wait for both goroutines to complete
	wg.Wait()
}

package main

import (
	"fmt"
	"sync"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	filePath := "graph.txt"

	// Call the ReadGraphFile function
	graph, err := ReadGraphFile(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the graph
	fmt.Println("Graph:")
	for key, value := range graph {
		fmt.Printf("%s -> %v\n", key, value)
	}

	// Create a new Node
	nodeA := Node{
		NodeID:        "a",
		NodeNeighbors: []string{"b", "d"},
		CoreNum:       2,
	}
	/*
		nodeB := Node{
			NodeID:        "b",
			NodeNeighbors: []string{"a", "d", "g"},
			CoreNum:       3,
		}*/
	nodeC := Node{
		NodeID:        "g",
		NodeNeighbors: []string{"b", "d", "f", "h"},
		CoreNum:       4,
	}

	// WaitGroup to wait for both Go routines to finish
	var wg sync.WaitGroup
	wg.Add(2)

	// Increment coreNum for nodeA in a separate goroutine
	go func() {
		defer wg.Done()
		nodeA.IncrementCoreNum()
		fmt.Printf("After increment, Node A: %+v\n", nodeA)
	}()

	// Increment coreNum for nodeC in a separate goroutine
	go func() {
		defer wg.Done()
		nodeC.IncrementCoreNum()
		fmt.Printf("After increment, Node C: %+v\n", nodeC)
	}()

	// Wait for both goroutines to complete
	wg.Wait()
}

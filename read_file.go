package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Edge struct {
	Weight int
	Target string
}

// ReadGraphFile reads a graph from a file and returns it as a hash map.
func ReadGraphFile(filename string) (map[string][]interface{}, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Create the hash map
	graph := make(map[string][]interface{})

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Split the line into parts
		parts := strings.Fields(line)

		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid line format: %s", line)
		}

		// Parse the weight and node
		weight, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid weight in line: %s", line)
		}
		fromNode := parts[0]
		toNode := parts[2]

		// Store in the hash map
		graph[fromNode] = append(graph[fromNode], Edge{Weight: weight, Target: toNode})
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return graph, nil
}

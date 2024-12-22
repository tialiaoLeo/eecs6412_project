package main

import (
	"fmt"
	"strconv"
)

type heart_beat_msg struct {
	from *Node
	to   *Node
}

type secure_msg struct {
	from   *Node  // sening node
	from_v string // Encrypted k-core
	to     *Node  // receiving node
	to_v   string // Encrypted CoreNum when replying
}

func (m *secure_msg) compare() bool {
	// Decrypt the CoreNum of sender
	fmt.Printf("Node %v decrypting from_v...\n", m.to.NodeID)
	fromVclear, err := DecryptWithPrivateKey(m.from_v, m.to.PrivateKey)
	if err != nil {
		fmt.Printf("Error decrypting from_v for Node %v: %v\n", m.to.NodeID, err)
		return false
	}

	// Decrypt the CoreNum of replier
	fmt.Printf("Node %v decrypting to_v...\n", m.to.NodeID)
	toVclear, err := DecryptWithPrivateKey(m.to_v, m.to.PrivateKey)
	if err != nil {
		fmt.Printf("Error decrypting to_v for Node %v: %v\n", m.to.NodeID, err)
		return false
	}

	// Convert decrypted values back to integers
	toInt, err1 := strconv.Atoi(toVclear)
	fromInt, err2 := strconv.Atoi(fromVclear)

	if err1 != nil || err2 != nil {
		fmt.Printf("Error converting decrypted values to integers for Node %v: %v, %v\n", m.to.NodeID, err1, err2)
		return false
	}

	// Print debug information
	fmt.Printf("Node %v decrypted to_v: %d, CoreNum: %d\n", m.to.NodeID, toInt, m.to.CoreNum)
	fmt.Printf("Node %v decrypted from_v: %d, CoreNum: %d\n", m.from.NodeID, fromInt, m.from.CoreNum)

	// Update the core-number if necessary
	if toInt != m.to.CoreNum {
		fmt.Printf("Node %v updating CoreNum from %d to %d\n", m.to.NodeID, m.to.CoreNum, toInt)
		toInt = m.to.CoreNum
	}
	if fromInt != m.from.CoreNum {
		fmt.Printf("Node %v updating CoreNum from %d to %d\n", m.from.NodeID, m.from.CoreNum, fromInt)
		fromInt = m.from.CoreNum
	}

	// Perform the comparison
	result := toInt >= fromInt
	fmt.Printf("Comparison result for Node %v (to) >= Node %v (from): %v\n", m.to.NodeID, m.from.NodeID, result)
	return result
}

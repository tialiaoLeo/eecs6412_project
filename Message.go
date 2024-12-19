package main

import (
	"fmt"
	"strconv"
)

type heart_beat_msg struct {
	id   int
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
	// decrypt the CoreNum of sender and replier
	fromVclear, err := DecryptWithPrivateKey(m.from_v, m.to.PrivateKey)
	if err != nil {
		fmt.Println("Error decrypting from_v:", err)
		return false
	}

	toVclear, err := DecryptWithPrivateKey(m.to_v, m.to.PrivateKey)
	if err != nil {
		fmt.Println("Error decrypting to_v:", err)
		return false
	}

	// converting back to an integer
	toInt, err1 := strconv.Atoi(toVclear)
	fromInt, err2 := strconv.Atoi(fromVclear)

	if err1 != nil || err2 != nil {
		fmt.Println("Error converting decrypted values to integers:", err1, err2)
		return false
	}

	// update the core-number
	if toInt != m.to.CoreNum {
		toInt = m.to.CoreNum
	}
	if fromInt != m.from.CoreNum {
		fromInt = m.from.CoreNum
	}

	// compare core number
	return toInt >= fromInt
}

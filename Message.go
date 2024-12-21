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
	from     *Node
	from_v   string
	to       *Node
	to_v     string
	from_key string
	to_key   string
}

func (m *secure_msg) compare() bool {
	fromVclear, _ := Decrypt(m.from_v, m.from_key)
	toVclear, _ := Decrypt(m.to_v, m.to_key)
	toInt, err1 := strconv.Atoi(toVclear)
	fromInt, err2 := strconv.Atoi(fromVclear)

	// Check for conversion errors
	if err1 != nil || err2 != nil {
		fmt.Println("Error converting strings to integers:", err1, err2)
	}
	if toInt != m.to.CoreNum {
		toInt = m.to.CoreNum
	}

	if fromInt != m.from.CoreNum {
		fromInt = m.from.CoreNum
	}
	// Perform the comparison
	if toInt >= fromInt {
		//fmt.Println("toVclear is greater than or equal to fromVclear")
		return true
	} else {
		//fmt.Println("toVclear is less than fromVclear")
		return false
	}
}

package main

import (
	"github.com/tuneinsight/lattigo/v4/bfv"
	"github.com/tuneinsight/lattigo/v4/rlwe"
)

type heart_beat_msg struct {
	from *Node
	to   *Node
}

type secure_msg struct {
	from      *Node
	from_v    *rlwe.Ciphertext
	to        *Node
	to_v      *rlwe.Ciphertext
	encoder   bfv.Encoder
	decryptor rlwe.Decryptor
	encryptor rlwe.Encryptor
	params    bfv.Parameters
}

func (m *secure_msg) compare() bool {
	ciphertext1 := m.from_v
	ciphertext2 := m.to_v
	cur1 := Encrypt(m.from.CoreNum, m)
	cur2 := Encrypt(m.to.CoreNum, m)
	if ciphertext1 != &cur1 {
		ciphertext1 = &cur1
	}
	if ciphertext2 != &cur2 {
		ciphertext2 = &cur2
	}
	// Step 5: Subtract ciphertext2 from ciphertext1
	evaluator := bfv.NewEvaluator(m.params, rlwe.EvaluationKey{})
	diffCiphertext := evaluator.SubNoModNew(ciphertext1, ciphertext2)

	// Step 6: Decrypt the result
	diffPlaintext := bfv.NewPlaintext(m.params, m.params.MaxLevel())
	m.decryptor.Decrypt(diffCiphertext, diffPlaintext)

	// Decode the result
	diff := m.encoder.DecodeUintNew(diffPlaintext)
	// Step 7: Determine comparison result
	if diff[0] == 0 || diff[0] > 55537 {
		//fmt.Println("yes") // value2 >= value1
		return true
	} else {
		//fmt.Println("no") // value2 < value1
		return false
	}
}

func Encoder() (bfv.Encoder, rlwe.Decryptor, rlwe.Encryptor, bfv.Parameters) {
	params, _ := bfv.NewParametersFromLiteral(bfv.PN12QP109)

	// Step 2: Key generation
	kgen := bfv.NewKeyGenerator(params)
	sk, pk := kgen.GenKeyPair()

	// Step 3: Create encryptor, decryptor, and encoder
	encryptor := bfv.NewEncryptor(params, pk)
	decryptor := bfv.NewDecryptor(params, sk)
	encoder := bfv.NewEncoder(params)
	return encoder, decryptor, encryptor, params
}

func Newsecure_msg(fromV int, from *Node, to *Node) *secure_msg {
	var encoder, decryptor, encryptor, params = Encoder()
	cur_msg := &secure_msg{from, nil, to, nil, encoder, decryptor, encryptor, params}
	c1 := Encrypt(fromV, cur_msg)
	cur_msg.from_v = &c1
	return cur_msg
}

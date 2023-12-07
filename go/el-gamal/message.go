package el_gamal

import "math/big"

type Message struct {
	x []*big.Int
	y []*big.Int
}

type EncryptedMessage struct {
	Ax []*big.Int
	Ay []*big.Int
	Bx []*big.Int
	By []*big.Int
}

func BytesToMessage(msg []byte) *Message {
	return nil
}

// Package el_gamal
// Copyright 2023 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
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

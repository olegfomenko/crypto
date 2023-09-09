// Copyright 2023 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package schnorr_bjj

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/stretchr/testify/assert"
)

func TestSchnorrSignature(t *testing.T) {
	var bytes [30]byte
	_, err := rand.Read(bytes[:])
	assert.NoError(t, err)

	prv := new(big.Int).SetBytes(bytes[:])

	pk := babyjub.NewPoint().Mul(prv, G)

	fmt.Println(pk.X.String())
	fmt.Println(pk.Y.String())

	message := new(big.Int).SetBytes(crypto.Keccak256([]byte("Hello world")))
	fmt.Println(message.String())

	sig, err := SignSchnorr(prv, pk, message)
	assert.NoError(t, err)

	fmt.Println(sig.R.X.String())
	fmt.Println(sig.R.Y.String())
	fmt.Println(sig.S.String())

	err = VerifySchnorr(sig, pk, message)
	assert.NoError(t, err)
}

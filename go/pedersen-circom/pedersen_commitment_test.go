// Package pedersen
// Copyright 2023 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package pedersen

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/olegfomenko/crypto/go/pedersen-circom/types"
	"github.com/stretchr/testify/assert"
)

func TestPedersenCommitment(t *testing.T) {
	prv := NewRandomPrivateKey()
	fmt.Println(prv)
	amount := types.Wrap(big.NewInt(100))

	com, proof, err := NewCommitment(prv, amount)
	assert.NoError(t, err)

	err = VerifyCommitmentProof(com, proof)
	assert.NoError(t, err)

	comJSON, err := json.MarshalIndent(com, "", " ")
	assert.NoError(t, err)

	proofJSON, err := json.MarshalIndent(proof, "", " ")
	assert.NoError(t, err)

	fmt.Println(string(comJSON))
	fmt.Println(string(proofJSON))
}

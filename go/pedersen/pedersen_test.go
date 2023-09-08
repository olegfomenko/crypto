// Package pedersen
// Copyright 2032 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package pedersen

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

func init() {
	_, G, _ = bn256.RandomG1(rand.Reader)
	_, H, _ = bn256.RandomG1(rand.Reader)
}

func TestBitRepresentation(t *testing.T) {
	fmt.Println(strconv.FormatUint(17, 2))
	fmt.Println(hexutil.Encode(bn256.Order.Bytes()))
}

func TestPedersenCommitment(t *testing.T) {
	proof, commitment, prv, err := CreatePedersenCommitment(10, 5)
	if err != nil {
		panic(err)
	}

	reconstructedCommitment := PedersenCommitment(big.NewInt(10), prv)
	fmt.Println("Constructed commitment with prv key: " + reconstructedCommitment.String())
	fmt.Println("Response commitment: " + commitment.String())

	fmt.Println("Private Key: " + hexutil.Encode(prv.Bytes()))

	if err = VerifyPedersenCommitment(commitment, proof); err != nil {
		panic(err)
	}
}

func TestPedersenCommitmentFails(t *testing.T) {
	_, _, _, err := CreatePedersenCommitment(128, 5)
	if err == nil {
		panic("Should fail")
	}
}

func TestSchnorrSignature(t *testing.T) {
	prv, err := rand.Int(rand.Reader, bn256.Order)
	if err != nil {
		panic(err)
	}

	pk := ScalarMul(G, prv)
	message := Hash([]byte("Hello world"))

	sig, err := SignSchnorr(prv, pk, message)
	if err != nil {
		panic(err)
	}

	if err := VerifySchnorr(sig, pk, message); err != nil {
		panic(err)
	}
}

func TestSchnorrSignatureAggregation(t *testing.T) {
	prvAlice, err := rand.Int(rand.Reader, bn256.Order)
	if err != nil {
		panic(err)
	}

	pubAlice := ScalarMul(G, prvAlice)

	prvBob, err := rand.Int(rand.Reader, bn256.Order)
	if err != nil {
		panic(err)
	}

	pubBob := ScalarMul(G, prvBob)

	pubCombined := Add(pubAlice, pubBob)

	message := Hash([]byte("Hello world"))

	sigAlice, err := SignSchnorr(prvAlice, pubCombined, message)
	if err != nil {
		panic(err)
	}

	sigBob, err := SignSchnorr(prvBob, pubCombined, message)
	if err != nil {
		panic(err)
	}

	rCom := Add(sigAlice.R, sigBob.R)
	sigCom := SchnorrSignature{
		S: add(sigAlice.S, sigBob.S),
		R: rCom,
	}

	if err := VerifySchnorr(sigCom, pubCombined, message); err != nil {
		panic(err)
	}
}

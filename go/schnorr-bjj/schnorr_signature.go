// Package schnorr_bjj
// Copyright 2023 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// File `schnorr_signature.go` implements Schnorr signature generation on the
// BabyJubjub (`https://eips.ethereum.org/EIPS/eip-2494`) elliptic curve.
package schnorr_bjj

import (
	"crypto/rand"
	"errors"
	goerr "errors"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/iden3/go-iden3-crypto/babyjub"
)

type HashF func(bytes ...[]byte) []byte

var Hash HashF = crypto.Keccak256

var ErrFailedRandom = goerr.New("dailed to generate secure random")

// SignSchnorr creates the Schnorr signature for the given public key and message.
// Can be used both for mono- and multi- signing.
// For multisignature `publicKey` should be an aggregated public key.
// For monosignature `publicKey` should be an elliptic point `prv*G`.
func SignSchnorr(prv *big.Int, publicKey *babyjub.Point, m *big.Int) (*SchnorrSignature, error) {
	var bytes [30]byte
	_, err := rand.Read(bytes[:])
	if err != nil {
		return nil, ErrFailedRandom
	}

	k := new(big.Int).SetBytes(bytes[:])

	kG := babyjub.NewPoint().Mul(k, G)

	hash := new(big.Int).SetBytes(Hash(m.Bytes(), publicKey.X.Bytes(), publicKey.Y.Bytes()))

	s := new(big.Int).Add(k, new(big.Int).Mul(hash, prv))

	return &SchnorrSignature{
		R: kG,
		S: s,
	}, nil
}

// VerifySchnorr verifies Schnorr signature validity.
// Can be used both for mono- and multi- signing.
// For multisignature `publicKey` should be an aggregated public key.
// For monosignature `publicKey` should be an elliptic point `prv*G`.
func VerifySchnorr(sig *SchnorrSignature, publicKey *babyjub.Point, m *big.Int) error {
	// s = k + hash*prv
	// r = kG

	// p2 = (k + hash*prv)*G
	// p1 = kG + hash*prv*G = (k + hash*prv)*Gs

	hash := new(big.Int).SetBytes(Hash(m.Bytes(), publicKey.X.Bytes(), publicKey.Y.Bytes()))

	p1 := babyjub.NewPoint().Mul(hash, publicKey)
	p1 = babyjub.NewPointProjective().Add(sig.R.Projective(), p1.Projective()).Affine()

	p2 := babyjub.NewPoint().Mul(sig.S, G)

	if !BJJPointEqual(p1, p2) {
		return errors.New("verification failed")
	}

	return nil
}

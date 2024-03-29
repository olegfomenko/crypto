// Package schnorr_bn256
// Copyright 2023 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// File `schnorr.go` implements Schnorr signature generation on the
// bn256 elliptic curve that compatible with alt_bn128 (https://eips.ethereum.org/EIPS/eip-1108)
package schnorr_bn256

import (
	"bytes"
	"crypto/rand"
	"github.com/cloudflare/bn256"
	"math/big"
)

func R(G *bn256.G1) (*big.Int, *bn256.G1, error) {
	r, err := rand.Int(rand.Reader, bn256.Order)
	if err != nil {
		return nil, nil, ErrFailedRandom
	}

	return r, new(bn256.G1).ScalarMult(G, r), nil
}

func Msg(bytes ...[]byte) *big.Int {
	return new(big.Int).Mod(new(big.Int).SetBytes(Hash(bytes...)), bn256.Order)
}

// MultiSigSchnorr creates the Schnorr signature for the given public key and message.
// For multi-signature `PubKeyCommon` should be an aggregated public key and `RCommon` should be aggregated R.
func MultiSigSchnorr(prv *big.Int, r *big.Int, PubKeyCommon *bn256.G1, RCommon *bn256.G1, m *big.Int) (*SchnorrSignature, error) {
	hash := Msg(m.Bytes(), PubKeyCommon.Marshal(), RCommon.Marshal())
	s := new(big.Int).Add(r, new(big.Int).Mul(hash, prv))

	return &SchnorrSignature{
		R: RCommon,
		S: s,
	}, nil
}

// SignSchnorr creates the Schnorr signature for the given public key and message.
// `PublicKey` should be an elliptic point `prv*G`.
func SignSchnorr(prv *big.Int, PublicKey *bn256.G1, G *bn256.G1, m *big.Int) (*SchnorrSignature, error) {
	r, R, err := R(G)
	if err != nil {
		return nil, err
	}

	hash := Msg(m.Bytes(), PublicKey.Marshal(), R.Marshal())

	s := new(big.Int).Add(r, new(big.Int).Mul(hash, prv))

	return &SchnorrSignature{
		R: R,
		S: s,
	}, nil
}

// VerifySchnorr verifies Schnorr signature validity.
// Can be used both for mono- and multi- signing.
// For multi-signature `PublicKey` should be an aggregated public key.
// For mono-signature `PublicKey` should be an elliptic point `prv*G`.
func VerifySchnorr(sig *SchnorrSignature, PublicKey *bn256.G1, G *bn256.G1, m *big.Int) bool {
	// s = r + hash*prv
	// R = rG

	// p2 = (r + hash*prv)*G
	// p1 = rG + hash*prv*G = (r + hash*prv)*G

	hash := Msg(m.Bytes(), PublicKey.Marshal(), sig.R.Marshal())

	p1 := new(bn256.G1).ScalarMult(PublicKey, hash)
	p1 = new(bn256.G1).Add(p1, sig.R)

	p2 := new(bn256.G1).ScalarMult(G, sig.S)

	return bytes.Equal(p1.Marshal(), p2.Marshal())
}

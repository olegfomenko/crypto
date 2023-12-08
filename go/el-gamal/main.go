// Package el_gamal
// Copyright 2023 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package el_gamal

import (
	"crypto/elliptic"
	"crypto/rand"
	"math/big"

	"github.com/olegfomenko/crypto/go/ec"
)

var Curve elliptic.Curve = ec.SECP256K1()

type PublicKey struct {
	X, Y *big.Int
}

type PrivateKey struct {
	*PublicKey
	D *big.Int
}

type Cypher struct {
	Ax, Ay *big.Int
	Bx, By *big.Int
}

func GeneratePrivateKey() (*PrivateKey, error) {
	d, err := rand.Int(rand.Reader, Curve.Params().N)
	if err != nil {
		return nil, err
	}

	x, y := Curve.ScalarBaseMult(d.Bytes())
	return &PrivateKey{
		PublicKey: &PublicKey{
			X: x,
			Y: y,
		},
		D: d,
	}, nil
}

func Encrypt(mx, my *big.Int, pub *PublicKey) (*Cypher, error) {
	k, err := rand.Int(rand.Reader, Curve.Params().N)
	if err != nil {
		return nil, err
	}

	ax, ay := Curve.ScalarBaseMult(k.Bytes())

	cx, cy := Curve.ScalarMult(pub.X, pub.Y, k.Bytes())
	bx, by := Curve.Add(mx, my, cx, cy)
	return &Cypher{ax, ay, bx, by}, nil
}

func Decrypt(cypher *Cypher, prv *PrivateKey) (mx *big.Int, my *big.Int) {
	x, y := Curve.ScalarMult(cypher.Ax, cypher.Ay, prv.D.Bytes())
	mx, my = Curve.Add(cypher.Bx, cypher.By, x, minus(y))
	return
}

func minus(val *big.Int) *big.Int {
	return new(big.Int).Mod(new(big.Int).Mul(val, big.NewInt(-1)), Curve.Params().P)
}

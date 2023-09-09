// Package deffie_hellman
// Copyright 2023 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package deffie_hellman

import (
	"crypto/rand"
	goerr "errors"
	"math/big"
)

var ErrInvalidParams = goerr.New("invalid params")

type Party struct {
	g      *big.Int
	p      *big.Int
	secret *big.Int
}

func NewParty(g, p *big.Int) (*Party, error) {
	if g == nil || p == nil {
		return nil, ErrInvalidParams
	}

	secret, err := rand.Int(rand.Reader, p)
	if err != nil {
		return nil, err
	}

	return &Party{
		g:      g,
		p:      p,
		secret: secret,
	}, nil
}

func (p *Party) GetShare() *big.Int {
	return new(big.Int).Exp(p.g, p.secret, p.p)
}

func (p *Party) ReceiveShare(share *big.Int) *big.Int {
	return new(big.Int).Exp(share, p.secret, p.p)
}

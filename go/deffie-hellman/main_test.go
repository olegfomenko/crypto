// Package deffie_hellman
// Copyright 2032 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package deffie_hellman

import (
	"crypto/rand"
	"math/big"
	"testing"
)

var g = new(big.Int).SetInt64(7)
var p, _ = rand.Prime(rand.Reader, 256)

func TestTwoParties(t *testing.T) {
	t.Logf("P = %s", p.String())

	alice, err := NewParty(g, p)
	if err != nil {
		panic(err)
	}

	bob, err := NewParty(g, p)
	if err != nil {
		panic(err)
	}

	aliceShare, bobShare := alice.GetShare(), bob.GetShare()
	aliceKey := alice.ReceiveShare(bobShare)
	bobKey := bob.ReceiveShare(aliceShare)

	t.Logf("Alice secret = %s", aliceKey.String())
	t.Logf("Bob secret = %s", bobKey.String())
	if aliceKey.Cmp(bobKey) != 0 {
		panic("is no equal")
	}
}

func TestFourParties(t *testing.T) {
	t.Logf("P = %s", p.String())

	const n = 4
	parties := make([]*Party, n)

	var err error
	for i := 0; i < n; i++ {
		parties[i], err = NewParty(g, p)
		if err != nil {
			panic(err)
		}
	}

	next := func(i int) int {
		return (i + 1) % n
	}

	for i := 0; i < n; i++ {
		cur := parties[i].GetShare()
		to := i
		for j := 0; j < n-1; j++ {
			to = next(to)
			cur = parties[to].ReceiveShare(cur)
		}
		t.Logf("%d secret = %s", i, cur.String())
	}
}

func TestTenParties(t *testing.T) {
	t.Logf("P = %s", p.String())

	const n = 10
	parties := make([]*Party, n)

	var err error
	for i := 0; i < n; i++ {
		parties[i], err = NewParty(g, p)
		if err != nil {
			panic(err)
		}
	}

	next := func(i int) int {
		return (i + 1) % n
	}

	for i := 0; i < n; i++ {
		cur := parties[i].GetShare()
		to := i
		for j := 0; j < n-1; j++ {
			to = next(to)
			cur = parties[to].ReceiveShare(cur)
		}
		t.Logf("%d secret = %s", i, cur.String())
	}
}

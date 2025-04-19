// Package parcs
// Copyright 2025 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package parcs

import (
	"math/big"
	"testing"
)

// helper: проверяет, что d ‑‑ нетривиальный делитель n.
func isValidFactor(n, d *big.Int) bool {
	one := big.NewInt(1)
	// 1 < d < n  и  n mod d == 0
	return d.Cmp(one) > 0 &&
		d.Cmp(n) < 0 &&
		new(big.Int).Rem(n, d).Sign() == 0
}

// --- ТЕСТ 1 -----------------------------------------------------------------
// 8051 = 83 · 97.
// Ожидаем любой из нетривиальных делителей (83 или 97).
func TestPollardRho_8051(t *testing.T) {
	n, _ := new(big.Int).SetString("8051", 10)

	d := Pollard(n)
	if !isValidFactor(n, d) {
		t.Fatalf("got %s, expected нетривиальный делитель 8051", d)
	}
}

// --- ТЕСТ 2 -----------------------------------------------------------------
// Чётное число 100.  pollardRho должен сразу вернуть 2.
func TestPollardRho_100(t *testing.T) {
	n := big.NewInt(100)

	d := Pollard(n)
	if !isValidFactor(n, d) {
		t.Fatalf("got %s, expected нетривиальный делитель 100", d)
	}
}

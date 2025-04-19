// Package parcs
// Copyright 2025 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package parcs

import (
	"fmt"
	"math/big"
)

const (
	ProbTestSteps = 20
	MaxDivs       = 1000
)

type ParcsPollrad struct {
	limit uint32
	sem   chan struct{}
}

func NewParcsPollrad(limit uint32) *ParcsPollrad {
	return &ParcsPollrad{
		limit: limit,
		sem:   make(chan struct{}, limit),
	}
}

func (p *ParcsPollrad) runTask(n *big.Int, result chan *big.Int) {
	p.sem <- struct{}{}
	go func() {
		fmt.Println("! Running task for n =", n)
		fmt.Println("Now task running:", len(p.sem))
		task(n, result, p.runTask)
		_ = <-p.sem
		fmt.Println("! Finished task for n =", n)
	}()
}

func task(n *big.Int, result chan *big.Int, spawn func(n *big.Int, result chan *big.Int)) {
	if n.ProbablyPrime(ProbTestSteps) {
		fmt.Println("Found prime n=", n)
		result <- n
		return
	}
	fmt.Println("Running Pollrad for n=", n)

	d := Pollrad(n)

	spawn(d, result)
	spawn(new(big.Int).Div(n, d), result)
}

func (p *ParcsPollrad) Run(n *big.Int) []*big.Int {
	result := make(chan *big.Int, MaxDivs)
	p.runTask(n, result)

	for len(p.sem) > 0 {
		// Waiting
	}

	fmt.Println("Collecting results")
	close(result)

	res := make([]*big.Int, 0, MaxDivs)
	for r := range result {
		res = append(res, r)
	}

	return res
}

/*
func factor(n *big.Int) ([]*big.Int, error) {
	if n.Cmp(big.NewInt(1)) == 0 {
		return nil, nil
	}
	// Если n — вероятно простое, возвращаем его.
	if n.ProbablyPrime(20) { // 20 раундов ≈ 2‑64 ложноположительных
		return []*big.Int{new(big.Int).Set(n)}, nil
	}

	// Ищем нетривиальный делитель.
	d, err := pollardRho(n)
	if err != nil {
		return nil, err
	}
	// Рекурсивно разлагаем d и n/d.
	left, err := factor(d)
	if err != nil {
		return nil, err
	}
	right, err := factor(new(big.Int).Div(n, d))
	if err != nil {
		return nil, err
	}
	return append(left, right...), nil
}
*/

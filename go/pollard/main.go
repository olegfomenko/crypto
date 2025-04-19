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

type ParcsPollard struct {
	limit uint32
	sem   chan struct{}
}

func NewParcsPollard(limit uint32) *ParcsPollard {
	return &ParcsPollard{
		limit: limit,
		sem:   make(chan struct{}, limit),
	}
}

func (p *ParcsPollard) runTask(n *big.Int, result chan *big.Int) {
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
	fmt.Println("Running Pollard for n=", n)

	d := Pollard(n)

	spawn(d, result)
	spawn(new(big.Int).Div(n, d), result)
}

func (p *ParcsPollard) Run(n *big.Int) []*big.Int {
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

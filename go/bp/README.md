# Bulletproofs: Short Proofs for Confidential Transactions

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This package contains the implementation of [Bulletproofs](https://eprint.iacr.org/2017/1066.pdf) - the logarithmic size
range proofs.

The [main.go](./main.go) contains the primary implementation of bulletproof and inner product argument proof that uses
Fiat-Shamir heuristic to make proof non-interactive.

The [main_test.go](./main_test.go) contains the example of usage of the primary implementation.

The [docs_test.go](./docs_test.go) contains several implementation of a word by word approach defined in 3-4.2
paragraphs of original doc.

## Usage

Explore [main_test.go](./main_test.go) `TestBulletProof` with example of usage.


# The zkSNARK (Pinocchio) repo on Golang

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## General

The following implementation works over bn128 paring field. Explore [Cloudflare bn256 Bilinear map implementation](https://github.com/cloudflare/bn256) where the implementation is defined. 
Some functions need to be provided both for G1 and G2 computations.
The proof contains the following fields:

```go
type Proof struct {
    G1_L       *bn256.G1
    G2_L       *bn256.G2
    G2_R       *bn256.G2
    G2_O       *bn256.G2
    G2_alpha_L *bn256.G2
    G2_alpha_R *bn256.G2
    G2_alpha_O *bn256.G2
    G2_h       *bn256.G2
}
````

The following methods expected to be used for calculations:
1. Setup function that creates setup params.
    ```go
        func Setup(l1 L1, l2 L2, r R2, o O2, n uint64) *SetupParams
    ```
   
2. Proof function
    ```go
        func MakeProof(params *SetupParams, bigL1 BigL1, bigL2 BigL2, bigR BigR2, bigO BigO2, h H2) *Proof
    ```
   
3. Verify function
    ```go
        func VerifyProof(params *SetupParams, proof *Proof) error
    ```

## Example

Explore [Test example](./example_test.go) to see example of proving that we know `a`, such that f(1, a, 2) = 8
```
    f(w,a,b) = w? (a * b) : (a + b)
```

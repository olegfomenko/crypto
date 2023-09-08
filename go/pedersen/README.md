# Back-Maxwell range proof for Pedersen Commitments on Go 

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Example implementation of [Back-Maxwell Rangeproof](https://blockstream.com/bitcoin17-final41.pdf) on Go 
for creating the Pedersen commitment with corresponding proof that committed value lies in [0..2^n-1] range.   
The implementation uses Ethereum bn128 G1 curve to produce commitments and proofs. 

## Usage
Explore [main_test.go](./main_test.go) `TestPedersenCommitment` with example of usage.

Note, that there are the following values defined in global space to be changed on your choice:

```go
var G *bn256.G1
var H *bn256.G1

// Hash function that should return the value in Curve.N field
var Hash func(...[]byte) *big.Int = defaultHash
```

## Schnorr Signature
Explore [main_test.go](./main_test.go) `TestSchnorrSignatureAggregation` with an example of Schnorr signature. 
It can be useful to sign the resulting C=C1-C2 commitment in transactions. 

It uses the scheme from [Schnorr Signature](https://mareknarozniak.com/2021/05/25/schnorr-signature/) article.

# Crypto

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Crypto library by Oleg Fomenko.

Includes:

- Golang implementations of different crypto algorithms:
    1. [Diffie Hellman](./go/deffie-hellman)
    2. [Dynamic Merkle tree (base on Treap)](./go/dynamic-merkle)
    3. [Pedersen commitment (with Back-Maxwell rangeproof)](./go/pedersen)
    4. [Schnorr signature over bn256 curve](./go/schnorr-bn256)
    5. [ZK-SNARK (Pinocchio protocol) basic implementation](./go/zk-snark)
    6. [Pedersen commitment on Circom circuits](./go/pedersen-circom)
    7. [Elliptic curve](./go/ec)
    8. [El-Gamal](./go/el-gamal)
    9. [Math (crypto)](./go/math)
    10. [RSA](./go/rsa) 
    11. [Bulletproofs](./go/bp)
    12. [Bulletproofs++](./go/bppp)
    13. [Pedersen commitment on GNARK circuits](./go/pedersen-gnark) 
    14. [FFT](./go/fft) 

- Circom circuits:
    1. [Schnorr signature](./circuits/schnorr)
    2. [Pedersen commitment](./circuits/pedersen)
    3. [Merkle tree](./circuits/merkle)

- SageMath implementations of different crypto algorithms:
    1. [ZK-STARK](./sage/zk-stark)
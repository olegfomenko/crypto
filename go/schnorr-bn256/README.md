# Schnorr signature over bn256 curve

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)


PubKey is `rG`.
Signing:

- `k = rand()`
- `R = kG`
- `s = k + Hash(msg|P|R)*r`
Sig = `<s, R>`

Verification:
- Check that `sG` equal to `R + Hash(msg|P|R)*P`

Description:
`s = k + hash*prv`, so `sG = G(k + hash*prv)`
`R = kG`, so `R + hash*P = kG + hash*prv*G = G (k + hash*prv)`

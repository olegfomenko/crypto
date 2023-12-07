# Elliptic curve

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This package contains Go (only) implementation of elliptic curve with equation: y^2 =x^3 + ax + b. It also implements
the standard
Go `elliptic.Curve` interface.

The `SECP256K1()` function returns curve instance with parameters
from [SEC 2: Recommended Elliptic Curve Domain Parameters](https://www.secg.org/SEC2-Ver-1.0.pdf).

In the [tests](./main_test.go) you can find the comparison of this implementation and popular Ethereum implementation.

[More about (Wiki)](https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication)
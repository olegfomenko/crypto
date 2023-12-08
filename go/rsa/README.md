# RSA

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This package implement the RSA cryptosystem. 
Private keys are generated using `math.TestPrime` for random value.

The following constant in [main.go](./main.go) file defines the size of p,q random primes in bytes:
```go
const size = 128
```

The size of n and phi(n) will be 2*n.

Example of usage:
```go
    key, err := GeneratePrivateKey()
    if err != nil {
        panic(err)
    }

    msg := new(big.Int).SetBytes([]byte("Hello World")) // should be in [1...n). Check for bytes count will be enaught.
	
	// Encryption
    cypher := Encrypt(msg, key.PublicKey)

    // Decryption
    text := Decrypt(cypher, key)
	
	
    fmt.Println("Original message:", string(msg.Bytes()))
    fmt.Println("Encrypted message:", cypher)
    fmt.Println("Decrypted message:", string(text.Bytes()))
```
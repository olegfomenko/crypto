package schnorr_bn256

import (
	"errors"
	"github.com/cloudflare/bn256"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

type (
	SchnorrSignature struct {
		R *bn256.G1
		S *big.Int
	}

	HashF func(bytes ...[]byte) []byte
)

var (
	Hash            HashF = crypto.Keccak256
	ErrFailedRandom       = errors.New("failed to generate secure random")
)

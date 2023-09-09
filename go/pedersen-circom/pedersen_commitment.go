// Package pedersen
// Copyright 2032 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package pedersen

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-rapidsnark/prover"
	zkptypes "github.com/iden3/go-rapidsnark/types"
	"github.com/iden3/go-rapidsnark/verifier"
	"github.com/iden3/go-rapidsnark/witness/v2"
	"github.com/iden3/go-rapidsnark/witness/wazero"
	"github.com/olegfomenko/crypto/go/pedersen-circom/circuit"
	"github.com/olegfomenko/crypto/go/pedersen-circom/types"
	"github.com/pkg/errors"
)

// NewRandomPrivateKey generates random private key using secure `crypto/rand`.
// The key length is 30 bytes because of circuit limitations.
func NewRandomPrivateKey() *types.BigInt {
	var prv [30]byte
	_, err := rand.Read(prv[:])
	if err != nil {
		panic(err)
	}

	return types.Wrap(new(big.Int).SetBytes(prv[:]))
}

// NewCommitment creates a new Pedersen commitment for the provided amount and private key.
// Returns commitment and ZK proof structures.
func NewCommitment(prv, amount *types.BigInt) (*types.Commitment, *zkptypes.ZKProof, error) {
	input := &types.InputData{
		R: prv,
		A: amount,
	}

	if err := input.Validate(); err != nil {
		return nil, nil, errors.Wrap(err, "input data validation failed")
	}

	proof, err := GenerateCommitmentProof(input)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to generate zk proof")
	}

	return types.NewCommitmentKeyFromInput(input), proof, nil
}

// GenerateCommitmentProof generates proof using circuit in `./circuit` (`.wasm` and `.zkey` files) for the provided input data.
// Uses `github.com/iden3/go-rapidsnark/witness/wazero` as a wasm engine.
// Using other engines from `go-rapidsnark` is not recommended because of possible memory leaks.
func GenerateCommitmentProof(input *types.InputData) (*zkptypes.ZKProof, error) {
	wasm, err := circuit.PedersenWASM.ReadFile(circuit.PedersenWASMFileName)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to parse: %s", circuit.PedersenWASMFileName))
	}

	zkey, err := circuit.PedersenZKEY.ReadFile(circuit.PedersenZKEYFileName)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to parse: %s", circuit.PedersenZKEYFileName))
	}

	wtnsCalculator, err := witness.NewCalculator(
		wasm,
		witness.WithWasmEngine(wazero.NewCircom2WZWitnessCalculator),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize witness calculator")
	}

	wtnsBytes, err := wtnsCalculator.CalculateWTNSBin(input.ToMap(), true)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate witness")
	}

	proof, err := prover.Groth16Prover(zkey, wtnsBytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate proof")
	}

	return proof, nil
}

// VerifyCommitmentProof verifies provided proof using verification key in `./circuit`.
// Also checks that public inputs corresponds to the provided commitment.
func VerifyCommitmentProof(commitment *types.Commitment, proof *zkptypes.ZKProof) error {
	if ln := len(proof.PubSignals); ln != 2 {
		return errors.New(fmt.Sprintf("invalid pub signals arr len: is %d should be %d", ln, 2))
	}

	point := &babyjub.Point{X: types.MustBigInt(proof.PubSignals[0]), Y: types.MustBigInt(proof.PubSignals[1])}

	if !BJJPointEqual(point, commitment.Public) {
		return errors.New(fmt.Sprintf("commitment does not correspond providede proof: invalid public input"))
	}

	verificationKey, err := circuit.VerificationKey.ReadFile(circuit.VerificationKeyFileName)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to parse: %s", circuit.VerificationKeyFileName))
	}

	if err := verifier.VerifyGroth16(*proof, verificationKey); err != nil {
		return errors.Wrap(err, "failed to verify generated proof")
	}
	return nil
}

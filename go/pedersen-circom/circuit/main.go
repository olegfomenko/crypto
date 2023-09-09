// Package circuit
// Copyright 2023 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package circuit

import "embed"

//go:embed pedersen.wasm
var PedersenWASM embed.FS

const PedersenWASMFileName = "pedersen.wasm"

//go:embed pedersen.zkey
var PedersenZKEY embed.FS

const PedersenZKEYFileName = "pedersen.zkey"

//go:embed verification_key.json
var VerificationKey embed.FS

const VerificationKeyFileName = "verification_key.json"

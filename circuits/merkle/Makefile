test: compile prove verify

compile: remove prepare build witness ptau setup verification_key solidity

remove:
	rm -rf ./generated

prepare:
	echo "Creating dir" && \
    mkdir generated

build:
	echo "Compiling scheme" && \
    circom ./merkle.circom --r1cs --wasm --sym -o=./generated && \
    echo "INFO:" && \
    snarkjs r1cs info ./generated/merkle.r1cs

witness:
	echo "Generating witness" && \
	node ./generated/merkle_js/generate_witness.js ./generated/merkle_js/merkle.wasm ./input.json ./generated/merkle_js/witness.wtns

ptau:
	echo "Generating powers of tau: phase #1" && \
	snarkjs powersoftau new bn128 12 ./generated/merkle_js/ptau_000.ptau && \
	snarkjs powersoftau contribute ./generated/merkle_js/ptau_000.ptau ./generated/merkle_js/ptau.ptau --name="First contribution"
	echo "Generating powers of tau: phase #2" && \
	snarkjs powersoftau prepare phase2 ./generated/merkle_js/ptau.ptau ./generated/merkle_js/final.ptau

setup:
	echo "Setup groth16" && \
	snarkjs groth16 setup ./generated/merkle.r1cs ./generated/merkle_js/final.ptau ./generated/merkle_js/merkle_000.zkey
	snarkjs zkey contribute ./generated/merkle_js/merkle_000.zkey ./generated/merkle_js/merkle.zkey --name="First contribution"


verification_key:
	echo "Exporting verification key" && \
	snarkjs zkey export verificationkey ./generated/merkle_js/merkle.zkey ./generated/merkle_js/verification_key.json

prove:
	echo "Prooving" && \
	snarkjs groth16 prove ./generated/merkle_js/merkle.zkey ./generated/merkle_js/witness.wtns ./generated/proof.json ./generated/public.json

verify:
	echo "Verifing" && \
	snarkjs groth16 verify ./generated/merkle_js/verification_key.json ./generated/public.json ./generated/proof.json

solidity:
	echo "Generating .sol" && \
	snarkjs zkey export solidityverifier ./generated/merkle_js/merkle.zkey verifier.sol
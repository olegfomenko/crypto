test: compile prove verify

compile: remove prepare build witness ptau setup verification_key solidity

remove:
	rm -rf ./generated

prepare:
	echo "Creating dir" && \
    mkdir generated

build:
	echo "Compiling scheme" && \
    circom ./pedersen.circom --r1cs --wasm --sym -o=./generated && \
    echo "INFO:" && \
    snarkjs r1cs info ./generated/pedersen.r1cs

witness:
	echo "Generating witness" && \
	node ./generated/pedersen_js/generate_witness.js ./generated/pedersen_js/pedersen.wasm ./input.json ./generated/pedersen_js/witness.wtns

ptau:
	echo "Generating powers of tau: phase #1" && \
	snarkjs powersoftau new bn128 12 ./generated/pedersen_js/ptau_000.ptau && \
	snarkjs powersoftau contribute ./generated/pedersen_js/ptau_000.ptau ./generated/pedersen_js/ptau.ptau --name="First contribution"
	echo "Generating powers of tau: phase #2" && \
	snarkjs powersoftau prepare phase2 ./generated/pedersen_js/ptau.ptau ./generated/pedersen_js/final.ptau

setup:
	echo "Setup groth16" && \
	snarkjs groth16 setup ./generated/pedersen.r1cs ./generated/pedersen_js/final.ptau ./generated/pedersen_js/pedersen_000.zkey
	snarkjs zkey contribute ./generated/pedersen_js/pedersen_000.zkey ./generated/pedersen_js/pedersen.zkey --name="First contribution"


verification_key:
	echo "Exporting verification key" && \
	snarkjs zkey export verificationkey ./generated/pedersen_js/pedersen.zkey ./generated/pedersen_js/verification_key.json

prove:
	echo "Prooving" && \
	snarkjs groth16 prove ./generated/pedersen_js/pedersen.zkey ./generated/pedersen_js/witness.wtns ./generated/proof.json ./generated/public.json

verify:
	echo "Verifing" && \
	snarkjs groth16 verify ./generated/pedersen_js/verification_key.json ./generated/public.json ./generated/proof.json

solidity:
	echo "Generating .sol" && \
	snarkjs zkey export solidityverifier ./generated/pedersen_js/pedersen.zkey verifier.sol
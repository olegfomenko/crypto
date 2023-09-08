test: compile prove verify

compile: remove prepare build witness ptau setup verification_key solidity

remove:
	rm -rf ./generated

prepare:
	echo "Creating dir" && \
    mkdir generated

build:
	echo "Compiling scheme" && \
    circom ./schnorr.circom --r1cs --wasm --sym -o=./generated && \
    echo "INFO:" && \
    snarkjs r1cs info ./generated/schnorr.r1cs

witness:
	echo "Generating witness" && \
	node ./generated/schnorr_js/generate_witness.js ./generated/schnorr_js/schnorr.wasm ./input.json ./generated/schnorr_js/witness.wtns

ptau:
	echo "Generating powers of tau: phase #1" && \
	snarkjs powersoftau new bn128 12 ./generated/schnorr_js/ptau_000.ptau && \
	snarkjs powersoftau contribute ./generated/schnorr_js/ptau_000.ptau ./generated/schnorr_js/ptau.ptau --name="First contribution"
	echo "Generating powers of tau: phase #2" && \
	snarkjs powersoftau prepare phase2 ./generated/schnorr_js/ptau.ptau ./generated/schnorr_js/final.ptau

setup:
	echo "Setup groth16" && \
	snarkjs groth16 setup ./generated/schnorr.r1cs ./generated/schnorr_js/final.ptau ./generated/schnorr_js/schnorr_000.zkey
	snarkjs zkey contribute ./generated/schnorr_js/schnorr_000.zkey ./generated/schnorr_js/schnorr.zkey --name="First contribution"


verification_key:
	echo "Exporting verification key" && \
	snarkjs zkey export verificationkey ./generated/schnorr_js/schnorr.zkey ./generated/schnorr_js/verification_key.json

prove:
	echo "Prooving" && \
	snarkjs groth16 prove ./generated/schnorr_js/schnorr.zkey ./generated/schnorr_js/witness.wtns ./generated/proof.json ./generated/public.json

verify:
	echo "Verifing" && \
	snarkjs groth16 verify ./generated/schnorr_js/verification_key.json ./generated/public.json ./generated/proof.json

solidity:
	echo "Generating .sol" && \
	snarkjs zkey export solidityverifier ./generated/schnorr_js/schnorr.zkey verifier.sol
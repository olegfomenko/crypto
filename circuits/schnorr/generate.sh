echo "Creating dir" && \
mkdir generated && \

echo "Compaling scheme" && \
circom schnorr.circom --r1cs --wasm --sym -o=./generated && \

echo "INFO:" && \
snarkjs r1cs info ./generated/schnorr.r1cs && \

echo "Generating witness" && \
node ./generated/schnorr_js/generate_witness.js ./generated/schnorr_js/schnorr.wasm ./input.json ./generated/schnorr_js/witness.wtns && \

echo "Generating powers of tau: phase #1" && \
snarkjs powersoftau new bn128 12 ./generated/schnorr_js/ptau.ptau && \

echo "Generating powers of tau: phase #2" && \
snarkjs powersoftau prepare phase2 ./generated/schnorr_js/ptau.ptau ./generated/schnorr_js/final.ptau && \

echo "Setup groth16" && \
snarkjs groth16 setup ./generated/schnorr.r1cs ./generated/schnorr_js/final.ptau ./generated/schnorr_js/circuit.zkey && \

echo "Exporting verification key" && \
snarkjs zkey export verificationkey ./generated/schnorr_js/circuit.zkey ./generated/schnorr_js/verification_key.json && \

echo "Prooving" && \
snarkjs groth16 prove ./generated/schnorr_js/circuit.zkey ./generated/schnorr_js/witness.wtns proof.json public.json && \

echo "Verifing" && \
snarkjs groth16 verify ./generated/schnorr_js/verification_key.json public.json proof.json 

echo "Generating .sol" && \
snarkjs zkey export solidityverifier ./generated/schnorr_js/circuit.zkey verifier.sol
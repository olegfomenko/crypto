pragma circom 2.1.0;

include "../../node_modules/circomlib/circuits/poseidon.circom";


// From https://github.com/tornadocash/tornado-core/blob/master/circuits/merkleTree.circom#L18C1-L26C2
// if s == 0 returns [in[0], in[1]]
// if s == 1 returns [in[1], in[0]]
template Selector() {
    signal input in[2];
    signal input s;
    signal output out[2];

    s * (1 - s) === 0;
    out[0] <== (in[1] - in[0])*s + in[0];
    out[1] <== (in[0] - in[1])*s + in[1];
}

template ProveMerkle(N) {
    // Merkle path
    signal input in[N];

    // 0 if Hash(current, in[i])
    // 1 if Hash(in[i], current)
    // indices[0] no difference; 
    signal input indices[N];

    // Resulting merkle root
    signal output root;

    component hash[N];
    component selector[N];
    var current = in[0];

    for(var i = 1; i < N; i++) {
        hash[i] = Poseidon(2);

        selector[i] = Selector();
        selector[i].in[0] <== current;
        selector[i].in[1] <== in[i];
        selector[i].s <== indices[i];

        hash[i].inputs[0] <== selector[i].out[0];
        hash[i].inputs[1] <== selector[i].out[1];
        current = hash[i].out;
    }

    root <== current;
}

component main = ProveMerkle(2);
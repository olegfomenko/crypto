pragma circom 2.1.0;

include "../../node_modules/circomlib/circuits/babyjub.circom";
include "../../node_modules/circomlib/circuits/bitify.circom";
include "../../node_modules/circomlib/circuits/escalarmulany.circom";
include "../../node_modules/circomlib/circuits/poseidon.circom";

template ECMulBase(n) {
    signal input scalar;
    signal output out[2];

    component scallarBits = Num2Bits(n);
    scallarBits.in <== scalar;

    var BASE[2] = [
        5299619240641551281634865583518297030282874472190772894086521144482721001553,
        16950150798460657717958625567821834550301663161624707787222815936182638968203
    ];

    component mulFix = EscalarMulFix(n, BASE);

    var i;
    for (i = 0; i < n; i++) {
        mulFix.e[i] <== scallarBits.out[i];
    }

    out[0] <== mulFix.out[0];
    out[1] <== mulFix.out[1];
}

template ECMul(n) {
    signal input scalar;
    signal input point[2];
    signal output out[2];
    
    component scallarBits = Num2Bits(n);
    scallarBits.in <== scalar;

    component bjjMul = EscalarMulAny(n);

    var i;
    for (i = 0; i < n; i++) {
        bjjMul.e[i] <== scallarBits.out[i];
    }

    bjjMul.p[0] <== point[0];
    bjjMul.p[1] <== point[1];

    out[0] <== bjjMul.out[0];
    out[1] <== bjjMul.out[1];
}

template SchnorrVerification(n) {
    // Schnorr signature
    signal input rx;
    signal input ry;
    signal input s;

    // message Poseidon hash
    signal input msg;

    // public key
    signal input pkx;
    signal input pky;

    signal output x;
    signal output y;

    component poseidon = Poseidon(5);

    poseidon.inputs[0] <== msg;
    poseidon.inputs[1] <== pkx;
    poseidon.inputs[2] <== pky;
    poseidon.inputs[3] <== rx;
    poseidon.inputs[4] <== ry;

    var hash = poseidon.out;

    component ecMul = ECMul(n);
    ecMul.scalar <== hash;
    ecMul.point[0] <== pkx;
    ecMul.point[1] <== pky;

    component ecAdd = BabyAdd();
    ecAdd.x1 <== rx;
    ecAdd.y1 <== ry;

    ecAdd.x2 <== ecMul.out[0];
    ecAdd.y2 <== ecMul.out[1];

    component ecMulBase = ECMulBase(n);
    ecMulBase.scalar <== s;

    var x1 = ecMulBase.out[0];
    var y1 = ecMulBase.out[1];

    var x2 = ecAdd.xout;
    var y2 = ecAdd.yout;

    x1 === x2;
    y1 === y2;

    x <==  x1 - x2;
    y <==  y1 - y2;
}

component main = SchnorrVerification(255);
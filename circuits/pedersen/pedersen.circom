pragma circom 2.1.0;

include "../../node_modules/circomlib/circuits/babyjub.circom";
include "../../node_modules/circomlib/circuits/bitify.circom";
include "../../node_modules/circomlib/circuits/escalarmulfix.circom";

template ECMulBase(N) {
    signal input scalar;
    signal output out[2];

    component scallarBits = Num2Bits(N);
    scallarBits.in <== scalar;

    var BASE[2] = [
        5299619240641551281634865583518297030282874472190772894086521144482721001553,
        16950150798460657717958625567821834550301663161624707787222815936182638968203
    ];

    component mulFix = EscalarMulFix(N, BASE);

    var i;
    for (i = 0; i < N; i++) {
        mulFix.e[i] <== scallarBits.out[i];
    }

    out[0] <== mulFix.out[0];
    out[1] <== mulFix.out[1];
}

template ECMulH(N) {
    signal input scalar;
    signal output out[2];

    component scallarBits = Num2Bits(N);
    scallarBits.in <== scalar;

    var H[2] = [
        15334330715717027115948243110556436026028216985345384579806128223314358448928,
        14640338696677432581567520324796424956409796398271990973432884194068091890885
    ];

    component mulFix = EscalarMulFix(N, H);

    var i;
    for (i = 0; i < N; i++) {
        mulFix.e[i] <== scallarBits.out[i];
    }

    out[0] <== mulFix.out[0];
    out[1] <== mulFix.out[1];
}

template PedersenCommitment(N) {
    signal input r;
    signal input a;

    var ETH = 1000000000000000000;
    var MAX = 1000000 * ETH;

    component less = LessThan(N);
    less.in[0] <== a;
    less.in[1] <== MAX;

    less.out === 1;

    signal output x;
    signal output y;

    component mul_aH = ECMulH(N);
    mul_aH.scalar <== a;

    component mul_rG = ECMulBase(N);
    mul_rG.scalar <== r;

    component add = BabyAdd();
    add.x1 <== mul_aH.out[0];
    add.y1 <== mul_aH.out[1];

    add.x2 <== mul_rG.out[0];
    add.y2 <== mul_rG.out[1];

    x <== add.xout;
    y <== add.yout;
}

component main = PedersenCommitment(252);
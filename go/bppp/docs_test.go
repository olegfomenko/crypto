package bppp

import "C"
import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	"github.com/cloudflare/bn256"
)

type ACPublic struct {
	Nm, Nl, Nv, Nw, No int // Nw = Nm + Nm + No (for L, R, O parts)
	K                  int
	G                  *bn256.G1
	GVec               []*bn256.G1 // Nm
	HVec               []*bn256.G1 // Nv+7

	Wm [][]*big.Int // Nm * Nw
	Wl [][]*big.Int // Nl * Nw

	Am []*big.Int // Nm
	Al []*big.Int // Nl

	Fl bool
	Fm bool

	// Commitments
	V []*bn256.G1
}

type PartitionF = func(typ int, index int) *int // typ = 1:lo, 2:ll, 3:lr, 4:no

type AcPrivate struct {
	v  [][]*big.Int // Nv
	sv []*big.Int   // k
	wl []*big.Int   // Nm
	wr []*big.Int   // Nm
	wo []*big.Int   // No
	f  PartitionF
}

// Creates commits Co and Cl, also map input witness using f partition func
func CommitOL(public *ACPublic, wo, wl []*big.Int, f PartitionF) (ro []*big.Int, rl []*big.Int, no []*big.Int, nl []*big.Int, lo []*big.Int, ll []*big.Int, Co *bn256.G1, Cl *bn256.G1) {
	ro_ := values(6)
	rl_ := values(5)

	// contains random values, except several positions (described in 5.2.4)
	ro = []*big.Int{ro_[0], ro_[1], ro_[2], ro_[3], big.NewInt(0), ro_[4], ro_[5], big.NewInt(0)}        // 8
	rl = []*big.Int{rl_[0], rl_[1], rl_[2], big.NewInt(0), rl_[3], rl_[4], big.NewInt(0), big.NewInt(0)} // 8

	// nl == wl and nr == wr (described in 5.2.1)
	nl = wl // Nm

	no = make([]*big.Int, public.Nm)
	for j := range no {
		no[j] = big.NewInt(0)

		if i := f(4, j); i != nil {
			no[j].Set(wo[*i])
		}
	}

	lo = make([]*big.Int, public.Nv)
	for j := range lo {
		lo[j] = big.NewInt(0)

		if i := f(1, j); i != nil {
			lo[j].Set(wo[*i])
		}
	}

	ll = make([]*big.Int, public.Nv)
	for j := range lo {
		ll[j] = big.NewInt(0)

		if i := f(2, j); i != nil {
			ll[j].Set(wo[*i])
		}
	}

	Co = new(bn256.G1).ScalarMult(public.G, ro[0])
	Co.Add(Co, vectorPointScalarMul(public.HVec, append(ro[1:], lo...)))
	Co.Add(Co, vectorPointScalarMul(public.GVec, no))

	Cl = new(bn256.G1).ScalarMult(public.G, rl[0])
	Cl.Add(Cl, vectorPointScalarMul(public.HVec, append(rl[1:], ll...)))
	Cl.Add(Cl, vectorPointScalarMul(public.GVec, nl))

	return
}

func CommitR(public *ACPublic, wo, wr []*big.Int, f PartitionF) (rr []*big.Int, nr []*big.Int, lr []*big.Int, Cr *bn256.G1) {
	rr_ := values(4) // 4

	// contains random values, except several positions (described in 5.2.4)
	rr = []*big.Int{rr_[0], rr_[1], big.NewInt(0), rr_[2], rr_[3], big.NewInt(0), big.NewInt(0), big.NewInt(0)}

	// nl == wl and nr == wr (described in 5.2.1)
	nr = wr

	// Creates commits Cr also map input witness using f partition func
	lr = make([]*big.Int, public.Nv)
	for j := range lr {
		lr[j] = big.NewInt(0)

		if i := f(3, j); i != nil {
			lr[j].Set(wo[*i])
		}
	}

	Cr = new(bn256.G1).ScalarMult(public.G, rr[0])
	Cr.Add(Cr, vectorPointScalarMul(public.HVec, append(rr[1:], lr...)))
	Cr.Add(Cr, vectorPointScalarMul(public.GVec, nr))
	return
}

func ArithmeticCircuitProtocol(public *ACPublic, private *AcPrivate) {
	ro, rl, no, nl, lo, ll, Co, Cl := CommitOL(public, private.wo, private.wl, private.f)

	rr, nr, lr, Cr := CommitR(public, private.wo, private.wr, private.f)

	InnerArithmeticCircuitProtocol(public, private,
		[][]*big.Int{rl, rr, ro},
		[][]*big.Int{nl, nr, no},
		[][]*big.Int{ll, lr, lo},
		[]*bn256.G1{Cl, Cr, Co},
	)
}

func InnerArithmeticCircuitProtocol(public *ACPublic, private *AcPrivate, r, n, l [][]*big.Int, C []*bn256.G1) {
	rl := r[0] // 8
	rr := r[1] // 8
	ro := r[2] // 8

	ll := l[0] // Nv
	lr := l[1] // Nv
	lo := l[2] // Nv

	nl := n[0] // Nm
	nr := n[1] // Nm
	no := n[2] // Nm

	Cl := C[0]
	Cr := C[1]
	Co := C[2]

	// Send Cl, Cr, Co to verifier

	// Verifier generates challenges and sends to prover
	ch_ro := values(1)[0]
	ch_lambda := values(1)[0]
	ch_beta := values(1)[0]
	ch_delta := values(1)[0]

	// Prover and Verifier computes

	var MlnL [][]*big.Int // Nl * Nm
	for i := 0; i < public.Nl; i++ {
		MlnL = append(MlnL, public.Wl[i][:public.Nm])
	}

	var MmnL [][]*big.Int // Nm * Nm
	for i := 0; i < public.Nm; i++ {
		MmnL = append(MmnL, public.Wm[i][:public.Nm])
	}

	var MlnR [][]*big.Int // Nl*Nm
	for i := 0; i < public.Nl; i++ {
		MlnR = append(MlnR, public.Wl[i][public.Nm:public.Nm*2])
	}

	var MmnR [][]*big.Int // Nm*Nm
	for i := 0; i < public.Nm; i++ {
		MmnR = append(MmnR, public.Wm[i][public.Nm:public.Nm*2])
	}

	var WlO [][]*big.Int // Nl*No
	for i := 0; i < public.Nl; i++ {
		WlO = append(WlO, public.Wl[i][public.Nm*2:])
	}

	var WmO [][]*big.Int // Nm*No
	for i := 0; i < public.Nm; i++ {
		WmO = append(WmO, public.Wm[i][public.Nm*2:])
	}

	//ManO, a = {l,m}

	var MlnO [][]*big.Int // Nl*Nm
	for i := 0; i < public.Nl; i++ {
		MlnO = append(MlnO, make([]*big.Int, public.Nm))

		for j := 0; j < public.Nm; j++ {
			if j_ := private.f(4, j); j_ != nil {
				MlnO[i][j] = WlO[i][*j_]
			}
		}
	}

	var MmnO [][]*big.Int // Nm*Nm
	for i := 0; i < public.Nm; i++ {
		MmnO = append(MmnO, make([]*big.Int, public.Nm))

		for j := 0; j < public.Nm; j++ {
			if j_ := private.f(4, j); j_ != nil {
				MmnO[i][j] = WmO[i][*j_]
			}
		}
	}

	// MalX, a = {l,m}, X = {L,R,O}

	// L
	var MllL [][]*big.Int // Nl*Nv
	for i := 0; i < public.Nl; i++ {
		MllL = append(MllL, make([]*big.Int, public.Nv))

		for j := 0; j < public.Nv; j++ {
			if j_ := private.f(2, j); j_ != nil {
				MllL[i][j] = WlO[i][*j_]
			}
		}
	}

	var MmlL [][]*big.Int // Nm*Nv
	for i := 0; i < public.Nm; i++ {
		MmlL = append(MmlL, make([]*big.Int, public.Nv))

		for j := 0; j < public.Nv; j++ {
			if j_ := private.f(2, j); j_ != nil {
				MmlL[i][j] = WmO[i][*j_]
			}
		}
	}

	// R
	var MllR [][]*big.Int // Nl*Nv
	for i := 0; i < public.Nl; i++ {
		MllR = append(MllR, make([]*big.Int, public.Nv))

		for j := 0; j < public.Nv; j++ {
			if j_ := private.f(3, j); j_ != nil {
				MllR[i][j] = WlO[i][*j_]
			}
		}
	}

	var MmlR [][]*big.Int // Nm*Nv
	for i := 0; i < public.Nm; i++ {
		MmlR = append(MmlR, make([]*big.Int, public.Nv))

		for j := 0; j < public.Nv; j++ {
			if j_ := private.f(3, j); j_ != nil {
				MmlR[i][j] = WmO[i][*j_]
			}
		}
	}

	// O
	var MllO [][]*big.Int // Nl*Nv
	for i := 0; i < public.Nl; i++ {
		MllO = append(MllO, make([]*big.Int, public.Nv))

		for j := 0; j < public.Nv; j++ {
			if j_ := private.f(1, j); j_ != nil {
				MllO[i][j] = WlO[i][*j_]
			}
		}
	}

	var MmlO [][]*big.Int // Nm*Nv
	for i := 0; i < public.Nm; i++ {
		MmlO = append(MmlO, make([]*big.Int, public.Nv))

		for j := 0; j < public.Nv; j++ {
			if j_ := private.f(1, j); j_ != nil {
				MmlO[i][j] = WmO[i][*j_]
			}
		}
	}

	ch_mu := mul(ch_ro, ch_ro)

	// Calculate linear combination of V

	var V_ *bn256.G1 = &bn256.G1{} // set infinite

	for i := 0; i < public.K; i++ {
		V_ = V_.Add(V_, new(bn256.G1).ScalarMult(
			public.V[i],
			add(
				mul(bbool(public.Fl), pow(ch_lambda, mul(bint(public.Nv), bint(i)))),
				mul(bbool(public.Fm), pow(ch_mu, add(mul(bint(public.Nv), bint(i)), bint(1)))),
			),
		))
	}

	V_.ScalarMult(V_, bint(2))

	// Calculate lambda vector (nl == nv * k)

	e_lambda_nl := powvector(ch_lambda, public.Nl)

	e_lambda_nv := powvector(ch_lambda, public.Nv)
	e_mu_nv := powvector(ch_mu, public.Nv)

	ch_mu_nv := pow(ch_mu, bint(public.Nv))
	ch_lambda_nv := pow(ch_lambda, bint(public.Nv))

	e_mu_nv_k := powvector(ch_mu_nv, public.K)
	e_lambda_nv_k := powvector(ch_lambda_nv, public.K)

	lambda := vectorAdd(
		vectorTensorMul(vectorMulOnScalar(e_lambda_nv, ch_mu), e_mu_nv_k),
		vectorTensorMul(e_mu_nv, e_lambda_nv_k),
	)

	lambda = vectorMulOnScalar(lambda, bbool(public.Fl && public.Fm))
	lambda = vectorSub(e_lambda_nl, lambda) //Nl

	// Calculate mu vector
	e_mu_nm := powvector(ch_mu, public.Nm)
	mu := vectorMulOnScalar(e_mu_nm, ch_mu) // Nm

	muDiagInv := diagInv(ch_mu, public.Nm) // Nm*Nm

	// Calculate coefficients clX, X = {L,R,O}
	cnL := vectorMulOnMatrix(vectorAdd(vectorMulOnMatrix(lambda, MlnL), vectorMulOnMatrix(mu, MmnL)), muDiagInv) // Nm
	cnR := vectorMulOnMatrix(vectorAdd(vectorMulOnMatrix(lambda, MlnR), vectorMulOnMatrix(mu, MmnR)), muDiagInv) // Nm
	cnO := vectorMulOnMatrix(vectorAdd(vectorMulOnMatrix(lambda, MlnO), vectorMulOnMatrix(mu, MmnO)), muDiagInv) // Nm

	clL := vectorAdd(vectorMulOnMatrix(lambda, MllL), vectorMulOnMatrix(mu, MmlL)) // Nv
	clR := vectorAdd(vectorMulOnMatrix(lambda, MllR), vectorMulOnMatrix(mu, MmlR)) // Nv
	clO := vectorAdd(vectorMulOnMatrix(lambda, MllO), vectorMulOnMatrix(mu, MmlO)) // Nv

	// Define pn(T) and ps(T) as vector of polynom coefficient
	pnT := map[int][]*big.Int{ // 4 * Nm
		0: zeros(public.Nm),                      // * T^0
		1: cnR,                                   // * T^1
		2: cnL,                                   // * T^2
		3: vectorMulOnScalar(cnO, inv(ch_delta)), // * T^3
	}

	psT := polyAdd(polyVectorMulWeight(pnT, pnT, ch_mu), map[int]*big.Int{ // 4
		0: bint(0),
		1: bint(0),
		2: bint(0),
		3: add(vectorMul(mu, public.Am), vectorMul(lambda, public.Al)),
	})

	// Prover computes
	ls := values(public.Nv) // Nv
	ns := values(public.Nm) // Nm

	// Calc linear combination of v[][0]
	var v_ = bint(0)

	for i := 0; i < public.K; i++ {
		v_ = add(v_, mul(
			private.v[i][0],
			add(
				mul(bbool(public.Fl), pow(ch_lambda, mul(bint(public.Nv), bint(i)))),
				mul(bbool(public.Fm), pow(ch_mu, add(mul(bint(public.Nv), bint(i)), bint(1)))),
			),
		))
	}

	v_ = mul(v_, bint(2))

	cl_T := map[int][]*big.Int{ // 4
		0: zeros(public.Nv),
		1: clR,
		2: clL,
		3: vectorMulOnScalar(clO, inv(ch_delta)),
	}

	if public.Fm {
		cl_T[0] = vectorAdd(cl_T[0], append(vectorMulOnScalar(powvector(ch_mu, public.Nv)[1:], ch_mu), bint(0)))
	}

	if public.Fl {
		cl_T[0] = vectorAdd(cl_T[0], append(powvector(ch_lambda, public.Nv)[1:], bint(0)))
	}

	l_T := map[int][]*big.Int{ // 5*Nv
		-1: ls,
		0:  vectorMulOnScalar(lo, ch_delta),
		1:  ll,
		2:  lr,
	}

	l_T_3 := func() []*big.Int {
		l_T3 := zeros(public.Nv - 1)

		for i := 0; i < public.K; i++ {
			l_T3 = vectorAdd(l_T3, vectorMulOnScalar(
				private.v[i][1:],
				add(
					mul(bbool(public.Fl), pow(ch_lambda, mul(bint(public.Nv), bint(i)))),
					mul(bbool(public.Fm), pow(ch_mu, add(mul(bint(public.Nv), bint(i)), bint(1)))),
				),
			))
		}

		l_T3 = vectorMulOnScalar(l_T3, bint(2))
		return append(l_T3, bint(0))
	}()

	l_T[3] = l_T_3

	n_T := map[int][]*big.Int{ // 4 * Nm
		-1: ns,
		0:  vectorMulOnScalar(no, ch_delta),
		1:  nl,
		2:  nr,
	}

	nT := polyVectorAdd(pnT, n_T)

	f_T := polySub(psT, polyVectorMulWeight(nT, nT, ch_mu)) // 8
	f_T[3] = add(f_T[3], v_)
	f_T = polySub(f_T, polyVectorMul(cl_T, l_T))

	rv := zeros(8) // 8
	rv[1] = func() *big.Int {
		rv1 := bint(0)

		for i := 0; i < public.K; i++ {
			rv1 = add(rv1, mul(
				private.sv[i],
				add(
					mul(bbool(public.Fl), pow(ch_lambda, mul(bint(public.Nv), bint(i)))),
					mul(bbool(public.Fm), pow(ch_mu, add(mul(bint(public.Nv), bint(i)), bint(1)))),
				),
			))
		}

		return mul(rv1, bint(2))
	}()

	sr := []*big.Int{
		mul(ch_beta, mul(ch_delta, ro[1])),
		bint(0),
		add(mul(inv(ch_delta), mul(ch_delta, ro[0])), rl[1]),
		add(add(mul(ch_delta, ro[2]), mul(inv(ch_beta), rl[0])), rr[1]),
		add(add(add(mul(ch_delta, ro[3]), rl[2]), rv[1]), mul(inv(ch_beta), rr[0])),
		add(rl[4], rr[3]),
		add(mul(ch_delta, ro[5]), rr[4]),
		add(mul(ch_delta, ro[6]), rl[5]),
	}

	f_ := []*big.Int{f_T[0], f_T[-1], f_T[1], f_T[2], f_T[3], f_T[5], f_T[6], f_T[7]}

	rs := vectorSub( // 8
		append([]*big.Int{f_[0]}, vectorMulOnScalar(f_[1:], inv(ch_beta))...), // TODO wtf is this shit?
		sr,
	)

	Cs := new(bn256.G1).ScalarMult(public.G, rs[0])
	Cs.Add(Cs, vectorPointScalarMul(public.HVec, append(rs[1:], ls...)))
	Cs.Add(Cs, vectorPointScalarMul(public.GVec, ns))

	// Prover sends Cs to verifier

	// Verifier selects random t and sends to Prover
	t := values(1)[0]

	// Prover computes

	tinv := inv(t)
	t2 := mul(t, t)
	t3 := mul(t2, t)

	rT := vectorMulOnScalar(rs, tinv) // 8
	rT = vectorAdd(rT, vectorMulOnScalar(ro, ch_delta))
	rT = vectorAdd(rT, vectorMulOnScalar(rl, t))
	rT = vectorAdd(rT, vectorMulOnScalar(rr, t2))
	rT = vectorAdd(rT, vectorMulOnScalar(rv, t3))

	vT := polyCalc(psT, t) // will not be used by prover, but should be used by verifier
	vT = add(vT, mul(v_, t3))
	vT = add(vT, rT[0])

	lT := append(rT[1:], polyVectorCalc(l_T, t)...)

	// Prover and verifier computes
	PT := new(bn256.G1).ScalarMult(public.G, polyCalc(psT, t))
	PT.Add(PT, vectorPointScalarMul(public.GVec, polyVectorCalc(pnT, t)))

	cr_T := []*big.Int{bint(1), mul(ch_beta, tinv), mul(ch_beta, t), mul(ch_beta, t2), mul(ch_beta, t3), mul(ch_beta, mul(t2, t3)), mul(ch_beta, mul(t3, t3)), mul(ch_beta, mul(mul(t3, t), t3))} // 8

	cT := append(cr_T[1:], polyVectorCalc(cl_T, t)...)

	CT := new(bn256.G1).Add(PT, new(bn256.G1).ScalarMult(Cs, tinv))
	CT.Add(CT, new(bn256.G1).ScalarMult(Co, ch_delta))
	CT.Add(CT, new(bn256.G1).ScalarMult(Cl, t))
	CT.Add(CT, new(bn256.G1).ScalarMult(Cr, t2))
	CT.Add(CT, new(bn256.G1).ScalarMult(V_, t3))

	wnla(public.G, public.GVec, public.HVec, cT, CT, ch_ro, ch_mu, lT, polyVectorCalc(nT, t))
}

func TestWNLA(t *testing.T) {
	// Public
	const N = 4

	g := points(1)[0]
	G := points(N)

	H := points(N)

	c := values(N)

	ro := values(1)[0]
	mu := mul(ro, ro)

	// Private
	l := []*big.Int{big.NewInt(4), big.NewInt(5), big.NewInt(10), big.NewInt(1)}
	n := []*big.Int{big.NewInt(2), big.NewInt(1), big.NewInt(2), big.NewInt(10)}

	// Com
	v := add(vectorMul(c, l), weightVectorMul(n, n, mu))
	C := new(bn256.G1).ScalarMult(g, v)
	C.Add(C, vectorPointScalarMul(H, l))
	C.Add(C, vectorPointScalarMul(G, n))

	wnla(g, G, H, c, C, ro, mu, l, n)
}

func wnla(g *bn256.G1, G, H []*bn256.G1, c []*big.Int, C *bn256.G1, ro, mu *big.Int, l, n []*big.Int) {
	roinv := inv(ro)
	fmt.Println("Running WNLA protocol")

	if len(l)+len(n) < 6 {

		// Prover sends l, n to Verifier
		// Next verifier computes:
		_v := add(vectorMul(c, l), weightVectorMul(n, n, mu))

		_C := new(bn256.G1).ScalarMult(g, _v)
		_C.Add(_C, vectorPointScalarMul(H, l))
		_C.Add(_C, vectorPointScalarMul(G, n))

		if !bytes.Equal(_C.Marshal(), C.Marshal()) {
			panic("Failed to verify!")
		}

		fmt.Println("Verified!")
		return
	}

	// Verifier selects random challenge
	y := values(1)[0]

	// Prover calculates new reduced values, vx and vr and sends X, R to verifier
	c0, c1 := reduceVector(c)
	l0, l1 := reduceVector(l)
	n0, n1 := reduceVector(n)
	G0, G1 := reducePoints(G)
	H0, H1 := reducePoints(H)

	l_ := vectorAdd(l0, vectorMulOnScalar(l1, y))
	n_ := vectorAdd(vectorMulOnScalar(n0, roinv), vectorMulOnScalar(n1, y))

	//v_ := add(vectorMul(c_, l_), weightVectorMul(n_, n_, mul(mu, mu)))

	vx := add(
		mul(weightVectorMul(n0, n1, mul(mu, mu)), mul(big.NewInt(2), roinv)),
		add(vectorMul(c0, l1), vectorMul(c1, l0)),
	)

	vr := add(weightVectorMul(n1, n1, mul(mu, mu)), vectorMul(c1, l1))

	X := new(bn256.G1).ScalarMult(g, vx)
	X.Add(X, vectorPointScalarMul(H0, l1))
	X.Add(X, vectorPointScalarMul(H1, l0))
	X.Add(X, vectorPointScalarMul(G0, vectorMulOnScalar(n1, ro)))
	X.Add(X, vectorPointScalarMul(G1, vectorMulOnScalar(n0, roinv)))

	R := new(bn256.G1).ScalarMult(g, vr)
	R.Add(R, vectorPointScalarMul(H1, l1))
	R.Add(R, vectorPointScalarMul(G1, n1))

	// Submit R, X to Verifier

	// Both computes
	H_ := vectorPointsAdd(H0, vectorPointMulOnScalar(H1, y))
	G_ := vectorPointsAdd(vectorPointMulOnScalar(G0, ro), vectorPointMulOnScalar(G1, y))
	c_ := vectorAdd(c0, vectorMulOnScalar(c1, y))

	ro_ := mu
	mu_ := mul(mu, mu)

	C_ := new(bn256.G1).Set(C)
	C_.Add(C_, new(bn256.G1).ScalarMult(X, y))
	C_.Add(C_, new(bn256.G1).ScalarMult(R, sub(mul(y, y), big.NewInt(1))))

	// Recursive run
	wnla(g, G_, H_, c_, C_, ro_, mu_, l_, n_)
}

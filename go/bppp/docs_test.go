package bppp

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
	v  [][]*big.Int // k*Nv
	sv []*big.Int   // k
	wl []*big.Int   // Nm
	wr []*big.Int   // Nm
	wo []*big.Int   // No
	f  PartitionF
}

func TestACProtocol(t *testing.T) {
	// Scheme to proof that we know such p, q that:
	// pq = r
	// for some public r.

	// r = 15, p = 3, q = 5

	p := bint(3)
	q := bint(5)

	// Challenge x = 10

	// Wv = [-10, -100] = [-z, -z^2]
	Wv := []*big.Int{bint(-10), bint(-100)}

	v := []*big.Int{p, q}

	// Wl*w = M(Wl*al + Wr*ar + Wo*ao)
	// fl*wv+al = v+al = -Wl*w = -M(Wl*al + Wr*ar + Wo*ao) = -M(Wv*v+c)
	// v+al = -M*(Wv*v) - M*c
	// if M such that -M(Wv*v) = v then al = -M*c

	// Corresponding matrix M such that -M(Wv*v) = v
	m := []*big.Int{frac(3, 530), frac(5, 530)}

	fmt.Println("Check -M(Wv*v) =", vectorMulOnScalar(vectorMulOnScalar(m, bint(-1)), vectorMul(Wv, v)))

	al := vectorMulOnScalar(m, bint(-15*1000)) // -m * c = -m * (r * z^3)

	// Wlw = M(Wl*al + Wr*ar + Wo*ao)
	// Wl*al + Wr*ar + Wo*ao = -30 - 500 + 15000 = 14470
	// M(Wl*al + Wr*ar + Wo*ao) = [1447/101, 14470/101]

	Wlw := vectorMulOnScalar(m, bint(14470)) // 2

	fmt.Println("Wl*w =", Wlw)

	// Wl = M(Wl*al + Wr*ar + Wo*ao) * w'
	// where w' - right inverse for w
	w := []*big.Int{bint(3), bint(5), bint(15)}

	// left inverse w = [3/259, 5/259, 15/259]
	wInv := []*big.Int{frac(3, 259), frac(5, 259), frac(15, 259)} // 3
	fmt.Println("Check w*w' =", vectorMul(w, wInv))               // PASS

	var Wl [][]*big.Int = make([][]*big.Int, 2)
	for i := range Wl {
		Wl[i] = make([]*big.Int, 3)

		for j := range Wl[i] {
			Wl[i][j] = mul(Wlw[i], wInv[j])
		}
	}

	{
		check := zeros(2)

		for i := range Wl {
			check[i] = vectorMul(Wl[i], w)
		}

		fmt.Println("Check circuit:", vectorAdd(check, vectorAdd(v, al)))
	}

	Wm := [][]*big.Int{
		{bint(0), bint(0), bint(1)},
	} // [0, 0, 1]

	fmt.Println("Wm*w =", vectorMul(w, Wm[0]))

	public := ACPublic{
		Nm:   1,
		Nl:   2,
		Nv:   2,
		Nw:   3,
		No:   1,
		K:    1,
		G:    points(1)[0],
		GVec: points(1),
		HVec: points(7 + 2),
		Wm:   Wm,
		Wl:   Wl,
		Am:   zeros(1),
		Al:   al,
		Fl:   true,
		Fm:   false,
	}

	private := AcPrivate{
		v:  [][]*big.Int{{p, q}},
		sv: values(1),
		wl: []*big.Int{p},
		wr: []*big.Int{q},
		wo: []*big.Int{mul(p, q)},
		f: func(typ int, index int) *int {
			if typ == 4 { // map all to no
				return &index
			}

			return nil
		},
	}

	ArithmeticCircuitProtocol(&public, &private)
}

func ArithmeticCircuitProtocol(public *ACPublic, private *AcPrivate) {
	public.V = make([]*bn256.G1, public.K)
	for i := range public.V {
		public.V[i] = Com(private.v[i], private.sv[i], public.G, public.HVec)
	}

	ro, rl, no, nl, lo, ll, Co, Cl := CommitOL(public, private.wo, private.wl, private.f)

	rr, nr, lr, Cr := CommitR(public, private.wo, private.wr, private.f)

	InnerArithmeticCircuitProtocol2(public, private,
		[][]*big.Int{rl, rr, ro},
		[][]*big.Int{nl, nr, no},
		[][]*big.Int{ll, lr, lo},
		[]*bn256.G1{Cl, Cr, Co},
	)
}

func Com(v []*big.Int, s *big.Int, G *bn256.G1, H []*bn256.G1) *bn256.G1 {
	res := new(bn256.G1).ScalarMult(G, v[0])
	res.Add(res, new(bn256.G1).ScalarMult(H[0], s))
	res.Add(res, vectorPointScalarMul(H[8:], v[1:]))
	return res
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

	no = make([]*big.Int, public.Nm) // Nm
	for j := range no {
		no[j] = big.NewInt(0)

		if i := f(4, j); i != nil {
			no[j].Set(wo[*i])
		}
	}

	lo = make([]*big.Int, public.Nv) // Nv
	for j := range lo {
		lo[j] = big.NewInt(0)

		if i := f(1, j); i != nil {
			lo[j].Set(wo[*i])
		}
	}

	ll = make([]*big.Int, public.Nv) // Nv
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
	rr = []*big.Int{rr_[0], rr_[1], big.NewInt(0), rr_[2], rr_[3], big.NewInt(0), big.NewInt(0), big.NewInt(0)} // 8

	// nl == wl and nr == wr (described in 5.2.1)
	nr = wr

	// Creates commits Cr also map input witness using f partition func
	lr = make([]*big.Int, public.Nv) // Nv
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

func InnerArithmeticCircuitProtocol2(public *ACPublic, private *AcPrivate, r, n, l [][]*big.Int, C []*bn256.G1) {
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
			MlnO[i][j] = big.NewInt(0)

			if j_ := private.f(4, j); j_ != nil {
				MlnO[i][j].Set(WlO[i][*j_])
			}
		}
	}

	var MmnO [][]*big.Int // Nm*Nm
	for i := 0; i < public.Nm; i++ {
		MmnO = append(MmnO, make([]*big.Int, public.Nm))

		for j := 0; j < public.Nm; j++ {
			MmnO[i][j] = big.NewInt(0)

			if j_ := private.f(4, j); j_ != nil {
				MmnO[i][j].Set(WmO[i][*j_])
			}
		}
	}

	// MalX, a = {l,m}, X = {L,R,O}

	// L
	var MllL [][]*big.Int // Nl*Nv
	for i := 0; i < public.Nl; i++ {
		MllL = append(MllL, make([]*big.Int, public.Nv))

		for j := 0; j < public.Nv; j++ {
			MllL[i][j] = big.NewInt(0)

			if j_ := private.f(2, j); j_ != nil {
				MllL[i][j].Set(WlO[i][*j_])
			}
		}
	}

	var MmlL [][]*big.Int // Nm*Nv
	for i := 0; i < public.Nm; i++ {
		MmlL = append(MmlL, make([]*big.Int, public.Nv))

		for j := 0; j < public.Nv; j++ {
			MmlL[i][j] = big.NewInt(0)

			if j_ := private.f(2, j); j_ != nil {
				MmlL[i][j].Set(WmO[i][*j_])
			}
		}
	}

	// R
	var MllR [][]*big.Int // Nl*Nv
	for i := 0; i < public.Nl; i++ {
		MllR = append(MllR, make([]*big.Int, public.Nv))

		for j := 0; j < public.Nv; j++ {
			MllR[i][j] = big.NewInt(0)

			if j_ := private.f(3, j); j_ != nil {
				MllR[i][j].Set(WlO[i][*j_])
			}
		}
	}

	var MmlR [][]*big.Int // Nm*Nv
	for i := 0; i < public.Nm; i++ {
		MmlR = append(MmlR, make([]*big.Int, public.Nv))

		for j := 0; j < public.Nv; j++ {
			MmlR[i][j] = big.NewInt(0)

			if j_ := private.f(3, j); j_ != nil {
				MmlR[i][j].Set(WmO[i][*j_])
			}
		}
	}

	// O
	var MllO [][]*big.Int // Nl*Nv
	for i := 0; i < public.Nl; i++ {
		MllO = append(MllO, make([]*big.Int, public.Nv))

		for j := 0; j < public.Nv; j++ {
			MllO[i][j] = big.NewInt(0)

			if j_ := private.f(1, j); j_ != nil {
				MllO[i][j].Set(WlO[i][*j_])
			}
		}
	}

	var MmlO [][]*big.Int // Nm*Nv
	for i := 0; i < public.Nm; i++ {
		MmlO = append(MmlO, make([]*big.Int, public.Nv))

		for j := 0; j < public.Nv; j++ {
			MmlO[i][j] = big.NewInt(0)

			if j_ := private.f(1, j); j_ != nil {
				MmlO[i][j].Set(WmO[i][*j_])
			}
		}
	}

	//////////

	// Check M matrix calculated ok
	Wlw := vectorAdd(matrixMulOnVector(lo, MllO), matrixMulOnVector(no, MlnO))
	Wlw = vectorAdd(Wlw, vectorAdd(matrixMulOnVector(ll, MllL), matrixMulOnVector(nl, MlnL)))
	Wlw = vectorAdd(Wlw, vectorAdd(matrixMulOnVector(lr, MllR), matrixMulOnVector(nr, MlnR)))
	fmt.Println("Wl*w =", Wlw)

	Wmw := vectorAdd(matrixMulOnVector(lo, MmlO), matrixMulOnVector(no, MmnO))
	Wmw = vectorAdd(Wmw, vectorAdd(matrixMulOnVector(ll, MmlL), matrixMulOnVector(nl, MmnL)))
	Wmw = vectorAdd(Wmw, vectorAdd(matrixMulOnVector(lr, MmlR), matrixMulOnVector(nr, MmnR)))
	fmt.Println("Wm*w =", Wmw)

	/////////

	ch_mu := mul(ch_ro, ch_ro)

	lcomb := func(i int) *big.Int {
		return add(
			mul(bbool(public.Fl), pow(ch_lambda, mul(bint(public.Nv), bint(i)))),
			mul(bbool(public.Fm), pow(ch_mu, add(mul(bint(public.Nv), bint(i)), bint(1)))),
		)
	}

	// Calculate linear combination of V
	V_ := func() *bn256.G1 {
		var V_ = new(bn256.G1).ScalarBaseMult(bint(0)) // set infinite

		for i := 0; i < public.K; i++ {
			V_ = V_.Add(V_, new(bn256.G1).ScalarMult(
				public.V[i],
				lcomb(i),
			))
		}

		return V_.ScalarMult(V_, bint(2))
	}()

	// Calculate lambda vector (nl == nv * k)
	lambda := vectorAdd(
		vectorTensorMul(vectorMulOnScalar(powvector(ch_lambda, public.Nv), ch_mu), powvector(pow(ch_mu, bint(public.Nv)), public.K)),
		vectorTensorMul(powvector(ch_mu, public.Nv), powvector(pow(ch_lambda, bint(public.Nv)), public.K)),
	)

	lambda = vectorMulOnScalar(lambda, bbool(public.Fl && public.Fm))
	lambda = vectorSub(powvector(ch_lambda, public.Nl), lambda) //Nl

	// Calculate mu vector
	mu := vectorMulOnScalar(powvector(ch_mu, public.Nm), ch_mu) // Nm
	fmt.Println("mu =", mu)

	/////////

	// Check Eq. 34

	c34 := vectorMul(lambda, vectorAdd(vectorAdd(Wlw, []*big.Int{bint(3), bint(5)}), public.Al))
	c34 = add(c34, weightVectorMul(private.wl, private.wr, ch_mu))
	c34 = sub(c34, vectorMul(mu, Wmw))
	fmt.Println("Check Eq. 34 =", c34)

	///////////

	// Calculate coefficients clX, X = {L,R,O}
	muDiagInv := diagInv(ch_mu, public.Nm) // Nm*Nm

	// TODO maybe - instead of +
	cnL := vectorMulOnMatrix(vectorSub(vectorMulOnMatrix(lambda, MlnL), vectorMulOnMatrix(mu, MmnL)), muDiagInv) // Nm
	cnR := vectorMulOnMatrix(vectorSub(vectorMulOnMatrix(lambda, MlnR), vectorMulOnMatrix(mu, MmnR)), muDiagInv) // Nm
	cnO := vectorMulOnMatrix(vectorSub(vectorMulOnMatrix(lambda, MlnO), vectorMulOnMatrix(mu, MmnO)), muDiagInv) // Nm

	clL := vectorSub(vectorMulOnMatrix(lambda, MllL), vectorMulOnMatrix(mu, MmlL)) // Nv
	clR := vectorSub(vectorMulOnMatrix(lambda, MllR), vectorMulOnMatrix(mu, MmlR)) // Nv
	clO := vectorSub(vectorMulOnMatrix(lambda, MllO), vectorMulOnMatrix(mu, MmlO)) // Nv

	// Prover computes
	ls := values(public.Nv) // Nv
	ns := values(public.Nm) // Nm

	// Calc linear combination of v[][0]
	v_ := func() *big.Int {
		v_ := bint(0)

		for i := 0; i < public.K; i++ {
			v_ = add(v_, mul(
				private.v[i][0],
				lcomb(i),
			))
		}

		return mul(v_, bint(2))
	}()

	rv := zeros(8)            // 8
	rv[0] = func() *big.Int { // TODO maybe 0 instead of 1
		rv1 := bint(0)

		for i := 0; i < public.K; i++ {
			rv1 = add(rv1, mul(
				private.sv[i],
				lcomb(i),
			))
		}

		return mul(rv1, bint(2))
	}()

	// Calc linear combination of v[][1:]
	v_1 := func() []*big.Int {
		var v_1 = zeros(1)

		for i := 0; i < public.K; i++ {
			v_1 = vectorAdd(v_1, vectorMulOnScalar(
				private.v[i][1:],
				lcomb(i),
			))
		}

		return vectorMulOnScalar(v_1, bint(2))
	}()

	// Check V_ correctness
	{
		check := new(bn256.G1).ScalarMult(public.G, v_)
		check.Add(check, vectorPointScalarMul(public.HVec, append(rv, v_1...)))
		fmt.Println("Check V_ correct:", bytes.Equal(V_.Marshal(), check.Marshal()))
	}

	// Define f'(t):
	f_ := make(map[int]*big.Int)

	// calc ps(t)
	f_[3] = mul(bint(2), add(vectorMul(lambda, public.Al), vectorMul(mu, public.Am)))

	// calc pn(t)^2 by mu
	f_[6] = add(f_[6], weightVectorMul(vectorMulOnScalar(cnO, inv(ch_delta)), vectorMulOnScalar(cnO, inv(ch_delta)), ch_mu)) // 3+3
	f_[5] = sub(f_[5], mul(bint(2), weightVectorMul(vectorMulOnScalar(cnO, inv(ch_delta)), cnL, ch_mu)))                     // 3+2
	f_[4] = add(f_[4], mul(bint(2), weightVectorMul(vectorMulOnScalar(cnO, inv(ch_delta)), cnR, ch_mu)))                     // 3+1

	f_[4] = add(f_[4], weightVectorMul(cnL, cnL, ch_mu))               // 2+2
	f_[3] = sub(f_[3], mul(bint(2), weightVectorMul(cnL, cnR, ch_mu))) // 2+1

	f_[2] = weightVectorMul(cnR, cnR, ch_mu) // 1+1

	// calc v_*T^3
	f_[3] = add(f_[3], v_)

	// calc <cl_, l_>

	f_[2] = sub(f_[2], mul(bint(2), mul(inv(ch_delta), vectorMul(clO, ls))))  // 3-1
	f_[3] = add(f_[3], mul(bint(2), vectorMul(clO, lo)))                      // 3+0
	f_[4] = sub(f_[4], mul(bint(2), mul(inv(ch_delta), vectorMul(clO, ll))))  // 3+1
	f_[5] = add(f_[5], mul(bint(2), mul(inv(ch_delta), vectorMul(clO, lr))))  // 3+2
	f_[6] = sub(f_[6], mul(bint(2), mul(inv(ch_delta), vectorMul(clO, v_1)))) // 3+3 // TODO check dimension correspondence

	f_[1] = add(f_[1], mul(bint(2), vectorMul(clL, ls)))                // 2-1
	f_[2] = sub(f_[2], mul(bint(2), mul(ch_delta, vectorMul(clL, lo)))) // 2+0
	f_[3] = add(f_[3], mul(bint(2), vectorMul(clL, ll)))                // 2+1
	f_[4] = sub(f_[4], mul(bint(2), vectorMul(clL, lr)))                // 2+2
	f_[5] = add(f_[5], mul(bint(2), vectorMul(clL, v_1)))               // 2+3 // TODO check dimension correspondence

	f_[0] = sub(f_[0], vectorMul(clR, ls))                              // 1-1
	f_[1] = add(f_[1], mul(bint(2), mul(ch_delta, vectorMul(clR, lo)))) // 1+0
	f_[2] = sub(f_[2], mul(bint(2), vectorMul(clR, ll)))                // 1+1
	f_[3] = add(f_[3], mul(bint(2), vectorMul(clR, lr)))                // 1+2
	f_[4] = sub(f_[4], mul(bint(2), vectorMul(clR, v_1)))               // 1+3 // TODO check dimension correspondence

	cl0 := vectorAdd(
		vectorMulOnScalar(vectorMulOnScalar(powvector(ch_mu, public.Nv)[1:], ch_mu), bbool(public.Fm)),
		vectorMulOnScalar(powvector(ch_lambda, public.Nv)[1:], bbool(public.Fl)),
	) // TODO check dimension correspondence

	f_[-1] = add(f_[-1], vectorMul(cl0, ls))              // 0-1
	f_[0] = sub(f_[0], mul(ch_delta, vectorMul(cl0, lo))) // 0+0
	f_[1] = add(f_[1], vectorMul(cl0, ll))                // 0+1
	f_[2] = sub(f_[2], vectorMul(cl0, lr))                // 0+2
	f_[3] = add(f_[3], vectorMul(cl0, v_1))               // 0+3 // TODO check dimension correspondence

	// calc weight norm |n(t)|^2 for mu

	// n(t) = pn(t) + n_(t) =
	// (t-1*ns - delta*no + t *nl - t2*nr) + (delta-1*T3*cnO - t2*cnl + t*cnr) =
	// t-1*ns - delta*no + t * (nl + cnR) - t2 * (nr + cnL) + delta-1*T3*cnO

	// -1: ns
	// 0: -delta*no
	// 1: nl+cnR
	// 2: -nr-cnL
	// 3: delta-1*cnO

	f_[-2] = sub(f_[-2], weightVectorMul(ns, ns, ch_mu)) // -1-1

	f_[-1] = add(f_[-1], mul(weightVectorMul(ns, no, ch_mu), mul(bint(2), ch_delta)))     // -1+0
	f_[0] = sub(f_[0], mul(weightVectorMul(ns, vectorAdd(nl, cnR), ch_mu), bint(2)))      // -1+1
	f_[1] = add(f_[1], mul(weightVectorMul(ns, vectorAdd(nr, cnL), ch_mu), bint(2)))      // -1+2
	f_[2] = sub(f_[2], mul(weightVectorMul(ns, cnO, ch_mu), mul(bint(2), inv(ch_delta)))) // -1+3

	f_[0] = sub(f_[0], weightVectorMul(vectorMulOnScalar(no, ch_delta), vectorMulOnScalar(no, ch_delta), ch_mu)) // 0 + 0
	f_[1] = add(f_[1], mul(weightVectorMul(no, vectorAdd(nl, cnR), ch_mu), mul(bint(2), ch_delta)))              // 0+1
	f_[2] = sub(f_[2], mul(weightVectorMul(no, vectorAdd(nr, cnL), ch_mu), mul(bint(2), ch_delta)))              // 0+2
	f_[3] = add(f_[3], mul(weightVectorMul(no, cnO, ch_mu), bint(2)))                                            // 0+3

	f_[2] = sub(f_[2], weightVectorMul(vectorAdd(nl, cnR), vectorAdd(nl, cnR), ch_mu))                    // 1 + 1
	f_[3] = add(f_[3], mul(weightVectorMul(vectorAdd(nl, cnR), vectorAdd(nr, cnL), ch_mu), bint(2)))      // 1+2
	f_[4] = sub(f_[4], mul(weightVectorMul(vectorAdd(nl, cnR), cnO, ch_mu), mul(bint(2), inv(ch_delta)))) // 1+3

	f_[4] = sub(f_[4], weightVectorMul(vectorAdd(nr, cnL), vectorAdd(nr, cnL), ch_mu))                    // 2 + 2
	f_[5] = add(f_[5], mul(weightVectorMul(vectorAdd(nr, cnL), cnO, ch_mu), mul(bint(2), inv(ch_delta)))) // 2+3

	f_[6] = sub(f_[6], weightVectorMul(vectorMulOnScalar(cnO, inv(ch_delta)), vectorMulOnScalar(cnO, inv(ch_delta)), ch_mu)) // 3 + 3

	fmt.Println("f'(T) =", f_)
	fmt.Println("f'(T)[3] =", f_[3])

	// TODO should be chosen later!!
	t := values(1)[0]
	tinv := inv(t)
	t2 := mul(t, t)
	t3 := mul(t2, t)

	// TODO calc without T
	sr := vectorMulOnScalar(ro, minus(mul(t, ch_delta)))
	sr = vectorAdd(sr, vectorMulOnScalar(rl, t2))
	sr = vectorSub(sr, vectorMulOnScalar(rr, t3))
	sr = vectorAdd(sr, vectorMulOnScalar(rv, mul(t3, t)))

	//sr := []*big.Int{
	//	mul(ch_beta, mul(ch_delta, ro[1])),
	//	bint(0),
	//	add(mul(inv(ch_beta), mul(ch_delta, ro[0])), rl[1]),
	//	add(add(mul(ch_delta, ro[2]), mul(inv(ch_beta), rl[0])), rr[1]),
	//	add(add(add(mul(ch_delta, ro[3]), rl[2]), rv[0]), mul(inv(ch_beta), rr[0])),
	//	add(rl[4], rr[3]),
	//	add(mul(ch_delta, ro[5]), rr[4]),
	//	add(mul(ch_delta, ro[6]), rl[5]),
	//}

	fcoef := []*big.Int{f_[-1], f_[-2], f_[0], f_[1], f_[2], f_[4], f_[5], f_[6]}

	rs := vectorSub( // 8
		append([]*big.Int{fcoef[0]}, vectorMulOnScalar(fcoef[1:], inv(ch_beta))...),
		sr,
	)

	Cs := new(bn256.G1).ScalarMult(public.G, rs[0])
	Cs.Add(Cs, vectorPointScalarMul(public.HVec, append(rs[1:], ls...)))
	Cs.Add(Cs, vectorPointScalarMul(public.GVec, ns))

	// Prover sends Cs to verifier

	// Verifier selects random t and sends to Prover
	// Chose t here

	// Prover computes

	r0 := mul(rs[0], tinv)
	r0 = sub(r0, mul(ro[0], ch_delta))
	r0 = add(r0, mul(rl[0], t))
	r0 = sub(r0, mul(rr[0], t2))

	lT := vectorMulOnScalar(append(rs[1:], ls...), tinv)
	lT = vectorSub(lT, vectorMulOnScalar(append(ro[1:], lo...), ch_delta))
	lT = vectorAdd(lT, vectorMulOnScalar(append(rl[1:], ll...), t))
	lT = vectorSub(lT, vectorMulOnScalar(append(rr[1:], lr...), t2))
	lT = vectorAdd(lT, vectorMulOnScalar(append(rv, v_1...), t3))

	pnT := vectorMulOnScalar(cnR, t)
	pnT = vectorSub(pnT, vectorMulOnScalar(cnL, t2))
	pnT = vectorAdd(pnT, vectorMulOnScalar(cnO, mul(inv(ch_delta), t3)))

	psT := weightVectorMul(pnT, pnT, ch_mu)
	psT = add(psT, mul(bint(2), mul(vectorMul(lambda, public.Al), t3)))
	psT = add(psT, mul(bint(2), mul(vectorMul(mu, public.Am), t3)))

	n_T := vectorMulOnScalar(ns, tinv)
	n_T = vectorSub(n_T, vectorMulOnScalar(no, ch_delta))
	n_T = vectorAdd(n_T, vectorMulOnScalar(nl, t))
	n_T = vectorSub(n_T, vectorMulOnScalar(nr, t2))

	nT := vectorAdd(pnT, n_T)

	// Prover and verifier computes
	PT := new(bn256.G1).ScalarMult(public.G, psT)
	PT.Add(PT, vectorPointScalarMul(public.GVec, pnT))

	cr_T := []*big.Int{
		bint(1),
		mul(ch_beta, tinv),
		mul(ch_beta, t),
		mul(ch_beta, t2),
		mul(ch_beta, t3),
		mul(ch_beta, mul(t2, t3)),
		mul(ch_beta, mul(t3, t3)),
		mul(ch_beta, mul(mul(t3, t), t3)),
	} // 8

	cl_T := vectorMulOnScalar(clO, mul(t3, inv(ch_delta)))
	cl_T = vectorSub(cl_T, vectorMulOnScalar(clL, t2))
	cl_T = vectorAdd(cl_T, vectorMulOnScalar(clR, t))
	cl_T = vectorMulOnScalar(cl_T, bint(2))
	cl_T = vectorSub(cl_T, cl0)

	cT := append(cr_T[1:], cl_T...)

	CT := new(bn256.G1).Add(PT, new(bn256.G1).ScalarMult(Cs, tinv))
	CT.Add(CT, new(bn256.G1).ScalarMult(Co, minus(ch_delta)))
	CT.Add(CT, new(bn256.G1).ScalarMult(Cl, t))
	CT.Add(CT, new(bn256.G1).ScalarMult(Cr, minus(t2)))
	CT.Add(CT, new(bn256.G1).ScalarMult(V_, t3))

	vT := add(psT, mul(v_, t3))
	vT = add(vT, r0)
	fmt.Println(r0)

	// Check that calculated commitment equals to v*G + <l,H> + <n,G>
	{
		CTPrv := new(bn256.G1).ScalarMult(public.G, vT)
		CTPrv.Add(CTPrv, vectorPointScalarMul(public.HVec, lT))
		CTPrv.Add(CTPrv, vectorPointScalarMul(public.GVec, nT))

		fmt.Println("Check C(t) = ComNormInnerArg(l, n, v):", bytes.Equal(CT.Marshal(), CTPrv.Marshal()))
	}

	// Check 45 eq from doc 510
	{
		CLeft := new(bn256.G1).ScalarMult(Cs, tinv)
		CLeft.Add(CLeft, new(bn256.G1).ScalarMult(Co, minus(ch_delta)))
		CLeft.Add(CLeft, new(bn256.G1).ScalarMult(Cl, t))
		CLeft.Add(CLeft, new(bn256.G1).ScalarMult(Cr, minus(t2)))
		CLeft.Add(CLeft, new(bn256.G1).ScalarMult(V_, t3))

		CTRight := new(bn256.G1).ScalarMult(public.G, add(mul(v_, t3), r0))
		CTRight.Add(CTRight, vectorPointScalarMul(public.HVec, lT))
		CTRight.Add(CTRight, vectorPointScalarMul(public.GVec, n_T))
		fmt.Println("Check (45):", bytes.Equal(CLeft.Marshal(), CTRight.Marshal()))
	}

	fmt.Println("Should be WNLA secret: ", vT)
	wnla(public.G, public.GVec, public.HVec, cT, CT, ch_ro, ch_mu, lT, nT)
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
	fmt.Println("Running WNLA protocol... WNLA secret: ", add(vectorMul(c, l), weightVectorMul(n, n, mu)))

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

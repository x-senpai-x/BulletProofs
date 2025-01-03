// +build pedersonCommitment

package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

func GenerateNumsPoint (curve elliptic.Curve,G ecdsa.PublicKey,)  (ecdsa.PublicKey){
	data := append(G.X.Bytes(), G.Y.Bytes()...)
	sha256Hash := sha256.Sum256(data)
	var Hx, Hy *big.Int
	found := false
	counter := 0
	for !found {
		hashWithCounter := append(sha256Hash[:], byte(counter)) // Append counter to hash
		xCoord := new(big.Int).SetBytes(hashWithCounter)        // Convert to big.Int
		Hx, Hy = curve.ScalarBaseMult(xCoord.Bytes())           // Use the curve to derive a point
		if curve.IsOnCurve(Hx, Hy) {                            // Check if the point is valid
			found = true
		} else {
			counter++ // Increment counter to try a new hash
		}
	}

	if !found {
		fmt.Println("Failed to generate a valid point H on the curve.")
		return ecdsa.PublicKey{}
	}

	fmt.Println("Valid point H found on the curve.")
	H:=ecdsa.PublicKey{
		Curve: curve,
		X: Hx,
		Y:Hy,
	}
	return H
}
func GenerateCommitment(curve elliptic.Curve, r *big.Int,a *big.Int ,H ecdsa.PublicKey,G ecdsa.PublicKey)(*ecdsa.PublicKey){
	rHx, rHy := curve.ScalarMult(H.X, H.Y, r.Bytes()) // rH
	aGx, aGy := curve.ScalarMult(G.X, G.Y, a.Bytes()) // aG
	Cx, Cy := curve.Add(rHx, rHy, aGx, aGy)         // C = rH + aG
	fmt.Printf("Commitment C: (%s, %s)\n", Cx.String(), Cy.String())
	return &ecdsa.PublicKey{
		Curve: curve,
		X: Cx,
		Y: Cy,
	}
}

// Commitment represents a Pedersen commitment
// type Commitment struct {
// 	C *Point
// 	R *big.Int
// 	X *big.Int
// }

// type Point struct {
// 	X, Y *big.Int
// }

// C= rH+ aG
// Here, C is the curve point we will use as a commitment (and give to some
// counterparty), a is the value we commit to (the secret value) (assume from now on it’s a number,
// not a string), r is the randomness which provides hiding, G is as already men-
// tioned the publically agreed generator of the elliptic curve, and H is another
// curve point, for which nobody knows the discrete logarithm q s.t. H= qG.
// this new commitment scheme does have a homomorphism:
// C(r1,a1) + C(r2,a2) = r1H+ a1G+ r2H+ a2G
// = (r1 + r2)H+ (a1 + a2)G
// = C(r1 + r2,a1 + a2)
// take the encoding of the point G, in binary, perhaps
// compressed or uncompressed form, take the SHA256 of that binary string and
// treat it as the x-coordinate of an elliptic curve point
// func perdersonCommitment()


// // Generate H by hashing G.X and G.Y and finding a valid point
// // 	Not all 256 bit hash digests will be such x-coordinates, but about half
// // of them are, so you can just use some simple iterative algorithm (e.g. append
// // byte “1”, “2”, ... after the encoding of G) to create not just one such NUMS
// // point H, but a large set of them like H1,H2,...,Hn. And indeed we will make
// // heavy use of such “sets of pre-determined curve points for which no one knows
// // the relative discrete logs” later.
// 	H:=GenerateNumsPoint(curve,G)
// 	// Generate randomness r and secret witness a
// 	r, err := rand.Int(rand.Reader, curve.Params().N)
// 	if err != nil {
// 		fmt.Println("Error generating randomness r:", err)
// 		return
// 	}
// 	a := big.NewInt(42) // Secret witness
// 	C:=GenerateCommitment(curve,r,a,H,G) 
// 	fmt.Println(C);

// 	m:=3//number of commitments
// 	commitments := make([]*Commitment, m)

// }


func main() {
	// Initialize the elliptic curve
	curve := elliptic.P256()
	// Generate a key pair
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		fmt.Println("Error generating private key:", err)
		return
	}
// 	Our argument of knowledge will come after we have generated a set of com-
// mitments for each of m vectors x1,x2,...,xm, each of the same dimension N
// (̸= m)). Explicitly:
// C1 = r1H+ x1G
// C2 = r2H+ x2G
// ...
// Cm = rmH+ xmG
// • P → V: C0 (a new commitment to a newly chosen random vector of
// 	dimension N)
// 	• V→ P: e (a random scalar)
// 	• P→ V: (z,s) (a single vector of dimension N, and another scalar)
// z = Σ e^i * x_i,  s = Σ e^i * r_i
//the verifer of course needs to verify whether the proof is valid.
// Σ e^i * C_i ? = sH + zG
	G := privateKey.PublicKey // Public key as the generator point
	H:=GenerateNumsPoint(curve,G)
	// Generate randomness r and secret witness a
	r0, err := rand.Int(rand.Reader, curve.Params().N)
	if err != nil {
		fmt.Println("Error generating randomness r:", err)
		return
	}
	x0 := big.NewInt(42) // Secret witness
	C0:=GenerateCommitment(curve,r0,x0,H,G) 
	fmt.Println("Prover sends commitment C0 to Verifier:", C0)
	m:=3//number of commitments
	//Verifier generates e
	e, err := rand.Int(rand.Reader, curve.Params().N)
	if err != nil {
		fmt.Println("Error generating e:", err)
		return
	}
	z:=big.NewInt(0);
	s:=big.NewInt(0);
	Commitments:=make([]*ecdsa.PublicKey,m)
	rS:=make([]*big.Int,m)
	xS:=make([]*big.Int,m)
	ep := big.NewInt(1)
	LHSX:=big.NewInt(0);
	LHSY:=big.NewInt(0);
	for i:=0;i<3;i++ {
		rS[i], err = rand.Int(rand.Reader, curve.Params().N)
		if err != nil {
			fmt.Println("Error generating randomness r:", err)
			return
		}
		xS[i] = new(big.Int).Set(x0) // Initialize xS[i] with the value of x0
		xS[i].Add(xS[i], rS[i])  
		Commitments[i]=GenerateCommitment(curve,rS[i],xS[i],H,G)
		ep.Exp(e, big.NewInt(int64(i+1)), curve.Params().N) // e^i
		z.Add(z,new(big.Int).Mul(ep, xS[i]))
		s.Add(s,new(big.Int).Mul(ep,rS[i]))
		eCiX,eCiY:=curve.ScalarMult(Commitments[i].X,Commitments[i].Y,ep.Bytes())
		LHSX,LHSY=curve.Add(LHSX,LHSY,eCiX,eCiY)
	}
	sHx, sHy := curve.ScalarMult(H.X, H.Y, s.Bytes()) // rH
	zGx, zGy := curve.ScalarMult(G.X, G.Y, z.Bytes()) // aG
	RHSX,RHSY:=curve.Add(sHx,sHy,zGx,zGy)
	fmt.Println("Verifier computes LHS:", LHSX, LHSY)
	fmt.Println("Verifier computes RHS:", RHSX, RHSY)

	if LHSX.Cmp(RHSX) == 0 && LHSY.Cmp(RHSY) == 0 {
		fmt.Println("Proof verified successfully")
	} else {
		fmt.Println("Proof verification failed")
	}
}
// func main() {
// 	// Initialize the elliptic curve
// 	curve := elliptic.P256()
// 	// Generate a key pair
// 	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
// 	if err != nil {
// 		fmt.Println("Error generating private key:", err)
// 		return
// 	}
// 	G := privateKey.PublicKey // Public key as the generator point

// +build extractor

// 3 goals
// Convincing the verifier ->Completeness
// Actially proving the truth of statement ->Soundness
// Prover reveals nothing else than the validity of the statement ->Zero Knowledge

package main
import (
	"crypto/ecdsa"
	"math/big"
	"crypto/rand"
	"crypto/sha256"
	"crypto/elliptic"
	"fmt"
	"errors"
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
type Transcript struct{
	Commitment *ecdsa.PublicKey
	Challenge *big.Int
	Z *big.Int
	S *big.Int	
}
func GenerateTranscript(curve elliptic.Curve, e,x0, r0 *big.Int, m int ,H,G ecdsa.PublicKey )(Transcript) {//return z and s and commitment 
	C0:=GenerateCommitment(curve,r0,x0,H,G) 
	z:=new(big.Int);
	s:=big.NewInt(0);
	Commitments:=make([]*ecdsa.PublicKey,m)
	rS:=make([]*big.Int,m)
	xS:=make([]*big.Int,m)
	ep := big.NewInt(1)
	LHSX:=big.NewInt(0);
	LHSY:=big.NewInt(0);
	for i:=0;i<m;i++ {
		var err error
		rS[i],err = rand.Int(rand.Reader, curve.Params().N)
		if err != nil {
			fmt.Println("Error generating randomness r:", err)
		}
		xS[i] = new(big.Int).Set(x0) // Initialize xS[i] with the value of x0
		xS[i].Add(xS[i], rS[i])  
		Commitments[i]=GenerateCommitment(curve,rS[i],xS[i],H,G)
		ep.Exp(e, big.NewInt(int64(i+1)), curve.Params().N) // e^i
		z.Add(z,new(big.Int).Mul(ep, xS[i]))
		s.Add(s,new(big.Int).Mul(ep,rS[i]))
		eCiX,eCiY:=curve.ScalarMult(Commitments[i].X,Commitments[i].Y,ep.Bytes())
		LHSX,LHSY=curve.Add(LHSX,LHSY,eCiX,eCiY)
		//LHS : Σ e^i * C_i
	}
	sHx, sHy := curve.ScalarMult(H.X, H.Y, s.Bytes()) // rH
	zGx, zGy := curve.ScalarMult(G.X, G.Y, z.Bytes()) // aG
	RHSX,RHSY:=curve.Add(sHx,sHy,zGx,zGy)
	//RHS: sH + zG
	//Verification LHS =RHS ? 
	if LHSX.Cmp(RHSX) == 0 && LHSY.Cmp(RHSY) == 0 {
		fmt.Println("Proof verified successfully")
	} else {
		fmt.Println("Proof verification failed")
	}
	return Transcript{
		Commitment: C0,
		Challenge: e,
		Z: z,
		S: s,
	}
}
func constructVandermondeMatrix(curve elliptic.Curve, challenges []*big.Int, m int) [][]*big.Int {
	matrix := make([][]*big.Int, m+1)
	for i := 0; i < m+1; i++ {
		matrix[i] = make([]*big.Int, m+1)
		for j := 0; j < m+1; j++ {
			matrix[i][j] = new(big.Int).Exp(challenges[i], big.NewInt(int64(j)), curve.Params().P)
		}
	}
	return matrix
}

func invertMatrix(curve elliptic.Curve, matrix [][]*big.Int) ([][]*big.Int, error) {
	n := len(matrix)
	Order := curve.Params().P

	// Initialize the inverse matrix
	inverse := make([][]*big.Int, n)
	for i := range inverse {
		inverse[i] = make([]*big.Int, n)
		for j := range inverse[i] {
			inverse[i][j] = new(big.Int) // Ensure every element is allocated
		}
	}

	// Initialize augmented matrix [matrix | I]
	augmented := make([][]*big.Int, n)
	for i := range augmented {
		augmented[i] = make([]*big.Int, 2*n)
		for j := 0; j < 2*n; j++ {
			augmented[i][j] = new(big.Int) // Allocate every element explicitly
			if j < n {
				augmented[i][j].Set(matrix[i][j])
			} else if j == n+i {
				augmented[i][j].SetInt64(1) // Identity matrix part
			}
		}
	}

	// Gaussian elimination with modular arithmetic
	for i := 0; i < n; i++ {
		// Find pivot
		pivot := augmented[i][i]
		if pivot.Sign() == 0 {
			return nil, errors.New("matrix is not invertible")
		}

		// Compute pivot inverse
		pivotInv := new(big.Int).ModInverse(pivot, Order)
		if pivotInv == nil {
			return nil, errors.New("pivot is not invertible")
		}

		// Normalize row i
		for j := 0; j < 2*n; j++ {
			augmented[i][j].Mul(augmented[i][j], pivotInv)
			augmented[i][j].Mod(augmented[i][j], Order)
		}

		// Eliminate column i
		for j := 0; j < n; j++ {
			if i != j {
				factor := new(big.Int).Set(augmented[j][i]) // Copy factor safely
				for k := 0; k < 2*n; k++ {
					temp := new(big.Int).Mul(factor, augmented[i][k])
					temp.Mod(temp, Order)
					augmented[j][k].Sub(augmented[j][k], temp)
					augmented[j][k].Mod(augmented[j][k], Order)
				}
			}
		}
	}

	// Extract inverse matrix
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			inverse[i][j].Set(augmented[i][n+j])
		}
	}

	return inverse, nil
}


//Extractor starts the prover
//Prover generates C0
//Extractor provides challenge e
//Obtains z,s( z = Σ e^i * x_i,  s = Σ e^i * r_i
//z1,s1,z2,s2,.... zm ,sm 
//Extractor creates vandermorte matrix using challenges
func main(){
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		fmt.Println("Error generating private key:", err)
		return
	}
	G:=privateKey.PublicKey
	H:=GenerateNumsPoint(curve,G)
	m:=3
	Transcripts:=make([]Transcript,m+1)
	challenges:=make([]*big.Int,m+1)
	secrets:=make([]*big.Int,m+1)
	randoms:=make([]*big.Int,m+1)
	for i:=0;i<m+1;i++ {
		secrets[i]=big.NewInt(int64(i*i))
		randoms[i],err=rand.Int(rand.Reader,curve.Params().P)
		challenges[i],err=rand.Int(rand.Reader,curve.Params().P)//given by extractor in reality
		//i.e challenge is known to extractor unlike secret and random
		Transcripts[i]=GenerateTranscript(curve,challenges[i],secrets[i],randoms[i],m,H,G)
	}
	V:=constructVandermondeMatrix(curve,challenges,m)
	I,err:=invertMatrix(curve,V)
	predictedSecrets:=make([]*big.Int,m+1)
	for i := 0; i < m+1; i++ {
		predictedSecrets[i] = big.NewInt(0)
		for j := 0; j < m+1; j++ {
			temp := new(big.Int).Mul(I[i][j], Transcripts[j].Z)
			temp.Mod(temp, curve.Params().P) // Ensure modular reduction at every step
			predictedSecrets[i].Add(predictedSecrets[i], temp)
			predictedSecrets[i].Mod(predictedSecrets[i], curve.Params().P)
		}
	}
	
	for i:=0;i<m+1;i++{
		fmt.Println("Actual Value :")
		fmt.Println(secrets[i])
		fmt.Println("Predicted Value Value :")
		fmt.Println(predictedSecrets[i])
	}
}
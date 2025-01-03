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
type Keys struct{
	p ecdsa.PrivateKey
	G ecdsa.PublicKey
}
func generatePvtAndPubKey(curve elliptic.Curve)(Keys){
	p,err:=ecdsa.GenerateKey(curve,rand.Reader)
	if err != nil {
		fmt.Println("Error generating private key:", err)
		return
	}
	G:=&p.PublicKey
	return Keys{
		p:*p,
		G:*G,
	}
}
func generateRandomScalar() *big.Int{
	randomScalar,err:=rand.Int(rand.Reader,curve.Params().P)
	if err!=nil{
		fmt.Println("Error generating random scalar")
	}
	return randomScalar
}
func GenerateRandomVector(dim int) []*big.Int {
	vector := make([]*big.Int, dim)
	for i := 0; i < dim; i++ {
		vector[i] = GenerateRandomScalar()
	}
	return vector
}
func DotProduct(a, b []*big.Int) *big.Int {
	result := big.NewInt(0)
	for i := 0; i < len(a); i++ {
		product := new(big.Int).Mul(a[i], b[i])
		result.Add(result, product)
	}
	return result
}
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
func GenerateCommitment(curve elliptic.Curve, r *big.Int,a []*big.Int ,H ecdsa.PublicKey,G ecdsa.PublicKey)(*ecdsa.PublicKey){
	commitmentX, commitmentY := big.NewInt(0), big.NewInt(0)
	rHx, rHy := curve.ScalarMult(H.X, H.Y, r.Bytes()) // rH
	for i := 0; i < len(a); i++ {
		// Scalar multiplication for each element in the vector
		aGx, aGy := curve.ScalarMult(G.X, G.Y, a[i].Bytes()) // v_i * G
		commitmentX.Add(commitmentX, aGx)
		commitmentY.Add(commitmentY, aGy)
	}
	Cx, Cy := curve.Add(rHx, rHy, commitmentX, commitmentY)         // C = rH + aG
	fmt.Printf("Commitment C: (%s, %s)\n", Cx.String(), Cy.String())
	return &ecdsa.PublicKey{
		Curve: curve,
		X: Cx,
		Y: Cy,
	}
}
func CommitmentStep(){

}
func main(){
	curve:=elliptic.P256()
	Keys:=generatePvtAndPubKey(curve)
	G:=Keys.G
	x := []*big.Int{big.NewInt(123456789), big.NewInt(987654321)} 
	y := []*big.Int{big.NewInt(1122334455), big.NewInt(9988776655)} 
	dx:=GenerateRandomVector(2)
	dy:=GenerateRandomVector(2)
	H:=GenerateNumsPoint(curve,G)
	rd:=generateRandomScalar()
	sd:=generateRandomScalar()
	t1:=generateRandomScalar()
	t0:=generateRandomScalar()
	Ad:=GenerateCommitment(curve,rd,dx,H,G)
	Bd:=GenerateCommitment(curve,sd,dy,H,G)
	temp0:=DotProduct(dx,dy)
	temp1:=DotProduct(x,dy).Add(DotProduct(y,dx))
	C1:=GenerateCommitment(curve,t1,temp1,H,G)
	C0:=GenerateCommitment(curve,t0,temp0,H,G)
}
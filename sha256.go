package main
import (
	"crypto/sha256"
	"encoding/hex"
	//"fmt"	//"crypto/ecdh"
)
// • Hiding - a commitment C does not reveal the value it commits to.
// • Binding - having made the commitment C(m) to m, you can’t change your mind and open it as a commitment to a diﬀerent message m′

//Hash functions are not homophorbic, meaning that if you change the input by even a single bit, the output will be completely different.
//example SHA256(a)+SHA256(b) != SHA256(a+b)
func hash (a string )string {
	sha256Hash :=sha256.Sum256([]byte(a))//returns the hash of the input 
	//SHA256 required input in byte form so we convert the string to byte using []byte<string>
	//sha256.Sum256() returns a fixed-size array of 32 bytes 
	hexVersion:=hex.EncodeToString(sha256Hash[:])//takes byte slice as input and converts to hex characters ([:] takes first 32 bytes here its already 32 bytes)
	return (hexVersion)
}

// func main() {
// 	a:="1234";
// 	b:="5678";
// 	d:=a+b;//"12345678"
// 	shaAstr:=hash(a);
// 	shaBstr:=hash(b);
// 	shaDstr:=hash(d);
// 	combinedHash:=shaAstr+shaBstr;
// 	fmt.Println("SHA256 of a: ",shaAstr);
// 	fmt.Println("SHA256 of b: ",shaBstr);
// 	fmt.Println("SHA256 of a+b: ",shaDstr);
// 	fmt.Println("SHA256 of a+SHA256 of b: ",combinedHash);
// 	//Clearly the hash of a+b is not equal to the hash of a+hash of b as expected
// }

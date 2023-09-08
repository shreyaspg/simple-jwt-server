package main

import (
	// "crypto/rand"
	"crypto/x509"
	// "crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

const privKeyLocation = "./private_key.pem"
const iss = "shrys-solutions@fortanix.iam.gserviceaccount.com"

func main() {

    // Read private key from file
    priv, err := ioutil.ReadFile(privKeyLocation)
    if err != nil {
        log.Println("Read error", err)
    }

    privPem, _ := pem.Decode(priv)
    privKey, err := x509.ParsePKCS8PrivateKey(privPem.Bytes)
    if err != nil {
        log.Println("Parse error", err)
    }

	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Minute).Unix()
	claims["authorized"] = true
	claims["iss"] = iss
	claims["scope"] = "https://www.googleapis.com/auth/cloud-platform"
	claims["aud"] = "https://www.googleapis.com/oauth2/v4/token"
	claims["iat"] = time.Now().Unix()

    log.Println(claims)
	tokenString, err := token.SignedString(privKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("tokenString:\n %s\n\n ", tokenString)
	token, err = jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodRSA)

		if !ok {
			fmt.Println("Unauthorized")
			return "Error", nil
		}
		fmt.Println("Authorized")
		return "Verified", nil
	})
}

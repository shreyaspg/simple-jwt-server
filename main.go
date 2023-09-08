package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"github.com/golang-jwt/jwt"
	"os"
	"time"
    "encoding/pem"
)

func main() {
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Minute)
	claims["authorized"] = true
	claims["iss"] = "test@service-account.gcloud.com"
	claims["scope"] = "https://www.googleapis.com/auth/cloud-platform"
	claims["aud"] = "https://www.googleapis.com/oauth2/v4/token"
	claims["iat"] = time.Now()

	privatekey, err := rsa.GenerateKey(rand.Reader, 4096)
    privateKeyPEM := &pem.Block{
        Type:  "RSA PRIVATE KEY",
        Bytes: x509.MarshalPKCS1PrivateKey(privatekey),
    }
    privateKeyFile, err := os.Create("private_key.pem")
    if err != nil {
        fmt.Println("Error creating private key file:", err)
        os.Exit(1)
    }
    pem.Encode(privateKeyFile, privateKeyPEM)
    privateKeyFile.Close()

    // Extract the public key from the private key
    publicKey := &privatekey.PublicKey

    // Encode the public key to the PEM format
    publicKeyPEM := &pem.Block{
        Type:  "RSA PUBLIC KEY",
        Bytes: x509.MarshalPKCS1PublicKey(publicKey),
    }
    publicKeyFile, err := os.Create("public_key.pem")
    if err != nil {
        fmt.Println("Error creating public key file:", err)
        os.Exit(1)
    }
    pem.Encode(publicKeyFile, publicKeyPEM)
    publicKeyFile.Close()
	// fmt.Printf("%x\n\n", privatekey.D.Bytes())
    // pub := &privatekey.PublicKey
    // fmt.Printf("%x", pub)
	tokenString, err := token.SignedString(privatekey)
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

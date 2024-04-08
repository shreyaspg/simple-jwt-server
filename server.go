package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"crypto/x509"
	"encoding/pem"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const privKeyLocation = "./private_key.pem"

var (
	PORT    string
	iss     string
	privKey any
)

type Jwt struct {
	Data string `json:"data"`
}

func issueGCPJwt(c echo.Context) error {
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Minute).Unix()
	claims["authorized"] = true
	claims["iss"] = iss
	claims["scope"] = "https://www.googleapis.com/auth/cloud-platform"
	claims["aud"] = "https://www.googleapis.com/oauth2/v4/token"
	claims["iat"] = time.Now().Unix()

	tokenString, err := token.SignedString(privKey)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, Jwt{Data: tokenString})
}

func main() {
	e := echo.New()

	// hide the startup banner
	e.HideBanner = true

	// set up loggin middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.GET("/upcheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "Up and running JWT issuer")
	})

	e.GET("/jwt/issue", issueGCPJwt)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", PORT)))
}

func init() {
	godotenv.Load()
	var isSet bool
	PORT, isSet = os.LookupEnv("PORT")
	if !isSet {
		log.Println("Using default PORT 1323")
		PORT = "1323"
	}

	iss, isSet = os.LookupEnv("ISS")
	if !isSet {
		log.Println("Missing ISS value")
		os.Exit(1)
	}

	// Read private key from file
	priv, err := os.ReadFile(privKeyLocation)
	if err != nil {
		log.Println("Read error", err)
		os.Exit(1)
	}

	privPem, _ := pem.Decode(priv)
	privKey, err = x509.ParsePKCS8PrivateKey(privPem.Bytes)
	if err != nil {
		log.Println("Parse error", err)
		os.Exit(1)
	}
}

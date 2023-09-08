package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"crypto/x509"
	"encoding/pem"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io/ioutil"
)

const privKeyLocation = "./private_key.pem"

var (
	PORT    string
	iss     string
	privKey any
)

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

	e.GET("/jwt/issue", handleJwtIssue)
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
	priv, err := ioutil.ReadFile(privKeyLocation)
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

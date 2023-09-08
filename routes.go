package main

import (
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type Jwt struct {
	Data string `json:"data"`
}

func handleJwtIssue(c echo.Context) error {
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
		return err
	}
	return c.JSON(http.StatusOK, Jwt{Data: tokenString})
}

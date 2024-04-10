package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(ttl time.Duration, payload interface{}, secretJWTKey string) (string,error){
    token := jwt.New(jwt.SigningMethodHS256)

    now := time.Now()

    claims := token.Claims.(jwt.MapClaims)
    claims["sub"] = payload
    claims["exp"] = now.Add(ttl).Unix()
    claims["iat"] = now.Unix()
    claims["nbf"] = now.Unix()

    tokenString, err := token.SignedString([]byte(secretJWTKey))
    if err != nil{
        return "", fmt.Errorf("failed to sign and generate jwt token: %w", err)
    }

    return tokenString, nil
}

func ValidateToken(tokenString, secretJWTKey string) (interface{}, error){
    tok, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
        }

        return []byte(secretJWTKey),nil
    })

    if err != nil{
        return nil, fmt.Errorf("failed to parse jwt token: %w", err)
    }

    claims,ok := tok.Claims.(jwt.MapClaims)
    if !ok || !tok.Valid{
        return nil, fmt.Errorf("invalid JWT token")
    }

    return claims["sub"], nil
}

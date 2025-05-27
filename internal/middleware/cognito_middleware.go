package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"math/big"
	"net/http"
	"os"
	"strings"
)

var (
	jwksCache map[string]interface{}
)

func Protected() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := parseAndValidateToken(tokenString)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if !token.Valid {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		ctx.Locals("userId", claims["sub"])
		ctx.Locals("username", claims["username"])
		ctx.Locals("cognitoClaims", claims)

		return ctx.Next()
	}
}

func parseAndValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid header not found in token")
		}
		jwk, err := getJWKS(kid)
		if err != nil {
			return nil, err
		}

		pem, err := convertJWKToPEM(jwk)
		if err != nil {
			return nil, err
		}

		return pem, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}

	expectedIssuer := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", os.Getenv("AWS_REGION"), os.Getenv("AWS_USER_POOL_ID"))
	if claims["iss"] != expectedIssuer {
		return nil, errors.New("invalid token issuer")
	}

	if claims["aud"] != os.Getenv("AWS_APP_CLIENT_ID") && !containsString(claims["client_id"], os.Getenv("AWS_APP_CLIENT_ID")) {
		return nil, errors.New("invalid token audience")
	}

	return token, nil
}

func getJWKS(kid string) (map[string]interface{}, error) {
	if jwksCache == nil {
		jwksCache = make(map[string]interface{})
	}

	if jwk, exists := jwksCache[kid]; exists {
		return jwk.(map[string]interface{}), nil
	}

	jwksURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", os.Getenv("AWS_REGION"), os.Getenv("AWS_USER_POOL_ID"))
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jwks struct {
		Keys []map[string]interface{} `json:"keys"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	for _, key := range jwks.Keys {
		if key["kid"] == kid {
			jwksCache[kid] = key
			return key, nil
		}
	}

	return nil, errors.New("no matching JWK found for kid")
}

func containsString(claim interface{}, str string) bool {
	switch v := claim.(type) {
	case string:
		return v == str
	case []string:
		for _, s := range v {
			if s == str {
				return true
			}
		}
	case []interface{}:
		for _, s := range v {
			if s, ok := s.(string); ok && s == str {
				return true
			}
		}
	}
	return false
}

func convertJWKToPEM(jwk map[string]interface{}) (interface{}, error) {
	nStr, ok := jwk["n"].(string)
	if !ok {
		return nil, errors.New("invalid JWK: missing or invalid modulus (n)")
	}

	eStr, ok := jwk["e"].(string)
	if !ok {
		return nil, errors.New("invalid JWK: missing or invalid exponent (e)")
	}

	n, err := base64.RawURLEncoding.DecodeString(nStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus: %v", err)
	}

	e, err := base64.RawURLEncoding.DecodeString(eStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent: %v", err)
	}

	modulus := new(big.Int)
	modulus.SetBytes(n)

	var exponent int
	for i := 0; i < len(e); i++ {
		exponent = exponent<<8 + int(e[i])
	}

	pubKey := &rsa.PublicKey{
		N: modulus,
		E: exponent,
	}

	return pubKey, nil
}

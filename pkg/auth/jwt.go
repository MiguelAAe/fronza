package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	_ "github.com/golang-jwt/jwt/v4"
)

var mySigningKey = []byte("portal-deliveries")

type CustomClaims struct {
	ID            int64
	Email         string
	Role          string
	Platform      string
	RegisterClaim *jwt.RegisteredClaims
}

func NewJWT(ID int64, email, role, platform string) (string, error) {

	// Create the claims
	claims := CustomClaims{
		ID:       ID,
		Email:    email,
		Role:     role,
		Platform: platform,
		RegisterClaim: &jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "test",
			Subject:   "somebody",
			ID:        "1",
			Audience:  []string{"somebody_else"},
		},
	}

	t := jwt.New(jwt.SigningMethodHS256)
	t.Claims = &claims

	token, err := t.SignedString(mySigningKey)

	return token, err
}

func ParseToken(tokenString string) (CustomClaims, bool) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	if err != nil {
		return CustomClaims{}, false
	}

	claims, ok := token.Claims.(*CustomClaims)
	if ok && token.Valid {
		return *claims, true
	}

	return CustomClaims{}, false
}

func (c CustomClaims) Valid() error {
	vErr := new(jwt.ValidationError)
	now := jwt.TimeFunc()

	// The claims below are optional, by default, so if they are set to the
	// default value in Go, let's not fail the verification for them.
	if !c.RegisterClaim.VerifyExpiresAt(now, false) {
		delta := now.Sub(c.RegisterClaim.ExpiresAt.Time)
		vErr.Inner = fmt.Errorf("%s by %v", delta, jwt.ErrTokenExpired)
		vErr.Errors |= jwt.ValidationErrorExpired
	}

	if !c.RegisterClaim.VerifyIssuedAt(now, false) {
		vErr.Inner = jwt.ErrTokenUsedBeforeIssued
		vErr.Errors |= jwt.ValidationErrorIssuedAt
	}

	if !c.RegisterClaim.VerifyNotBefore(now, false) {
		vErr.Inner = jwt.ErrTokenNotValidYet
		vErr.Errors |= jwt.ValidationErrorNotValidYet
	}

	if vErr.Errors == 0 {
		return nil
	}

	return vErr
}

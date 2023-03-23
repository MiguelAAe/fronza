package auth

import (
	"fmt"
	"testing"

	"github.com/golang-jwt/jwt/v4"
)

func TestParseToken(t *testing.T) {
	tokenString, err := NewJWT(1, "user", "", "ios")
	if err != nil {
		t.Error(err)
		return
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	claims, ok := token.Claims.(*CustomClaims)
	if ok && token.Valid {
		t.Logf("%d %s", claims.ID, claims.RegisterClaim.Issuer)
	} else {
		t.Logf("inside else")
		t.Error(err)
	}

	fmt.Println(claims)
}

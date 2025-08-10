package auth

import "github.com/golang-jwt/jwt/v5"

var EmptyClaims = &Claims{}

type Claims struct {
	jwt.RegisteredClaims
}

func (c *Claims) UserID() string {
	subject, err := c.GetSubject()
	if err != nil {
		return ""
	}
	return subject
}

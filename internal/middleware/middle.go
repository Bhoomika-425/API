package middleware

import (
	"fmt"
	"project/internal/auth"
)

type Mid struct {
	auth auth.UserAuth
}

func NewMiddleware(a auth.UserAuth) (Mid, error) {
	if a == nil {
		return Mid{}, fmt.Errorf("auth cant be null")
	}
	return Mid{
		auth: a,
	}, nil
}

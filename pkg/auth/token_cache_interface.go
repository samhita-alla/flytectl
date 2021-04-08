package auth

import "golang.org/x/oauth2"

type TokenCache interface {
	SaveToken(token oauth2.Token) error
	GetToken() (oauth2.Token, error)
}

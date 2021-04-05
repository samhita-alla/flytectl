package auth

import (
	"context"

	goauth "golang.org/x/oauth2"
)

func GenerateClientConfig(ctx context.Context) (goauth.Config, error) {
	var clientConfigFromAdmin ClientConfigFromAdmin
	var serverConfigFromAdmin ServerConfigFromAdmin
	var err error
	if clientConfigFromAdmin, err = GetClientConfigFromAdmin(ctx); err != nil {
		return goauth.Config{}, err
	}
	if serverConfigFromAdmin, err = GetAuthServerConfigFromAdmin(ctx); err != nil {
		return goauth.Config{}, err
	}
	genClientConf := goauth.Config{
		ClientID: clientConfigFromAdmin.ClientId,
		RedirectURL: clientConfigFromAdmin.RedirectUri,
		Scopes: serverConfigFromAdmin.ScopesSupported,
		Endpoint: goauth.Endpoint{
			TokenURL: serverConfigFromAdmin.TokenEndpoint,
			AuthURL: serverConfigFromAdmin.AuthorizationEndpoint,
		},
	}
	return genClientConf, nil
}
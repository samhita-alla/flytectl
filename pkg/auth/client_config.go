package auth

import (
	"context"
	"golang.org/x/oauth2"
)
// Initialize this config once and use it in different flows.
var clientConf *oauth2.Config

func GenerateClientConfig(ctx context.Context) (*oauth2.Config, error) {
	if clientConf != nil {
		return clientConf, nil
	}
	var clientConfigFromAdmin ClientConfigFromAdmin
	var serverConfigFromAdmin ServerConfigFromAdmin
	var err error
	if clientConfigFromAdmin, err = GetClientConfigFromAdmin(ctx); err != nil {
		return nil, err
	}
	if serverConfigFromAdmin, err = GetAuthServerConfigFromAdmin(ctx); err != nil {
		return nil, err
	}
	clientConf = &oauth2.Config{
		ClientID: clientConfigFromAdmin.ClientId,
		RedirectURL: clientConfigFromAdmin.RedirectUri,
		Scopes: serverConfigFromAdmin.ScopesSupported,
		Endpoint: oauth2.Endpoint{
			TokenURL: serverConfigFromAdmin.TokenEndpoint,
			AuthURL: serverConfigFromAdmin.AuthorizationEndpoint,
		},
	}
	return clientConf, nil
}
package auth

import (
	"golang.org/x/oauth2"
)

// Initialize this config once and use it in different flows.
var clientConf *oauth2.Config

func GenerateClientConfig() (*oauth2.Config, error) {
	if clientConf != nil {
		return clientConf, nil
	}
	var clientConfigFromAdmin ClientConfigFromAdmin
	var serverConfigFromAdmin ServerConfigFromAdmin
	var err error
	if clientConfigFromAdmin, err = GetClientConfigFromAdmin(); err != nil {
		return nil, err
	}
	if serverConfigFromAdmin, err = GetAuthServerConfigFromAdmin(); err != nil {
		return nil, err
	}
	clientConf = &oauth2.Config{
		ClientID:    clientConfigFromAdmin.ClientID,
		RedirectURL: clientConfigFromAdmin.RedirectURI,
		Scopes:      serverConfigFromAdmin.ScopesSupported,
		Endpoint: oauth2.Endpoint{
			TokenURL: serverConfigFromAdmin.TokenEndpoint,
			AuthURL:  serverConfigFromAdmin.AuthorizationEndpoint,
		},
	}
	return clientConf, nil
}

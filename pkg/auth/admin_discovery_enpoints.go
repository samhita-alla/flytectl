package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var (
	clientConfigURL     = "/config/v1/flyte_client"
	authServerConfigURL = "/.well-known/oauth-authorization-server"
)

type ClientConfigFromAdmin struct {
	ClientID                 string   `json:"client_id"`
	RedirectURI              string   `json:"redirect_uri"`
	Scopes                   []string `json:"scopes"`
	AuthorizationMetadataKey string   `json:"authorization_metadata_key"`
}

type ServerConfigFromAdmin struct {
	Issuer                            string   `json:"issuer"`
	AuthorizationEndpoint             string   `json:"authorization_endpoint"`
	TokenEndpoint                     string   `json:"token_endpoint"`
	ResponseTypesSupported            []string `json:"response_types_supported"`
	ScopesSupported                   []string `json:"scopes_supported"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
	JSONWebKeysURI                    string   `json:"jwks_uri"`
	CodeChallengeMethodsSupported     []string `json:"code_challenge_methods_supported"`
	// optional
	GrantTypesSupported []string `json:"grant_types_supported"`
}

func GetClientConfigFromAdmin() (ClientConfigFromAdmin, error) {
	// adminEndpoint := admin.GetConfig(ctx).Endpoint
	// @TODO : fix this before checkin. We need GRPC endpoint for it
	resp, err := http.Get("http://localhost:8088" + clientConfigURL)
	clientConfigFromAdmin := ClientConfigFromAdmin{}
	if err != nil {
		return clientConfigFromAdmin, err
	}
	defer resp.Body.Close()
	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return clientConfigFromAdmin, err
	}
	if jsonErr := json.Unmarshal(body, &clientConfigFromAdmin); jsonErr != nil {
		return clientConfigFromAdmin, jsonErr
	}
	return clientConfigFromAdmin, nil
}

func GetAuthServerConfigFromAdmin() (ServerConfigFromAdmin, error) {
	// adminEndpoint := admin.GetConfig(ctx).Endpoint
	// @TODO : fix this before checkin. We need GRPC endpoint for it
	resp, err := http.Get("http://localhost:8088" + authServerConfigURL)
	serverConfigFromAdmin := ServerConfigFromAdmin{}
	if err != nil {
		return serverConfigFromAdmin, err
	}
	defer resp.Body.Close()
	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return serverConfigFromAdmin, err
	}
	if jsonErr := json.Unmarshal(body, &serverConfigFromAdmin); jsonErr != nil {
		return serverConfigFromAdmin, jsonErr
	}
	return serverConfigFromAdmin, nil
}

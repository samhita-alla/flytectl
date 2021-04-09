package auth

import (
	"encoding/json"
	"io"
	"net/http"
)

var (
	clientConfigUrl     = "/config/v1/flyte_client"
	authServerConfigUrl = "/.well-known/oauth-authorization-server"
)

type ClientConfigFromAdmin struct {
	ClientId                 string   `json:"client_id"`
	RedirectUri              string   `json:"redirect_uri"`
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
	JSONWebKeysUri                    string   `json:"jwks_uri"`
	CodeChallengeMethodsSupported     []string `json:"code_challenge_methods_supported"`
	// optional
	GrantTypesSupported []string `json:"grant_types_supported"`
}

func GetClientConfigFromAdmin() (ClientConfigFromAdmin, error) {
	// adminEndpoint := admin.GetConfig(ctx).Endpoint
	// @TODO : fix this before checkin. We need GRPC endpoint for it
	resp, err := http.Get("http://localhost:8088" + clientConfigUrl)
	clientConfigFromAdmin := ClientConfigFromAdmin{}
	if err != nil {
		return clientConfigFromAdmin, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if jsonErr := json.Unmarshal(body, &clientConfigFromAdmin); jsonErr != nil {
		return clientConfigFromAdmin, jsonErr
	}
	return clientConfigFromAdmin, nil
}

func GetAuthServerConfigFromAdmin() (ServerConfigFromAdmin, error) {
	// adminEndpoint := admin.GetConfig(ctx).Endpoint
	// @TODO : fix this before checkin. We need GRPC endpoint for it
	resp, err := http.Get("http://localhost:8088" + authServerConfigUrl)
	serverConfigFromAdmin := ServerConfigFromAdmin{}
	if err != nil {
		return serverConfigFromAdmin, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if jsonErr := json.Unmarshal(body, &serverConfigFromAdmin); jsonErr != nil {
		return serverConfigFromAdmin, jsonErr
	}
	return serverConfigFromAdmin, nil
}

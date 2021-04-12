package auth

import (
	"context"

	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/service"

	"golang.org/x/oauth2"
)

// Initialize this config once and use it in different flows.
var clientConf *oauth2.Config

func GenerateClientConfig(ctx context.Context, authClient service.AuthServiceClient) (*oauth2.Config, error) {
	var flyteClientResp *service.FlyteClientResponse
	var err error
	if flyteClientResp, err = authClient.FlyteClient(ctx, &service.FlyteClientRequest{}); err != nil {
		return nil, err
	}
	var flyteOauthMetaResp *service.OAuth2MetadataResponse
	if flyteOauthMetaResp, err = authClient.OAuth2Metadata(ctx, &service.OAuth2MetadataRequest{}); err != nil {
		return nil, err
	}
	clientConf = &oauth2.Config{
		ClientID:    flyteClientResp.ClientId,
		RedirectURL: flyteClientResp.RedirectUri,
		Scopes:      flyteClientResp.Scopes,
		Endpoint: oauth2.Endpoint{
			TokenURL: flyteOauthMetaResp.TokenEndpoint,
			AuthURL:  flyteOauthMetaResp.AuthorizationEndpoint,
		},
	}
	return clientConf, nil
}

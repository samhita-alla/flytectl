package auth

import (
	"context"
	"fmt"
	"github.com/flyteorg/flyteidl/clients/go/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/service"
	"github.com/flyteorg/flytestdlib/logger"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"sync"
)

// Initialize this config once and use it in different flows.
var clientConf *oauth2.Config


func InitializeAuthClientFromConfig(ctx context.Context) (service.AuthServiceClient, error) {
	cfg := admin.GetConfig(ctx)

	if cfg == nil {
		return nil, fmt.Errorf("retrieved Nil config for admin")
	}
	return InitializeAuthClient(ctx, *cfg), nil
}


var (
	once            = sync.Once{}
	adminConnection *grpc.ClientConn
)

// Create an AuthServiceClient with a shared Admin connection for the process
func InitializeAuthClient(ctx context.Context, cfg admin.Config) service.AuthServiceClient {
	once.Do(func() {
		var err error
		adminConnection, err = admin.NewAdminConnection(ctx, cfg)
		if err != nil {
			logger.Panicf(ctx, "failed to initialize Admin connection. Err: %s", err.Error())
		}
	})
	return NewAuthClient(ctx, adminConnection)
}


func NewAuthClient(ctx context.Context, conn *grpc.ClientConn) service.AuthServiceClient {
	logger.Infof(ctx, "Initialized Auth client")
	return service.NewAuthServiceClient(conn)
}


func GenerateClientConfig(ctx context.Context) (*oauth2.Config, error) {
	authServiceClient := InitializeAuthClient(ctx, *admin.GetConfig(ctx))
	authServiceClient :=
	var flyteClientResp *service.FlyteClientResponse
	var err error
	if flyteClientResp, err = authServiceClient.FlyteClient(ctx, &service.FlyteClientRequest{}); err != nil {
		return nil, err
	}
	var flyteOauthMetaResp *service.OAuth2MetadataResponse
	if flyteOauthMetaResp, err =authServiceClient.OAuth2Metadata(ctx, &service.OAuth2MetadataRequest{}); err != nil {
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

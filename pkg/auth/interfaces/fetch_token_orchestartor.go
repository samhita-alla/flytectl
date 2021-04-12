package interfaces

import (
	"context"

	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/service"

	"golang.org/x/oauth2"
)

//go:generate mockery -all -case=underscore

type FetchTokenOrchestrator interface {
	RefreshTheToken(ctx context.Context, clientConf *oauth2.Config, token *oauth2.Token) *oauth2.Token
	FetchTokenFromCacheOrRefreshIt(ctx context.Context, authClient service.AuthServiceClient) *oauth2.Token
	FetchTokenFromAuthFlow(ctx context.Context, authClient service.AuthServiceClient) (*oauth2.Token, error)
}

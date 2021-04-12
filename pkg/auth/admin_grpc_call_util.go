package auth

import (
	"context"

	"github.com/flyteorg/flyteidl/clients/go/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/service"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/status"
)

type AdminGrpcAPICallContext func(ctx context.Context, callOptions []grpc.CallOption) error

type AdminGrpcCallOptions func(ctx context.Context, callOptions []grpc.CallOption) ([]grpc.CallOption, error)

func callOptionForToken(ctx context.Context, token *oauth2.Token) grpc.CallOption {
	var callOption grpc.CallOption
	accessToken := FlyteCtlInsecureTokenSource{
		flyteCtlToken: token,
	}
	if admin.GetConfig(ctx).UseInsecureConnection {
		callOption = grpc.PerRPCCredsCallOption{Creds: InsecurePerRPCCredentials{TokenSource: &accessToken}}
	} else {
		callOption = grpc.PerRPCCredsCallOption{Creds: oauth.TokenSource{TokenSource: &accessToken}}
	}
	return callOption
}

func updateWithNewToken(ctx context.Context, authClient service.AuthServiceClient, callOptions []grpc.CallOption) ([]grpc.CallOption, error) {
	var newToken *oauth2.Token
	var err error
	if newToken, err = FetchTokenFromAuthFlow(ctx, authClient); err != nil {
		return nil, err
	}
	return append(callOptions, callOptionForToken(ctx, newToken)), nil
}

func updateWithCachedOrRefreshedToken(ctx context.Context, authClient service.AuthServiceClient, callOptions []grpc.CallOption) []grpc.CallOption {
	var cachedOrRefreshedToken *oauth2.Token
	if cachedOrRefreshedToken = FetchTokenFromCacheOrRefreshIt(ctx, authClient); cachedOrRefreshedToken == nil {
		return callOptions
	}
	return append(callOptions, callOptionForToken(ctx, cachedOrRefreshedToken))
}

func Do(ctx context.Context, authClient service.AuthServiceClient, grpcAPICallContext AdminGrpcAPICallContext, callOptions []grpc.CallOption, useAuth bool) error {
	// Fetch from the cache only when usAuth is enabled.
	if useAuth {
		callOptions = updateWithCachedOrRefreshedToken(ctx, authClient, callOptions)
	}
	if grpcStatusError := grpcAPICallContext(ctx, callOptions); grpcStatusError != nil {
		if grpcStatus, ok := status.FromError(grpcStatusError); ok && grpcStatus.Code() == codes.Unauthenticated && useAuth {
			var err error
			if callOptions, err = updateWithNewToken(ctx, authClient, callOptions); err != nil {
				return err
			}
			return grpcAPICallContext(ctx, callOptions)
		}
		return grpcStatusError
	}
	return nil
}

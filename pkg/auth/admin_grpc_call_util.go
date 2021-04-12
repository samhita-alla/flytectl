package auth

import (
	"context"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"

	"github.com/flyteorg/flyteidl/clients/go/admin"

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

func updateWithNewToken(ctx context.Context, cmdCtx cmdCore.CommandContext,callOptions []grpc.CallOption) ([]grpc.CallOption, error) {
	var newToken *oauth2.Token
	var err error
	if newToken, err = FetchTokenFromAuthFlow(ctx, cmdCtx); err != nil {
		return nil, err
	}
	return append(callOptions, callOptionForToken(ctx, newToken)), nil
}

func updateWithCachedOrRefreshedToken(ctx context.Context, cmdCtx cmdCore.CommandContext, callOptions []grpc.CallOption) []grpc.CallOption {
	var cachedOrRefreshedToken *oauth2.Token
	if cachedOrRefreshedToken = FetchTokenFromCacheOrRefreshIt(ctx, cmdCtx); cachedOrRefreshedToken == nil {
		return callOptions
	}
	return append(callOptions, callOptionForToken(ctx, cachedOrRefreshedToken))
}

func Do(ctx context.Context, cmdCtx cmdCore.CommandContext, grpcAPICallContext AdminGrpcAPICallContext, callOptions []grpc.CallOption, useAuth bool) error {
	// Fetch from the cache only when usAuth is enabled.
	if useAuth {
		callOptions = updateWithCachedOrRefreshedToken(ctx, cmdCtx, callOptions)
	}
	if grpcStatusError := grpcAPICallContext(ctx, callOptions); grpcStatusError != nil {
		if grpcStatus, ok := status.FromError(grpcStatusError); ok && grpcStatus.Code() == codes.Unauthenticated && useAuth {
			var err error
			if callOptions, err = updateWithNewToken(ctx, cmdCtx, callOptions); err != nil {
				return err
			}
			return grpcAPICallContext(ctx, callOptions)
		}
		return grpcStatusError
	}
	return nil
}

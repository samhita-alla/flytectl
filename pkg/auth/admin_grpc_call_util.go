package auth

import (
	"context"
	"golang.org/x/oauth2"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AdminGrpcApiCallContext func(ctx context.Context, callOptions []grpc.CallOption) error

type AdminGrpcCallOptions func(ctx context.Context, callOptions []grpc.CallOption) ([]grpc.CallOption, error)


func callOptionForToken(token *oauth2.Token) grpc.CallOption {
	var callOption grpc.CallOption
	accessToken := FlyteCtlTokenSource{
		flyteCtlToken: token,
	}
	callOption = grpc.PerRPCCredsCallOption{Creds: InsecurePerRPCCredentials{TokenSource: &accessToken}}
	return callOption
}

func updateWithNewToken(ctx context.Context, callOptions []grpc.CallOption) ([]grpc.CallOption, error) {
	var newToken *oauth2.Token
	var err error
	if newToken, err = FetchTokenFromAuthFlow(ctx); err != nil {
		return  nil, err
	}
	return append(callOptions, callOptionForToken(newToken)), nil
}

func updateWithCachedOrRefreshedToken(ctx context.Context, callOptions []grpc.CallOption) []grpc.CallOption {
	var cachedOrRefreshedToken *oauth2.Token
	if cachedOrRefreshedToken = FetchTokenFromCacheOrRefreshIt(ctx); cachedOrRefreshedToken == nil {
		return  callOptions
	}
	return append(callOptions, callOptionForToken(cachedOrRefreshedToken))
}

func Do(grpcApiCallContext AdminGrpcApiCallContext, ctx context.Context, callOptions []grpc.CallOption, useAuth bool) error {
	callOptions = updateWithCachedOrRefreshedToken(ctx, callOptions)
	if grpcStatusError := grpcApiCallContext(ctx, callOptions); grpcStatusError != nil {
		if grpcStatus, ok := status.FromError(grpcStatusError) ; ok && grpcStatus.Code() == codes.Unauthenticated && useAuth {
			var err error
			if callOptions, err = updateWithNewToken(ctx, callOptions); err != nil {
				return err
			}
			return grpcApiCallContext(ctx, callOptions)
		}
		return grpcStatusError
	}
	return nil
}

package auth

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AdminGrpcApiCallContext func(ctx context.Context, callOptions []grpc.CallOption) error

type AdminGrpcCallOptions func(ctx context.Context, callOptions []grpc.CallOption) ([]grpc.CallOption, error)

var grpcCallOptions = func(ctx context.Context, callOptions []grpc.CallOption) ([]grpc.CallOption, error) {
	var authFlowCallOption grpc.CallOption
	var err error
	if authFlowCallOption, err = StartAuthFlow(ctx); err != nil {
		return  nil, err
	}
	return append(callOptions, authFlowCallOption), nil
}

func Do(grpcApiCallContext AdminGrpcApiCallContext, ctx context.Context, callOptions []grpc.CallOption, useAuth bool) error {
	if grpcStatusError := grpcApiCallContext(ctx, callOptions); grpcStatusError != nil {
		if grpcStatus, ok := status.FromError(grpcStatusError) ; ok && grpcStatus.Code() == codes.Unauthenticated && useAuth {
			var err error
			if callOptions, err = grpcCallOptions(ctx, callOptions); err != nil {
				return err
			}
			return grpcApiCallContext(ctx, callOptions)
		}
		return grpcStatusError
	}
	return nil
}

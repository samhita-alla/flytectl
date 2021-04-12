package auth

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/flyteorg/flytectl/cmd/testutils"
	"github.com/flyteorg/flytectl/pkg/auth/interfaces/mocks"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var setup = testutils.Setup

var (
	ctx                   context.Context
	authClient            service.AuthServiceClient
	grpcAPICallContext    AdminGrpcAPICallContext
	callOptions           []grpc.CallOption
	useAuth               bool
	mockTokenOrchestrator *mocks.FetchTokenOrchestrator
)

func GrpcCallUtilSetup() {
	ctx = testutils.Ctx
	authClient = testutils.MockAuthClient
	useAuth = false
	mockTokenOrchestrator = &mocks.FetchTokenOrchestrator{}
	defaultTokenOrchestrator = mockTokenOrchestrator
	grpcAPICallContext = func(ctx context.Context, callOptions []grpc.CallOption) error {
		return nil
	}
}

func TestDo(t *testing.T) {
	setup()
	GrpcCallUtilSetup()
	err := Do(ctx, authClient, grpcAPICallContext, callOptions, useAuth)
	assert.Nil(t, err)
}

func TestDoWithNonAuthError(t *testing.T) {
	setup()
	GrpcCallUtilSetup()
	grpcAPICallContext = func(ctx context.Context, callOptions []grpc.CallOption) error {
		return fmt.Errorf("nonAuthError")
	}
	err := Do(ctx, authClient, grpcAPICallContext, callOptions, useAuth)
	assert.NotNil(t, err)
	assert.Equal(t, err, fmt.Errorf("nonAuthError"))
}

func TestDoWithAuthErrorWithClientAuthDisabled(t *testing.T) {
	setup()
	GrpcCallUtilSetup()
	grpcAPICallContext = func(ctx context.Context, callOptions []grpc.CallOption) error {
		return status.New(codes.Unauthenticated, "empty identity").Err()
	}
	err := Do(ctx, authClient, grpcAPICallContext, callOptions, useAuth)
	assert.NotNil(t, err)
	assert.Equal(t, err, status.New(codes.Unauthenticated, "empty identity").Err())
}

func TestDoWithAuthErrorWithClientAuthEnabled(t *testing.T) {
	setup()
	GrpcCallUtilSetup()
	useAuth = true
	grpcAPICallContext = func(ctx context.Context, callOptions []grpc.CallOption) error {
		return status.New(codes.Unauthenticated, "empty identity").Err()
	}

	token := &oauth2.Token{
		AccessToken: "fakeAccessToken",
		Expiry:      time.Now().Add(time.Minute * 30),
	}
	mockTokenOrchestrator.OnFetchTokenFromAuthFlowMatch(mock.Anything, mock.Anything, mock.Anything).Return(token, nil)
	mockTokenOrchestrator.OnFetchTokenFromCacheOrRefreshItMatch(mock.Anything, mock.Anything).Return(token)
	mockTokenOrchestrator.OnRefreshTheTokenMatch(mock.Anything, mock.Anything, mock.Anything).Return(token)
	err := Do(ctx, authClient, grpcAPICallContext, callOptions, useAuth)
	assert.NotNil(t, err)
	assert.Equal(t, err, status.New(codes.Unauthenticated, "empty identity").Err())
}

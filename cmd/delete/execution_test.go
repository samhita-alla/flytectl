package delete

import (
	"context"
	"errors"
	"testing"

	"github.com/flyteorg/flytectl/cmd/config"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/cmd/testutils"
	"github.com/flyteorg/flyteidl/clients/go/admin/mocks"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"

	"github.com/stretchr/testify/assert"
)

var (
	ctx                   context.Context
	mockClient            *mocks.AdminServiceClient
	args                  []string
	cmdCtx                cmdCore.CommandContext
	terminateExecRequests []*admin.ExecutionTerminateRequest
	terminateExecResponse *admin.ExecutionTerminateResponse
)

var setup = testutils.Setup

func deleteExecutionSetup() {
	ctx = testutils.Ctx
	cmdCtx = testutils.CmdCtx
	mockClient = testutils.MockAdminClient
	args = append(args, "exec1", "exec2")
	terminateExecRequests = []*admin.ExecutionTerminateRequest{
		{Id: &core.WorkflowExecutionIdentifier{
			Name:    "exec1",
			Project: config.GetConfig().Project,
			Domain:  config.GetConfig().Domain,
		}},
		{Id: &core.WorkflowExecutionIdentifier{
			Name:    "exec2",
			Project: config.GetConfig().Project,
			Domain:  config.GetConfig().Domain,
		}},
	}
	terminateExecResponse = &admin.ExecutionTerminateResponse{}
}

func TestTerminateExecutionFunc(t *testing.T) {
	setup()
	deleteExecutionSetup()
	mockClient.OnTerminateExecutionMatch(ctx, terminateExecRequests[0]).Return(terminateExecResponse, nil)
	mockClient.OnTerminateExecutionMatch(ctx, terminateExecRequests[1]).Return(terminateExecResponse, nil)
	err := terminateExecutionFunc(ctx, args, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "TerminateExecution", ctx, terminateExecRequests[0])
	mockClient.AssertCalled(t, "TerminateExecution", ctx, terminateExecRequests[1])
}

func TestTerminateExecutionFuncWithError(t *testing.T) {
	setup()
	deleteExecutionSetup()
	mockClient.OnTerminateExecutionMatch(ctx, terminateExecRequests[0]).Return(nil, errors.New("failed to terminate"))
	mockClient.OnTerminateExecutionMatch(ctx, terminateExecRequests[1]).Return(terminateExecResponse, nil)
	err := terminateExecutionFunc(ctx, args, cmdCtx)
	assert.Equal(t, errors.New("failed to terminate"), err)
	mockClient.AssertCalled(t, "TerminateExecution", ctx, terminateExecRequests[0])
	mockClient.AssertNotCalled(t, "TerminateExecution", ctx, terminateExecRequests[1])
}

func TestTerminateExecutionFuncWithPartialSuccess(t *testing.T) {
	setup()
	deleteExecutionSetup()
	mockClient.OnTerminateExecutionMatch(ctx, terminateExecRequests[0]).Return(terminateExecResponse, nil)
	mockClient.OnTerminateExecutionMatch(ctx, terminateExecRequests[1]).Return(nil, errors.New("failed to terminate"))
	err := terminateExecutionFunc(ctx, args, cmdCtx)
	assert.Equal(t, errors.New("failed to terminate"), err)
	mockClient.AssertCalled(t, "TerminateExecution", ctx, terminateExecRequests[0])
	mockClient.AssertCalled(t, "TerminateExecution", ctx, terminateExecRequests[1])
}

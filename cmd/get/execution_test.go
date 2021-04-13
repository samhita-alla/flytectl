package get

import (
	"errors"
	"testing"

	"github.com/flyteorg/flytectl/cmd/testutils"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"

	"github.com/stretchr/testify/assert"
)

var (
	execGetRequest    *admin.WorkflowExecutionGetRequest
	execListRequest   *admin.ResourceListRequest
	executionResponse *admin.Execution
	executionList     *admin.ExecutionList
)

func GetExecutionSetup() {
	ctx = testutils.Ctx
	cmdCtx = testutils.CmdCtx
	mockClient = testutils.MockAdminClient
	execGetRequest = &admin.WorkflowExecutionGetRequest{
		Id: &core.WorkflowExecutionIdentifier{
			Project: projectValue,
			Domain:  domainValue,
			Name:    executionNameValue,
		},
	}
	execListRequest = &admin.ResourceListRequest{
		Limit: 100,
		Id: &admin.NamedEntityIdentifier{
			Project: projectValue,
			Domain:  domainValue,
		},
	}
	executionResponse = &admin.Execution{
		Id: &core.WorkflowExecutionIdentifier{
			Project: projectValue,
			Domain:  domainValue,
			Name:    executionNameValue,
		},
		Spec: &admin.ExecutionSpec{
			LaunchPlan: &core.Identifier{
				Project: projectValue,
				Domain:  domainValue,
				Name:    launchPlanNameValue,
				Version: launchPlanVersionValue,
			},
		},
		Closure: &admin.ExecutionClosure{
			WorkflowId: &core.Identifier{
				Project: projectValue,
				Domain:  domainValue,
				Name:    workflowNameValue,
				Version: workflowVersionValue,
			},
			Phase: core.WorkflowExecution_SUCCEEDED,
		},
	}
	executions := []*admin.Execution{executionResponse}
	executionList = &admin.ExecutionList{
		Executions: executions,
	}
}

func TestListExecutionFunc(t *testing.T) {
	setup()
	GetExecutionSetup()
	var args []string
	mockClient.OnListExecutionsMatch(ctx, execListRequest).Return(executionList, nil)
	err := getExecutionFunc(ctx, args, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "ListExecutions", ctx, execListRequest)
}

func TestListExecutionFuncWithError(t *testing.T) {
	setup()
	GetExecutionSetup()
	var args []string
	mockClient.OnListExecutionsMatch(ctx, execListRequest).Return(nil, errors.New("executions NotFound"))
	err := getExecutionFunc(ctx, args, cmdCtx)
	assert.NotNil(t, err)
	assert.Equal(t, err, errors.New("executions NotFound"))
	mockClient.AssertCalled(t, "ListExecutions", ctx, execListRequest)
}

func TestGetExecutionFunc(t *testing.T) {
	setup()
	GetExecutionSetup()
	args := []string{executionNameValue}
	mockClient.OnGetExecutionMatch(ctx, execGetRequest).Return(executionResponse, nil)
	err := getExecutionFunc(ctx, args, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "GetExecution", ctx, execGetRequest)
}

func TestGetExecutionFuncWithError(t *testing.T) {
	setup()
	GetExecutionSetup()
	args := []string{executionNameValue}
	mockClient.OnGetExecutionMatch(ctx, execGetRequest).Return(nil, errors.New("execution NotFound"))
	err := getExecutionFunc(ctx, args, cmdCtx)
	assert.NotNil(t, err)
	assert.Equal(t, err, errors.New("execution NotFound"))
	mockClient.AssertCalled(t, "GetExecution", ctx, execGetRequest)
}

package get

import (
	"errors"
	"testing"

	"github.com/flyteorg/flytectl/cmd/testutils"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	workflowListRequest       *admin.ResourceListRequest
	workflowList              *admin.WorkflowList
	workflowIdentifierList    *admin.NamedEntityIdentifierList
	workflowIdentifierListReq *admin.NamedEntityIdentifierListRequest
	args                      []string
)

func GetWorkflowSetup() {
	ctx = testutils.Ctx
	cmdCtx = testutils.CmdCtx
	mockClient = testutils.MockAdminClient
	var workflowName string
	if len(args) > 0 {
		workflowName = args[0]
	}
	workflowIdentifierListReq = &admin.NamedEntityIdentifierListRequest{
		Project: projectValue,
		Domain:  domainValue,
		SortBy: &admin.Sort{
			Key:       "name",
			Direction: admin.Sort_ASCENDING,
		},
		Limit: 100,
	}
	var entities []*admin.NamedEntityIdentifier
	id1 := &admin.NamedEntityIdentifier{
		Project: projectValue,
		Domain:  domainValue,
		Name:    "worflow1",
	}
	id2 := &admin.NamedEntityIdentifier{
		Project: projectValue,
		Domain:  domainValue,
		Name:    "worflow2",
	}
	entities = append(entities, id1, id2)
	workflowIdentifierList = &admin.NamedEntityIdentifierList{
		Entities: entities,
	}

	workflowListRequest = &admin.ResourceListRequest{
		Id: &admin.NamedEntityIdentifier{
			Project: projectValue,
			Domain:  domainValue,
			Name:    workflowName,
		},
		SortBy: &admin.Sort{
			Key:       "created_at",
			Direction: admin.Sort_DESCENDING,
		},
		Limit: 100,
	}
	var workflowArray []*admin.Workflow
	workflow1 := &admin.Workflow{
		Id: &core.Identifier{
			Name:    "workflow1",
			Version: "v2",
		},
		Closure: &admin.WorkflowClosure{
			CreatedAt: &timestamppb.Timestamp{Seconds: 1, Nanos: 0},
		},
	}
	workflowArray = append(workflowArray, workflow1)
	workflowList = &admin.WorkflowList{
		Workflows: workflowArray,
	}
}

func TestListWorkflowsFunc(t *testing.T) {
	setup()
	GetWorkflowSetup()
	mockClient.OnListWorkflowIdsMatch(ctx, workflowIdentifierListReq).Return(workflowIdentifierList, nil)
	err = getWorkflowFunc(ctx, args, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "ListWorkflowIds", ctx, workflowIdentifierListReq)
}

func TestListWorkflowsFuncWithError(t *testing.T) {
	setup()
	GetWorkflowSetup()
	mockClient.OnListWorkflowIdsMatch(ctx, workflowIdentifierListReq).Return(nil, errors.New("workflows NotFound"))
	err = getWorkflowFunc(ctx, args, cmdCtx)
	assert.NotNil(t, err)
	assert.Equal(t, err, errors.New("workflows NotFound"))
	mockClient.AssertCalled(t, "ListWorkflowIds", ctx, workflowIdentifierListReq)
}

func TestGetWorkflowsFunc(t *testing.T) {
	setup()
	args = append(args, "workflow1")
	GetWorkflowSetup()
	mockClient.OnListWorkflowsMatch(ctx, workflowListRequest).Return(workflowList, nil)
	err = getWorkflowFunc(ctx, args, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "ListWorkflows", ctx, workflowListRequest)
}

func TestGetWorkflowsFuncWithError(t *testing.T) {
	setup()
	args = append(args, "workflow1")
	GetWorkflowSetup()
	mockClient.OnListWorkflowsMatch(ctx, workflowListRequest).Return(nil, errors.New("workflows NotFound"))
	err = getWorkflowFunc(ctx, args, cmdCtx)
	assert.NotNil(t, err)
	assert.Equal(t, err, errors.New("workflows NotFound"))
	mockClient.AssertCalled(t, "ListWorkflows", ctx, workflowListRequest)
}

//
//func TestListProjectsFuncError(t *testing.T) {
//	setup()
//	GetProjectSetup()
//	var args []string
//	mockClient.OnListProjectsMatch(ctx, projectListRequest).Return(nil, errors.New("projects NotFound"))
//	err = getProjectsFunc(ctx, args, cmdCtx)
//	assert.NotNil(t, err)
//	assert.Equal(t, err, errors.New("projects NotFound"))
//	mockClient.AssertCalled(t, "ListProjects", ctx, projectListRequest)
//}
//
//func TestGetProjectFunc(t *testing.T) {
//	setup()
//	GetProjectSetup()
//	args := []string{projectValue}
//	mockClient.OnListProjectsMatch(ctx, projectListRequest).Return(projects, nil)
//	err = getProjectsFunc(ctx, args, cmdCtx)
//	assert.Nil(t, err)
//	mockClient.AssertCalled(t, "ListProjects", ctx, projectListRequest)
//}
//
//func TestGetProjectFuncError(t *testing.T) {
//	setup()
//	GetProjectSetup()
//	args := []string{projectValue}
//	mockClient.OnListProjectsMatch(ctx, projectListRequest).Return(nil, errors.New("projects NotFound"))
//	err = getProjectsFunc(ctx, args, cmdCtx)
//	assert.NotNil(t, err)
//	assert.Equal(t, err, errors.New("projects NotFound"))
//	mockClient.AssertCalled(t, "ListProjects", ctx, projectListRequest)
//}
//
//func TestGetProjectFuncNotFound(t *testing.T) {
//	setup()
//	GetProjectSetup()
//	args := []string{"notFoundProject"}
//	mockClient.OnListProjectsMatch(ctx, projectListRequest).Return(projects, nil)
//	err = getProjectsFunc(ctx, args, cmdCtx)
//	assert.Nil(t, err)
//	mockClient.AssertCalled(t, "ListProjects", ctx, projectListRequest)
//}

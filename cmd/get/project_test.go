package get

import (
	"errors"
	"testing"

	"github.com/flyteorg/flytectl/cmd/testutils"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"

	"github.com/stretchr/testify/assert"
)

var (
	projectListRequest *admin.ProjectListRequest
	projects           *admin.Projects
)

func GetProjectSetup() {
	ctx = testutils.Ctx
	cmdCtx = testutils.CmdCtx
	mockClient = testutils.MockAdminClient
	projectListRequest = &admin.ProjectListRequest{}
	var projectsArray []*admin.Project
	var domains []*admin.Domain
	domain1 := &admin.Domain{Id: domainValue, Name: domainValue}
	domains = append(domains, domain1)
	project1 := &admin.Project{
		Id:      projectValue,
		Name:    projectValue,
		Domains: domains,
	}
	projectsArray = append(projectsArray, project1)
	projects = &admin.Projects{
		Projects: projectsArray,
	}
}

func TestListProjectsFunc(t *testing.T) {
	setup()
	GetProjectSetup()
	var args []string
	mockClient.OnListProjectsMatch(ctx, projectListRequest).Return(projects, nil)
	err = getProjectsFunc(ctx, args, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "ListProjects", ctx, projectListRequest)
}

func TestListProjectsFuncError(t *testing.T) {
	setup()
	GetProjectSetup()
	var args []string
	mockClient.OnListProjectsMatch(ctx, projectListRequest).Return(nil, errors.New("projects NotFound"))
	err = getProjectsFunc(ctx, args, cmdCtx)
	assert.NotNil(t, err)
	assert.Equal(t, err, errors.New("projects NotFound"))
	mockClient.AssertCalled(t, "ListProjects", ctx, projectListRequest)
}

func TestGetProjectFunc(t *testing.T) {
	setup()
	GetProjectSetup()
	args := []string{projectValue}
	mockClient.OnListProjectsMatch(ctx, projectListRequest).Return(projects, nil)
	err = getProjectsFunc(ctx, args, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "ListProjects", ctx, projectListRequest)
}

func TestGetProjectFuncError(t *testing.T) {
	setup()
	GetProjectSetup()
	args := []string{projectValue}
	mockClient.OnListProjectsMatch(ctx, projectListRequest).Return(nil, errors.New("projects NotFound"))
	err = getProjectsFunc(ctx, args, cmdCtx)
	assert.NotNil(t, err)
	assert.Equal(t, err, errors.New("projects NotFound"))
	mockClient.AssertCalled(t, "ListProjects", ctx, projectListRequest)
}

func TestGetProjectFuncNotFound(t *testing.T) {
	setup()
	GetProjectSetup()
	args := []string{"notFoundProject"}
	mockClient.OnListProjectsMatch(ctx, projectListRequest).Return(projects, nil)
	err = getProjectsFunc(ctx, args, cmdCtx)
	assert.Nil(t, err)
	mockClient.AssertCalled(t, "ListProjects", ctx, projectListRequest)
}

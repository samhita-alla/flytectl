package get

import (
	"context"
	"fmt"

	"github.com/flyteorg/flytectl/cmd/config"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/pkg/auth"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"

	"google.golang.org/grpc"
)

// Reads the task config to drive fetching the correct tasks.
func FetchTaskForName(ctx context.Context, name string, project string, domain string, cmdCtx cmdCore.CommandContext) ([]*admin.Task, error) {
	var tasks []*admin.Task
	var err error
	var task *admin.Task
	if taskConfig.Latest {
		if task, err = FetchTaskLatestVersion(ctx, name, project, domain, cmdCtx); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	} else if taskConfig.Version != "" {
		if task, err = FetchTaskVersion(ctx, name, taskConfig.Version, project, domain, cmdCtx); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	} else {
		tasks, err = FetchAllVerOfTask(ctx, name, project, domain, cmdCtx)
		if err != nil {
			return nil, err
		}
	}
	if taskConfig.ExecFile != "" {
		// There would be atleast one task object when code reaches here and hence the length assertion is not required.
		task = tasks[0]
		// Only write the first task from the tasks object.
		if err = CreateAndWriteExecConfigForTask(task, taskConfig.ExecFile); err != nil {
			return nil, err
		}
	}
	return tasks, nil
}

func FetchAllVerOfTask(ctx context.Context, name string, project string, domain string, cmdCtx cmdCore.CommandContext) ([]*admin.Task, error) {
	var callOptions []grpc.CallOption
	var tList *admin.TaskList
	grpcAPICall := func(_ctx context.Context, _callOptions []grpc.CallOption) error {
		var err error
		tList, err = cmdCtx.AdminClient().ListTasks(_ctx, &admin.ResourceListRequest{
			Id: &admin.NamedEntityIdentifier{
				Project: project,
				Domain:  domain,
				Name:    name,
			},
			SortBy: &admin.Sort{
				Key:       "created_at",
				Direction: admin.Sort_DESCENDING,
			},
			Limit: 100,
		}, callOptions...)
		if err != nil {
			return err
		}
		return nil
	}
	err := auth.Do(ctx, cmdCtx.AuthClient(), grpcAPICall, callOptions, config.GetConfig().UseAuth)
	if err != nil {
		return nil, err
	}
	if len(tList.Tasks) == 0 {
		return nil, fmt.Errorf("no tasks retrieved for %v", name)
	}
	return tList.Tasks, nil
}

func FetchTaskLatestVersion(ctx context.Context, name string, project string, domain string, cmdCtx cmdCore.CommandContext) (*admin.Task, error) {
	var t *admin.Task
	var err error
	// Fetch the latest version of the task.
	var taskVersions []*admin.Task
	taskVersions, err = FetchAllVerOfTask(ctx, name, project, domain, cmdCtx)
	if err != nil {
		return nil, err
	}
	t = taskVersions[0]
	return t, nil
}

func FetchTaskVersion(ctx context.Context, name string, version string, project string, domain string, cmdCtx cmdCore.CommandContext) (*admin.Task, error) {
	var callOptions []grpc.CallOption
	var t *admin.Task
	grpcAPICall := func(_ctx context.Context, _callOptions []grpc.CallOption) error {
		var err error
		t, err = cmdCtx.AdminClient().GetTask(_ctx, &admin.ObjectGetRequest{
			Id: &core.Identifier{
				ResourceType: core.ResourceType_TASK,
				Project:      project,
				Domain:       domain,
				Name:         name,
				Version:      version,
			},
		}, _callOptions...)
		if err != nil {
			return err
		}
		return nil
	}
	err := auth.Do(ctx, cmdCtx.AuthClient(), grpcAPICall, callOptions, config.GetConfig().UseAuth)
	if err != nil {
		return nil, err
	}
	return t, nil
}

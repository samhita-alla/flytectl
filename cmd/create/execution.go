package create

import (
	"context"
	"fmt"
	"github.com/flyteorg/flytectl/pkg/auth"
	"google.golang.org/grpc"

	"github.com/flyteorg/flytectl/cmd/config"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
)

const (
	executionShort = "Create execution resources"
	executionLong  = `
Create the executions for given workflow/task in a project and domain.

There are three steps in generating an execution.

- Generate the execution spec file using the get command.
- Update the inputs for the execution if needed.
- Run the execution by passing in the generated yaml file.

The spec file should be generated first and then run the execution using the spec file.
You can reference the flytectl get task for more details

::

 flytectl get tasks -d development -p flytectldemo core.advanced.run_merge_sort.merge  --version v2 --execFile execution_spec.yaml

The generated file would look similar to this

.. code-block:: yaml

	 iamRoleARN: ""
	 inputs:
	   sorted_list1:
	   - 0
	   sorted_list2:
	   - 0
	 kubeServiceAcct: ""
	 targetDomain: ""
	 targetProject: ""
	 task: core.advanced.run_merge_sort.merge
	 version: "v2"


The generated file can be modified to change the input values.

.. code-block:: yaml

	 iamRoleARN: 'arn:aws:iam::12345678:role/defaultrole'
	 inputs:
	   sorted_list1:
	   - 2
	   - 4
	   - 6
	   sorted_list2:
	   - 1
	   - 3
	   - 5
	 kubeServiceAcct: ""
	 targetDomain: ""
	 targetProject: ""
	 task: core.advanced.run_merge_sort.merge
	 version: "v2"

And then can be passed through the command line.
Notice the source and target domain/projects can be different.
The root project and domain flags of -p and -d should point to task/launch plans project/domain.

::

 flytectl create execution --execFile execution_spec.yaml -p flytectldemo -d development --targetProject flytesnacks

Usage
`
)

//go:generate pflags ExecutionConfig --default-var executionConfig

// ExecutionConfig hold configuration for create execution flags and configuration of the actual task or workflow  to be launched.
type ExecutionConfig struct {
	// pflag section
	ExecFile        string `json:"execFile,omitempty" pflag:",file for the execution params.If not specified defaults to <<workflow/task>_name>.execution_spec.yaml"`
	TargetDomain    string `json:"targetDomain" pflag:",project where execution needs to be created.If not specified configured domain would be used."`
	TargetProject   string `json:"targetProject" pflag:",project where execution needs to be created.If not specified configured project would be used."`
	KubeServiceAcct string `json:"kubeServiceAcct" pflag:",kubernetes service account AuthRole for launching execution."`
	IamRoleARN      string `json:"iamRoleARN" pflag:",iam role ARN AuthRole for launching execution."`
	// Non plfag section is read from the execution config generated by get task/launchplan
	Workflow string                 `json:"workflow,omitempty"`
	Task     string                 `json:"task,omitempty"`
	Version  string                 `json:"version"`
	Inputs   map[string]interface{} `json:"inputs"`
}

type ExecutionParams struct {
	name   string
	isTask bool
}

var (
	executionConfig = &ExecutionConfig{}
)

func createExecutionCommand(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	var execParams ExecutionParams
	var err error
	sourceProject := config.GetConfig().Project
	sourceDomain := config.GetConfig().Domain
	if execParams, err = readConfigAndValidate(config.GetConfig().Project, config.GetConfig().Domain); err != nil {
		return err
	}
	var executionRequest *admin.ExecutionCreateRequest
	if execParams.isTask {
		if executionRequest, err = createExecutionRequestForTask(ctx, execParams.name, sourceProject, sourceDomain, cmdCtx); err != nil {
			return err
		}
	} else {
		if executionRequest, err = createExecutionRequestForWorkflow(ctx, execParams.name, sourceProject, sourceDomain, cmdCtx); err != nil {
			return err
		}
	}
	var callOptions []grpc.CallOption
	var exec *admin.ExecutionCreateResponse
	grpcApiCall := func(_ctx context.Context, _callOptions []grpc.CallOption) error {
		var err error
		exec, err = cmdCtx.AdminClient().CreateExecution(_ctx, executionRequest, _callOptions...)
		return err
	}
	err = auth.Do(grpcApiCall, ctx, callOptions, config.GetConfig().UseAuth)
	if err != nil {
		return err
	}
	fmt.Printf("execution identifier %v\n", exec.Id)
	return nil
}

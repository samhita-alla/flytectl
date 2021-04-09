package get

import (
	"context"

	"github.com/flyteorg/flytectl/cmd/config"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/pkg/auth"
	"github.com/flyteorg/flytectl/pkg/printer"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"
	"github.com/flyteorg/flytestdlib/logger"
	"github.com/golang/protobuf/proto"

	"google.golang.org/grpc"
)

const (
	executionShort = "Gets execution resources"
	executionLong  = `
Retrieves all the executions within project and domain.(execution,executions can be used interchangeably in these commands)
::

 bin/flytectl get execution -p flytesnacks -d development

Retrieves execution by name within project and domain.

::

 bin/flytectl get execution -p flytesnacks -d development oeh94k9r2r

Retrieves execution by filters
::

 Not yet implemented

Retrieves all the execution within project and domain in yaml format

::

 bin/flytectl get execution -p flytesnacks -d development -o yaml

Retrieves all the execution within project and domain in json format.

::

 bin/flytectl get execution -p flytesnacks -d development -o json

Usage
`
)

var executionColumns = []printer.Column{
	{Header: "Name", JSONPath: "$.id.name"},
	{Header: "Launch Plan Name", JSONPath: "$.spec.launchPlan.name"},
	{Header: "Type", JSONPath: "$.spec.launchPlan.resourceType"},
	{Header: "Phase", JSONPath: "$.closure.phase"},
	{Header: "Started", JSONPath: "$.closure.startedAt"},
	{Header: "Elapsed Time", JSONPath: "$.closure.duration"},
}

func ExecutionToProtoMessages(l []*admin.Execution) []proto.Message {
	messages := make([]proto.Message, 0, len(l))
	for _, m := range l {
		messages = append(messages, m)
	}
	return messages
}

func getExecutionFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	adminPrinter := printer.Printer{}
	var executions []*admin.Execution
	var callOptions []grpc.CallOption
	if len(args) > 0 {
		name := args[0]
		var execution *admin.Execution
		grpcAPICall := func(_ctx context.Context, _callOptions []grpc.CallOption) error {
			var err error
			execution, err = cmdCtx.AdminClient().GetExecution(
				_ctx, &admin.WorkflowExecutionGetRequest{
					Id: &core.WorkflowExecutionIdentifier{
						Project: config.GetConfig().Project,
						Domain:  config.GetConfig().Domain,
						Name:    name,
					},
				}, _callOptions...)
			return err
		}
		err := auth.Do(ctx, grpcAPICall, callOptions, config.GetConfig().UseAuth)
		if err != nil {
			return err
		}
		executions = append(executions, execution)
	} else {
		var executionList *admin.ExecutionList
		grpcAPICallListExecs := func(_ctx context.Context, _callOptions []grpc.CallOption) error {
			var err error
			executionList, err = cmdCtx.AdminClient().ListExecutions(
				_ctx, &admin.ResourceListRequest{
					Limit: 100,
					Id: &admin.NamedEntityIdentifier{
						Project: config.GetConfig().Project,
						Domain:  config.GetConfig().Domain,
					},
				}, _callOptions...)
			return err
		}
		err := auth.Do(ctx, grpcAPICallListExecs, callOptions, config.GetConfig().UseAuth)
		if err != nil {
			return err
		}
		executions = executionList.Executions
	}
	logger.Infof(ctx, "Retrieved %v executions", len(executions))
	err := adminPrinter.Print(config.GetConfig().MustOutputFormat(), executionColumns, ExecutionToProtoMessages(executions)...)
	if err != nil {
		return err
	}
	return nil
}

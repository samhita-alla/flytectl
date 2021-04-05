package get

import (
	"context"
	"github.com/flyteorg/flytectl/pkg/auth"
	"google.golang.org/grpc"

	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flytestdlib/logger"
	"github.com/golang/protobuf/proto"

	"github.com/flyteorg/flytectl/cmd/config"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/pkg/printer"
)

const (
	projectShort = "Gets project resources"
	projectLong  = `
Retrieves all the projects.(project,projects can be used interchangeably in these commands)
::

 bin/flytectl get project

Retrieves project by name

::

 bin/flytectl get project flytesnacks

Retrieves project by filters
::

 Not yet implemented

Retrieves all the projects in yaml format

::

 bin/flytectl get project -o yaml

Retrieves all the projects in json format

::

 bin/flytectl get project -o json

Usage
`
)

var projectColumns = []printer.Column{
	{Header: "ID", JSONPath: "$.id"},
	{Header: "Name", JSONPath: "$.name"},
	{Header: "Description", JSONPath: "$.description"},
}

func ProjectToProtoMessages(l []*admin.Project) []proto.Message {
	messages := make([]proto.Message, 0, len(l))
	for _, m := range l {
		messages = append(messages, m)
	}
	return messages
}

func getProjectsFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	adminPrinter := printer.Printer{}
	var projects *admin.Projects
	var callOptions []grpc.CallOption
	grpcApiCall := func(_ctx context.Context, _callOptions []grpc.CallOption) error {
		var err error
		projects, err = cmdCtx.AdminClient().ListProjects(ctx, &admin.ProjectListRequest{}, _callOptions...)
		if err != nil {
			return err
		}
		return nil
	}
	// useAuth will be controlled by a flag.
	if err := auth.Do(grpcApiCall, ctx, callOptions, true); err != nil {
		return err
	}
	if len(args) == 1 {
		name := args[0]
		logger.Debugf(ctx, "Retrieved %v projects", len(projects.Projects))
		for _, v := range projects.Projects {
			if v.Name == name {
				err := adminPrinter.Print(config.GetConfig().MustOutputFormat(), projectColumns, v)
				if err != nil {
					return err
				}
				return nil
			}
		}
		return nil
	}
	logger.Debugf(ctx, "Retrieved %v projects", len(projects.Projects))
	return adminPrinter.Print(config.GetConfig().MustOutputFormat(), projectColumns, ProjectToProtoMessages(projects.Projects)...)
}

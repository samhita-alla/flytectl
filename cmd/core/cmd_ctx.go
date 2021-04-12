package cmdcore

import (
	"io"

	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/service"
)

type CommandContext struct {
	adminClient service.AdminServiceClient
	authClient  service.AuthServiceClient
	in          io.Reader
	out         io.Writer
}

func NewCommandContext(adminClient service.AdminServiceClient, authClient service.AuthServiceClient, out io.Writer) CommandContext {
	return CommandContext{adminClient: adminClient, authClient: authClient, out: out}
}

func (c CommandContext) AdminClient() service.AdminServiceClient {
	return c.adminClient
}

func (c CommandContext) AuthClient() service.AuthServiceClient {
	return c.authClient
}

func (c CommandContext) OutputPipe() io.Writer {
	return c.out
}

func (c CommandContext) InputPipe() io.Reader {
	return c.in
}

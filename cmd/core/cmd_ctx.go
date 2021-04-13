package cmdcore

import (
	"github.com/flyteorg/flyteidl/clients/go/admin"
	"io"

	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/service"
)

type CommandContext struct {
	clientSet   *admin.Clientset
	in          io.Reader
	out         io.Writer
}

func NewCommandContext(clientSet *admin.Clientset, out io.Writer) CommandContext {
	return CommandContext{clientSet: clientSet, out: out}
}

func (c CommandContext) AdminClient() service.AdminServiceClient {
	if c.clientSet == nil {
		return nil
	}
	return c.clientSet.AdminClient()
}

func (c CommandContext) AuthClient() service.AuthServiceClient {
	if c.clientSet == nil {
		return nil
	}
	return c.clientSet.AuthClient()
}

func (c CommandContext) OutputPipe() io.Writer {
	return c.out
}

func (c CommandContext) InputPipe() io.Reader {
	return c.in
}

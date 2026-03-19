package e2e

import (
	"bytes"
	"io"
	"log/slog"
	"testing"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	client "github.com/metal-stack/api/go/client"
	"github.com/metal-stack/cli/cmd"
	"github.com/metal-stack/cli/cmd/completion"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

type TestClientConfig[Request, Response any] struct {
	WantRequest  Request  // for client expectation
	WantResponse Response // for client return

	FsMocks   func(fs *afero.Afero)
	MockStdin *bytes.Buffer
}

func NewRootCmd[Request, Response any](t *testing.T, c *TestClientConfig[Request, Response]) NewRootCmdFunc {
	interceptors := []connect.Interceptor{
		validate.NewInterceptor(),
		&testClientInterceptor[Request, Response]{
			t:        t,
			response: c.WantResponse,
			request:  c.WantRequest,
		},
	}

	cl, err := client.New(&client.DialConfig{
		BaseURL:      "http://this-is-just-for-testing",
		Interceptors: interceptors,
		UserAgent:    "cli-test",
		Log:          slog.Default(),
	})
	require.NoError(t, err)

	fs := afero.NewMemMapFs()
	if c.FsMocks != nil {
		c.FsMocks(&afero.Afero{
			Fs: fs,
		})
	}

	var in io.Reader
	if c.MockStdin != nil {
		in = bytes.NewReader(c.MockStdin.Bytes())
	}

	return func() (*cobra.Command, *bytes.Buffer) {
		var out bytes.Buffer

		return cmd.NewRootCmd(&config.Config{
			Fs:         fs,
			Out:        &out,
			In:         in,
			PromptOut:  io.Discard,
			Completion: &completion.Completion{},
			Client:     cl,
		}), &out
	}
}

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

type TestClientConfig struct {
	FsMocks   func(fs *afero.Afero)
	MockStdin *bytes.Buffer
}

func NewRootCmd(t *testing.T, c *TestClientConfig, calls ...ClientCall) NewRootCmdFunc {
	return func() (*cobra.Command, *bytes.Buffer) {
		interceptors := []connect.Interceptor{
			&testClientInterceptor{
				t:     t,
				calls: calls,
				count: 0,
			},
			validate.NewInterceptor(),
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

		var out bytes.Buffer

		return cmd.NewRootCmd(&config.Config{
			Fs:        fs,
			Out:       &out,
			In:        in,
			PromptOut: io.Discard,
			Completion: &completion.Completion{
				Client: cl,
			},
			Client: cl,
		}), &out
	}
}

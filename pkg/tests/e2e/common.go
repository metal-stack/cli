package e2e

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"

	"slices"

	"buf.build/go/protoyaml"
	"connectrpc.com/connect"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	client "github.com/metal-stack/api/go/client"
	"github.com/metal-stack/cli/cmd"
	"github.com/metal-stack/cli/cmd/completion"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/testcommon"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoimpl"
)

type Test[Request, Response any] struct {
	Name string
	Cmd  func() []string

	Client    client.Client
	FsMocks   func(fs afero.Fs)
	MockStdin *bytes.Buffer

	// DisableMockClient bool // can switch off mock client creation

	WantErr       error
	WantRequest   Request
	WantResponse  Response      // for client return and json and yaml
	WantObject    proto.Message // domain object for yaml/json structural comparison
	WantTable     *string       // for table printer
	WantWideTable *string       // for wide table printer
	Template      *string       // for template printer
	WantTemplate  *string       // for template printer
	WantMarkdown  *string       // for markdown printer
}

func (c *Test[Request, Response]) TestCmd(t *testing.T) {
	require.NotEmpty(t, c.Name, "test name must not be empty")
	require.NotEmpty(t, c.Cmd, "cmd must not be empty")

	if c.WantErr != nil {
		_, _, conf := c.newCmdConfig(t)

		cmd := cmd.NewRootCmd(conf)
		os.Args = append([]string{config.BinaryName}, c.Cmd()...)

		err := cmd.Execute()
		if diff := cmp.Diff(c.WantErr, err, testcommon.IgnoreUnexported(), testcommon.ErrorStringComparer()); diff != "" {
			t.Errorf("error diff (+got -want):\n %s", diff)
		}
	}

	for _, format := range outputFormats(c) {
		t.Run(fmt.Sprintf("%v", format.Args()), func(t *testing.T) {
			_, out, conf := c.newCmdConfig(t)

			cmd := cmd.NewRootCmd(conf)
			os.Args = append([]string{config.BinaryName}, c.Cmd()...)
			os.Args = append(os.Args, format.Args()...)

			err := cmd.Execute()
			require.NoError(t, err)

			format.Validate(t, out.Bytes())
		})
	}
}

func (c *Test[Request, Response]) newCmdConfig(t *testing.T) (any, *bytes.Buffer, *config.Config) {
	interceptors := []connect.Interceptor{
		&testClientInterceptor[Request, Response]{
			t:        t,
			response: c.WantResponse,
			request:  c.WantRequest,
		},
		// validate.NewInterceptor(),
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
		c.FsMocks(fs)
	}

	var in io.Reader
	if c.MockStdin != nil {
		in = bytes.NewReader(c.MockStdin.Bytes())
	}

	var (
		out    bytes.Buffer
		config = &config.Config{
			Fs:         fs,
			Out:        &out,
			In:         in,
			PromptOut:  io.Discard,
			Completion: &completion.Completion{},
			Client:     cl,
		}
	)

	return nil, &out, config
}

func AssertExhaustiveArgs(t *testing.T, args []string, exclude ...string) {
	assertContainsPrefix := func(ss []string, prefix string) error {
		for _, s := range ss {
			if strings.HasPrefix(s, prefix) {
				return nil
			}
		}
		return fmt.Errorf("not exhaustive: does not contain %q", prefix)
	}

	root := cmd.NewRootCmd(&config.Config{})
	cmd, args, err := root.Find(args)
	require.NoError(t, err)

	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		if slices.Contains(exclude, f.Name) {
			return
		}
		require.NoError(t, assertContainsPrefix(args, "--"+f.Name), "please ensure you all available args are used in order to increase coverage or exclude them explicitly")
	})
}

func outputFormats[Request, Response any](c *Test[Request, Response]) []outputFormat[Response] {
	var formats []outputFormat[Response]

	if c.WantObject != nil {
		formats = append(formats,
			&protoYAMLOutputFormat[Response]{want: c.WantObject},
			&protoJSONOutputFormat[Response]{want: c.WantObject},
		)
	}

	if c.WantTable != nil {
		formats = append(formats, &tableOutputFormat[Response]{table: *c.WantTable})
	}

	if c.WantWideTable != nil {
		formats = append(formats, &wideTableOutputFormat[Response]{table: *c.WantWideTable})
	}

	if c.Template != nil && c.WantTemplate != nil {
		formats = append(formats, &templateOutputFormat[Response]{template: *c.Template, templateOutput: *c.WantTemplate})
	}

	if c.WantMarkdown != nil {
		formats = append(formats, &markdownOutputFormat[Response]{table: *c.WantMarkdown})
	}

	return formats
}

type outputFormat[R any] interface {
	Args() []string
	Validate(t *testing.T, output []byte)
}

type protoYAMLOutputFormat[R any] struct {
	want proto.Message
}

func (o *protoYAMLOutputFormat[R]) Args() []string {
	return []string{"-o", "yaml"}
}

func (o *protoYAMLOutputFormat[R]) Validate(t *testing.T, output []byte) {
	t.Logf("got following yaml output:\n\n%s\n\nconsider using this for test comparison if it looks correct.", string(output))

	got := proto.Clone(o.want)
	proto.Reset(got)

	err := protoyaml.Unmarshal(output, got)
	require.NoError(t, err)

	if diff := cmp.Diff(o.want, got, testcommon.IgnoreUnexported(), cmpopts.IgnoreTypes(protoimpl.MessageState{})); diff != "" {
		t.Errorf("diff (+got -want):\n %s", diff)
	}
}

type protoJSONOutputFormat[R any] struct {
	want proto.Message
}

func (o *protoJSONOutputFormat[R]) Args() []string {
	return []string{"-o", "json"}
}

func (o *protoJSONOutputFormat[R]) Validate(t *testing.T, output []byte) {
	t.Logf("got following json output:\n\n%s\n\nconsider using this for test comparison if it looks correct.", string(output))

	got := proto.Clone(o.want)
	proto.Reset(got)

	err := protojson.Unmarshal(output, got)
	require.NoError(t, err)

	if diff := cmp.Diff(o.want, got, testcommon.IgnoreUnexported(), cmpopts.IgnoreTypes(protoimpl.MessageState{})); diff != "" {
		t.Errorf("diff (+got -want):\n %s", diff)
	}
}

type tableOutputFormat[R any] struct {
	table string
}

func (o *tableOutputFormat[R]) Args() []string {
	return []string{"-o", "table"}
}

func (o *tableOutputFormat[R]) Validate(t *testing.T, output []byte) {
	validateTableRows(t, o.table, string(output))
}

type wideTableOutputFormat[R any] struct {
	table string
}

func (o *wideTableOutputFormat[R]) Args() []string {
	return []string{"-o", "wide"}
}

func (o *wideTableOutputFormat[R]) Validate(t *testing.T, output []byte) {
	validateTableRows(t, o.table, string(output))
}

type templateOutputFormat[R any] struct {
	template       string
	templateOutput string
}

func (o *templateOutputFormat[R]) Args() []string {
	return []string{"-o", "template", "--template", o.template}
}

func (o *templateOutputFormat[R]) Validate(t *testing.T, output []byte) {
	t.Logf("got following template output:\n\n%s\n\nconsider using this for test comparison if it looks correct.", string(output))

	assert.Equal(t, strings.TrimSpace(o.templateOutput), strings.TrimSpace(string(output)))
	// if diff := cmp.Diff(strings.TrimSpace(o.templateOutput), strings.TrimSpace(string(output))); diff != "" {
	// 	t.Errorf("diff (+got -want):\n %s", diff)
	// }
}

type markdownOutputFormat[R any] struct {
	table string
}

func (o *markdownOutputFormat[R]) Args() []string {
	return []string{"-o", "markdown"}
}

func (o *markdownOutputFormat[R]) Validate(t *testing.T, output []byte) {
	validateTableRows(t, o.table, string(output))
}

func validateTableRows(t *testing.T, want, got string) {
	trimAll := func(ss []string) []string {
		var res []string
		for _, s := range ss {
			res = append(res, strings.TrimSpace(s))
		}
		return res
	}

	var (
		trimmedWant = strings.TrimSpace(want)
		trimmedGot  = strings.TrimSpace(string(got))

		wantRows = trimAll(strings.Split(trimmedWant, "\n"))
		gotRows  = trimAll(strings.Split(trimmedGot, "\n"))
	)

	t.Logf("got following table output:\n\n%s\n\nconsider using this for test comparison if it looks correct.", trimmedGot)

	// t.Log(cmp.Diff(trimmedWant, trimmedGot))

	require.Equal(t, len(wantRows), len(gotRows), "tables have different lengths")

	for i := range wantRows {
		wantFields := trimAll(strings.Split(wantRows[i], " "))
		gotFields := trimAll(strings.Split(gotRows[i], " "))

		require.Equal(t, len(wantFields), len(gotFields), "table fields have different lengths")

		for i := range wantFields {
			assert.Equal(t, wantFields[i], gotFields[i])
		}
	}
}

type testClientInterceptor[Request, Response any] struct {
	t        *testing.T
	request  Request
	response Response
}

func (t *testClientInterceptor[Request, Response]) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, ar connect.AnyRequest) (connect.AnyResponse, error) {
		assert.Equal(t.t, &t.request, ar.Any())
		return connect.NewResponse(&t.response), nil
	}
}

func (t *testClientInterceptor[Request, Response]) WrapStreamingClient(connect.StreamingClientFunc) connect.StreamingClientFunc {
	t.t.Errorf("streaming not supported")
	return nil
}

func (t *testClientInterceptor[Request, Response]) WrapStreamingHandler(connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	t.t.Errorf("streaming not supported")
	return nil
}

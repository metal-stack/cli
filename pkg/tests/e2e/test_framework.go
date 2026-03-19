package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"slices"

	"buf.build/go/protoyaml"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/metal-stack/cli/cmd"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/testcommon"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoimpl"
	"sigs.k8s.io/yaml"
)

// NewRootCmdFunc returns the root command for the cli and an output buffer which returns the output after command execution
type NewRootCmdFunc func() (rootCmd *cobra.Command, out *bytes.Buffer)

type Test[Response, Object any] struct {
	Name string

	NewRootCmd NewRootCmdFunc
	CmdArgs    []string
	Out        *bytes.Buffer

	// output format tests
	WantObject      Object        // object for rawyaml/rawjson structural comparison
	WantProtoObject proto.Message // object for yaml/json structural comparison
	WantTable       *string       // for table printer
	WantWideTable   *string       // for wide table printer
	WantMarkdown    *string       // for markdown printer
	WantTemplate    *string       // for template printer
	Template        *string       // for template printer

	WantErr error
}

func (c *Test[Response, Object]) TestCmd(t *testing.T) {
	require.NotEmpty(t, c.Name, "test name must not be empty")
	require.NotEmpty(t, c.CmdArgs, "cmd must not be empty")

	rootCmd, out := c.NewRootCmd()

	if c.WantErr != nil {
		os.Args = append([]string{rootCmd.Use}, c.CmdArgs...)

		err := rootCmd.Execute()
		if diff := cmp.Diff(c.WantErr, err, testcommon.IgnoreUnexported(), testcommon.ErrorStringComparer()); diff != "" {
			t.Errorf("error diff (+got -want):\n %s", diff)
		}
	}

	formats := outputFormats(c)

	if len(formats) == 0 {
		t.Errorf("at least one want section for output formats must be specified, otherwise no command is getting executed")
		return
	}

	for _, format := range formats {
		succeeded := t.Run(fmt.Sprintf("%v", format.Args()), func(t *testing.T) {
			out.Reset()

			os.Args = append([]string{rootCmd.Use}, c.CmdArgs...)
			os.Args = append(os.Args, format.Args()...)

			err := rootCmd.Execute()
			require.NoError(t, err)

			format.Validate(t, out.Bytes())
		})

		if !succeeded {
			t.FailNow()
		}
	}
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

func outputFormats[Response, Object any](c *Test[Response, Object]) []outputFormat {
	var formats []outputFormat

	if c.WantProtoObject != nil {
		formats = append(formats,
			&protoYAMLOutputFormat[Response]{want: c.WantProtoObject},
			&rawYamlOutputFormat[Object]{want: c.WantObject},
			&protoJSONOutputFormat[Response]{want: c.WantProtoObject},
			&rawJsonOutputFormat[Object]{want: c.WantObject},
		)
	}

	if c.WantTable != nil {
		formats = append(formats, &tableOutputFormat{table: *c.WantTable})
	}

	if c.WantWideTable != nil {
		formats = append(formats, &wideTableOutputFormat{table: *c.WantWideTable})
	}

	if c.Template != nil && c.WantTemplate != nil {
		formats = append(formats, &templateOutputFormat{template: *c.Template, templateOutput: *c.WantTemplate})
	}

	if c.WantMarkdown != nil {
		formats = append(formats, &markdownOutputFormat{table: *c.WantMarkdown})
	}

	return formats
}

type outputFormat interface {
	Args() []string
	Validate(t *testing.T, output []byte)
}

type rawYamlOutputFormat[R any] struct {
	want R
}

func (o *rawYamlOutputFormat[R]) Args() []string {
	return []string{"-o", "yamlraw"}
}

func (o *rawYamlOutputFormat[R]) Validate(t *testing.T, output []byte) {
	t.Logf("got following yamlraw output:\n\n%s\n", string(output))

	var got R

	err := yaml.Unmarshal(output, &got)
	require.NoError(t, err)

	if diff := cmp.Diff(o.want, got, testcommon.IgnoreUnexported(), cmpopts.IgnoreTypes(protoimpl.MessageState{})); diff != "" {
		t.Errorf("diff (+got -want):\n %s", diff)
	}
}

type rawJsonOutputFormat[R any] struct {
	want R
}

func (o *rawJsonOutputFormat[R]) Args() []string {
	return []string{"-o", "jsonraw"}
}

func (o *rawJsonOutputFormat[R]) Validate(t *testing.T, output []byte) {
	t.Logf("got following jsonraw output:\n\n%s\n", string(output))

	var got R

	err := json.Unmarshal(output, &got)
	require.NoError(t, err)

	if diff := cmp.Diff(o.want, got, testcommon.IgnoreUnexported(), cmpopts.IgnoreTypes(protoimpl.MessageState{})); diff != "" {
		t.Errorf("diff (+got -want):\n %s", diff)
	}
}

type protoYAMLOutputFormat[R any] struct {
	want proto.Message
}

func (o *protoYAMLOutputFormat[R]) Args() []string {
	return []string{"-o", "yaml"}
}

func (o *protoYAMLOutputFormat[R]) Validate(t *testing.T, output []byte) {
	t.Logf("got following yaml output:\n\n%s\n", string(output))

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
	t.Logf("got following json output:\n\n%s\n", string(output))

	got := proto.Clone(o.want)
	proto.Reset(got)

	err := protojson.Unmarshal(output, got)
	require.NoError(t, err)

	if diff := cmp.Diff(o.want, got, testcommon.IgnoreUnexported(), cmpopts.IgnoreTypes(protoimpl.MessageState{})); diff != "" {
		t.Errorf("diff (+got -want):\n %s", diff)
	}
}

type tableOutputFormat struct {
	table string
}

func (o *tableOutputFormat) Args() []string {
	return []string{"-o", "table"}
}

func (o *tableOutputFormat) Validate(t *testing.T, output []byte) {
	validateTableRows(t, o.table, string(output))
}

type wideTableOutputFormat struct {
	table string
}

func (o *wideTableOutputFormat) Args() []string {
	return []string{"-o", "wide"}
}

func (o *wideTableOutputFormat) Validate(t *testing.T, output []byte) {
	validateTableRows(t, o.table, string(output))
}

type templateOutputFormat struct {
	template       string
	templateOutput string
}

func (o *templateOutputFormat) Args() []string {
	return []string{"-o", "template", "--template", o.template}
}

func (o *templateOutputFormat) Validate(t *testing.T, output []byte) {
	t.Logf("got following template output:\n\n%s\n\nconsider using this for test comparison if it looks correct.", string(output))

	if diff := cmp.Diff(strings.TrimSpace(o.templateOutput), strings.TrimSpace(string(output))); diff != "" {
		t.Errorf("diff (+got -want):\n %s", diff)
	}
}

type markdownOutputFormat struct {
	table string
}

func (o *markdownOutputFormat) Args() []string {
	return []string{"-o", "markdown"}
}

func (o *markdownOutputFormat) Validate(t *testing.T, output []byte) {
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

	// somehow this diff does not look nice anymore. :(
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

package e2e

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/metal-stack/metal-lib/pkg/testcommon"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/testing/protocmp"
)

type testClientInterceptor struct {
	t     *testing.T
	calls []ClientCall
	count int
}

type ClientCall struct {
	WantRequest  any
	WantResponse func() connect.AnyResponse
	WantError    *connect.Error
}

func (t *testClientInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, ar connect.AnyRequest) (connect.AnyResponse, error) {
		defer func() { t.count++ }()

		if t.count >= len(t.calls) {
			t.t.Errorf("received an unexpected client call of type %T: %v", ar.Any(), ar.Any())
			t.t.FailNow()
		}

		call := t.calls[t.count]

		if diff := cmp.Diff(call.WantRequest, ar.Any(), protocmp.Transform(), testcommon.IgnoreUnexported(), cmpopts.IgnoreTypes(protoimpl.MessageState{})); diff != "" {
			t.t.Errorf("request diff (+got -want):\n %s", diff)
			t.t.FailNow()
		}

		if call.WantError != nil {
			return nil, call.WantError
		}

		return call.WantResponse(), nil
	}
}

func (t *testClientInterceptor) WrapStreamingClient(connect.StreamingClientFunc) connect.StreamingClientFunc {
	t.t.Errorf("streaming not supported")
	return nil
}

func (t *testClientInterceptor) WrapStreamingHandler(connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	t.t.Errorf("streaming not supported")
	return nil
}

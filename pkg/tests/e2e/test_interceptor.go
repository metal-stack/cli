package e2e

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/metal-stack/metal-lib/pkg/testcommon"
	"google.golang.org/protobuf/runtime/protoimpl"
)

type testClientInterceptor[Request, Response any] struct {
	t        *testing.T
	request  Request
	response Response
}

func (t *testClientInterceptor[Request, Response]) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, ar connect.AnyRequest) (connect.AnyResponse, error) {
		if diff := cmp.Diff(&t.request, ar.Any(), testcommon.IgnoreUnexported(), cmpopts.IgnoreTypes(protoimpl.MessageState{})); diff != "" {
			t.t.Errorf("request diff (+got -want):\n %s", diff)
			t.t.FailNow()
		}

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

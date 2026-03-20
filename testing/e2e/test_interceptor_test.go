package e2e

import (
	"log/slog"
	"testing"

	"connectrpc.com/connect"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	client "github.com/metal-stack/api/go/client"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/testcommon"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestInterceptor(t *testing.T) {
	cl, err := client.New(&client.DialConfig{
		BaseURL: "http://this-is-just-for-testing",
		Interceptors: []connect.Interceptor{
			&testClientInterceptor{
				t: t,
				calls: []ClientCall{
					{
						WantRequest: &apiv2.IPServiceGetRequest{
							Ip: "1.2.3.4",
						},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&apiv2.IPServiceGetResponse{
								Ip: &apiv2.IP{Ip: "1.2.3.4"},
							})
						},
					},
				},
			},
		},
		UserAgent: "cli-test",
		Log:       slog.Default(),
	})
	require.NoError(t, err)

	resp, err := cl.Apiv2().IP().Get(t.Context(), &apiv2.IPServiceGetRequest{
		Ip: "1.2.3.4",
	})
	require.NoError(t, err)

	if diff := cmp.Diff(&apiv2.IPServiceGetResponse{
		Ip: &apiv2.IP{
			Ip: "1.2.3.4",
		},
	}, resp, protocmp.Transform(), testcommon.IgnoreUnexported(), cmpopts.IgnoreTypes(protoimpl.MessageState{})); diff != "" {
		t.Errorf("diff = %s", diff)
	}
}

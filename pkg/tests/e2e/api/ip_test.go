package api_e2e

import (
	"testing"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/pkg/tests/e2e"
	"github.com/metal-stack/metal-lib/pkg/tag"
)

var (
	ip1 = func() *apiv2.IP {
		return &apiv2.IP{
			Uuid:        "2e0144a2-09ef-42b7-b629-4263295db6e8",
			Ip:          "1.1.1.1",
			Name:        "a",
			Description: "a description",
			Project:     "ce19a655-7933-4745-8f3e-9592b4a90488",
			Type:        apiv2.IPType_IP_TYPE_STATIC,
			Meta: &apiv2.Meta{
				Labels: &apiv2.Labels{
					Labels: map[string]string{
						tag.ClusterServiceFQN: "<cluster>/default/ingress-nginx",
					},
				},
			},
		}
	}
	ip2 = func() *apiv2.IP {
		return &apiv2.IP{
			Uuid:        "9cef40ec-29c6-4dfa-aee8-47ee1f49223d",
			Ip:          "4.3.2.1",
			Name:        "b",
			Description: "b description",
			Project:     "46bdfc45-9c8d-4268-b359-b40e3079d384",
			Type:        apiv2.IPType_IP_TYPE_EPHEMERAL,
			Meta: &apiv2.Meta{
				Labels: &apiv2.Labels{
					Labels: map[string]string{
						"a": "b",
					},
				},
			},
		}
	}
)

func Test_IPCmd_List(t *testing.T) {
	tests := []*e2e.Test[apiv2.IPServiceListRequest, apiv2.IPServiceListResponse]{
		{
			Name: "list",
			Cmd: func() []string {
				return []string{"ip", "list", "--project", "a"}
			},
			WantRequest: apiv2.IPServiceListRequest{
				Project: "a",
			},
			WantResponse: apiv2.IPServiceListResponse{
				Ips: []*apiv2.IP{
					ip1(),
					ip2(),
				},
			},
			WantTable: new(`
IP       PROJECT                               ID                                    TYPE       NAME  ATTACHED SERVICE
4.3.2.1  46bdfc45-9c8d-4268-b359-b40e3079d384  9cef40ec-29c6-4dfa-aee8-47ee1f49223d  ephemeral  b
1.1.1.1  ce19a655-7933-4745-8f3e-9592b4a90488  2e0144a2-09ef-42b7-b629-4263295db6e8  static     a
`),
			WantWideTable: new(`
IP       PROJECT                               ID                                    TYPE       NAME  DESCRIPTION    LABELS
4.3.2.1  46bdfc45-9c8d-4268-b359-b40e3079d384  9cef40ec-29c6-4dfa-aee8-47ee1f49223d  ephemeral  b     b description  a=b
1.1.1.1  ce19a655-7933-4745-8f3e-9592b4a90488  2e0144a2-09ef-42b7-b629-4263295db6e8  static     a     a description  cluster.metal-stack.io/id/namespace/service=<cluster>/default/ingress-nginx
`),
			Template: new("{{ .ip }} {{ .project }}"),
			WantTemplate: new(`
4.3.2.1 46bdfc45-9c8d-4268-b359-b40e3079d384
1.1.1.1 ce19a655-7933-4745-8f3e-9592b4a90488
			`),
			WantMarkdown: new(`
| IP      | PROJECT                              | ID                                   | TYPE      | NAME | ATTACHED SERVICE |
|---------|--------------------------------------|--------------------------------------|-----------|------|------------------|
| 4.3.2.1 | 46bdfc45-9c8d-4268-b359-b40e3079d384 | 9cef40ec-29c6-4dfa-aee8-47ee1f49223d | ephemeral | b    |                  |
| 1.1.1.1 | ce19a655-7933-4745-8f3e-9592b4a90488 | 2e0144a2-09ef-42b7-b629-4263295db6e8 | static    | a    |                  |
`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_IPCmd_Describe(t *testing.T) {
	ip1 := ip1()

	tests := []*e2e.Test[apiv2.IPServiceGetRequest, apiv2.IPServiceGetResponse]{
		{
			Name: "describe",
			Cmd: func() []string {
				return []string{"ip", "describe", "--project", ip1.Project, ip1.Ip}
			},
			WantRequest: apiv2.IPServiceGetRequest{
				Ip:      ip1.Ip,
				Project: ip1.Project,
			},
			WantResponse: apiv2.IPServiceGetResponse{
				Ip: ip1,
			},
			WantObject: ip1,
			WantTable: new(`
IP       PROJECT                               ID                                    TYPE    NAME  ATTACHED SERVICE
1.1.1.1  ce19a655-7933-4745-8f3e-9592b4a90488  2e0144a2-09ef-42b7-b629-4263295db6e8  static  a
`),
			WantWideTable: new(`
IP       PROJECT                               ID                                    TYPE    NAME  DESCRIPTION    LABELS
1.1.1.1  ce19a655-7933-4745-8f3e-9592b4a90488  2e0144a2-09ef-42b7-b629-4263295db6e8  static  a     a description  cluster.metal-stack.io/id/namespace/service=<cluster>/default/ingress-nginx
`),
			Template: new("{{ .ip }} {{ .project }}"),
			WantTemplate: new(`
1.1.1.1 ce19a655-7933-4745-8f3e-9592b4a90488
			`),
			WantMarkdown: new(`
| IP      | PROJECT                              | ID                                   | TYPE   | NAME | ATTACHED SERVICE |
|---------|--------------------------------------|--------------------------------------|--------|------|------------------|
| 1.1.1.1 | ce19a655-7933-4745-8f3e-9592b4a90488 | 2e0144a2-09ef-42b7-b629-4263295db6e8 | static | a    |                  |
`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

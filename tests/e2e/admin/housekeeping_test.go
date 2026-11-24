package admin_e2e

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/metal-stack/api/go/client"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	e2erootcmd "github.com/metal-stack/cli/testing/e2e"
	"github.com/metal-stack/cli/tests/e2e/testresources"
	"github.com/metal-stack/metal-lib/pkg/genericcli/e2e"
)

func Test_AdminHousekeepingCmd_MachineLiveliness(t *testing.T) {
	tests := []*e2e.Test[adminv2.HousekeepingServiceMachineLivelinessResponse, *adminv2.HousekeepingServiceMachineLivelinessResponse]{
		{
			Name:    "machine-liveliness",
			CmdArgs: []string{"admin", "housekeeping", "machine-liveliness"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.HousekeepingServiceMachineLivelinessRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.HousekeepingServiceMachineLivelinessResponse{
								ModifiedMachines: []*adminv2.MachineWithLiveliness{
									testresources.MachineWithLiveliness1,
									testresources.MachineWithLiveliness2,
								},
							})
						},
					},
				},
			}),
			WantProtoObject: &adminv2.HousekeepingServiceMachineLivelinessResponse{
				ModifiedMachines: []*adminv2.MachineWithLiveliness{
					testresources.MachineWithLiveliness1,
					testresources.MachineWithLiveliness2,
				},
			},
			Template: new(`{{ range .modified_machines }}{{ .uuid }} {{ .liveliness }}{{ "\n" }}{{ end }}`),
			WantTemplate: new(`
5fa2bbe1-407c-4142-92d5-e4419daf9646 1
673fc473-63ca-4ea4-b9dd-b45cb2127a6f 2
			`),
			WantDefault: new(`
modifiedMachines:
    - uuid: 5fa2bbe1-407c-4142-92d5-e4419daf9646
      liveliness: MACHINE_LIVELINESS_ALIVE
    - uuid: 673fc473-63ca-4ea4-b9dd-b45cb2127a6f
      liveliness: MACHINE_LIVELINESS_DEAD
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminHousekeepingCmd_MachineVPNConnected(t *testing.T) {
	tests := []*e2e.Test[adminv2.HousekeepingServiceMachineVPNConnectedResponse, apiv2.VPNNode]{
		{
			Name:    "machine-vpn-connected",
			CmdArgs: []string{"admin", "housekeeping", "machine-vpn-connected"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.HousekeepingServiceMachineVPNConnectedRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.HousekeepingServiceMachineVPNConnectedResponse{
								ModifiedMachines: []*apiv2.VPNNode{
									testresources.VpnNode2(),
									testresources.VpnNode1(),
								},
							})
						},
					},
				},
			}),
			Template: new("{{ .id }} {{ .project }}"),
			WantTemplate: new(`
2 f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c
1 0d81bca7-73f6-4da3-8397-4a8c52a0c583
			`),
			WantTable: new(`
            ID  NAME    PROJECT                               IPS              LAST SEEN
            2   node-2  f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  2.3.4.5          2000-01-01 00:01:00
            1   node-1  0d81bca7-73f6-4da3-8397-4a8c52a0c583  1.2.3.4,1.2.3.5  2000-01-01 00:00:00
			`),
			WantWideTable: new(`
            ID  NAME    PROJECT                               IPS              LAST SEEN
            2   node-2  f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  2.3.4.5          2000-01-01 00:01:00
            1   node-1  0d81bca7-73f6-4da3-8397-4a8c52a0c583  1.2.3.4,1.2.3.5  2000-01-01 00:00:00
			`),
			WantMarkdown: new(`
            | ID | NAME   | PROJECT                              | IPS             | LAST SEEN           |
            |----|--------|--------------------------------------|-----------------|---------------------|
            | 2  | node-2 | f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c | 2.3.4.5         | 2000-01-01 00:01:00 |
            | 1  | node-1 | 0d81bca7-73f6-4da3-8397-4a8c52a0c583 | 1.2.3.4,1.2.3.5 | 2000-01-01 00:00:00 |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminHousekeepingCmd_NetworksConsistent(t *testing.T) {
	response := &adminv2.HousekeepingServiceNetworksConsistentResponse{
		NetworksMissingInIpam:  []*apiv2.Network{testresources.Network1()},
		NetworksMissingInMetal: []*apiv2.Network{testresources.Network2()},
		IpsMissingInIpam:       []*apiv2.IP{testresources.IP1()},
		IpsMissingInMetal:      []*apiv2.IP{testresources.IP2()},
	}

	tests := []*e2e.Test[adminv2.HousekeepingServiceNetworksConsistentResponse, *adminv2.HousekeepingServiceNetworksConsistentResponse]{
		{
			Name:    "networks-consistent",
			CmdArgs: []string{"admin", "housekeeping", "networks-consistent"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.HousekeepingServiceNetworksConsistentRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(response)
						},
					},
				},
			}),
			WantDefault: new(`
networksMissingInIpam:
    - id: 6988ebb0-9531-4f9b-a893-d7868258e2ef
      meta:
        labels:
            labels:
                cluster.metal-stack.io/id/namespace/service: <cluster>/default/ingress-nginx
        createdAt: "2000-01-01T00:00:00Z"
      name: internet
      description: internet network
      partition: partition-1
      project: 0d81bca7-73f6-4da3-8397-4a8c52a0c583
      prefixes:
        - 10.0.0.0/16
        - 2001:db8::/32
      type: NETWORK_TYPE_EXTERNAL
      natType: NAT_TYPE_NONE
networksMissingInMetal:
    - id: d83ffb0a-7aa6-4a66-8e03-0b5ee8b718a0
      meta:
        labels:
            labels:
                a: b
        createdAt: "2000-01-01T00:00:00Z"
      name: private
      description: private network
      partition: partition-1
      project: f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c
      prefixes:
        - 192.168.1.0/24
      type: NETWORK_TYPE_CHILD
      natType: NAT_TYPE_IPV4_MASQUERADE
      vrf: 100
      parentNetwork: 6988ebb0-9531-4f9b-a893-d7868258e2ef
ipsMissingInIpam:
    - uuid: 2e0144a2-09ef-42b7-b629-4263295db6e8
      meta:
        labels:
            labels:
                cluster.metal-stack.io/id/namespace/service: <cluster>/default/ingress-nginx
      ip: 1.1.1.1
      name: aa
      description: a description
      network: internet
      project: ce19a655-7933-4745-8f3e-9592b4a90488
      type: IP_TYPE_STATIC
ipsMissingInMetal:
    - uuid: 9cef40ec-29c6-4dfa-aee8-47ee1f49223d
      meta:
        labels:
            labels:
                a: b
      ip: 4.3.2.1
      name: bb
      description: b description
      network: internet
      project: 46bdfc45-9c8d-4268-b359-b40e3079d384
      type: IP_TYPE_EPHEMERAL
run the following command inside of a pod of the ipam deployment to get the list of existing prefixes:
kubectl exec -it ipam-xxx-yyy -- /cli prefix list

then select one of the networks listed in networksMissingInMetal and delete it in ipam with
kubectl exec -it ipam-xxx-yyy -- /cli prefix delete --cidr 1.2.3.0/24 [--namespace optional-namespace]


run the following command inside of a pod of the ipam deployment to get the list of existing ips:
kubectl exec -it ipam-xxx-yyy -- /cli ip list --prefix prefix-of-the-given-metal-network

then select one of the ips listed in ipsMissingInMetal and delete it in ipam with
kubectl exec -it ipam-xxx-yyy -- /cli ip release --ip 1.2.3.4 --prefix 1.2.3.0/24 [--namespace optional-namespace]

			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

func Test_AdminHousekeepingCmd_MachineResurrect(t *testing.T) {
	tests := []*e2e.Test[adminv2.HousekeepingServiceMachineResurrectResponse, apiv2.Machine]{
		{
			Name:    "machine-resurrect",
			CmdArgs: []string{"admin", "housekeeping", "machine-resurrect"},
			NewRootCmd: e2erootcmd.NewRootCmd(t, &e2erootcmd.TestConfig{
				ClientCalls: []client.ClientCall{
					{
						WantRequest: &adminv2.HousekeepingServiceMachineResurrectRequest{},
						WantResponse: func() connect.AnyResponse {
							return connect.NewResponse(&adminv2.HousekeepingServiceMachineResurrectResponse{
								ModifiedMachines: []*apiv2.Machine{
									testresources.Machine2(),
									testresources.Machine1(),
								},
							})
						},
					},
				},
			}),
			Template: new("{{ .allocation.uuid }} {{ .allocation.project }}"),
			WantTemplate: new(`
4f94e87b-b08f-4f82-b053-9b8305de60ad f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c
<no value> <no value>
			`),
			WantTable: new(`
            ID                                       LAST EVENT   WHEN  AGE  HOSTNAME   PROJECT                               SIZE           IMAGE         PARTITION    RACK
            673fc473-63ca-4ea4-b9dd-b45cb2127a6f  🛡  Phoned Home  1m    1m   machine-2  f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  v1-medium-x86  Ubuntu 24.04  partition-2  rack-1
            5fa2bbe1-407c-4142-92d5-e4419daf9646     Alive        1m                                                          v1-medium-x86                partition-1  rack-1
			`),
			WantWideTable: new(`
            ID                                    LAST EVENT   WHEN  AGE  DESCRIPTION  NAME       HOSTNAME   PROJECT                               IPS        SIZE           IMAGE         PARTITION    RACK    STARTED               TAGS  STATE
            673fc473-63ca-4ea4-b9dd-b45cb2127a6f  Phoned Home  1m    1m   machine 2    machine-2  machine-2  f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c  4.5.6.7    v1-medium-x86  Ubuntu 24.04  partition-2  rack-1  1999-12-31T23:59:00Z  c=d   available
                                                                                                                                                    192.1.1.1
            5fa2bbe1-407c-4142-92d5-e4419daf9646  Alive        1m                                                                                             v1-medium-x86                partition-1  rack-1                        a=b   available
			`),
			WantMarkdown: new(`
            | ID                                   |   | LAST EVENT  | WHEN | AGE | HOSTNAME  | PROJECT                              | SIZE          | IMAGE        | PARTITION   | RACK   |
            |--------------------------------------|---|-------------|------|-----|-----------|--------------------------------------|---------------|--------------|-------------|--------|
            | 673fc473-63ca-4ea4-b9dd-b45cb2127a6f | 🛡 | Phoned Home | 1m   | 1m  | machine-2 | f3b4e6a1-2c8d-4e5f-a7b9-1d3e5f7a9b0c | v1-medium-x86 | Ubuntu 24.04 | partition-2 | rack-1 |
            | 5fa2bbe1-407c-4142-92d5-e4419daf9646 |   | Alive       | 1m   |     |           |                                      | v1-medium-x86 |              | partition-1 | rack-1 |
			`),
		},
	}
	for _, tt := range tests {
		tt.TestCmd(t)
	}
}

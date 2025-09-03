// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package vms

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/vmware/govmomi/cli"
	"github.com/vmware/govmomi/cli/flags"
	"github.com/vmware/govmomi/vapi/cis/tasks"
	"github.com/vmware/govmomi/vapi/esx/settings/clusters/vms"
	"github.com/vmware/govmomi/vim25/types"
)

type set struct {
	*flags.ClientFlag
	*flags.ClusterFlag
	*flags.ResourcePoolFlag
	*flags.NetworkFlag

	vms.SolutionSpec
}

func init() {
	cli.Register("esx.settings.clusters.vms.solutions.set", &set{})
}

func (cmd *set) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.ClientFlag, ctx = flags.NewClientFlag(ctx)
	cmd.ClientFlag.Register(ctx, f)

	cmd.ClusterFlag, ctx = flags.NewClusterFlag(ctx)
	cmd.ClusterFlag.Register(ctx, f)

	cmd.ResourcePoolFlag, ctx = flags.NewResourcePoolFlag(ctx)
	cmd.ResourcePoolFlag.Register(ctx, f)

	cmd.NetworkFlag, ctx = flags.NewNetworkFlag(ctx)
	cmd.NetworkFlag.Register(ctx, f)
}

func (cmd *set) Usage() string {
	return "CLUSTER"
}

func (cmd *set) Description() string {
	return `set CLUSTER in datacenter.

Lorem ipsum

Examples:
  govc esx.settings.clusters.vms.solutions ClusterA`
}

func (cmd *set) Process(ctx context.Context) error {
	if err := cmd.ClientFlag.Process(ctx); err != nil {
		return err
	}

	if err := cmd.ClusterFlag.Process(ctx); err != nil {
		return err
	}

	if err := cmd.ResourcePoolFlag.Process(ctx); err != nil {
		return err
	}

	if err := cmd.NetworkFlag.Process(ctx); err != nil {
		return err
	}

	return nil
}

func (cmd *set) Run(ctx context.Context, f *flag.FlagSet) error {
	if f.NArg() != 1 {
		return flag.ErrHelp
	}

	network, err := cmd.Network()
	if err != nil {
		return err
	}

	netInfo, err := network.EthernetCardBackingInfo(ctx)
	if err != nil {
		return err
	}

	seven := int32(7)

	netDevice := &types.VirtualVmxnet3{
		VirtualVmxnet: types.VirtualVmxnet{
			VirtualEthernetCard: types.VirtualEthernetCard{
				//							AddressType: "Generated",

				AddressType: string(types.VirtualEthernetCardMacTypeGenerated),
				VirtualDevice: types.VirtualDevice{
					Key:        0,
					UnitNumber: &seven,
					Backing:    netInfo,
					DeviceInfo: &types.Description{
						Label:   "Network 1",
						Summary: "VM Network",
					},
				},
			},
		},
	}

	_ = netDevice

	deviceCfg := &types.VirtualDeviceConfigSpec{
		Operation:     types.VirtualDeviceConfigSpecOperationAdd,
		FileOperation: types.VirtualDeviceConfigSpecFileOperationCreate,
		Device:        netDevice,
	}

	cfg := &types.VirtualMachineConfigSpec{
		DeviceChange: []types.BaseVirtualDeviceConfigSpec{
			deviceCfg,
		},
	}

	var buf bytes.Buffer
	types.NewJSONEncoder(&buf).Encode(cfg)
	if err != nil {
		return err
	}

	s := &vms.SolutionSpec{
		DeploymentType: vms.ClusterVmSet,
		DisplayName:    "food",
		DisplayVersion: "1",
		VmNameTemplate: vms.VmNameTemplate{
			Prefix: "templatename",
			Suffix: vms.Counter,
		},
		ClusterSolutionSpec: vms.ClusterSolutionSpec{
			VmCount:             1,
			VmPlacementPolicies: []vms.VmPlacementPolicy{vms.VmVmAntiAffinity},
			RemediationPolicy:   vms.Sequential,
			//Devices:             cfg,
			Devices: buf.Bytes(),
		},
		VmCloneConfig:      vms.NoClones,
		HookConfigurations: map[vms.LifecycleState]vms.LifecycleHookConfig{vms.PostProvisioning: {}},
		OvfResource: vms.OvfResource{
			LocationType:             vms.RemoteFile,
			Url:                      "https://lvn-dvm-10-161-123-247.dvm.lvn.broadcom.net:5480/wcpagent/photon-ova.ovf",
			AuthenticationScheme:     vms.None,
			SslCertificateValidation: vms.SslCertificateValidationDisabled,
		},
		VmResourceSpec: vms.VmResourceSpec{
			OvfDeploymentOption: "small",
		},
		VmStoragePolicy:    vms.Profile,
		VmStorageProfiles:  []string{"dbcb3f01-0930-4899-aa76-bf4c05476c34"},
		VmDiskType:         vms.DiskTypeThick,
		VmResourcePool:     "resgroup-1043",
		VmFolder:           "group-v1044",
		RedeploymentPolicy: vms.BlueGreen,
	}

	b, err := json.Marshal(s)
	fmt.Println(string(b))
	//os.Exit(0)

	rc, err := cmd.RestClient()
	if err != nil {
		return err
	}

	cluster, err := cmd.Cluster()
	if err != nil {
		return err
	}

	fmt.Println("setting solution")
	solutionId := "solution1"
	m := vms.Manager{rc}
	if err := m.Set(ctx, cluster.Reference(), solutionId, s); err != nil {
		return err
	}

	fmt.Println("applying solution")
	spec := &vms.ApplySpec{
		ClusterSolutions: &vms.ClusterSolutionsApplyFilterSpec{
			Solutions: []string{solutionId},
		},
	}

	taskId, err := m.Apply(ctx, cluster.Reference(), spec)
	if err != nil {
		return err
	}

	fmt.Println("waiting")

	if _, err = tasks.NewManager(rc).WaitForRunningOrError(ctx, taskId); err != nil {
		return err
	}

	info, err := m.Get(ctx, cluster.Reference(), solutionId)
	spew.Dump(info)
	spew.Dump(err)

	info, err = m.Get(ctx, cluster.Reference(), "solutionId")
	spew.Dump(info)
	spew.Dump(err)

	//m.List(ctx, cluster.Reference())
	//m.List(ctx, types.ManagedObjectReference{Value: "domain-c52"})

	fmt.Println("checking compliance")
	compliances, err := m.CheckCompliance(ctx, cluster.Reference(), solutionId, &vms.CheckComplianceFilterSpec{Solutions: []string{solutionId}})
	if err != nil {
		return err
	}

	spew.Dump(compliances)

	return nil
}

/*
 * Copyright (c) 2019 SUSE LLC.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package cluster

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/SUSE/skuba/pkg/skuba"
	cluster "github.com/SUSE/skuba/pkg/skuba/actions/cluster/init"
)

type initOptions struct {
	ControlPlane      string
	KubernetesVersion string
	CloudProvider     string
	StrictCapDefaults bool
	CniPlugin         string
}

// NewInitCmd creates a new `skuba cluster init` cobra command
func NewInitCmd() *cobra.Command {
	initOptions := initOptions{}

	cmd := &cobra.Command{
		Use:   "init <cluster-name> --control-plane <IP/FQDN>",
		Short: "Initialize skuba structure for cluster deployment",
		Run: func(cmd *cobra.Command, args []string) {
			initConfig, err := cluster.NewInitConfiguration(
				args[0],
				initOptions.CloudProvider,
				initOptions.ControlPlane,
				initOptions.KubernetesVersion,
				initOptions.StrictCapDefaults,
				initOptions.CniPlugin)
			if err != nil {
				klog.Fatalf("init failed due to error: %s", err)
			}

			if err = cluster.Init(initConfig); err != nil {
				klog.Fatalf("init failed due to error: %s", err)
			}
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("accepts 1 arg, received %d", len(args))
			}
			if len(args[0]) == 0 {
				return fmt.Errorf("Invalid cluster name: expected non empty cluster name")
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&initOptions.ControlPlane, "control-plane", "", "The control plane location (IP/FQDN) that will load balance the master nodes (required)")
	if skuba.BuildType == "development" {
		cmd.Flags().StringVar(&initOptions.KubernetesVersion, "kubernetes-version", "", "The kubernetes version to bootstrap with (only in development build)")
	}
	cmd.Flags().StringVar(&initOptions.CloudProvider, "cloud-provider", "", "Enable cloud provider integration with the chosen cloud. Valid values: aws, azure, openstack, vsphere")
	_ = cmd.MarkFlagRequired("control-plane")

	cmd.Flags().BoolVar(&initOptions.StrictCapDefaults, "strict-capability-defaults", false, "All the containers will start with CRI-O default capabilities")

	cmd.Flags().StringVar(&initOptions.CniPlugin, "cni-plugin", "cilium", "Specify the CNI plugin to be used across the cluster. Valid values: cilium")

	return cmd
}

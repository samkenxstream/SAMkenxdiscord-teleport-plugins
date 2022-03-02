/*
Copyright 2015-2021 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package test

import (
	"time"

	"github.com/gravitational/teleport/api/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

func (s *TerraformSuite) TestClusterNetworkingConfig() {
	name := "teleport_cluster_networking_config.test"

	resource.Test(s.T(), resource.TestCase{
		ProtoV6ProviderFactories:  s.terraformProviders,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: s.getFixture("networking_config_0_set.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "kind", "cluster_networking_config"),
					resource.TestCheckResourceAttr(name, "metadata.labels.example", "yes"),
					resource.TestCheckResourceAttr(name, "spec.client_idle_timeout", "30m"),
				),
			},
			{
				Config:   s.getFixture("networking_config_0_set.tf"),
				PlanOnly: true,
			},
			{
				Config: s.getFixture("networking_config_1_update.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "kind", "cluster_networking_config"),
					resource.TestCheckResourceAttr(name, "metadata.labels.example", "no"),
					resource.TestCheckResourceAttr(name, "spec.client_idle_timeout", "1h"),
				),
			},
			{
				Config:   s.getFixture("networking_config_1_update.tf"),
				PlanOnly: true,
			},
		},
	})
}

func (s *TerraformSuite) TestImportClusterNetworkingConfig() {
	name := "teleport_cluster_networking_config.test"

	clusterNetworkingConfig := &types.ClusterNetworkingConfigV2{
		Metadata: types.Metadata{},
		Spec: types.ClusterNetworkingConfigSpecV2{
			ClientIdleTimeout: types.Duration(30 * time.Second),
		},
	}
	err := clusterNetworkingConfig.CheckAndSetDefaults()
	require.NoError(s.T(), err)

	err = s.client.SetClusterNetworkingConfig(s.Context(), clusterNetworkingConfig)
	require.NoError(s.T(), err)

	resource.Test(s.T(), resource.TestCase{
		ProtoV6ProviderFactories: s.terraformProviders,
		Steps: []resource.TestStep{
			{
				Config:        s.terraformConfig + "\n" + `resource "teleport_cluster_networking_config" "test" { }`,
				ResourceName:  name,
				ImportState:   true,
				ImportStateId: "test",
				ImportStateCheck: func(state []*terraform.InstanceState) error {
					require.Equal(s.T(), state[0].Attributes["kind"], "cluster_networking_config")
					require.Equal(s.T(), state[0].Attributes["spec.client_idle_timeout"], "30s")

					return nil
				},
			},
		},
	})
}

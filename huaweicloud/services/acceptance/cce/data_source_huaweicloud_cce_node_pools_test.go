package cce

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataNodePools_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		all = "data.huaweicloud_cce_node_pools.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byShowDefault   = "data.huaweicloud_cce_node_pools.with_default"
		dcByShowDefault = acceptance.InitDataSourceCheck(byShowDefault)

		byWithoutDefault   = "data.huaweicloud_cce_node_pools.without_default"
		dcByWithoutDefault = acceptance.InitDataSourceCheck(byWithoutDefault)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataNodePools_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "node_pools.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "node_pools.0.id"),
					resource.TestCheckResourceAttrPair(all, "node_pools.0.id",
						"huaweicloud_cce_node_pool.test", "id"),
					resource.TestCheckResourceAttrPair(all, "node_pools.0.name",
						"huaweicloud_cce_node_pool.test", "name"),
					resource.TestCheckResourceAttrPair(all, "node_pools.0.initial_node_count",
						"huaweicloud_cce_node_pool.test", "initial_node_count"),
					resource.TestCheckResourceAttrPair(all, "node_pools.0.flavor_id",
						"huaweicloud_cce_node_pool.test", "flavor_id"),
					resource.TestCheckResourceAttrPair(all, "node_pools.0.type",
						"huaweicloud_cce_node_pool.test", "type"),
					resource.TestCheckResourceAttrPair(all, "node_pools.0.availability_zone",
						"huaweicloud_cce_node_pool.test", "availability_zone"),
					resource.TestCheckResourceAttrPair(all, "node_pools.0.os",
						"huaweicloud_cce_node_pool.test", "os"),
					resource.TestCheckResourceAttrPair(all, "node_pools.0.key_pair",
						"huaweicloud_cce_node_pool.test", "key_pair"),
					resource.TestCheckResourceAttrPair(all, "node_pools.0.subnet_id",
						"huaweicloud_cce_node_pool.test", "subnet_id"),
					resource.TestCheckResourceAttrPair(all, "node_pools.0.runtime",
						"huaweicloud_cce_node_pool.test", "runtime"),
					resource.TestCheckResourceAttrPair(all, "node_pools.0.labels.test_key",
						"huaweicloud_cce_node_pool.test", "labels.test_key"),
					resource.TestCheckResourceAttrPair(all, "node_pools.0.tags.foo",
						"huaweicloud_cce_node_pool.test", "tags.foo"),
					resource.TestCheckResourceAttrPair(all, "node_pools.0.taints.0.key",
						"huaweicloud_cce_node_pool.test", "taints.0.key"),
					resource.TestCheckResourceAttrPair(all, "node_pools.0.taints.0.value",
						"huaweicloud_cce_node_pool.test", "taints.0.value"),
					resource.TestCheckResourceAttrPair(all, "node_pools.0.taints.0.effect",
						"huaweicloud_cce_node_pool.test", "taints.0.effect"),
					resource.TestCheckResourceAttrPair(all, "node_pools.0.extend_params.0.max_pods",
						"huaweicloud_cce_node_pool.test", "extend_params.0.max_pods"),
					resource.TestCheckResourceAttrPair(all, "node_pools.0.extend_params.0.docker_base_size",
						"huaweicloud_cce_node_pool.test", "extend_params.0.docker_base_size"),
					resource.TestCheckResourceAttrPair(all, "node_pools.0.hostname_config.0.type",
						"huaweicloud_cce_node_pool.test", "hostname_config.0.type"),
					resource.TestCheckResourceAttrPair(all, "node_pools.0.root_volume.0.size",
						"huaweicloud_cce_node_pool.test", "root_volume.0.size"),
					resource.TestCheckResourceAttrPair(all, "node_pools.0.root_volume.0.volumetype",
						"huaweicloud_cce_node_pool.test", "root_volume.0.volumetype"),
					resource.TestCheckResourceAttrPair(all, "node_pools.0.data_volumes.0.size",
						"huaweicloud_cce_node_pool.test", "data_volumes.0.size"),
					resource.TestCheckResourceAttrPair(all, "node_pools.0.data_volumes.0.volumetype",
						"huaweicloud_cce_node_pool.test", "data_volumes.0.volumetype"),
					resource.TestCheckResourceAttrPair(all, "node_pools.0.extension_scale_groups.0.metadata.0.name",
						"huaweicloud_cce_node_pool.test", "extension_scale_groups.0.metadata.0.name"),
					resource.TestCheckResourceAttrPair(all, "node_pools.0.extension_scale_groups.0.spec.0.flavor",
						"huaweicloud_cce_node_pool.test", "extension_scale_groups.0.spec.0.flavor"),
					resource.TestCheckResourceAttrPair(all, "node_pools.0.extension_scale_groups.0.spec.0.az",
						"huaweicloud_cce_node_pool.test", "extension_scale_groups.0.spec.0.az"),
					resource.TestCheckResourceAttrSet(all, "node_pools.0.storage.#"),
					resource.TestCheckResourceAttrSet(all, "node_pools.0.storage.0.selectors.#"),
					resource.TestCheckResourceAttrSet(all, "node_pools.0.storage.0.groups.#"),
					resource.TestCheckOutput("is_storage_set", "true"),

					// filter by show_default_node_pool
					dcByShowDefault.CheckResourceExists(),
					dcByWithoutDefault.CheckResourceExists(),
					resource.TestCheckOutput("is_show_default_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataNodePools_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[2]s"
}

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "Huawei Cloud EulerOS 2.0"
  flavor_id                = data.huaweicloud_compute_flavors.test.ids[0]
  initial_node_count       = 0
  availability_zone        = data.huaweicloud_availability_zones.test.names[0]
  key_pair                 = huaweicloud_kps_keypair.test.name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"
  runtime                  = "containerd"
  subnet_id                = huaweicloud_vpc_subnet.test.id
  security_groups          = [huaweicloud_networking_secgroup.test.id]
  initialized_conditions   = ["CCEInitial"]

  tag_policy_on_existing_nodes   = "ignore"
  label_policy_on_existing_nodes = "ignore"
  taint_policy_on_existing_nodes = "ignore"

  labels = {
    test_key = "test_value"
  }

  tags = {
    foo = "bar"
  }

  taints {
    key    = "test_key"
    value  = "test_value"
    effect = "NoSchedule"
  }

  extend_params {
    max_pods         = 50
    docker_base_size = 20
  }

  hostname_config {
    type = "cceNodeName"
  }

  root_volume {
    size       = 40
    volumetype = "SSD"
  }

  data_volumes {
    size       = 100
    volumetype = "SSD"
  }

  storage {
    selectors {
      name              = "cceUse"
      type              = "evs"
      match_label_size  = "100"
      match_label_count = "1"
    }

    groups {
      name           = "vgpaas"
      selector_names = ["cceUse"]
      cce_managed    = true

      virtual_spaces {
        name        = "kubernetes"
        size        = "10%%"
        lvm_lv_type = "linear"
      }

      virtual_spaces {
        name            = "runtime"
        size            = "90%%"
        runtime_lv_type = "linear"
      }
    }
  }

  extension_scale_groups {
    metadata {
      name = "%[2]s-group1"
    }

    spec {
      flavor = data.huaweicloud_compute_flavors.test.ids[0]
      az     = data.huaweicloud_availability_zones.test.names[1]

      autoscaling {
        extension_priority = 1
        enable             = true
      }
    }
  }
}

data "huaweicloud_cce_node_pools" "test" {
  depends_on = [
    huaweicloud_cce_node_pool.test,
  ]

  cluster_id = huaweicloud_cce_cluster.test.id
}

# Filter by show_default_node_pool=true
data "huaweicloud_cce_node_pools" "with_default" {
  depends_on = [
    huaweicloud_cce_node_pool.test,
  ]

  cluster_id             = huaweicloud_cce_cluster.test.id
  show_default_node_pool = "true"
}

# Filter by show_default_node_pool=false
data "huaweicloud_cce_node_pools" "without_default" {
  depends_on = [
    huaweicloud_cce_node_pool.test,
  ]

  cluster_id             = huaweicloud_cce_cluster.test.id
  show_default_node_pool = "false"
}

output "is_show_default_filter_useful" {
  value = (
    length(data.huaweicloud_cce_node_pools.with_default.node_pools) >= length(data.huaweicloud_cce_node_pools.without_default.node_pools) &&
    length(data.huaweicloud_cce_node_pools.with_default.node_pools) > 0
  )
}

locals {
  target_pools = [
    for p in data.huaweicloud_cce_node_pools.test.node_pools : p
    if p.id == huaweicloud_cce_node_pool.test.id
  ]
  target_pool = length(local.target_pools) > 0 ? local.target_pools[0] : null
}

output "is_storage_set" {
  value = local.target_pool != null && length(local.target_pool.storage) > 0
}
`, testAccNodePool_base(name), name)
}

package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccCCENodePoolV3DataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "data.huaweicloud_cce_node_pool.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodePoolV3DataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodePoolV3DataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "extension_scale_groups.#"),
					resource.TestCheckResourceAttrSet(resourceName, "extension_scale_groups.0.metadata.#"),
					resource.TestCheckResourceAttrSet(resourceName, "extension_scale_groups.0.metadata.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "extension_scale_groups.0.metadata.0.uid"),
					resource.TestCheckResourceAttrSet(resourceName, "extension_scale_groups.0.spec.#"),
					resource.TestCheckResourceAttrSet(resourceName, "extension_scale_groups.0.spec.0.flavor"),
					resource.TestCheckResourceAttrSet(resourceName, "extension_scale_groups.0.spec.0.az"),
					resource.TestCheckResourceAttrSet(resourceName, "extension_scale_groups.0.spec.0.autoscaling.#"),
					resource.TestCheckResourceAttrSet(resourceName, "extension_scale_groups.0.spec.0.autoscaling.0.extension_priority"),
					resource.TestCheckResourceAttrSet(resourceName, "extension_scale_groups.0.spec.0.autoscaling.0.enable"),
				),
			},
		},
	})
}

func testAccCheckCCENodePoolV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find node pools data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("Node pool data source ID not set ")
		}

		return nil
	}
}

func testAccCCENodePoolV3DataSource_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.9"
  flavor_id                = data.huaweicloud_compute_flavors.test.ids[0]
  initial_node_count       = 1
  availability_zone        = data.huaweicloud_availability_zones.test.names[0]
  key_pair                 = huaweicloud_kps_keypair.test.name
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
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

data "huaweicloud_cce_node_pool" "test" {
  depends_on = [huaweicloud_cce_node_pool.test]

  cluster_id = huaweicloud_cce_cluster.test.id
  name       = huaweicloud_cce_node_pool.test.name
}
`, testAccNodePool_base(name), name)
}

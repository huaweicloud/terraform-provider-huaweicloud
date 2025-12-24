package cceautopilot

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccClusterUpgrade_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterUpgrade_basic(name),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testAccClusterUpgrade_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_autopilot_cluster" "test" {
  name        = "%[2]s"
  flavor      = "cce.autopilot.cluster"
  description = "created by terraform"
  version     = "v1.28"

  host_network {
    vpc    = huaweicloud_vpc.test.id
    subnet = huaweicloud_vpc_subnet.test.id
  }

  container_network {
    mode = "eni"
  }

  eni_network {
    subnets {
      subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id
    }
  }

  tags = {
    "foo" = "bar"
    "key" = "value"
  }
}
`, common.TestVpc(name), name)
}

func testAccClusterUpgrade_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_autopilot_cluster_upgrade" "test" {
  cluster_id     = huaweicloud_cce_autopilot_cluster.test.id
  target_version = "v1.31"

  strategy {
    type = "inPlaceRollingUpdate"

    in_place_rolling_update {
      user_defined_step = 20
    }
  }
}
`, testAccClusterUpgrade_base(name), name)
}

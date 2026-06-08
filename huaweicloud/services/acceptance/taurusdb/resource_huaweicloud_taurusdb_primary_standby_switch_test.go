package taurusdb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTaurusDBPrimaryStandbySwitch_basic(t *testing.T) {
	resourceName := "huaweicloud_taurusdb_primary_standby_switch.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBPrimaryStandbySwitch_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttrSet(resourceName, "node_id"),
				),
			},
		},
	})
}

func testAccTaurusDBPrimaryStandbySwitch_basic() string {
	return `
data "huaweicloud_taurusdb_instances" "test" {}

locals {
  instance_id = data.huaweicloud_taurusdb_instances.test.instances.0.id
  slave_node_ids = [for node in data.huaweicloud_taurusdb_instances.test.instances.0.nodes : node.id if node.type == "slave"]
  slave_node_id = length(local.slave_node_ids) > 0 ? local.slave_node_ids[0] : null
}

resource "huaweicloud_taurusdb_primary_standby_switch" "test" {
  instance_id = local.instance_id
  node_id     = local.slave_node_id
}
`
}

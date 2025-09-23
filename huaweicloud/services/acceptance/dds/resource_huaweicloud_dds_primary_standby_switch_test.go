package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDDSPrimaryStandbySwitch_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDDSPrimaryStandbySwitch_Instance(rName),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccDDSPrimaryStandbySwitch_Instance(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dds_primary_standby_switch" "test" {
  instance_id = huaweicloud_dds_instance.instance.id
}`, testAccDDSInstanceReplicaSetBasic(rName))
}

func TestAccDDSPrimaryStandbySwitch_node(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDDSPrimaryStandbySwitch_Node(rName),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccDDSPrimaryStandbySwitch_Node(rName string) string {
	return fmt.Sprintf(`
%s

locals {
  node_ids = [for node in huaweicloud_dds_instance.instance.nodes: node.id if node.role == "Secondary"]
}

resource "huaweicloud_dds_primary_standby_switch" "test" {
  instance_id = huaweicloud_dds_instance.instance.id
  node_id     = local.node_ids[0]

  lifecycle {
    ignore_changes = [
      node_id,
    ]
  }
}`, testAccDDSInstanceV3Config_basic(rName, 8800))
}

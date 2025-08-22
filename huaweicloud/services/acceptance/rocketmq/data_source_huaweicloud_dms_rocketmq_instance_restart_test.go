package rocketmq

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDmsRocketMQInstanceRestart_basic(t *testing.T) {
	var (
		dataSourceName = "huaweicloud_dms_rocketmq_instance_restart.test"
		instanceID     = acceptance.HW_DMS_ROCKETMQ_INSTANCE_ID
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSRocketMQInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRocketMQInstanceRestart_basic(instanceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "instance_id", instanceID),
					resource.TestCheckResourceAttr(dataSourceName, "nodes.#", "2"),
					resource.TestCheckResourceAttrSet(dataSourceName, "nodes.0"),
				),
			},
		},
	})
}

func testAccDmsRocketMQInstanceRestart_basic(instanceID string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_rocketmq_instance_restart" "test" {
  instance_id = data.huaweicloud_dms_rocketmq_instance_nodes.test.instance_id

  nodes = [
    try(data.huaweicloud_dms_rocketmq_instance_nodes.test.nodes[0].id),
    try(data.huaweicloud_dms_rocketmq_instance_nodes.test.nodes[1].id)
  ]
}
`, testAccDmsRocketMQInstanceNodes_basic(instanceID))
}

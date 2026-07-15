package rocketmq

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccNodeBatchRestart_basic(t *testing.T) {
	var (
		dataSourceName = "huaweicloud_dms_rocketmq_node_batch_restart.test"
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
				Config: testAccNodeBatchRestart_basic(instanceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "instance_id", instanceID),
					resource.TestCheckResourceAttr(dataSourceName, "nodes.#", "2"),
					resource.TestCheckResourceAttrSet(dataSourceName, "nodes.0"),
				),
			},
		},
	})
}

func testAccNodeBatchRestart_basic(instanceID string) string {
	return fmt.Sprintf(`
data "huaweicloud_dms_rocketmq_instance_nodes" "test" {
  instance_id = "%[1]s"
}

resource "huaweicloud_dms_rocketmq_node_batch_restart" "test" {
  instance_id = "%[1]s"

  nodes = [
    try(data.huaweicloud_dms_rocketmq_instance_nodes.test.nodes[0].id),
    try(data.huaweicloud_dms_rocketmq_instance_nodes.test.nodes[1].id)
  ]
}
`, instanceID)
}

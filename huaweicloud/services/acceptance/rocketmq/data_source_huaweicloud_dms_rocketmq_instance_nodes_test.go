package rocketmq

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDmsRocketMQInstanceNodes_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_dms_rocketmq_instance_nodes.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
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
				Config: testAccDmsRocketMQInstanceNodes_basic(instanceID),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "nodes.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "nodes.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "nodes.0.broker_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "nodes.0.broker_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "nodes.0.address"),
				),
			},
		},
	})
}

func testAccDmsRocketMQInstanceNodes_basic(instanceID string) string {
	return fmt.Sprintf(`
data "huaweicloud_dms_rocketmq_instance_nodes" "test" {
  instance_id = "%s"
}
`, instanceID)
}

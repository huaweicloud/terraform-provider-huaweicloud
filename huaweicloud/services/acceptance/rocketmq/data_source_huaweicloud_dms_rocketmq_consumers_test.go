package rocketmq

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDmsRocketmqConsumers_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dms_rocketmq_consumers.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSRocketMQInstanceID(t)
			acceptance.TestAccPreCheckDMSRocketMQGroupName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDmsRocketmqConsumers_basic(true),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "clients.#"),
					resource.TestCheckResourceAttrSet(dataSource, "clients.0.language"),
					resource.TestCheckResourceAttrSet(dataSource, "clients.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "clients.0.client_id"),
					resource.TestCheckResourceAttrSet(dataSource, "clients.0.client_address"),
					resource.TestCheckResourceAttrSet(dataSource, "clients.0.subscriptions.0.topic"),
					resource.TestCheckResourceAttrSet(dataSource, "clients.0.subscriptions.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "clients.0.subscriptions.0.expression"),
					resource.TestCheckResourceAttrSet(dataSource, "online"),
					resource.TestCheckResourceAttrSet(dataSource, "subscription_consistency"),
				),
			},
			{
				Config: testDataSourceDmsRocketmqConsumers_basic(false),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "clients.#"),
					resource.TestCheckResourceAttrSet(dataSource, "clients.0.language"),
					resource.TestCheckResourceAttrSet(dataSource, "clients.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "clients.0.client_id"),
					resource.TestCheckResourceAttrSet(dataSource, "clients.0.client_address"),
					resource.TestCheckResourceAttrSet(dataSource, "online"),
					resource.TestCheckResourceAttrSet(dataSource, "subscription_consistency"),
				),
			},
		},
	})
}

func testDataSourceDmsRocketmqConsumers_basic(detail bool) string {
	return fmt.Sprintf(`
data "huaweicloud_dms_rocketmq_consumers" "test" {
  instance_id = "%[1]s"
  group       = "%[2]s"
  is_detail   = %t
}
`, acceptance.HW_DMS_ROCKETMQ_INSTANCE_ID, acceptance.HW_DMS_ROCKETMQ_GROUP_NAME, detail)
}

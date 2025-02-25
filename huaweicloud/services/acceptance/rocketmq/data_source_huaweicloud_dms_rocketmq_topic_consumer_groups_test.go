package rocketmq

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDmsRocketmqTopicConsumerGroups_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dms_rocketmq_topic_consumer_groups.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSRocketMQInstanceID(t)
			acceptance.TestAccPreCheckDMSRocketMQTopicName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDmsRocketmqTopicConsumerGroups_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "groups.#"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDmsRocketmqTopicConsumerGroups_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dms_rocketmq_topic_consumer_groups" "test" {
  instance_id = "%s"
  topic_name  = "%s"
}
`, acceptance.HW_DMS_ROCKETMQ_INSTANCE_ID, acceptance.HW_DMS_ROCKETMQ_TOPIC_NAME)
}

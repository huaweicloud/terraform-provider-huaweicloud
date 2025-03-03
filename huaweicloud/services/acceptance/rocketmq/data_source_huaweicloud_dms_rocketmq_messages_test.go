package rocketmq

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDmsRocketMQMessages_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dms_rocketmq_messages.test"
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
				Config: testDataSourceDmsRocketMQMessages_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "messages.#"),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.message_id"),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.store_time"),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.born_time"),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.reconsume_times"),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.body_crc"),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.store_size"),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.born_host"),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.store_host"),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.queue_id"),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.queue_offset"),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.property_list.#"),

					resource.TestCheckOutput("msg_id_validation", "true"),
				),
			},
		},
	})
}

func testDataSourceDmsRocketMQMessages_basic() string {
	currentTime := time.Now()
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).UnixMilli()
	endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location()).UnixMilli()
	return fmt.Sprintf(`
data "huaweicloud_dms_rocketmq_messages" "test" {
  instance_id = "%[1]s"
  topic       = "%[2]s"
  start_time  = %[3]v
  end_time    = %[4]v
}

locals {
  msg_id = data.huaweicloud_dms_rocketmq_messages.test.messages.0.message_id
}

data "huaweicloud_dms_rocketmq_messages" "msg_id" {
  instance_id = "%[1]s"
  topic       = "%[2]s"
  message_id  = try(local.msg_id, "")
}

output "msg_id_validation" {
  value = alltrue([for v in data.huaweicloud_dms_rocketmq_messages.msg_id.messages[*].message_id : v == local.msg_id])
}
`, acceptance.HW_DMS_ROCKETMQ_INSTANCE_ID, acceptance.HW_DMS_ROCKETMQ_TOPIC_NAME, startTime, endTime)
}

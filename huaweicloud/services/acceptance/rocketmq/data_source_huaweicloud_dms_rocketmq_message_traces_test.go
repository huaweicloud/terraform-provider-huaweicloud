package rocketmq

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDmsRocketmqMessageTraces_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dms_rocketmq_message_traces.test"
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
				Config: testDataSourceDmsRocketmqMessageTraces_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "traces.#"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.message_id"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.body_length"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.client_host"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.consume_status"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.cost_time"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.from_transaction_check"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.group_name"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.keys"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.retry_times"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.store_host"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.success"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.time"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.trace_type"),
				),
			},
		},
	})
}

func testDataSourceDmsRocketmqMessageTraces_basic() string {
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

data "huaweicloud_dms_rocketmq_message_traces" "test" {
  instance_id = "%[1]s"
  message_id  = try(local.msg_id, "")
}
`, acceptance.HW_DMS_ROCKETMQ_INSTANCE_ID, acceptance.HW_DMS_ROCKETMQ_TOPIC_NAME, startTime, endTime)
}

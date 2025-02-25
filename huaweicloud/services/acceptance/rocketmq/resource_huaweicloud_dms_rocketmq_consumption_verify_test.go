package rocketmq

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRocketMQConsumptionVerify_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSRocketMQInstanceID(t)
			acceptance.TestAccPreCheckDMSRocketMQGroupName(t)
			acceptance.TestAccPreCheckDMSRocketMQTopicName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRocketMQConsumptionVerify_basic(),
			},
		},
	})
}

func testAccRocketMQConsumptionVerify_basic() string {
	currentTime := time.Now()
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).UnixMilli()
	endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location()).UnixMilli()
	return fmt.Sprintf(`
data "huaweicloud_dms_rocketmq_messages" "test" {
  instance_id = "%[1]s"
  topic       = "%[3]s"
  start_time  = %[4]v
  end_time    = %[5]v
}

data "huaweicloud_dms_rocketmq_consumers" "test" {
  instance_id = "%[1]s"
  group       = "%[2]s"
  is_detail   = true
}

locals {
  msg_id    = data.huaweicloud_dms_rocketmq_messages.test.messages.0.message_id
  client_id = data.huaweicloud_dms_rocketmq_consumers.test.clients.0.client_id
}

resource "huaweicloud_dms_rocketmq_consumption_verify" "test" {
  instance_id     = "%[1]s"
  group           = "%[2]s"
  topic           = "%[3]s"
  client_id       = local.client_id
  message_id_list = [local.msg_id]
}`, acceptance.HW_DMS_ROCKETMQ_INSTANCE_ID, acceptance.HW_DMS_ROCKETMQ_GROUP_NAME, acceptance.HW_DMS_ROCKETMQ_TOPIC_NAME,
		startTime, endTime)
}

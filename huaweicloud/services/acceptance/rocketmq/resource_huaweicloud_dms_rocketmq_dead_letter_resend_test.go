package rocketmq

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRocketMQDeadLetterResend_basic(t *testing.T) {
	resourceName := "huaweicloud_dms_rocketmq_dead_letter_resend.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSRocketMQInstanceID(t)
			acceptance.TestAccPreCheckDMSRocketMQGroupName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRocketMQDeadLetterResend_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "resend_results.#"),
					resource.TestCheckResourceAttrSet(resourceName, "resend_results.0.message_id"),
					resource.TestCheckResourceAttrSet(resourceName, "resend_results.0.error_code"),
					resource.TestCheckResourceAttrSet(resourceName, "resend_results.0.error_message"),
				),
			},
		},
	})
}

func testAccRocketMQDeadLetterResend_basic() string {
	currentTime := time.Now()
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()).UnixMilli()
	endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, currentTime.Location()).UnixMilli()
	return fmt.Sprintf(`
data "huaweicloud_dms_rocketmq_messages" "test" {
  instance_id = "%[1]s"
  topic       = urlencode("%%DLQ%%%[2]s")
  start_time  = %[3]v
  end_time    = %[4]v
}

locals {
  msg_id = data.huaweicloud_dms_rocketmq_messages.test.messages.0.message_id
}

resource "huaweicloud_dms_rocketmq_dead_letter_resend" "test" {
  instance_id     = "%[1]s"
  topic           = "%%DLQ%%%[2]s"
  message_id_list = [local.msg_id]
}`, acceptance.HW_DMS_ROCKETMQ_INSTANCE_ID, acceptance.HW_DMS_ROCKETMQ_GROUP_NAME, startTime, endTime)
}

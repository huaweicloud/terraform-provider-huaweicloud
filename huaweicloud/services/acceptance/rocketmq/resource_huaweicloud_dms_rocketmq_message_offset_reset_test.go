package rocketmq

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRocketMQMessageOffsetReset_basic(t *testing.T) {
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
				Config: testAccRocketMQMessageOffsetReset_basic(),
			},
		},
	})
}

func testAccRocketMQMessageOffsetReset_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_rocketmq_message_offset_reset" "test" {
  instance_id = "%[1]s"
  group       = "%[2]s"
  topic       = "%[3]s"
  timestamp   = 0
}`, acceptance.HW_DMS_ROCKETMQ_INSTANCE_ID, acceptance.HW_DMS_ROCKETMQ_GROUP_NAME, acceptance.HW_DMS_ROCKETMQ_TOPIC_NAME)
}

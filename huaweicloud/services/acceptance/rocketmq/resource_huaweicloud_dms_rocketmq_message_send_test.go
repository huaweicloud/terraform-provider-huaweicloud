package rocketmq

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRocketMQMessageSend_basic(t *testing.T) {
	resourceName := "huaweicloud_dms_rocketmq_message_send.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSRocketMQInstanceID(t)
			acceptance.TestAccPreCheckDMSRocketMQTopicName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		//  This is a one-time action resource not delete logic.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      testAccRocketMQMessageSend_instanceNotFound(),
				ExpectError: regexp.MustCompile(`This DMS instance does not exist`),
			},
			{
				Config:      testAccRocketMQMessageSend_topicNotFound(),
				ExpectError: regexp.MustCompile(`Topic does not exist`),
			},
			{
				Config: testAccRocketMQMessageSend_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "property_list.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "msg_id"),
					resource.TestCheckResourceAttrSet(resourceName, "queue_id"),
					resource.TestCheckResourceAttrSet(resourceName, "broker_name"),
					// `queue_offset` value may be 0, so we don't check it.
				),
			},
		},
	})
}

func testAccRocketMQMessageSend_instanceNotFound() string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
resource "huaweicloud_dms_rocketmq_message_send" "test" {
  instance_id = "%[1]s"
  topic       = "%[2]s"
  body        = "tf terraform script test"
}
`, randomId, acceptance.HW_DMS_ROCKETMQ_TOPIC_NAME)
}

func testAccRocketMQMessageSend_topicNotFound() string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_rocketmq_message_send" "test" {
  instance_id = "%s"
  topic       = "not_found"
  body        = "tf terraform script test"
}
`, acceptance.HW_DMS_ROCKETMQ_INSTANCE_ID)
}

func testAccRocketMQMessageSend_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_rocketmq_message_send" "test" {
  instance_id = "%[1]s"
  topic       = "%[2]s"
  body        = "tf terraform script test"

  property_list {
    name  = "KEYS"
    value = "owner"
  }
  property_list {
    name  = "TAGS"
    value = "terraform"
  }
}
`, acceptance.HW_DMS_ROCKETMQ_INSTANCE_ID, acceptance.HW_DMS_ROCKETMQ_TOPIC_NAME)
}

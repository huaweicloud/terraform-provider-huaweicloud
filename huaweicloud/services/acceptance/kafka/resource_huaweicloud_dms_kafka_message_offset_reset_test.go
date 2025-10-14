package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this test, ensure the consumer group status is `EMPTY` and has consumed the topic.
// The topic has produced messages.
func TestAccMessageOffsetReset_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
			acceptance.TestAccPreCheckDMSKafkaConsumerGroupName(t)
			acceptance.TestAccPreCheckDMSKafkaTopicName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config:      testAccMessageOffsetReset_instanceNotFound(),
				ExpectError: regexp.MustCompile(`This DMS instance does not exist`),
			},
			{
				Config: testAccMessageOffsetReset_basic_step1(),
			},
			{
				Config: testAccMessageOffsetReset_basic_step2(),
			},
		},
	})
}

func testAccMessageOffsetReset_instanceNotFound() string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_message_offset_reset" "not_found" {
  instance_id = "%[1]s"
  group       = "%[2]s"
  partition   = -1
  timestamp   = "0"
}`, randomId, acceptance.HW_DMS_KAFKA_CONSUMER_GROUP_NAME)
}

// Reset message offset for all topic with timestamp.
func testAccMessageOffsetReset_basic_step1() string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_message_offset_reset" "test" {
  instance_id = "%[1]s"
  group       = "%[2]s"
  partition   = -1
  timestamp   = "0"
}`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, acceptance.HW_DMS_KAFKA_CONSUMER_GROUP_NAME)
}

// Reset message offset for all partition under specific topic with message offset.
func testAccMessageOffsetReset_basic_step2() string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_message_offset_reset" "test2" {
  instance_id    = "%[1]s"
  group          = "%[2]s"
  topic          = "%[3]s"
  partition      = 0
  message_offset = "1"
}`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, acceptance.HW_DMS_KAFKA_CONSUMER_GROUP_NAME, acceptance.HW_DMS_KAFKA_TOPIC_NAME)
}

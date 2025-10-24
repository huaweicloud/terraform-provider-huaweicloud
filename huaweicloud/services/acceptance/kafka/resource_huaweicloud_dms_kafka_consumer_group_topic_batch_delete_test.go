package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running the test, ensure that the consumer group status is `EMPTY` and has consumed the topic.
func TestAccConsumerGroupTopicBatchDelete_basic(t *testing.T) {
	rName := "huaweicloud_dms_kafka_consumer_group_topic_batch_delete.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
			acceptance.TestAccPreCheckDMSKafkaConsumerGroupName(t)
			acceptance.TestAccPreCheckDMSKafkaTopicName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      testAccConsumerGroupTopicBatchDelete_instanceNotFound(),
				ExpectError: regexp.MustCompile(`This DMS instance does not exist`),
			},
			{
				Config:      testAccConsumerGroupTopicBatchDelete_consumerGroupNotFound(),
				ExpectError: regexp.MustCompile(`The consumer group does not exist`),
			},
			{
				Config: testAccConsumerGroupTopicBatchDelete_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "result.#", "1"),
					resource.TestCheckResourceAttr(rName, "result.0.name", acceptance.HW_DMS_KAFKA_TOPIC_NAME),
					resource.TestCheckResourceAttr(rName, "result.0.success", "true"),
				),
			},
		},
	})
}

func testAccConsumerGroupTopicBatchDelete_instanceNotFound() string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_consumer_group_topic_batch_delete" "test" {
  instance_id = "%[1]s"
  group       = "%[2]s"
  topics      = ["%[3]s"]
}
`, randomId, acceptance.HW_DMS_KAFKA_CONSUMER_GROUP_NAME, acceptance.HW_DMS_KAFKA_TOPIC_NAME)
}

func testAccConsumerGroupTopicBatchDelete_consumerGroupNotFound() string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_consumer_group_topic_batch_delete" "test" {
  instance_id = "%[1]s"
  group       = "not_found"
  topics      = ["%[2]s"]
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, acceptance.HW_DMS_KAFKA_TOPIC_NAME)
}

func testAccConsumerGroupTopicBatchDelete_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_consumer_group_topic_batch_delete" "test" {
  instance_id = "%[1]s"
  group       = "%[2]s"
  topics      = ["%[3]s"]
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, acceptance.HW_DMS_KAFKA_CONSUMER_GROUP_NAME, acceptance.HW_DMS_KAFKA_TOPIC_NAME)
}

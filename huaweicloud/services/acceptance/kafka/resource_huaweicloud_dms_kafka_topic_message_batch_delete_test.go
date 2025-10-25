package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTopicMessageBatchDelete_basic(t *testing.T) {
	rName := "huaweicloud_dms_kafka_topic_message_batch_delete.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This is a one-time action resource not delete logic.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      testAccTopicMessageBatchDelete_instanceNotFound(),
				ExpectError: regexp.MustCompile(`This DMS instance does not exist`),
			},
			{
				Config:      testAccTopicMessageBatchDelete_topicNotFound(),
				ExpectError: regexp.MustCompile(`Invalid topic in the request`),
			},
			{
				Config: testAccTopicMessageBatchDelete_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "result.#", "1"),
					resource.TestCheckResourceAttr(rName, "result.0.partition", "0"),
					resource.TestCheckResourceAttr(rName, "result.0.result", "success"),
				),
			},
		},
	})
}

func testAccTopicMessageBatchDelete_instanceNotFound() string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_topic_message_batch_delete" "test" {
  instance_id = "%[1]s"
  topic       = "instance_not_found"

  partitions {
    partition = 0
    offset    = 1
  }
}
`, randomId)
}

func testAccTopicMessageBatchDelete_topicNotFound() string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_topic_message_batch_delete" "test" {
  instance_id = "%[1]s"
  topic       = "topic_not_found"

  partitions {
    partition = 0
    offset    = 1
  }
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID)
}

func testAccTopicMessageBatchDelete_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_topic" "test2" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  partitions  = 2
}

resource "huaweicloud_dms_kafka_message_produce" "test" {
  instance_id = "%[1]s"
  topic       = huaweicloud_dms_kafka_topic.test2.name
  body        = "tf test"

  property_list {
    name  = "PARTITION"
    value = "0"
  }
}

resource "huaweicloud_dms_kafka_topic_message_batch_delete" "test" {
  instance_id = "%[1]s"
  topic       = huaweicloud_dms_kafka_topic.test2.name

  partitions {
    partition = 0
    offset    = 1
  }

  depends_on = [huaweicloud_dms_kafka_message_produce.test]
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, acceptance.RandomAccResourceName())
}

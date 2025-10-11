package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this test, make sure that the consumer group is online (status is STABLE).
func TestAccDataConsumerGroupMessageOffsets_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dms_kafka_consumer_group_message_offsets.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
			acceptance.TestAccPreCheckDMSKafkaConsumerGroupName(t)
			acceptance.TestAccPreCheckDMSKafkaTopicName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataConsumerGroupMessageOffsets_instanceNotFound(),
				ExpectError: regexp.MustCompile(`This DMS instance does not exist`),
			},
			{
				Config: testAccDataConsumerGroupMessageOffsets_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "message_offsets.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "message_offsets.0.consumer_id"),
					resource.TestCheckResourceAttrSet(all, "message_offsets.0.host"),
					resource.TestCheckResourceAttrSet(all, "message_offsets.0.client_id"),
					resource.TestCheckOutput("message_offset_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccDataConsumerGroupMessageOffsets_instanceNotFound() string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dms_kafka_consumer_group_message_offsets" "test" {
  instance_id = "%[1]s"
  group       = "%[2]s"
  topic       = "%[3]s"
}
`, randomId, acceptance.HW_DMS_KAFKA_CONSUMER_GROUP_NAME, acceptance.HW_DMS_KAFKA_TOPIC_NAME)
}

func testAccDataConsumerGroupMessageOffsets_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_message_produce" "test" {
  instance_id = "%[1]s"
  topic       = "%[3]s"
  body        = "terraform test"

  property_list {
    name  = "PARTITION"
    value = "0"
  }
}

data "huaweicloud_dms_kafka_consumer_group_message_offsets" "test" {
  instance_id = "%[1]s"
  group       = "%[2]s"
  topic       = "%[3]s"
}

locals {
  message_offset = try([for v in data.huaweicloud_dms_kafka_consumer_group_message_offsets.test.message_offsets : v
  if v.partition == 0][0], {})
}

output "message_offset_validation_pass" {
  value = (
    lookup(local.message_offset, "partition", null) == 0 &&
    lookup(local.message_offset, "message_current_offset", null) > 0 &&
    lookup(local.message_offset, "message_log_end_offset", null) > 0
  )
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, acceptance.HW_DMS_KAFKA_CONSUMER_GROUP_NAME, acceptance.HW_DMS_KAFKA_TOPIC_NAME)
}

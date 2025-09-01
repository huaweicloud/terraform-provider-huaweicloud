package rocketmq

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this test, make sure that the subscribed topic exists under the consumer group and
// that messages have been sent to the topic.
func TestAccDataConsumerGroupTopics_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dms_rocketmq_consumer_group_topics.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byTopic   = "data.huaweicloud_dms_rocketmq_consumer_group_topics.filter_by_topic"
		dcByTopic = acceptance.InitDataSourceCheck(byTopic)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSRocketMQInstanceID(t)
			acceptance.TestAccPreCheckDMSRocketMQGroupName(t)
			acceptance.TestAccPreCheckDMSRocketMQTopicName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataConsumerGroupTopics_instanceNotFound(),
				ExpectError: regexp.MustCompile(`This DMS instance does not exist`),
			},
			{
				Config:      testAccDataConsumerGroupTopics_groupNotFound(),
				ExpectError: regexp.MustCompile(`Group does not exist`),
			},
			{
				Config: testAccDataConsumerGroupTopics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "topics.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByTopic.CheckResourceExists(),
					resource.TestCheckOutput("is_topic_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataConsumerGroupTopics_instanceNotFound() string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dms_rocketmq_consumer_group_topics" "test" {
  instance_id = "%[1]s"
  group       = "%[2]s"
}
`, randomId, acceptance.HW_DMS_ROCKETMQ_GROUP_NAME)
}

func testAccDataConsumerGroupTopics_groupNotFound() string {
	return fmt.Sprintf(`
data "huaweicloud_dms_rocketmq_consumer_group_topics" "test" {
  instance_id = "%[1]s"
  group       = "non_found_consumer_group"
}
`, acceptance.HW_DMS_ROCKETMQ_INSTANCE_ID)
}

func testAccDataConsumerGroupTopics_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dms_rocketmq_consumer_group_topics" "test" {
  instance_id = "%[1]s"
  group       = "%[2]s"
}

data "huaweicloud_dms_rocketmq_consumer_group_topics" "filter_by_topic" {
  instance_id = "%[1]s"
  group       = "%[2]s"
  topic       = "%[3]s"
}

locals {
  filter_result = data.huaweicloud_dms_rocketmq_consumer_group_topics.filter_by_topic
}

output "is_topic_filter_useful" {
  value = (
    local.filter_result.max_offset > 0 &&
    local.filter_result.consumer_offset > 0 &&
    length(local.filter_result.brokers) > 0 &&
    length(local.filter_result.brokers[0].broker_name) != "" &&
    length(local.filter_result.brokers[0].queues) > 0
  )
}
`, acceptance.HW_DMS_ROCKETMQ_INSTANCE_ID, acceptance.HW_DMS_ROCKETMQ_GROUP_NAME, acceptance.HW_DMS_ROCKETMQ_TOPIC_NAME)
}

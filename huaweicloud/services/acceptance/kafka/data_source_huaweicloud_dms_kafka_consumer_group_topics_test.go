package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this test, make sure that the subscribed topic exists under the consumer group.
func TestAccDataConsumerGroupTopics_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dms_kafka_consumer_group_topics.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byTopic   = "data.huaweicloud_dms_kafka_consumer_group_topics.filter_by_topic"
		dcByTopic = acceptance.InitDataSourceCheck(byTopic)

		bySortAsc   = "data.huaweicloud_dms_kafka_consumer_group_topics.filter_by_sort_asc"
		dcBySortAsc = acceptance.InitDataSourceCheck(bySortAsc)

		bySortDesc   = "data.huaweicloud_dms_kafka_consumer_group_topics.filter_by_sort_desc"
		dcBySortDesc = acceptance.InitDataSourceCheck(bySortDesc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
			acceptance.TestAccPreCheckDMSKafkaConsumerGroupName(t)
			acceptance.TestAccPreCheckDMSKafkaTopicName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"null": {
				Source:            "hashicorp/null",
				VersionConstraint: "3.2.1",
			},
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccDataConsumerGroupTopics_instanceNotFound(),
				ExpectError: regexp.MustCompile(`This DMS instance does not exist`),
			},
			{
				Config: testAccDataConsumerGroupTopics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "topics.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "topics.0.topic"),
					resource.TestMatchResourceAttr(all, "topics.0.partitions", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByTopic.CheckResourceExists(),
					resource.TestCheckOutput("is_topic_filter_useful", "true"),
					dcBySortAsc.CheckResourceExists(),
					dcBySortDesc.CheckResourceExists(),
					resource.TestCheckOutput("is_sort_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataConsumerGroupTopics_instanceNotFound() string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dms_kafka_consumer_group_topics" "test" {
  instance_id = "%[1]s"
  group       = "%[2]s"
}
`, randomId, acceptance.HW_DMS_KAFKA_CONSUMER_GROUP_NAME)
}

func testAccDataConsumerGroupTopics_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_message_produce" "test" {
  instance_id = "%[1]s"
  topic       = "%[3]s"
  body        = "terraform test"
}

# Wait for the message to be produced.
resource "null_resource" "test" {
  provisioner "local-exec" {
    command = "sleep 20"
  }

  depends_on = [huaweicloud_dms_kafka_message_produce.test]
}

data "huaweicloud_dms_kafka_consumer_group_topics" "test" {
  instance_id = "%[1]s"
  group       = "%[2]s"

  depends_on = [null_resource.test]
}

# Fuzzy search.
data "huaweicloud_dms_kafka_consumer_group_topics" "filter_by_topic" {
  instance_id = "%[1]s"
  group       = "%[2]s"
  topic       = "%[3]s"

  depends_on = [null_resource.test]
}

locals {
  filter_by_topic_result = [for v in data.huaweicloud_dms_kafka_consumer_group_topics.filter_by_topic.topics : strcontains(v.topic, "%[3]s")]
}

output "is_topic_filter_useful" {
  value = length(local.filter_by_topic_result) > 0 && alltrue(local.filter_by_topic_result)
}

data "huaweicloud_dms_kafka_consumer_group_topics" "filter_by_sort_asc" {
  instance_id = "%[1]s"
  group       = "%[2]s"
  sort_key    = "topic"
  sort_dir    = "asc"

  depends_on = [null_resource.test]
}

data "huaweicloud_dms_kafka_consumer_group_topics" "filter_by_sort_desc" {
  instance_id = "%[1]s"
  group       = "%[2]s"
  sort_key    = "topic"

  depends_on = [null_resource.test]
}

locals {
  filter_by_sort_asc_result  = data.huaweicloud_dms_kafka_consumer_group_topics.filter_by_sort_asc.topics[*].topic
  filter_by_sort_desc_result = data.huaweicloud_dms_kafka_consumer_group_topics.filter_by_sort_desc.topics[*].topic
}

output "is_sort_filter_useful" {
  value = (
    length(local.filter_by_sort_asc_result) > 0
    && length(local.filter_by_sort_asc_result) == length(local.filter_by_sort_desc_result)
    && local.filter_by_sort_asc_result[0] == local.filter_by_sort_desc_result[length(local.filter_by_sort_desc_result) - 1]
  )
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, acceptance.HW_DMS_KAFKA_CONSUMER_GROUP_NAME, acceptance.HW_DMS_KAFKA_TOPIC_NAME)
}

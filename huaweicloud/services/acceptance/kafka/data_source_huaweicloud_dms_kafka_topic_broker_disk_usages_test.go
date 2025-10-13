package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataTopicBrokerDiskUsages_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dms_kafka_topic_broker_disk_usages.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byTop   = "data.huaweicloud_dms_kafka_topic_broker_disk_usages.filter_by_top"
		dcByTop = acceptance.InitDataSourceCheck(byTop)

		byMinSize   = "data.huaweicloud_dms_kafka_topic_broker_disk_usages.filter_by_min_size"
		dcByMinSize = acceptance.InitDataSourceCheck(byMinSize)

		byPercentage   = "data.huaweicloud_dms_kafka_topic_broker_disk_usages.filter_by_percentage"
		dcByPercentage = acceptance.InitDataSourceCheck(byPercentage)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataTopicBrokerDiskUsages_instanceNotFound(),
				ExpectError: regexp.MustCompile("This DMS instance does not exist"),
			},
			{
				Config: testAccDataTopicBrokerDiskUsages_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "disk_usages.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "disk_usages.0.broker_name"),
					resource.TestCheckResourceAttrSet(dataSource, "disk_usages.0.data_disk_size"),
					resource.TestCheckResourceAttrSet(dataSource, "disk_usages.0.data_disk_use"),
					resource.TestCheckResourceAttrSet(dataSource, "disk_usages.0.data_disk_free"),
					resource.TestCheckResourceAttrSet(dataSource, "disk_usages.0.data_disk_use_percentage"),
					resource.TestCheckResourceAttrSet(dataSource, "disk_usages.0.status"),
					resource.TestMatchResourceAttr(dataSource, "disk_usages.0.topics.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "disk_usages.0.topics.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "disk_usages.0.topics.0.topic_name"),
					resource.TestCheckResourceAttrSet(dataSource, "disk_usages.0.topics.0.topic_partition"),
					resource.TestCheckResourceAttrSet(dataSource, "disk_usages.0.topics.0.percentage"),
					dcByTop.CheckResourceExists(),
					resource.TestCheckOutput("is_top_useful", "true"),
					dcByMinSize.CheckResourceExists(),
					resource.TestCheckOutput("is_min_size_useful", "true"),
					dcByPercentage.CheckResourceExists(),
					resource.TestCheckOutput("is_percentage_useful", "true"),
				),
			},
		},
	})
}

func testAccDataTopicBrokerDiskUsages_instanceNotFound() string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dms_kafka_topic_broker_disk_usages" "test" {
  instance_id = "%s"
}
`, randomId)
}

func testAccDataTopicBrokerDiskUsages_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dms_kafka_topic_broker_disk_usages" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_dms_kafka_topic_broker_disk_usages" "filter_by_top" {
  instance_id = "%[1]s"
  top         = 1
}

locals {
  top_filter_result = [for v in data.huaweicloud_dms_kafka_topic_broker_disk_usages.filter_by_top.disk_usages : length(v.topics) <= 1]
}

output "is_top_useful" {
  value = length(local.top_filter_result) >= 0 && alltrue(local.top_filter_result)
}

data "huaweicloud_dms_kafka_topic_broker_disk_usages" "filter_by_min_size" {
  instance_id = "%[1]s"
  min_size    = "1K"
}

locals {
  min_size_filter_result = [for v in data.huaweicloud_dms_kafka_topic_broker_disk_usages.filter_by_min_size.disk_usages : length(v.topics) > 0]
}

output "is_min_size_useful" {
  value = length(local.min_size_filter_result) >= 0 && alltrue(local.min_size_filter_result)
}

data "huaweicloud_dms_kafka_topic_broker_disk_usages" "filter_by_percentage" {
  instance_id = "%[1]s"
  percentage  = floor(data.huaweicloud_dms_kafka_topic_broker_disk_usages.test.disk_usages[0].topics[0].percentage)
}

locals {
  percentage_filter_result = [for v in data.huaweicloud_dms_kafka_topic_broker_disk_usages.filter_by_percentage.disk_usages : length(v.topics) > 0]
}

output "is_percentage_useful" {
  value = length(local.percentage_filter_result) >= 0 && alltrue(local.percentage_filter_result)
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID)
}

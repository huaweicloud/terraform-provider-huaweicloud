package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this test, make sure that the consumer group is `STABLE` status.
func TestAccDataInstanceMonitoringDimensions_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dms_kafka_instance_monitoring_dimensions.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
			acceptance.TestAccPreCheckDMSKafkaTopicName(t)
			acceptance.TestAccPreCheckDMSKafkaConsumerGroupName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataInstanceMonitoringDimensions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "dimensions.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "dimensions.0.name"),
					resource.TestMatchResourceAttr(dataSourceName, "dimensions.0.metrics.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestMatchResourceAttr(dataSourceName, "dimensions.0.key_name.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestMatchResourceAttr(dataSourceName, "dimensions.0.dim_router.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestMatchResourceAttr(dataSourceName, "instance_ids.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "instance_ids.0.name"),
					resource.TestMatchResourceAttr(dataSourceName, "nodes.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "nodes.0.name"),
					resource.TestMatchResourceAttr(dataSourceName, "queues.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "queues.0.name"),
					resource.TestMatchResourceAttr(dataSourceName, "queues.0.partitions.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "queues.0.partitions.0.name"),
					resource.TestMatchResourceAttr(dataSourceName, "groups.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.name"),
					resource.TestMatchResourceAttr(dataSourceName, "groups.0.queues.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.queues.0.name"),
					resource.TestMatchResourceAttr(dataSourceName, "groups.0.queues.0.partitions.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.queues.0.partitions.0.name"),
					resource.TestCheckOutput("queues_validation", "true"),
					resource.TestCheckOutput("groups_validation", "true"),
				),
			},
		},
	})
}

func testAccDataInstanceMonitoringDimensions_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dms_kafka_instance_monitoring_dimensions" "test" {
  instance_id = "%[1]s"
}

locals {
  queues_filter_result = [for v in data.huaweicloud_dms_kafka_instance_monitoring_dimensions.test.queues : v.name == "%[2]s"]
  group_filter_result  = [for v in data.huaweicloud_dms_kafka_instance_monitoring_dimensions.test.groups : v.name == "%[3]s"]
}

output "queues_validation" {
  value = length(local.queues_filter_result) > 0 && alltrue(local.queues_filter_result)
}

output "groups_validation" {
  value = length(local.group_filter_result) > 0 && alltrue(local.group_filter_result)
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, acceptance.HW_DMS_KAFKA_TOPIC_NAME, acceptance.HW_DMS_KAFKA_CONSUMER_GROUP_NAME)
}

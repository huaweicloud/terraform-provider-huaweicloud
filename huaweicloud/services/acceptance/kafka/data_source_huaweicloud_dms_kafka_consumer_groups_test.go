package kafka

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDmsKafkaConsumerGroups_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dms_kafka_consumer_groups.all"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDmsKafkaConsumerGroups_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "groups.#"),
					resource.TestCheckOutput("name_validation", "true"),
					resource.TestCheckOutput("state_validation", "true"),
					resource.TestCheckOutput("description_validation", "true"),
					resource.TestCheckOutput("lag_validation", "true"),
					resource.TestCheckOutput("coordinator_id_validation", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDmsKafkaConsumerGroups_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_kafka_consumer_groups" "all" {
  depends_on = [huaweicloud_dms_kafka_consumer_group.test]

  instance_id = huaweicloud_dms_kafka_instance.test.id
}

data "huaweicloud_dms_kafka_consumer_groups" "test" {
  instance_id    = huaweicloud_dms_kafka_instance.test.id
  name           = local.test_refer.name
  description    = local.test_refer.description
  lag            = local.test_refer.lag
  coordinator_id = local.test_refer.coordinator_id
  state          = local.test_refer.state
}

locals {
  test_refer   = huaweicloud_dms_kafka_consumer_group.test
  test_results = data.huaweicloud_dms_kafka_consumer_groups.test
}

output "name_validation" {
  value = alltrue([for v in local.test_results.groups[*].name : strcontains(v, local.test_refer.name)])
}

output "description_validation" {
  value = alltrue([for v in local.test_results.groups[*].description : strcontains(v, local.test_refer.description)])
}

output "lag_validation" {
  value = alltrue([for v in local.test_results.groups[*].lag : v == local.test_refer.lag])
}

output "coordinator_id_validation" {
  value = alltrue([for v in local.test_results.groups[*].coordinator_id : v == local.test_refer.coordinator_id])
}

output "state_validation" {
  value = alltrue([for v in local.test_results.groups[*].state : v == local.test_refer.state])
}
`, testDmsKafkaConsumerGroup_basic(name))
}

func TestAccDataSourceDmsKafkaConsumerGroups_consumers(t *testing.T) {
	dataSource := "data.huaweicloud_dms_kafka_consumer_groups.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDmsKafkaConsumerGroups_consumers(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "groups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.members.#"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.group_message_offsets.#"),
				),
			},
		},
	})
}

func testDataSourceDmsKafkaConsumerGroups_consumers() string {
	return fmt.Sprintf(`
data "huaweicloud_dms_kafka_consumer_groups" "test" {
  instance_id = "%[1]s"
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID)
}

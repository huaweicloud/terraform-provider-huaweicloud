package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this test, make sure that the `HW_DMS_KAFKA_CONSUMER_GROUP_NAME` is online (status is online).
func TestAccDataConsumerGroups_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dms_kafka_consumer_groups.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByName   = "data.huaweicloud_dms_kafka_consumer_groups.filter_by_name"
		dcFilterByName = acceptance.InitDataSourceCheck(filterByName)

		filterByDescription   = "data.huaweicloud_dms_kafka_consumer_groups.filter_by_description"
		dcFilterByDescription = acceptance.InitDataSourceCheck(filterByDescription)

		filterByLag   = "data.huaweicloud_dms_kafka_consumer_groups.filter_by_lag"
		dcFilterByLag = acceptance.InitDataSourceCheck(filterByLag)

		filterByCoordinatorId   = "data.huaweicloud_dms_kafka_consumer_groups.filter_by_coordinator_id"
		dcFilterByCoordinatorId = acceptance.InitDataSourceCheck(filterByCoordinatorId)

		filterByState   = "data.huaweicloud_dms_kafka_consumer_groups.filter_by_state"
		dcFilterByState = acceptance.InitDataSourceCheck(filterByState)

		online   = "data.huaweicloud_dms_kafka_consumer_groups.online"
		dcOnline = acceptance.InitDataSourceCheck(online)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
			acceptance.TestAccPreCheckDMSKafkaConsumerGroupName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataConsumerGroups_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "groups.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by 'name' parameter.
					dcFilterByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					// Filter by 'description' parameter.
					dcFilterByDescription.CheckResourceExists(),
					resource.TestCheckOutput("is_description_filter_useful", "true"),
					// Filter by 'lag' parameter.
					dcFilterByLag.CheckResourceExists(),
					resource.TestCheckOutput("is_lag_filter_useful", "true"),
					// Filter by 'coordinator_id' parameter.
					dcFilterByCoordinatorId.CheckResourceExists(),
					resource.TestCheckOutput("is_coordinator_id_filter_useful", "true"),
					// Filter by 'state' parameter.
					dcFilterByState.CheckResourceExists(),
					resource.TestCheckOutput("is_state_filter_useful", "true"),
					// Use the online consumer group to filter.
					dcOnline.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(online, "groups.0.assignment_strategy"),
					resource.TestMatchResourceAttr(online, "groups.0.members.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
				),
			},
		},
	})
}

func testAccDataConsumerGroups_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_consumer_group" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  description = "Created by Terraform script"
}

# Without any filter parameters.
data "huaweicloud_dms_kafka_consumer_groups" "all" {
  instance_id = "%[1]s"
  depends_on  = [huaweicloud_dms_kafka_consumer_group.test]
}

# Filter by 'name' parameter.
locals {
  name = huaweicloud_dms_kafka_consumer_group.test.name
}

data "huaweicloud_dms_kafka_consumer_groups" "filter_by_name" {
  instance_id = "%[1]s"
  name        = local.name
  depends_on  = [huaweicloud_dms_kafka_consumer_group.test]
}

locals {
  name_filter_result = [for v in data.huaweicloud_dms_kafka_consumer_groups.filter_by_name.groups[*].name :
    strcontains(v, local.name)
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by 'description' parameter.
locals {
  description = huaweicloud_dms_kafka_consumer_group.test.description
}

data "huaweicloud_dms_kafka_consumer_groups" "filter_by_description" {
  instance_id = "%[1]s"
  description = local.description
  depends_on  = [huaweicloud_dms_kafka_consumer_group.test]
}

locals {
  description_filter_result = [
    for v in data.huaweicloud_dms_kafka_consumer_groups.filter_by_description.groups[*].description :
    v == local.description
  ]
}

output "is_description_filter_useful" {
  value = length(local.description_filter_result) > 0 && alltrue(local.description_filter_result)
}

# Filter by 'lag' parameter.
locals {
  lag = try(data.huaweicloud_dms_kafka_consumer_groups.all.groups[0].lag, null)
}

data "huaweicloud_dms_kafka_consumer_groups" "filter_by_lag" {
  instance_id = "%[1]s"
  lag         = local.lag
}

locals {
  lag_filter_result = [for v in data.huaweicloud_dms_kafka_consumer_groups.filter_by_lag.groups[*].lag :
    v == local.lag
  ]
}

output "is_lag_filter_useful" {
  value = length(local.lag_filter_result) > 0 && alltrue(local.lag_filter_result)
}

# Filter by 'coordinator_id' parameter.
locals {
  coordinator_id = huaweicloud_dms_kafka_consumer_group.test.coordinator_id
}

data "huaweicloud_dms_kafka_consumer_groups" "filter_by_coordinator_id" {
  instance_id    = "%[1]s"
  coordinator_id = local.coordinator_id
}

locals {
  coordinator_id_filter_result = [
    for v in data.huaweicloud_dms_kafka_consumer_groups.filter_by_coordinator_id.groups[*].coordinator_id :
    v == local.coordinator_id
  ]
}

output "is_coordinator_id_filter_useful" {
  value = length(local.coordinator_id_filter_result) > 0 && alltrue(local.coordinator_id_filter_result)
}

# Filter by 'state' parameter.
locals {
  state = huaweicloud_dms_kafka_consumer_group.test.state
}

data "huaweicloud_dms_kafka_consumer_groups" "filter_by_state" {
  instance_id = "%[1]s"
  state       = local.state
}

locals {
  state_filter_result = [
    for v in data.huaweicloud_dms_kafka_consumer_groups.filter_by_state.groups[*].state :
    v == local.state
  ]
}

output "is_state_filter_useful" {
  value = length(local.state_filter_result) > 0 && alltrue(local.state_filter_result)
}

# Use the online consumer group to filter.
data "huaweicloud_dms_kafka_consumer_groups" "online" {
  instance_id = "%[1]s"
  name        = "%[3]s"
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, name, acceptance.HW_DMS_KAFKA_CONSUMER_GROUP_NAME)
}

package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAsNotifications_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_as_notifications.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byName   = "data.huaweicloud_as_notifications.filter_by_topic_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNameNotFound   = "data.huaweicloud_as_notifications.not_found"
		dcByNameNotFound = acceptance.InitDataSourceCheck(byNameNotFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare the AS group containing the notifications in advance and configure the AS group ID into
			// the environment variable.
			acceptance.TestAccPreCheckASScalingGroupID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAsNotifications_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "topics.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "topics.0.topic_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "topics.0.topic_urn"),
					resource.TestCheckResourceAttrSet(dataSourceName, "topics.0.events.#"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("topic_name_filter_is_useful", "true"),

					dcByNameNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found", "true"),
				),
			},
		},
	})
}

func testDataSourceAsNotifications_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_as_notifications" "test" {
  scaling_group_id = "%[1]s"
}

# Filter by topic_name
locals {
  topic_name = data.huaweicloud_as_notifications.test.topics[0].topic_name
}

data "huaweicloud_as_notifications" "filter_by_topic_name" {
  scaling_group_id = "%[1]s"
  topic_name       = local.topic_name
}

locals {
  topic_name_filter_result = [
    for v in data.huaweicloud_as_notifications.filter_by_topic_name.topics[*].topic_name : v == local.topic_name
  ]
}

output "topic_name_filter_is_useful" {
  value = alltrue(local.topic_name_filter_result) && length(local.topic_name_filter_result) > 0
}

data "huaweicloud_as_notifications" "not_found" {
  scaling_group_id = "%[1]s"
  topic_name       = "not_found"
}

output "is_not_found" {
  value = length(data.huaweicloud_as_notifications.not_found.topics) == 0
}
`, acceptance.HW_AS_SCALING_GROUP_ID)
}

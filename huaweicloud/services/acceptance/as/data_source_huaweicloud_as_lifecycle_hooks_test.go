package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceLifecycleHooks_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_as_lifecycle_hooks.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byName   = "data.huaweicloud_as_lifecycle_hooks.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byType   = "data.huaweicloud_as_lifecycle_hooks.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byDefaultResult   = "data.huaweicloud_as_lifecycle_hooks.filter_by_default_result"
		dcByDefaultResult = acceptance.InitDataSourceCheck(byDefaultResult)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare the AS group containing the lifecycle hooks in advance and configure the AS group ID into
			// the environment variable.
			acceptance.TestAccPreCheckASScalingGroupID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccLifecycleHooks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "lifecycle_hooks.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "lifecycle_hooks.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "lifecycle_hooks.0.default_result"),
					resource.TestCheckResourceAttrSet(dataSourceName, "lifecycle_hooks.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "lifecycle_hooks.0.notification_topic_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "lifecycle_hooks.0.notification_topic_urn"),
					resource.TestCheckResourceAttrSet(dataSourceName, "lifecycle_hooks.0.timeout"),
					resource.TestCheckResourceAttrSet(dataSourceName, "lifecycle_hooks.0.type"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),

					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),

					dcByDefaultResult.CheckResourceExists(),
					resource.TestCheckOutput("is_default_result_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccLifecycleHooks_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_as_lifecycle_hooks" "test" {
  scaling_group_id = "%[1]s"
}

# Filter by name
locals {
  name = data.huaweicloud_as_lifecycle_hooks.test.lifecycle_hooks[0].name
}

data "huaweicloud_as_lifecycle_hooks" "filter_by_name" {
  scaling_group_id = "%[1]s"
  name             = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_as_lifecycle_hooks.filter_by_name.lifecycle_hooks[*].name : v == local.name
  ]
}

output "is_name_filter_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

# Filter by type
locals {
  type = data.huaweicloud_as_lifecycle_hooks.test.lifecycle_hooks[0].type
}

data "huaweicloud_as_lifecycle_hooks" "filter_by_type" {
  scaling_group_id = "%[1]s"
  type             = local.type
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_as_lifecycle_hooks.filter_by_type.lifecycle_hooks[*].type : v == local.type
  ]
}

output "is_type_filter_useful" {
  value = alltrue(local.type_filter_result) && length(local.type_filter_result) > 0
}

# Filter by default_result
locals {
  default_result = data.huaweicloud_as_lifecycle_hooks.test.lifecycle_hooks[0].default_result
}

data "huaweicloud_as_lifecycle_hooks" "filter_by_default_result" {
  scaling_group_id = "%[1]s"
  default_result   = local.default_result
}

locals {
  default_result_filter_result = [
    for v in data.huaweicloud_as_lifecycle_hooks.filter_by_default_result.lifecycle_hooks[*].default_result : v == local.default_result
  ]
}

output "is_default_result_filter_useful" {
  value = alltrue(local.default_result_filter_result) && length(local.default_result_filter_result) > 0
}
`, acceptance.HW_AS_SCALING_GROUP_ID)
}

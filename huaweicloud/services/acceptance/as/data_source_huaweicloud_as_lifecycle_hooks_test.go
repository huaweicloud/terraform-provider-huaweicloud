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
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckASScalingGroupID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccLifecycleHooks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "scaling_group_id"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("default_result_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccLifecycleHooks_basic() string {
	rName := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_as_lifecycle_hook" "test1" {
  name                   = "%[2]s1"
  type                   = "ADD"
  default_result         = "ABANDON"
  scaling_group_id       = huaweicloud_as_group.acc_as_group.id
  notification_topic_urn = huaweicloud_smn_topic.test.topic_urn
  notification_message   = "This is a test message"
}

resource "huaweicloud_as_lifecycle_hook" "test2" {
  name                   = "%[2]s2"
  type                   = "REMOVE"
  default_result         = "CONTINUE"
  scaling_group_id       = huaweicloud_as_group.acc_as_group.id
  notification_topic_urn = huaweicloud_smn_topic.test.topic_urn
  notification_message   = "This is a test message"
}

resource "huaweicloud_as_lifecycle_hook" "test3" {
  name                   = "%[2]s3"
  type                   = "ADD"
  default_result         = "CONTINUE"
  scaling_group_id       = huaweicloud_as_group.acc_as_group.id
  notification_topic_urn = huaweicloud_smn_topic.test.topic_urn
  notification_message   = "This is a test message"
}

data "huaweicloud_as_lifecycle_hooks" "test" {
  depends_on = [
	huaweicloud_as_lifecycle_hook.test1,
	huaweicloud_as_lifecycle_hook.test2,
	huaweicloud_as_lifecycle_hook.test3,
  ]

  scaling_group_id = huaweicloud_as_group.acc_as_group.id
}

locals {
  name = data.huaweicloud_as_lifecycle_hooks.test.lifecycle_hooks[0].name
}

data "huaweicloud_as_lifecycle_hooks" "filter_by_name" {
  scaling_group_id = huaweicloud_as_group.acc_as_group.id
  name             = local.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_as_lifecycle_hooks.filter_by_name.lifecycle_hooks) > 0 && alltrue( 
    [for v in data.huaweicloud_as_lifecycle_hooks.filter_by_name.lifecycle_hooks[*].name : v == local.name]
  )  
}

locals {
  type = data.huaweicloud_as_lifecycle_hooks.test.lifecycle_hooks[0].type
}

data "huaweicloud_as_lifecycle_hooks" "filter_by_type" {
  scaling_group_id = huaweicloud_as_group.acc_as_group.id
  type             = local.type
}

output "type_filter_is_useful" {
  value = length(data.huaweicloud_as_lifecycle_hooks.filter_by_type.lifecycle_hooks) > 0 && alltrue( 
    [for v in data.huaweicloud_as_lifecycle_hooks.filter_by_type.lifecycle_hooks[*].type : v == local.type]
  )  
}

locals {
  default_result = data.huaweicloud_as_lifecycle_hooks.test.lifecycle_hooks[0].default_result
}

data "huaweicloud_as_lifecycle_hooks" "filter_by_default_result" {
  scaling_group_id = huaweicloud_as_group.acc_as_group.id
  default_result   = local.default_result
}

output "default_result_filter_is_useful" {
  value = length(data.huaweicloud_as_lifecycle_hooks.filter_by_default_result.lifecycle_hooks) > 0 && alltrue(
    [for v in data.huaweicloud_as_lifecycle_hooks.filter_by_default_result.lifecycle_hooks[*].default_result : v == local.default_result]
  )  
}
`, testASLifecycleHook_base(rName), rName)
}

package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceV2Policies_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_asv2_policies.test"
		name           = acceptance.RandomAccResourceName()
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byPolicyId   = "data.huaweicloud_asv2_policies.filter_by_policy_id"
		dcByPolicyId = acceptance.InitDataSourceCheck(byPolicyId)

		byPolicyName   = "data.huaweicloud_asv2_policies.filter_by_policy_name"
		dcByPolicyName = acceptance.InitDataSourceCheck(byPolicyName)

		byPolicyType   = "data.huaweicloud_asv2_policies.filter_by_policy_type"
		dcByPolicyType = acceptance.InitDataSourceCheck(byPolicyType)

		byResourceType   = "data.huaweicloud_asv2_policies.filter_by_resource_type"
		dcByResourceType = acceptance.InitDataSourceCheck(byResourceType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceV2Policies_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "scaling_policies.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "scaling_policies.0.scaling_policy_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "scaling_policies.0.scaling_policy_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "scaling_policies.0.policy_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "scaling_policies.0.scaling_policy_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "scaling_policies.0.scaling_resource_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "scaling_policies.0.scaling_policy_action.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "scaling_policies.0.scheduled_policy.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "scaling_policies.0.scheduled_policy.0.launch_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "scaling_policies.0.scheduled_policy.0.recurrence_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "scaling_policies.0.scheduled_policy.0.recurrence_value"),
					resource.TestCheckResourceAttrSet(dataSourceName, "scaling_policies.0.cool_down_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "scaling_policies.0.create_time"),

					dcByPolicyId.CheckResourceExists(),
					resource.TestCheckOutput("scaling_policy_id_filter_useful", "true"),

					dcByPolicyName.CheckResourceExists(),
					resource.TestCheckOutput("scaling_policy_name_filter_useful", "true"),

					dcByPolicyType.CheckResourceExists(),
					resource.TestCheckOutput("scaling_policy_type_filter_useful", "true"),

					dcByResourceType.CheckResourceExists(),
					resource.TestCheckOutput("scaling_resource_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceV2Policies_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_bandwidth" "test" {
  name = "%[1]s"
  size = 5
}

resource "huaweicloud_as_bandwidth_policy" "test" {
  scaling_policy_name = "%[1]s"
  scaling_policy_type = "RECURRENCE"
  bandwidth_id        = huaweicloud_vpc_bandwidth.test.id
  cool_down_time      = 600
  action              = "pause"

  scaling_policy_action {
    operation = "ADD"
    size      = 1
  }

  scheduled_policy {
    launch_time      = "07:00"
    recurrence_type  = "Weekly"
    recurrence_value = "1,3,5"
    start_time       = "2022-09-30T12:00Z"
    end_time         = "2122-12-30T12:00Z"
  }
}

data "huaweicloud_asv2_policies" "test" {
  depends_on = [huaweicloud_as_bandwidth_policy.test]
}

locals {
  scaling_policy_id = data.huaweicloud_asv2_policies.test.scaling_policies[0].scaling_policy_id
}

data "huaweicloud_asv2_policies" "filter_by_policy_id" {
  scaling_policy_id = local.scaling_policy_id
}

output "scaling_policy_id_filter_useful" {
  value = length(data.huaweicloud_asv2_policies.filter_by_policy_id.scaling_policies) > 0 && alltrue(
    [for v in data.huaweicloud_asv2_policies.filter_by_policy_id.scaling_policies[*].scaling_policy_id : v == local.scaling_policy_id]
  )
}

locals {
  scaling_policy_name = data.huaweicloud_asv2_policies.test.scaling_policies[0].scaling_policy_name
}

data "huaweicloud_asv2_policies" "filter_by_policy_name" {
  scaling_policy_name = local.scaling_policy_name
}

output "scaling_policy_name_filter_useful" {
  value = length(data.huaweicloud_asv2_policies.filter_by_policy_name.scaling_policies) > 0 && alltrue(
    [for v in data.huaweicloud_asv2_policies.filter_by_policy_name.scaling_policies[*].scaling_policy_name : v == local.scaling_policy_name]
  )
}

locals {
  scaling_policy_type = data.huaweicloud_asv2_policies.test.scaling_policies[0].scaling_policy_type
}

data "huaweicloud_asv2_policies" "filter_by_policy_type" {
  scaling_policy_type = local.scaling_policy_type
}

output "scaling_policy_type_filter_useful" {
  value = length(data.huaweicloud_asv2_policies.filter_by_policy_type.scaling_policies) > 0 && alltrue(
    [for v in data.huaweicloud_asv2_policies.filter_by_policy_type.scaling_policies[*].scaling_policy_type : v == local.scaling_policy_type]
  )
}

locals {
  scaling_resource_type = data.huaweicloud_asv2_policies.test.scaling_policies[0].scaling_resource_type
}

data "huaweicloud_asv2_policies" "filter_by_resource_type" {
  scaling_resource_type = local.scaling_resource_type
}

output "scaling_resource_type_filter_useful" {
  value = length(data.huaweicloud_asv2_policies.filter_by_resource_type.scaling_policies) > 0 && alltrue(
    [for v in data.huaweicloud_asv2_policies.filter_by_resource_type.scaling_policies[*].scaling_resource_type : v == local.scaling_resource_type]
  )
}
`, name)
}

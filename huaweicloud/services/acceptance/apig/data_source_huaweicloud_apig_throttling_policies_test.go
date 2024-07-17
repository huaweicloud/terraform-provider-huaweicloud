package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceThrottlingPolicies_basic(t *testing.T) {
	var (
		rName = "data.huaweicloud_apig_throttling_policies.test"
		dc    = acceptance.InitDataSourceCheck(rName)

		byId   = "data.huaweicloud_apig_throttling_policies.filter_by_id"
		dcById = acceptance.InitDataSourceCheck(byId)

		byName   = "data.huaweicloud_apig_throttling_policies.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNotFoundName   = "data.huaweicloud_apig_throttling_policies.filter_by_not_found_name"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)

		byType   = "data.huaweicloud_apig_throttling_policies.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceThrottlingPolicies_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(rName, "policies.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcById.CheckResourceExists(),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckResourceAttr(byId, "policies.0.type", "API-based"),
					resource.TestCheckResourceAttrSet(byId, "policies.0.name"),
					resource.TestCheckResourceAttrSet(byId, "policies.0.period"),
					resource.TestCheckResourceAttrSet(byId, "policies.0.period_unit"),
					resource.TestCheckResourceAttrSet(byId, "policies.0.max_api_requests"),
					resource.TestCheckResourceAttrSet(byId, "policies.0.created_at"),
					resource.TestCheckResourceAttr(byId, "policies.0.app_throttles.#", "1"),
					resource.TestCheckResourceAttr(byId, "policies.0.app_throttles.0.%", "4"),
					resource.TestCheckResourceAttrSet(byId, "policies.0.app_throttles.0.id"),
					resource.TestCheckResourceAttrSet(byId, "policies.0.app_throttles.0.max_api_requests"),
					resource.TestCheckResourceAttrSet(byId, "policies.0.app_throttles.0.throttling_object_id"),
					resource.TestCheckResourceAttrSet(byId, "policies.0.app_throttles.0.throttling_object_name"),
					resource.TestCheckResourceAttr(byId, "policies.0.user_throttles.#", "1"),
					resource.TestCheckResourceAttr(byId, "policies.0.user_throttles.0.%", "4"),
					resource.TestCheckResourceAttrSet(byId, "policies.0.user_throttles.0.id"),
					resource.TestCheckResourceAttrSet(byId, "policies.0.user_throttles.0.max_api_requests"),
					resource.TestCheckResourceAttrSet(byId, "policies.0.user_throttles.0.throttling_object_id"),
					resource.TestCheckResourceAttrSet(byId, "policies.0.user_throttles.0.throttling_object_name"),
					resource.TestCheckResourceAttr(byId, "policies.0.bind_num", "1"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_not_found_filter_useful", "true"),
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceThrottlingPolicies_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_apig_throttling_policies" "test" {
  depends_on = [
    huaweicloud_apig_throttling_policy_associate.test,
  ]
  instance_id = local.instance_id
}

# Filter by ID
locals {
  policy_id = huaweicloud_apig_throttling_policy.test.id
}

data "huaweicloud_apig_throttling_policies" "filter_by_id" {
  depends_on = [
    huaweicloud_apig_throttling_policy_associate.test,
  ]
  instance_id = local.instance_id
  policy_id   = local.policy_id
}

locals {
  id_filter_result = [
    for v in data.huaweicloud_apig_throttling_policies.filter_by_id.policies[*].id : v == local.policy_id
  ]
}

output "is_id_filter_useful" {
  value = length(local.id_filter_result) > 0 && alltrue(local.id_filter_result)
}

# Filter by name
locals {
  policy_name = huaweicloud_apig_throttling_policy.test.name
}

data "huaweicloud_apig_throttling_policies" "filter_by_name" {
  depends_on = [
    huaweicloud_apig_throttling_policy_associate.test,
  ]
  instance_id = local.instance_id
  name        = local.policy_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_apig_throttling_policies.filter_by_name.policies[*].name : v == local.policy_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by name (not found)
locals {
  not_found_name = "not_found"
}

data "huaweicloud_apig_throttling_policies" "filter_by_not_found_name" {
  depends_on = [
    huaweicloud_apig_throttling_policy_associate.test,
  ]
  instance_id = local.instance_id
  name        = local.not_found_name
}

locals {
  not_found_name_filter_result = [
    for v in data.huaweicloud_apig_throttling_policies.filter_by_not_found_name.policies[*].name : strcontains(v, local.not_found_name)
  ]
}

output "is_name_not_found_filter_useful" {
  value = length(local.not_found_name_filter_result) == 0
}

# Filter by type
locals {
  policy_type = huaweicloud_apig_throttling_policy.test.type
}

data "huaweicloud_apig_throttling_policies" "filter_by_type" {
  depends_on = [
    huaweicloud_apig_throttling_policy_associate.test,
  ]
  instance_id = local.instance_id
  type        = local.policy_type
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_apig_throttling_policies.filter_by_type.policies[*].type : v == local.policy_type
  ]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}
`, testAccDataSourceApiAssociatedThrottlingPolicies_base())
}

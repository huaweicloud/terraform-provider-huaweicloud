package organizations

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataPolicies_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_organizations_policies.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byBuildType   = "data.huaweicloud_organizations_policies.filter_by_build_type"
		dcByBuildType = acceptance.InitDataSourceCheck(byBuildType)

		byName   = "data.huaweicloud_organizations_policies.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byType   = "data.huaweicloud_organizations_policies.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataPolicies_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "policies.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by 'build_type' parameter.
					dcByBuildType.CheckResourceExists(),
					resource.TestCheckOutput("is_build_type_filter_useful", "true"),
					// Filter by 'name' parameter.
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(byName, "policies.0.id"),
					resource.TestCheckResourceAttrSet(byName, "policies.0.name"),
					resource.TestCheckResourceAttrSet(byName, "policies.0.type"),
					resource.TestCheckResourceAttrSet(byName, "policies.0.urn"),
					resource.TestCheckResourceAttrSet(byName, "policies.0.description"),
					resource.TestCheckResourceAttrSet(byName, "policies.0.build_type"),
					// Filter by 'type' parameter.
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataPolicies_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_organizations_policy" "test" {
  name        = "%[1]s"
  type        = "service_control_policy"
  description = "Created by terraform script"
  content = jsonencode({
    Version = "5.0",
    Statement = [
      {
        Effect = "Deny"
        Action = []
      }
    ]
  })

  tags = {
    "foo" = "bar"
  }
}
`, name)
}

func testAccDataPolicies_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameters.
data "huaweicloud_organizations_policies" "test" {
  depends_on = [huaweicloud_organizations_policy.test]
}

# Filter by 'build_type' parameter.
data "huaweicloud_organizations_policies" "filter_by_build_type" {
  build_type = "custom"
  depends_on = [huaweicloud_organizations_policy.test]
}

locals {
  build_type_filter_result = [for v in data.huaweicloud_organizations_policies.filter_by_build_type.policies :
  v.build_type == "custom"]
}

output "is_build_type_filter_useful" {
  value = length(local.build_type_filter_result) > 0 && alltrue(local.build_type_filter_result)
}

# Filter by 'name' parameter.
locals {
  policy_name = huaweicloud_organizations_policy.test.name
}

data "huaweicloud_organizations_policies" "filter_by_name" {
  name       = local.policy_name
  depends_on = [huaweicloud_organizations_policy.test]
}

locals {
  name_filter_result = [for v in data.huaweicloud_organizations_policies.filter_by_name.policies : v.name == local.policy_name]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by 'type' parameter.
locals {
  policy_type = huaweicloud_organizations_policy.test.type
}

data "huaweicloud_organizations_policies" "filter_by_type" {
  type       = local.policy_type
  depends_on = [huaweicloud_organizations_policy.test]
}

locals {
  type_filter_result = [for v in data.huaweicloud_organizations_policies.filter_by_type.policies : v.type == local.policy_type]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}
`, testAccDataPolicies_base(name))
}

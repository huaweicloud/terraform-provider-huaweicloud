package organizations

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourcePolicies_basic(t *testing.T) {
	rName := "data.huaweicloud_organizations_policies.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePolicies_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "policies.0.id"),
					resource.TestCheckResourceAttrSet(rName, "policies.0.name"),
					resource.TestCheckResourceAttrSet(rName, "policies.0.type"),
					resource.TestCheckResourceAttrSet(rName, "policies.0.urn"),
					resource.TestCheckResourceAttrSet(rName, "policies.0.description"),
					resource.TestCheckResourceAttrSet(rName, "policies.0.build_type"),

					resource.TestCheckOutput("build_type_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourcePolicies_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_organizations_policies" "test" {}

data "huaweicloud_organizations_policies" "name_filter" {
  name = "%[2]s"

  depends_on = [huaweicloud_organizations_policy.test]
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_organizations_policies.name_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_organizations_policies.name_filter.policies[*].name : v == "%[2]s"]
  )  
}

data "huaweicloud_organizations_policies" "build_type_filter" {
  build_type = "custom"

  depends_on = [data.huaweicloud_organizations_policies.name_filter]
}
output "build_type_filter_is_useful" {
  value = length(data.huaweicloud_organizations_policies.build_type_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_organizations_policies.build_type_filter.policies[*].build_type : v == "custom"]
  )  
}

data "huaweicloud_organizations_policies" "type_filter" {
  type = "service_control_policy"

  depends_on = [data.huaweicloud_organizations_policies.name_filter]
}
output "type_filter_is_useful" {
  value = length(data.huaweicloud_organizations_policies.type_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_organizations_policies.type_filter.policies[*].type : v == "service_control_policy"]
  )  
}
`, testPolicy_basic(name), name)
}

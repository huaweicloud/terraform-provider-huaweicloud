package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataEnterpriseProjectGroups_basic(t *testing.T) {
	all := "data.huaweicloud_identity_enterprise_project_groups.all"
	dc := acceptance.InitDataSourceCheck(all)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEnterpriseProjectGroups_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataEnterpriseProjectGroups_basic() string {
	return fmt.Sprintf(`
# All
locals {
  enterprise_project_id = "%[1]s"
}

data "huaweicloud_identity_enterprise_project_groups" "all" {
  enterprise_project_id = local.enterprise_project_id
}

locals {
  enterprise_project_id_filter_result = [
    for v in data.huaweicloud_identity_enterprise_project_groups.all.groups[*].enterprise_project_id : v == local.enterprise_project_id
  ]
}

output "is_enterprise_project_id_filter_useful" {
  value = length(local.enterprise_project_id_filter_result) > 0 && alltrue(local.enterprise_project_id_filter_result)
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

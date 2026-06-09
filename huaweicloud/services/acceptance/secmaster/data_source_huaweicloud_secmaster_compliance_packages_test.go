package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCompliancePackages_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_compliance_packages.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCompliancePackages_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "builtin_compliance_num"),
					resource.TestCheckResourceAttrSet(dataSource, "customized_compliance_num"),
					resource.TestCheckResourceAttrSet(dataSource, "disabled_compliance_num"),
					resource.TestCheckResourceAttrSet(dataSource, "enabled_compliance_num"),
					resource.TestCheckResourceAttrSet(dataSource, "compliance_packages.#"),
					resource.TestCheckResourceAttrSet(dataSource, "compliance_packages.0.uuid"),
					resource.TestCheckResourceAttrSet(dataSource, "compliance_packages.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "compliance_packages.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "compliance_packages.0.owner"),
					resource.TestCheckResourceAttrSet(dataSource, "compliance_packages.0.spec_catalog_vo_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "compliance_packages.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "compliance_packages.0.classify"),
					resource.TestCheckResourceAttrSet(dataSource, "compliance_packages.0.areas"),
					resource.TestCheckResourceAttrSet(dataSource, "compliance_packages.0.region"),
					resource.TestCheckResourceAttrSet(dataSource, "compliance_packages.0.state"),
					resource.TestCheckResourceAttrSet(dataSource, "compliance_packages.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "compliance_packages.0.check_items_num"),
					resource.TestCheckResourceAttrSet(dataSource, "compliance_packages.0.has_auto_check_items"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCompliancePackages_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_compliance_packages" "test" {
  workspace_id = "%[1]s"
}

locals {
  name = data.huaweicloud_secmaster_compliance_packages.test.compliance_packages[0].name
}

data "huaweicloud_secmaster_compliance_packages" "filter_by_name" {
  workspace_id = "%[1]s"
  name         = local.name
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_secmaster_compliance_packages.filter_by_name.compliance_packages) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_compliance_packages.filter_by_name.compliance_packages[*] : v.name == local.name]
  )
}

locals {
  type = data.huaweicloud_secmaster_compliance_packages.test.compliance_packages[0].type
}

data "huaweicloud_secmaster_compliance_packages" "filter_by_type" {
  workspace_id = "%[1]s"
  type         = local.type
}

output "is_type_filter_useful" {
  value = length(data.huaweicloud_secmaster_compliance_packages.filter_by_type.compliance_packages) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_compliance_packages.filter_by_type.compliance_packages[*] : v.type == local.type]
  )
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}

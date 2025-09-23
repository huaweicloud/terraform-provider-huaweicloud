package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running the test case, please ensure that there is at least one WAF instance in the current region.
func TestAccDataSourceReferenceTables_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		dataSourceName = "data.huaweicloud_waf_reference_tables.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byName   = "data.huaweicloud_waf_reference_tables.name_filter"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNonExist   = "data.huaweicloud_waf_reference_tables.non_exist_filter"
		dcByNonExist = acceptance.InitDataSourceCheck(byNonExist)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWafReferenceTables_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "tables.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tables.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tables.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tables.0.conditions.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tables.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tables.0.creation_time"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),

					dcByNonExist.CheckResourceExists(),
					resource.TestCheckOutput("non_exist_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccWafReferenceTables_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_waf_reference_tables" "test" {
  enterprise_project_id = "%[2]s"

  depends_on = [
    huaweicloud_waf_reference_table.test
  ]
}

# Filter by name
locals {
  name = data.huaweicloud_waf_reference_tables.test.tables.0.name
}

data "huaweicloud_waf_reference_tables" "name_filter" {
  name                  = local.name
  enterprise_project_id = "%[2]s"
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_waf_reference_tables.name_filter.tables[*].name : v == local.name
  ]
}

output "name_filter_is_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)  
}

# Filter by non-exist name
data "huaweicloud_waf_reference_tables" "non_exist_filter" {
  name                  = "non-exist"
  enterprise_project_id = "%[2]s"
}

locals {
  non_exist_filter_result = [
    for v in data.huaweicloud_waf_reference_tables.non_exist_filter.tables[*].name : v == local.name
  ]
}

output "non_exist_filter_is_useful" {
  value = length(local.non_exist_filter_result) == 0
}
`, testAccWafReferenceTable_basic(name), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

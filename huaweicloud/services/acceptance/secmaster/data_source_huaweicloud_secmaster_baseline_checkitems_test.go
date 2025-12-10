package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBaselineCheckitems_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_baseline_checkitems.test"
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
				Config: testDataSourceBaselineCheckitems_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "builtin_checkitem_num"),
					resource.TestCheckResourceAttrSet(dataSource, "checkitem_num"),
					resource.TestCheckResourceAttrSet(dataSource, "checkitems.#"),
					resource.TestCheckResourceAttrSet(dataSource, "checkitems.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "checkitems.0.level"),
					resource.TestCheckResourceAttrSet(dataSource, "checkitems.0.method"),
					resource.TestCheckResourceAttrSet(dataSource, "checkitems.0.source"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_condition_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceBaselineCheckitems_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_baseline_checkitems" "test" {
  workspace_id = "%[1]s"
}

locals {
  name = data.huaweicloud_secmaster_baseline_checkitems.test.checkitems[0].name
}

data "huaweicloud_secmaster_baseline_checkitems" "filter_by_name" {
  workspace_id = "%[1]s"
  name         = local.name
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_secmaster_baseline_checkitems.filter_by_name.checkitems) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_baseline_checkitems.filter_by_name.checkitems[*] : v.name == local.name]
  )
}

# condition parameter filtering is only used to test whether API calls are successful.
data "huaweicloud_secmaster_baseline_checkitems" "filter_by_condition" {
  workspace_id = "%[1]s"

  condition {
    conditions {
      name = "test"
      data = ["data_test1", "data_test"]
    }

    logics = ["value1", "value2"]

  }
}

output "is_condition_filter_useful" {
  value = length(data.huaweicloud_secmaster_baseline_checkitems.filter_by_name.checkitems) > 0
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}

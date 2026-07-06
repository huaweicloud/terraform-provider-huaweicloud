package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDataclassTypes_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_dataclass_types.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterDataClassID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataclassTypes_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.dataclass_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.category"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.category_code"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.sub_category"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.sub_category_code"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.is_built_in"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.update_time"),

					resource.TestCheckOutput("filter_by_category_code_is_useful", "true"),
					resource.TestCheckOutput("filter_by_enabled_is_useful", "true"),
					resource.TestCheckOutput("filter_by_is_built_in_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataclassTypes_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_dataclass_types" "test" {
  workspace_id = "%[1]s"
  dataclass_id = "%[2]s"
}

locals {
  dataclass_category_code = data.huaweicloud_secmaster_dataclass_types.test.data[0].category_code
}

data "huaweicloud_secmaster_dataclass_types" "filter_by_category_code" {
  workspace_id  = "%[1]s"
  dataclass_id  = "%[2]s"
  category_code = local.dataclass_category_code
}

output "filter_by_category_code_is_useful" {
  value = length(data.huaweicloud_secmaster_dataclass_types.filter_by_category_code.data) > 0 && alltrue(
    [for info in data.huaweicloud_secmaster_dataclass_types.filter_by_category_code.data :
      info.category_code == local.dataclass_category_code]
  )
}

locals {
  dataclass_enabled = data.huaweicloud_secmaster_dataclass_types.test.data[0].enabled
}

data "huaweicloud_secmaster_dataclass_types" "filter_by_enabled" {
  workspace_id = "%[1]s"
  dataclass_id = "%[2]s"
  enabled      = local.dataclass_enabled
}

output "filter_by_enabled_is_useful" {
  value = length(data.huaweicloud_secmaster_dataclass_types.filter_by_enabled.data) > 0 && alltrue(
    [for info in data.huaweicloud_secmaster_dataclass_types.filter_by_enabled.data :
      info.enabled == local.dataclass_enabled]
  )
}

locals {
  dataclass_is_built_in = data.huaweicloud_secmaster_dataclass_types.test.data[0].is_built_in
}

data "huaweicloud_secmaster_dataclass_types" "filter_by_is_built_in" {
  workspace_id = "%[1]s"
  dataclass_id = "%[2]s"
  is_built_in  = local.dataclass_is_built_in
}

output "filter_by_is_built_in_is_useful" {
  value = length(data.huaweicloud_secmaster_dataclass_types.filter_by_is_built_in.data) > 0 && alltrue(
    [for info in data.huaweicloud_secmaster_dataclass_types.filter_by_is_built_in.data :
      info.is_built_in == local.dataclass_is_built_in]
  )
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_DATACLASS_ID)
}

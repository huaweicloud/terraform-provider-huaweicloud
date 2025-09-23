package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecmasterDataClasses_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_data_classes.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceSecmasterDataClasses_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_classes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_classes.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_classes.0.business_code"),
					resource.TestCheckResourceAttrSet(dataSource, "data_classes.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "data_classes.0.is_built_in"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_business_code_filter_useful", "true"),
					resource.TestCheckOutput("is_description_filter_useful", "true"),
					resource.TestCheckOutput("is_built_in_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceSecmasterDataClasses_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_data_classes" "test" {
  workspace_id = "%[1]s"
}

locals {
  name          = data.huaweicloud_secmaster_data_classes.test.data_classes[0].name
  business_code = data.huaweicloud_secmaster_data_classes.test.data_classes[0].business_code
  description   = data.huaweicloud_secmaster_data_classes.test.data_classes[0].description
  is_built_in   = data.huaweicloud_secmaster_data_classes.test.data_classes[0].is_built_in
}

data "huaweicloud_secmaster_data_classes" "filter_by_name" {
  workspace_id = "%[1]s"
  name         = local.name
}

data "huaweicloud_secmaster_data_classes" "filter_by_business_code" {
  workspace_id  = "%[1]s"
  business_code = local.business_code
}

data "huaweicloud_secmaster_data_classes" "filter_by_description" {
  workspace_id = "%[1]s"
  description  = local.description
}

data "huaweicloud_secmaster_data_classes" "filter_by_is_built_in" {
  workspace_id = "%[1]s"
  is_built_in  = tostring(local.is_built_in)
}

locals {
  list_by_name          = data.huaweicloud_secmaster_data_classes.filter_by_name.data_classes
  list_by_business_code = data.huaweicloud_secmaster_data_classes.filter_by_business_code.data_classes
  list_by_description   = data.huaweicloud_secmaster_data_classes.filter_by_description.data_classes
  list_by_is_built_in   = data.huaweicloud_secmaster_data_classes.filter_by_is_built_in.data_classes
}

output "is_name_filter_useful" {
  value = length(local.list_by_name) > 0 && alltrue(
    [for v in local.list_by_name[*].name : strcontains(v, local.name)]
  )
}

output "is_business_code_filter_useful" {
  value = length(local.list_by_business_code) > 0 && alltrue(
    [for v in local.list_by_business_code[*].business_code : strcontains(v, local.business_code)]
  )
}

output "is_description_filter_useful" {
  value = length(local.list_by_description) > 0 && alltrue(
    [for v in local.list_by_description[*].description : strcontains(v, local.description)]
  )
}

output "is_built_in_filter_useful" {
  value = length(local.list_by_is_built_in) > 0 && alltrue(
    [for v in local.list_by_is_built_in[*].is_built_in : v == local.is_built_in]
  )
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}

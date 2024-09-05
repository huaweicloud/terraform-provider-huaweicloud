package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecmasterDataClassFields_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_data_class_fields.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceSecmasterDataClassFields_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "fields.#"),
					resource.TestCheckResourceAttrSet(dataSource, "fields.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "fields.0.is_built_in"),
					resource.TestCheckResourceAttrSet(dataSource, "fields.0.mapping"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_built_in_filter_useful", "true"),
					resource.TestCheckOutput("is_mapping_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceSecmasterDataClassFields_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_data_classes" "test" {
  workspace_id = "%[1]s"
}

locals {
  data_class_id = data.huaweicloud_secmaster_data_classes.test.data_classes[0].id
}

data "huaweicloud_secmaster_data_class_fields" "test" {
  workspace_id  = "%[1]s"
  data_class_id = local.data_class_id
}

locals {
  name        = data.huaweicloud_secmaster_data_class_fields.test.fields[0].name
  is_built_in = data.huaweicloud_secmaster_data_class_fields.test.fields[0].is_built_in
  mapping     = data.huaweicloud_secmaster_data_class_fields.test.fields[0].mapping
}

data "huaweicloud_secmaster_data_class_fields" "filter_by_name" {
  workspace_id  = "%[1]s"
  data_class_id = local.data_class_id
  name          = local.name
}

data "huaweicloud_secmaster_data_class_fields" "filter_by_is_built_in" {
  workspace_id  = "%[1]s"
  data_class_id = local.data_class_id
  is_built_in   = tostring(local.is_built_in)
}

data "huaweicloud_secmaster_data_class_fields" "filter_by_mapping" {
  workspace_id  = "%[1]s"
  data_class_id = local.data_class_id
  mapping       = tostring(local.mapping)
}

locals {
  list_by_name        = data.huaweicloud_secmaster_data_class_fields.filter_by_name.fields
  list_by_is_built_in = data.huaweicloud_secmaster_data_class_fields.filter_by_is_built_in.fields
  list_by_mapping     = data.huaweicloud_secmaster_data_class_fields.filter_by_mapping.fields
}

output "is_name_filter_useful" {
  value = length(local.list_by_name) > 0 && alltrue(
    [for v in local.list_by_name[*].name : strcontains(v, local.name)]
  )
}

output "is_built_in_filter_useful" {
  value = length(local.list_by_is_built_in) > 0 && alltrue(
    [for v in local.list_by_is_built_in[*].is_built_in : v == local.is_built_in]
  )
}

output "is_mapping_filter_useful" {
  value = length(local.list_by_mapping) > 0 && alltrue(
    [for v in local.list_by_mapping[*].mapping : v == local.mapping]
  )
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}

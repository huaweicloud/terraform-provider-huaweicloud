package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceLayoutFields_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_layout_fields.test"
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
				Config: testAccDataSourceLayoutFields_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "fields.#"),
					resource.TestCheckResourceAttrSet(dataSource, "fields.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "fields.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "fields.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "fields.0.creator_id"),
					resource.TestCheckResourceAttrSet(dataSource, "fields.0.creator_name"),
					resource.TestCheckResourceAttrSet(dataSource, "fields.0.field_key"),
					resource.TestCheckResourceAttrSet(dataSource, "fields.0.is_built_in"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_built_in_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceLayoutFields_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_layout_fields" "test" {
  workspace_id  = "%[1]s"
  business_code = "Incident"
}

locals {
  name = data.huaweicloud_secmaster_layout_fields.test.fields[0].name
}

data "huaweicloud_secmaster_layout_fields" "name_filter" {
  workspace_id  = "%[1]s"
  business_code = "Incident"
  name          = local.name
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_secmaster_layout_fields.name_filter.fields) > 0
}

data "huaweicloud_secmaster_layout_fields" "built_in_filter" {
  workspace_id  = "%[1]s"
  business_code = "Incident"
  is_built_in   = "false"
}

output "is_built_in_filter_useful" {
  value = length(data.huaweicloud_secmaster_layout_fields.built_in_filter.fields) > 0
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}

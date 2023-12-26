package dataarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceTemplateOptionalFields_basic(t *testing.T) {
	rName := "data.huaweicloud_dataarts_architecture_ds_template_optionals.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceTemplateOptionalFields_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "optional_fields.0.fd_name"),
					resource.TestCheckResourceAttrSet(rName, "optional_fields.0.description"),
					resource.TestCheckResourceAttrSet(rName, "optional_fields.0.description_en"),
					resource.TestCheckResourceAttrSet(rName, "optional_fields.0.required"),
					resource.TestCheckResourceAttrSet(rName, "optional_fields.0.searchable"),

					resource.TestCheckOutput("fd_name_filter_is_useful", "true"),
					resource.TestCheckOutput("required_filter_is_useful", "true"),
					resource.TestCheckOutput("searchable_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceTemplateOptionalFields_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dataarts_architecture_ds_template_optionals" "test" {
  workspace_id = "%[1]s"
}

data "huaweicloud_dataarts_architecture_ds_template_optionals" "fd_name_filter" {
  workspace_id = "%[1]s"
  fd_name      = data.huaweicloud_dataarts_architecture_ds_template_optionals.test.optional_fields.0.fd_name
}
locals{
  fd_name = data.huaweicloud_dataarts_architecture_ds_template_optionals.test.optional_fields.0.fd_name
}
output "fd_name_filter_is_useful" {
  value = length(data.huaweicloud_dataarts_architecture_ds_template_optionals.fd_name_filter.optional_fields) > 0 && alltrue(
    [for v in data.huaweicloud_dataarts_architecture_ds_template_optionals.fd_name_filter.optional_fields[*].fd_name : 
    v == local.fd_name]
  )  
}

data "huaweicloud_dataarts_architecture_ds_template_optionals" "required_filter" {
  workspace_id = "%[1]s"
  required     = data.huaweicloud_dataarts_architecture_ds_template_optionals.test.optional_fields.0.required
}
locals{
  required = data.huaweicloud_dataarts_architecture_ds_template_optionals.test.optional_fields.0.required
}
output "required_filter_is_useful" {
  value = length(data.huaweicloud_dataarts_architecture_ds_template_optionals.required_filter.optional_fields) > 0 && alltrue(
    [for v in data.huaweicloud_dataarts_architecture_ds_template_optionals.required_filter.optional_fields[*].required : 
    v == local.required]
  )  
}

data "huaweicloud_dataarts_architecture_ds_template_optionals" "searchable_filter" {
  workspace_id = "%[1]s"
  searchable   = data.huaweicloud_dataarts_architecture_ds_template_optionals.test.optional_fields.0.searchable
}
locals{
  searchable = data.huaweicloud_dataarts_architecture_ds_template_optionals.test.optional_fields.0.searchable
}
output "searchable_filter_is_useful" {
  value = length(data.huaweicloud_dataarts_architecture_ds_template_optionals.searchable_filter.optional_fields) > 0 && alltrue(
    [for v in data.huaweicloud_dataarts_architecture_ds_template_optionals.searchable_filter.optional_fields[*].searchable : 
    v == local.searchable]
  )  
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID)
}

package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecmasterWorkflows_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_workflows.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSecmasterWorkflows_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "workflows.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "workflows.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "workflows.0.data_class_id"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("desc_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("data_class_id_filter_is_useful", "true"),
					resource.TestCheckOutput("enabled_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSecmasterWorkflows_basic() string {
	randDescription := acceptance.RandomAccResourceNameWithDash()
	return fmt.Sprintf(`
data "huaweicloud_secmaster_workflows" "test" {
  workspace_id = "%[1]s"
}

locals {
  name          = data.huaweicloud_secmaster_workflows.test.workflows[0].name
  type          = data.huaweicloud_secmaster_workflows.test.workflows[0].type
  description   = data.huaweicloud_secmaster_workflows.test.workflows[0].description
  data_class_id = data.huaweicloud_secmaster_workflows.test.workflows[0].data_class_id
  enabled       = tostring(data.huaweicloud_secmaster_workflows.test.workflows[0].enabled)
}

data "huaweicloud_secmaster_workflows" "filter_by_name" {
  workspace_id = "%[1]s"
  name         = local.name
}

data "huaweicloud_secmaster_workflows" "filter_by_desc" {
  workspace_id = "%[1]s"
  description  = local.description
}

data "huaweicloud_secmaster_workflows" "filter_null_by_desc" {
  workspace_id = "%[1]s"
  description  = "%[2]s"
}

data "huaweicloud_secmaster_workflows" "filter_by_type" {
  workspace_id = "%[1]s"
  type         = local.type
}

data "huaweicloud_secmaster_workflows" "filter_by_data_class_id" {
  workspace_id  = "%[1]s"
  data_class_id = local.data_class_id
}

data "huaweicloud_secmaster_workflows" "filter_by_enabled" {
  workspace_id = "%[1]s"
  enabled      = local.enabled
}

locals {
  list_by_name          = data.huaweicloud_secmaster_workflows.filter_by_name.workflows
  list_by_desc          = data.huaweicloud_secmaster_workflows.filter_by_desc.workflows
  list_null_by_desc     = data.huaweicloud_secmaster_workflows.filter_null_by_desc.workflows
  list_by_type          = data.huaweicloud_secmaster_workflows.filter_by_type.workflows
  list_by_enabled       = data.huaweicloud_secmaster_workflows.filter_by_enabled.workflows
  list_by_data_class_id = data.huaweicloud_secmaster_workflows.filter_by_data_class_id.workflows
}

output "name_filter_is_useful" {
  value = length(local.list_by_name) > 0 && alltrue(
    [for v in local.list_by_name[*].name : v == local.name]
  )
}

output "desc_filter_is_useful" {
  value = length(local.list_by_desc) >= 1 && length(local.list_null_by_desc) == 0
}

output "type_filter_is_useful" {
  value = length(local.list_by_type) > 0 && alltrue(
    [for v in local.list_by_type[*].type : v == local.type]
  )
}

output "data_class_id_filter_is_useful" {
  value = length(local.list_by_data_class_id) > 0 && alltrue(
    [for v in local.list_by_data_class_id[*].data_class_id : v == local.data_class_id]
  )
}

output "enabled_filter_is_useful" {
  value = length(local.list_by_enabled) > 0 && alltrue(
    [for v in local.list_by_enabled[*].enabled : tostring(v) == local.enabled]
  )
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, randDescription)
}

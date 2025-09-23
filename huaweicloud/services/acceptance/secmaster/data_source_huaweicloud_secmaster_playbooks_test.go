package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecmasterPlaybooks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_playbooks.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSecmasterPlaybooks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "playbooks.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "playbooks.0.enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "playbooks.0.data_class_name"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("desc_filter_is_useful", "true"),
					resource.TestCheckOutput("data_class_name_filter_is_useful", "true"),
					resource.TestCheckOutput("enabled_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSecmasterPlaybooks_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_playbooks" "test" {
  workspace_id = "%[1]s"
}

locals {
  name            = data.huaweicloud_secmaster_playbooks.test.playbooks[0].name
  description     = data.huaweicloud_secmaster_playbooks.test.playbooks[0].description
  enabled         = tostring(data.huaweicloud_secmaster_playbooks.test.playbooks[0].enabled)
  data_class_name = data.huaweicloud_secmaster_playbooks.test.playbooks[0].data_class_name
}
	
data "huaweicloud_secmaster_playbooks" "filter_by_name" {
  workspace_id = "%[1]s"
  name         = local.name
}

data "huaweicloud_secmaster_playbooks" "filter_by_desc" {
  workspace_id = "%[1]s"
  description  = local.description
}

data "huaweicloud_secmaster_playbooks" "filter_by_enabled" {
  workspace_id = "%[1]s"
  enabled      = local.enabled
}

data "huaweicloud_secmaster_playbooks" "filter_by_data_class_name" {
  workspace_id    = "%[1]s"
  data_class_name = local.data_class_name
}
	
locals {
  list_by_name            = data.huaweicloud_secmaster_playbooks.filter_by_name.playbooks
  list_by_desc            = data.huaweicloud_secmaster_playbooks.filter_by_desc.playbooks
  list_by_enabled         = data.huaweicloud_secmaster_playbooks.filter_by_enabled.playbooks
  list_by_data_class_name = data.huaweicloud_secmaster_playbooks.filter_by_data_class_name.playbooks
}
	
output "name_filter_is_useful" {
  value = length(local.list_by_name) > 0 && alltrue(
    [for v in local.list_by_name[*].name : v == local.name]
  )
}
	
output "desc_filter_is_useful" {
  value = length(local.list_by_desc) > 0 && alltrue(
    [for v in local.list_by_desc[*].description : v == local.description]
  )
}
	
output "enabled_filter_is_useful" {
  value = length(local.list_by_enabled) > 0 && alltrue(
    [for v in local.list_by_enabled[*].enabled : tostring(v) == local.enabled]
  )
}
	
output "data_class_name_filter_is_useful" {
  value = length(local.list_by_data_class_name) > 0 && alltrue(
    [for v in local.list_by_data_class_name[*].data_class_name : v == local.data_class_name]
  )
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}

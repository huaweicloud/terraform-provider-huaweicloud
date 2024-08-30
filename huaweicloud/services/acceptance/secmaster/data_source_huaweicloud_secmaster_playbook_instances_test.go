package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecmasterPlaybookInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_playbook_instances.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSecmasterPlaybookInstances_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.trigger_type"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.playbook.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.data_class.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.data_class.0.name"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_trigger_type_filter_useful", "true"),
					resource.TestCheckOutput("is_data_class_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSecmasterPlaybookInstances_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_playbook_instances" "test" {
  workspace_id = "%[1]s"
}

locals {
  status          = data.huaweicloud_secmaster_playbook_instances.test.instances[0].status
  trigger_type    = data.huaweicloud_secmaster_playbook_instances.test.instances[0].trigger_type
  data_class_name = data.huaweicloud_secmaster_playbook_instances.test.instances[0].data_class[0].name
}

data "huaweicloud_secmaster_playbook_instances" "filter_by_status" {
  workspace_id = "%[1]s"
  status       = local.status
}

data "huaweicloud_secmaster_playbook_instances" "filter_by_trigger_type" {
  workspace_id = "%[1]s"
  trigger_type = local.trigger_type
}

data "huaweicloud_secmaster_playbook_instances" "filter_by_data_class_name" {
  workspace_id    = "%[1]s"
  data_class_name = local.data_class_name
}

locals {
  list_by_status          = data.huaweicloud_secmaster_playbook_instances.filter_by_status.instances
  list_by_trigger_type    = data.huaweicloud_secmaster_playbook_instances.filter_by_trigger_type.instances
  list_by_data_class_name = data.huaweicloud_secmaster_playbook_instances.filter_by_data_class_name.instances
}

output "is_status_filter_useful" {
  value = length(local.list_by_status) > 0 && alltrue(
    [for v in local.list_by_status[*].status : v == local.status]
  )
}

output "is_trigger_type_filter_useful" {
  value = length(local.list_by_trigger_type) > 0 && alltrue(
    [for v in local.list_by_trigger_type[*].trigger_type : v == local.trigger_type]
  )
}

output "is_data_class_name_filter_useful" {
  value = length(local.list_by_data_class_name) > 0 && alltrue(
    [for v in local.list_by_data_class_name[*].data_class[0].name : v == local.data_class_name]
  )
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}

package secmaster

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func TestAccDataSourceWorkflowInstances_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_workflow_instances.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
		startTime  = utils.FormatTimeStampRFC3339(time.Now().Unix()-24*60*60, false, "2006-01-02T15:04:05.000Z+0800")
		endTime    = utils.FormatTimeStampRFC3339(time.Now().Unix(), false, "2006-01-02T15:04:05.000Z+0800")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceWorkflowInstances_basic(startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instance.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instance.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instance.0.workflow.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instance.0.workflow.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance.0.workflow.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instance.0.workflow.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "instance.0.dataclass.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instance.0.dataclass.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance.0.dataclass.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instance.0.playbook.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instance.0.playbook.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance.0.playbook.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instance.0.trigger_type"),
					resource.TestCheckResourceAttrSet(dataSource, "instance.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "instance.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "instance.0.end_time"),

					resource.TestCheckOutput("is_workflow_id_filter_useful", "true"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_dataclass_id_filter_useful", "true"),
					resource.TestCheckOutput("is_playbook_id_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_trigger_type_filter_useful", "true"),
					resource.TestCheckOutput("is_time_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceWorkflowInstances_basic(startTime, endTime string) string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_workflow_instances" "test" {
  workspace_id  = "%[1]s"
}

locals {
  workflow_id = data.huaweicloud_secmaster_workflow_instances.test.instance[0].workflow[0].id
}

data "huaweicloud_secmaster_workflow_instances" "workflow_id_filter" {
  workspace_id = "%[1]s"
  workflow_id  = local.workflow_id
}

output "is_workflow_id_filter_useful" {
  value = length(data.huaweicloud_secmaster_workflow_instances.workflow_id_filter.instance) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_workflow_instances.workflow_id_filter.instance[*].workflow[0].id : v == local.workflow_id]
  )
}

locals {
  id = data.huaweicloud_secmaster_workflow_instances.test.instance[0].id
}

data "huaweicloud_secmaster_workflow_instances" "id_filter" {
  workspace_id = "%[1]s"
  id           = local.id
}

output "is_id_filter_useful" {
  value = length(data.huaweicloud_secmaster_workflow_instances.id_filter.instance) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_workflow_instances.id_filter.instance[*].id : v == local.id]
  )
}

locals {
  name = data.huaweicloud_secmaster_workflow_instances.test.instance[0].name
}


data "huaweicloud_secmaster_workflow_instances" "name_filter" {
  workspace_id = "%[1]s"
  name         = local.name
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_secmaster_workflow_instances.name_filter.instance) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_workflow_instances.name_filter.instance[*].name : v == local.name]
  )
}

locals {
  dataclass_id = data.huaweicloud_secmaster_workflow_instances.test.instance[0].dataclass[0].id
}

data "huaweicloud_secmaster_workflow_instances" "dataclass_id_filter" {
  workspace_id = "%[1]s"
  dataclass_id = local.dataclass_id
}

output "is_dataclass_id_filter_useful" {
  value = length(data.huaweicloud_secmaster_workflow_instances.dataclass_id_filter.instance) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_workflow_instances.dataclass_id_filter.instance[*].dataclass[0].id : v == local.dataclass_id]
  )
}

locals {
  playbook_id = data.huaweicloud_secmaster_workflow_instances.test.instance[0].playbook[0].id
}

data "huaweicloud_secmaster_workflow_instances" "playbook_id_filter" {
  workspace_id = "%[1]s"
  playbook_id  = local.playbook_id
}

output "is_playbook_id_filter_useful" {
  value = length(data.huaweicloud_secmaster_workflow_instances.playbook_id_filter.instance) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_workflow_instances.playbook_id_filter.instance[*].playbook[0].id : v == local.playbook_id]
  )
}

locals {
  status = data.huaweicloud_secmaster_workflow_instances.test.instance[0].status
}

data "huaweicloud_secmaster_workflow_instances" "status_filter" {
  workspace_id  = "%[1]s"
  status        = local.status
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_secmaster_workflow_instances.status_filter.instance) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_workflow_instances.status_filter.instance[*].status : v == local.status]
  )
}

locals {
  trigger_type = data.huaweicloud_secmaster_workflow_instances.test.instance[0].trigger_type
}

data "huaweicloud_secmaster_workflow_instances" "trigger_type_filter" {
  workspace_id = "%[1]s"
  trigger_type = local.trigger_type
}

output "is_trigger_type_filter_useful" {
  value = length(data.huaweicloud_secmaster_workflow_instances.trigger_type_filter.instance) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_workflow_instances.trigger_type_filter.instance[*].trigger_type : v == local.trigger_type]
  )
}

data "huaweicloud_secmaster_workflow_instances" "time_filter" {
  workspace_id = "%[1]s"
  from_date    = "%[2]s"
  to_date      = "%[3]s"
}

output "is_time_filter_useful" {
  value = length(data.huaweicloud_secmaster_workflow_instances.time_filter.instance) > 0
}

`, acceptance.HW_SECMASTER_WORKSPACE_ID, startTime, endTime)
}

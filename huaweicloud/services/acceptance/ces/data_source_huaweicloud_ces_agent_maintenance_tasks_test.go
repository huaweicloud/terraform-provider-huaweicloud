package ces

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCesAgentMaintenanceTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ces_agent_maintenance_tasks.test"
	dc := acceptance.InitDataSourceCheck(dataSource)
	basicConfig := testResourceAgentMaintenanceTask_base()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCesAgentMaintenanceTasks_basic(basicConfig),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "invocations.0.invocation_id"),
					resource.TestCheckResourceAttrSet(dataSource, "invocations.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "invocations.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "invocations.0.instance_type"),
					resource.TestCheckResourceAttrSet(dataSource, "invocations.0.invocation_type"),
					resource.TestCheckResourceAttrSet(dataSource, "invocations.0.invocation_status"),
					resource.TestCheckResourceAttrSet(dataSource, "invocations.0.invocation_target"),
					resource.TestCheckResourceAttrSet(dataSource, "invocations.0.target_version"),
					resource.TestCheckResourceAttrSet(dataSource, "invocations.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "invocations.0.update_time"),
					resource.TestCheckOutput("is_default_filter_useful", "true"),
					resource.TestCheckOutput("is_instance_id_filter_useful", "true"),
					resource.TestCheckOutput("is_invocation_id_filter_useful", "true"),
					resource.TestCheckOutput("is_invocation_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCesAgentMaintenanceTasks_basic(config string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_ces_agent_maintenance_tasks" "test" {
  depends_on = [huaweicloud_ces_agent_maintenance_task.test]
}

output "is_default_filter_useful" {
  value = length(data.huaweicloud_ces_agent_maintenance_tasks.test.invocations) >= 1 
}

locals {
  instance_id     = huaweicloud_compute_instance.test.id
  invocation_id   = huaweicloud_ces_agent_maintenance_task.test.id
  invocation_type = huaweicloud_ces_agent_maintenance_task.test.invocation_type
}

data "huaweicloud_ces_agent_maintenance_tasks" "filter_by_instance_id" {
  instance_id = local.instance_id
 
  depends_on = [huaweicloud_ces_agent_maintenance_task.test]
}

output "is_instance_id_filter_useful" {
  value = length(data.huaweicloud_ces_agent_maintenance_tasks.filter_by_instance_id.invocations) >= 1 && alltrue([
    for i in data.huaweicloud_ces_agent_maintenance_tasks.filter_by_instance_id.invocations[*] : i.instance_id == local.instance_id
  ])
}

data "huaweicloud_ces_agent_maintenance_tasks" "filter_by_invocation_id" {
  invocation_id = local.invocation_id
 
  depends_on = [huaweicloud_ces_agent_maintenance_task.test]
}

output "is_invocation_id_filter_useful" {
  value = length(data.huaweicloud_ces_agent_maintenance_tasks.filter_by_invocation_id.invocations) >= 1 && alltrue([
    for i in data.huaweicloud_ces_agent_maintenance_tasks.filter_by_invocation_id.invocations[*] : i.invocation_id == local.invocation_id
  ])
}

data "huaweicloud_ces_agent_maintenance_tasks" "filter_by_invocation_type" {
  invocation_type = local.invocation_type
 
  depends_on = [huaweicloud_ces_agent_maintenance_task.test]
}

output "is_invocation_type_filter_useful" {
  value = length(data.huaweicloud_ces_agent_maintenance_tasks.filter_by_invocation_type.invocations) >= 1 && alltrue([
    for i in data.huaweicloud_ces_agent_maintenance_tasks.filter_by_invocation_type.invocations[*] : i.invocation_type == local.invocation_type
  ])
}
`, testResourceCesAgentMaintenanceTask_basic(config))
}

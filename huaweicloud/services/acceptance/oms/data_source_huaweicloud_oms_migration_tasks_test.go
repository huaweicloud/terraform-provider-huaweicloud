package oms

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceMigrationTasks_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_oms_migration_tasks.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byGroupId   = "data.huaweicloud_oms_migration_tasks.group_id_filter"
		dcByGroupId = acceptance.InitDataSourceCheck(byGroupId)

		byStatus   = "data.huaweicloud_oms_migration_tasks.status_filter"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOmsInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMigrationTasks_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "tasks.#", regexp.MustCompile(`^[1-9][0-9]*$`)),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.src_node.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.src_node.0.bucket"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.src_node.0.cloud_type"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.src_node.0.region"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.dst_node.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.dst_node.0.bucket"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.dst_node.0.region"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.task_type"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.group_type"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.enable_metadata_migration"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.enable_restore"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.object_overwrite_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.consistency_check"),

					dcByGroupId.CheckResourceExists(),
					resource.TestCheckOutput("group_id_filter_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceMigrationTasks_basic = `
data "huaweicloud_oms_migration_tasks" "test" {}

locals {
  group_id = data.huaweicloud_oms_migration_tasks.test.tasks[0].group_id
}

data "huaweicloud_oms_migration_tasks" "group_id_filter" {
  group_id = local.group_id
}

output "group_id_filter_useful" {
  value = length(data.huaweicloud_oms_migration_tasks.group_id_filter.tasks) > 0 && alltrue(
    [for v in data.huaweicloud_oms_migration_tasks.group_id_filter.tasks[*].group_id : v == local.group_id]
  )
}

locals {
  status = data.huaweicloud_oms_migration_tasks.test.tasks[0].status
}

data "huaweicloud_oms_migration_tasks" "status_filter" {
  status = local.status
}

output "status_filter_useful" {
  value = length(data.huaweicloud_oms_migration_tasks.status_filter.tasks) > 0 && alltrue(
    [for v in data.huaweicloud_oms_migration_tasks.status_filter.tasks[*].status : v == local.status]
  )
}
`

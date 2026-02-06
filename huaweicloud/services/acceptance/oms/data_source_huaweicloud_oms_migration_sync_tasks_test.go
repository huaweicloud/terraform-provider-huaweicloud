package oms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceMigrationSyncTasks_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_oms_migration_sync_tasks.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOmsInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMigrationSyncTasks_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.sync_task_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.src_cloud_type"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.src_region"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.src_bucket"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.dst_bucket"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.dst_region"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.object_overwrite_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.consistency_check"),

					resource.TestCheckOutput("status_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceMigrationSyncTasks_basic = `
data "huaweicloud_oms_migration_sync_tasks" "test" {}

locals {
  status = data.huaweicloud_oms_migration_sync_tasks.test.tasks[0].status
}

data "huaweicloud_oms_migration_sync_tasks" "status_filter" {
  status = local.status
}

output "status_filter_useful" {
  value = length(data.huaweicloud_oms_migration_sync_tasks.status_filter.tasks) > 0 && alltrue(
    [for v in data.huaweicloud_oms_migration_sync_tasks.status_filter.tasks[*].status : v == local.status]
  )
}
`

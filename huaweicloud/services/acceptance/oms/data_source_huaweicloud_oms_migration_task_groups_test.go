package oms

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceMigrationTaskGroups_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_oms_migration_task_groups.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byStatus   = "data.huaweicloud_oms_migration_task_groups.status_filter"
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
				Config: testAccDataSourceMigrationTaskGroups_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "taskgroups.#", regexp.MustCompile(`^[1-9][0-9]*$`)),
					resource.TestCheckResourceAttrSet(dataSource, "taskgroups.0.group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "taskgroups.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "taskgroups.0.src_node.#"),
					resource.TestCheckResourceAttrSet(dataSource, "taskgroups.0.src_node.0.bucket"),
					resource.TestCheckResourceAttrSet(dataSource, "taskgroups.0.src_node.0.cloud_type"),
					resource.TestCheckResourceAttrSet(dataSource, "taskgroups.0.src_node.0.region"),
					resource.TestCheckResourceAttrSet(dataSource, "taskgroups.0.dst_node.#"),
					resource.TestCheckResourceAttrSet(dataSource, "taskgroups.0.dst_node.0.bucket"),
					resource.TestCheckResourceAttrSet(dataSource, "taskgroups.0.dst_node.0.region"),
					resource.TestCheckResourceAttrSet(dataSource, "taskgroups.0.task_type"),
					resource.TestCheckResourceAttrSet(dataSource, "taskgroups.0.enable_metadata_migration"),
					resource.TestCheckResourceAttrSet(dataSource, "taskgroups.0.enable_restore"),
					resource.TestCheckResourceAttrSet(dataSource, "taskgroups.0.object_overwrite_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "taskgroups.0.consistency_check"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceMigrationTaskGroups_basic = `
data "huaweicloud_oms_migration_task_groups" "test" {}

locals {
  status = data.huaweicloud_oms_migration_task_groups.test.taskgroups[0].status
}

data "huaweicloud_oms_migration_task_groups" "status_filter" {
  status = local.status
}

output "status_filter_useful" {
  value = length(data.huaweicloud_oms_migration_task_groups.status_filter.taskgroups) > 0 && alltrue(
    [for v in data.huaweicloud_oms_migration_task_groups.status_filter.taskgroups[*].status : v == local.status]
  )
}
`

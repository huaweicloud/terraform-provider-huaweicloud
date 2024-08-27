package sfsturbo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDataTasks_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_sfs_turbo_data_tasks.test"
		rName      = acceptance.RandomAccResourceName()
		randInt    = acctest.RandInt()
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byType   = "data.huaweicloud_sfs_turbo_data_tasks.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byStatus   = "data.huaweicloud_sfs_turbo_data_tasks.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBSEndpoint(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataTasks_basic(rName, randInt),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.src_target"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.dest_target"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.src_prefix"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.dest_prefix"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.end_time"),

					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("type_filter_is_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataTasks_basic(name string, randInt int) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_sfs_turbo_data_tasks" "test" {
  depends_on = [
    huaweicloud_sfs_turbo_data_task.test
  ]

  share_id = huaweicloud_sfs_turbo.test.id
}

locals {
  type = data.huaweicloud_sfs_turbo_data_tasks.test.tasks[0].type
}

data "huaweicloud_sfs_turbo_data_tasks" "filter_by_type" {
  share_id = huaweicloud_sfs_turbo.test.id
  type     = local.type
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_sfs_turbo_data_tasks.filter_by_type.tasks[*].type : v == local.type
  ]
}

output "type_filter_is_useful" {
  value = alltrue(local.type_filter_result) && length(local.type_filter_result) > 0
}

locals {
  status = data.huaweicloud_sfs_turbo_data_tasks.test.tasks[0].status
}

data "huaweicloud_sfs_turbo_data_tasks" "filter_by_status" {
  share_id = huaweicloud_sfs_turbo.test.id
  status   = local.status
}

locals {
  status_filter_result = [ 
    for v in data.huaweicloud_sfs_turbo_data_tasks.filter_by_status.tasks[*].status : v == local.status
  ]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}
`, testAccDataTask_basic(name, randInt))
}

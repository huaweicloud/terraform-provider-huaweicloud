package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDdsInstantTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dds_instant_tasks.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSTimeRange(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDdsInstantTasks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.ended_at"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDdsInstantTasks_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dds_instant_tasks" "test" {
  start_time = "%[1]s"
  end_time   = "%[2]s"
}

locals {
  test_results = data.huaweicloud_dds_instant_tasks.test
}

// filter by name
data "huaweicloud_dds_instant_tasks" "filter_by_name" {
  start_time = "%[1]s"
  end_time   = "%[2]s"
  name       = local.test_results.jobs[0].name
}

locals {
  filter_result_by_name = [for v in data.huaweicloud_dds_instant_tasks.filter_by_name.jobs[*].name : 
    v == local.test_results.jobs[0].name]
}

output "is_name_filter_useful" {
  value = length(local.filter_result_by_name) > 0 && alltrue(local.filter_result_by_name) 
}

// filter by status
data "huaweicloud_dds_instant_tasks" "filter_by_status" {
  start_time = "%[1]s"
  end_time   = "%[2]s"
  status     = local.test_results.jobs[0].status
}

locals {
  filter_result_by_status = [for v in data.huaweicloud_dds_instant_tasks.filter_by_status.jobs[*].status : 
    v == local.test_results.jobs[0].status]
}

output "is_status_filter_useful" {
  value = length(local.filter_result_by_status) > 0 && alltrue(local.filter_result_by_status) 
}
`, acceptance.HW_DDS_START_TIME, acceptance.HW_DDS_END_TIME)
}

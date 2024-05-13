package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcsDiagnosisTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dcs_diagnosis_tasks.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDCSInstanceID(t)
			acceptance.TestAccPreCheckDcsTimeRange(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDcsDiagnosisTasks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "diagnosis_tasks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "diagnosis_tasks.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "diagnosis_tasks.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "diagnosis_tasks.0.begin_time"),
					resource.TestCheckResourceAttrSet(dataSource, "diagnosis_tasks.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "diagnosis_tasks.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "diagnosis_tasks.0.node_num"),
					resource.TestCheckResourceAttrSet(dataSource, "diagnosis_tasks.0.abnormal_item_sum"),
					resource.TestCheckResourceAttrSet(dataSource, "diagnosis_tasks.0.failed_item_sum"),

					resource.TestCheckOutput("task_id_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("begin_time_filter_is_useful", "true"),
					resource.TestCheckOutput("end_time_filter_is_useful", "true"),
					resource.TestCheckOutput("node_num_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDcsDiagnosisTasks_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dcs_diagnosis_tasks" "test" {
  depends_on = [huaweicloud_dcs_diagnosis_task.test]

  instance_id = "%[2]s"
}

locals{
  task_id = huaweicloud_dcs_diagnosis_task.test.id
}
data "huaweicloud_dcs_diagnosis_tasks" "task_id_filter" {
  depends_on = [huaweicloud_dcs_diagnosis_task.test]

  instance_id = "%[2]s"
  task_id     = huaweicloud_dcs_diagnosis_task.test.id
}
output "task_id_filter_is_useful" {
  value = length(data.huaweicloud_dcs_diagnosis_tasks.task_id_filter.diagnosis_tasks) > 0 && alltrue(
  [for v in data.huaweicloud_dcs_diagnosis_tasks.task_id_filter.diagnosis_tasks[*].id : v == local.task_id]  
  )
}

locals{
  status = "finished"
}
data "huaweicloud_dcs_diagnosis_tasks" "status_filter" {
  depends_on = [huaweicloud_dcs_diagnosis_task.test]

  instance_id = "%[2]s"
  status      = "finished"
}
output "status_filter_is_useful" {
  value = length(data.huaweicloud_dcs_diagnosis_tasks.status_filter.diagnosis_tasks) > 0 && alltrue(
  [for v in data.huaweicloud_dcs_diagnosis_tasks.status_filter.diagnosis_tasks[*].status : v == local.status]  
  )
}

locals{
  begin_time = "%[3]s"
}
data "huaweicloud_dcs_diagnosis_tasks" "begin_time_filter" {
  depends_on = [huaweicloud_dcs_diagnosis_task.test]

  instance_id = "%[2]s"
  begin_time  = "%[3]s"
}
output "begin_time_filter_is_useful" {
  value = length(data.huaweicloud_dcs_diagnosis_tasks.begin_time_filter.diagnosis_tasks) > 0 && alltrue(
  [for v in data.huaweicloud_dcs_diagnosis_tasks.begin_time_filter.diagnosis_tasks[*].begin_time : v == local.begin_time]  
  )
}

locals{
  end_time = "%[4]s"
}
data "huaweicloud_dcs_diagnosis_tasks" "end_time_filter" {
  depends_on = [huaweicloud_dcs_diagnosis_task.test]

  instance_id = "%[2]s"
  end_time    = "%[4]s"
}
output "end_time_filter_is_useful" {
  value = length(data.huaweicloud_dcs_diagnosis_tasks.end_time_filter.diagnosis_tasks) > 0 && alltrue(
  [for v in data.huaweicloud_dcs_diagnosis_tasks.end_time_filter.diagnosis_tasks[*].end_time : v == local.end_time]  
  )
}

locals{
  node_num = length(huaweicloud_dcs_diagnosis_task.test.node_ip_list)
}
data "huaweicloud_dcs_diagnosis_tasks" "node_num_filter" {
  depends_on = [huaweicloud_dcs_diagnosis_task.test]

  instance_id = "%[2]s"
  node_num    = length(huaweicloud_dcs_diagnosis_task.test.node_ip_list)
}
output "node_num_filter_is_useful" {
  value = length(data.huaweicloud_dcs_diagnosis_tasks.node_num_filter.diagnosis_tasks) > 0 && alltrue(
  [for v in data.huaweicloud_dcs_diagnosis_tasks.node_num_filter.diagnosis_tasks[*].node_num : v == local.node_num]  
  )
}
`, testAccDiagnosisTask_basic(), acceptance.HW_DCS_INSTANCE_ID, acceptance.HW_DCS_BEGIN_TIME, acceptance.HW_DCS_END_TIME)
}

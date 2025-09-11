package coc

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocDiagnosisTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_diagnosis_tasks.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocDiagnosisTasks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.user_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.progress"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.instance_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.region"),
					resource.TestCheckOutput("task_id_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("region_filter_is_useful", "true"),
					resource.TestCheckOutput("creator_filter_is_useful", "true"),
					resource.TestCheckOutput("start_time_filter_is_useful", "true"),
					resource.TestCheckOutput("end_time_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocDiagnosisTasks_basic() string {
	currentTime := time.Now()
	tenMinutesAgo := currentTime.Add(-10*time.Minute).Unix() * 1e3
	tenMinutesLater := currentTime.Add(10*time.Minute).Unix() * 1e3
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_coc_diagnosis_tasks" "test" {
  task_id = huaweicloud_coc_diagnosis_task.test.id
}

output "task_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_diagnosis_tasks.test.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_diagnosis_tasks.test.data[*].id : v == huaweicloud_coc_diagnosis_task.test.id]
  )
}

data "huaweicloud_coc_diagnosis_tasks" "type_filter" {
  type = huaweicloud_coc_diagnosis_task.test.type
}

output "type_filter_is_useful" {
  value = length(data.huaweicloud_coc_diagnosis_tasks.type_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_diagnosis_tasks.type_filter.data[*].type :
      v == huaweicloud_coc_diagnosis_task.test.type]
  )
}

data "huaweicloud_coc_diagnosis_tasks" "status_filter" {
  status = huaweicloud_coc_diagnosis_task.test.status
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_coc_diagnosis_tasks.status_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_diagnosis_tasks.status_filter.data[*].status :
      v == huaweicloud_coc_diagnosis_task.test.status]
  )
}

data "huaweicloud_coc_diagnosis_tasks" "region_filter" {
  region = huaweicloud_coc_diagnosis_task.test.region
}

output "region_filter_is_useful" {
  value = length(data.huaweicloud_coc_diagnosis_tasks.region_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_diagnosis_tasks.region_filter.data[*].region :
      v == huaweicloud_coc_diagnosis_task.test.region]
  )
}

data "huaweicloud_coc_diagnosis_tasks" "creator_filter" {
  creator = huaweicloud_coc_diagnosis_task.test.user_name
}

output "creator_filter_is_useful" {
  value = length(data.huaweicloud_coc_diagnosis_tasks.creator_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_diagnosis_tasks.creator_filter.data[*].user_name :
      v == huaweicloud_coc_diagnosis_task.test.user_name]
  )
}

data "huaweicloud_coc_diagnosis_tasks" "start_time_filter" {
  start_time = %[2]v
}

output "start_time_filter_is_useful" {
  value = length(data.huaweicloud_coc_diagnosis_tasks.start_time_filter.data) > 0
}

data "huaweicloud_coc_diagnosis_tasks" "end_time_filter" {
  end_time = %[3]v
}

output "end_time_filter_is_useful" {
  value = length(data.huaweicloud_coc_diagnosis_tasks.end_time_filter.data) > 0
}
`, testCocDiagnosisTask_basic(), tenMinutesAgo, tenMinutesLater)
}

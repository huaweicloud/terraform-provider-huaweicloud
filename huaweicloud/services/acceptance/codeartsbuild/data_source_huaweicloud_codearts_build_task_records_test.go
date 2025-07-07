package codeartsbuild

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCodeArtsBuildTaskRecords_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_build_task_records.test"
	dc := acceptance.InitDataSourceCheck(dataSource)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testBuildTaskAction_basic(name),
				// wait seconds for delay
				Check: func(_ *terraform.State) error {
					// lintignore:R018
					time.Sleep(10 * time.Second)
					return nil
				},
			},
			{
				Config: testAccDataSourceCodeArtsBuildTaskRecords_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.#"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.status_code"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.schedule_time"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.queued_time"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.finish_time"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.duration"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.build_duration"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.pending_duration"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.display_name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.trigger_name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.execution_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.repository"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.build_yml_path"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.build_yml_url"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.daily_build_number"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.trigger_type"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.scm_type"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.scm_web_url"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.user_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.build_no"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.daily_build_no"),
					resource.TestCheckOutput("is_branches_filter_useful", "true"),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
					resource.TestCheckOutput("is_date_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceCodeArtsBuildTaskRecords_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_build_tasks" "test" {
  depends_on = [huaweicloud_codearts_build_task_action.stop]

  project_id = huaweicloud_codearts_project.test.id
}

data "huaweicloud_codearts_build_task_records" "test" {
  build_project_id = data.huaweicloud_codearts_build_tasks.test.tasks[0].build_project_id
}

// filter by triggers
data "huaweicloud_codearts_build_task_records" "filter_by_trigger" {
  build_project_id = data.huaweicloud_codearts_build_tasks.test.tasks[0].build_project_id
  triggers         = [data.huaweicloud_codearts_build_task_records.test.records[0].trigger_name]
}

output "is_triggers_filter_useful" {
  value = length(data.huaweicloud_codearts_build_task_records.filter_by_trigger.records) > 0
}

// filter by branches
data "huaweicloud_codearts_build_task_records" "filter_by_branches" {
  build_project_id = data.huaweicloud_codearts_build_tasks.test.tasks[0].build_project_id
  branches         = ["master"]
}

output "is_branches_filter_useful" {
  value = length(data.huaweicloud_codearts_build_task_records.filter_by_branches.records) > 0
}

// filter by tags
data "huaweicloud_codearts_build_task_records" "filter_by_tags" {
  build_project_id = data.huaweicloud_codearts_build_tasks.test.tasks[0].build_project_id
  tags             = ["test"]
}

output "is_tags_filter_useful" {
  value = length(data.huaweicloud_codearts_build_task_records.filter_by_tags.records) > 0
}

// filter by date
data "huaweicloud_codearts_build_task_records" "filter_by_date" {
  build_project_id = data.huaweicloud_codearts_build_tasks.test.tasks[0].build_project_id
  from_date        = "2025-01-01 00:00:00"
  to_date          = "2025-12-31 23:59:59"
}

output "is_date_filter_useful" {
  value = length(data.huaweicloud_codearts_build_task_records.filter_by_date.records) > 0
}
`, testBuildTaskAction_basic(name))
}

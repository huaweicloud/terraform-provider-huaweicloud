package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataQualityConsistencyTasks_basic(t *testing.T) {
	var (
		// Without any filter parameters.
		all = "data.huaweicloud_dataarts_quality_consistency_tasks.all"
		dc  = acceptance.InitDataSourceCheck(all)

		// Filter by name.
		byName   = "data.huaweicloud_dataarts_quality_consistency_tasks.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		// Filter by not found name.
		byNotFoundName   = "data.huaweicloud_dataarts_quality_consistency_tasks.filter_by_not_found_name"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)

		// Filter by category ID.
		byCategoryId   = "data.huaweicloud_dataarts_quality_consistency_tasks.filter_by_category_id"
		dcByCategoryId = acceptance.InitDataSourceCheck(byCategoryId)

		// Filter by creator.
		byCreator   = "data.huaweicloud_dataarts_quality_consistency_tasks.filter_by_creator"
		dcByCreator = acceptance.InitDataSourceCheck(byCreator)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataQualityConsistencyTasks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "tasks.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(all, "tasks.0.id"),
					resource.TestCheckResourceAttrSet(all, "tasks.0.name"),
					resource.TestCheckResourceAttrSet(all, "tasks.0.category_id"),
					resource.TestCheckResourceAttrSet(all, "tasks.0.schedule_status"),
					resource.TestCheckResourceAttrSet(all, "tasks.0.create_time"),
					resource.TestCheckResourceAttrSet(all, "tasks.0.creator"),
					resource.TestMatchResourceAttr(all, "tasks.0.create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),

					// Filter by name.
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),

					// Filter by not found name.
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found_name_filter_useful", "true"),

					// Filter by category ID.
					dcByCategoryId.CheckResourceExists(),
					resource.TestCheckOutput("is_category_id_filter_useful", "true"),

					// Filter by creator.
					dcByCreator.CheckResourceExists(),
					resource.TestCheckOutput("is_creator_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataQualityConsistencyTasks_basic() string {
	return fmt.Sprintf(`
// Without any filter parameters.
data "huaweicloud_dataarts_quality_consistency_tasks" "all" {
  workspace_id = "%[1]s"
}

// Filter by name.
locals {
  task_name = data.huaweicloud_dataarts_quality_consistency_tasks.all.tasks[0].name
}

data "huaweicloud_dataarts_quality_consistency_tasks" "filter_by_name" {
  workspace_id = "%[1]s"

  name = local.task_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dataarts_quality_consistency_tasks.filter_by_name.tasks[*].name : v == local.task_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

// Filter by not found name.
data "huaweicloud_dataarts_quality_consistency_tasks" "filter_by_not_found_name" {
  workspace_id = "%[1]s"

  name = "not_found_task"
}

output "is_not_found_name_filter_useful" {
  value = length(data.huaweicloud_dataarts_quality_consistency_tasks.filter_by_not_found_name.tasks) == 0
}

// Filter by category ID.
locals {
  category_id = data.huaweicloud_dataarts_quality_consistency_tasks.all.tasks[0].category_id
}

data "huaweicloud_dataarts_quality_consistency_tasks" "filter_by_category_id" {
  workspace_id = "%[1]s"

  category_id = local.category_id
}

locals {
  category_id_filter_result = [
    for v in data.huaweicloud_dataarts_quality_consistency_tasks.filter_by_category_id.tasks[*].category_id :
    v == local.category_id
  ]
}

output "is_category_id_filter_useful" {
  value = length(local.category_id_filter_result) > 0 && alltrue(local.category_id_filter_result)
}

// Filter by creator.
locals {
  task_creator = data.huaweicloud_dataarts_quality_consistency_tasks.all.tasks[0].creator
}

data "huaweicloud_dataarts_quality_consistency_tasks" "filter_by_creator" {
  workspace_id = "%[1]s"

  creator = local.task_creator
}

locals {
  creator_filter_result = [
    for v in data.huaweicloud_dataarts_quality_consistency_tasks.filter_by_creator.tasks[*].creator :
    v == local.task_creator
  ]
}

output "is_creator_filter_useful" {
  value = length(local.creator_filter_result) > 0 && alltrue(local.creator_filter_result)
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID)
}

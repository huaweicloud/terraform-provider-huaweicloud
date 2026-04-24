package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataArchitectureProcess_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dataarts_architecture_process.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_dataarts_architecture_process.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byParentId   = "data.huaweicloud_dataarts_architecture_process.filter_by_parent_id"
		dcByParentId = acceptance.InitDataSourceCheck(byParentId)

		byOwner   = "data.huaweicloud_dataarts_architecture_process.filter_by_owner"
		dcByOwner = acceptance.InitDataSourceCheck(byOwner)

		byCreator   = "data.huaweicloud_dataarts_architecture_process.filter_by_creator"
		dcByCreator = acceptance.InitDataSourceCheck(byCreator)

		byTimeRange   = "data.huaweicloud_dataarts_architecture_process.filter_by_time_range"
		dcByTimeRange = acceptance.InitDataSourceCheck(byTimeRange)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataArchitectureProcess_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "processes.#", regexp.MustCompile(`^[1-9]([0-9]*)$`)),
					resource.TestCheckResourceAttrSet(all, "processes.0.id"),
					resource.TestCheckResourceAttrSet(all, "processes.0.name"),
					resource.TestCheckResourceAttrSet(all, "processes.0.name_en"),
					resource.TestCheckResourceAttrSet(all, "processes.0.description"),
					resource.TestCheckResourceAttrSet(all, "processes.0.guid"),
					resource.TestCheckResourceAttrSet(all, "processes.0.owner"),
					resource.TestCheckResourceAttrSet(all, "processes.0.parent_id"),
					resource.TestCheckResourceAttrSet(all, "processes.0.qualified_id"),
					resource.TestCheckResourceAttrSet(all, "processes.0.create_by"),
					resource.TestCheckResourceAttrSet(all, "processes.0.update_by"),
					resource.TestCheckResourceAttrSet(all, "processes.0.create_time"),
					resource.TestCheckResourceAttrSet(all, "processes.0.update_time"),
					resource.TestCheckResourceAttrSet(all, "processes.0.bizmetric_num"),
					resource.TestCheckResourceAttrSet(all, "processes.0.children_num"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByParentId.CheckResourceExists(),
					resource.TestCheckOutput("is_parent_id_filter_useful", "true"),
					dcByOwner.CheckResourceExists(),
					resource.TestCheckOutput("is_owner_filter_useful", "true"),
					dcByCreator.CheckResourceExists(),
					resource.TestCheckOutput("is_creator_filter_useful", "true"),
					dcByTimeRange.CheckResourceExists(),
					resource.TestCheckOutput("is_time_range_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataArchitectureProcess_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_process" "parent" {
  workspace_id = "%[1]s"
  name         = "%[2]s_parent"
  owner        = "owner_parent"
  description  = "Parent process"
}

resource "huaweicloud_dataarts_architecture_process" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  owner        = "owner"
  description  = "Child process"
  parent_id    = huaweicloud_dataarts_architecture_process.parent.id
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testAccDataArchitectureProcess_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Query all process architectures without any filter
data "huaweicloud_dataarts_architecture_process" "test" {
  workspace_id = "%[2]s"

  depends_on = [huaweicloud_dataarts_architecture_process.test]
}

# Filter by name
data "huaweicloud_dataarts_architecture_process" "filter_by_name" {
  workspace_id = "%[2]s"
  name         = huaweicloud_dataarts_architecture_process.test.name

  depends_on = [huaweicloud_dataarts_architecture_process.test]
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_process.filter_by_name.processes[*].name:
    strcontains(v, huaweicloud_dataarts_architecture_process.test.name)
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by parent_id
data "huaweicloud_dataarts_architecture_process" "filter_by_parent_id" {
  workspace_id = "%[2]s"
  parent_id    = huaweicloud_dataarts_architecture_process.parent.id

  depends_on = [huaweicloud_dataarts_architecture_process.test]
}

locals {
  parent_id_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_process.filter_by_parent_id.processes[*].parent_id:
    v == huaweicloud_dataarts_architecture_process.parent.id
  ]
}

output "is_parent_id_filter_useful" {
  value = length(local.parent_id_filter_result) > 0 && alltrue(local.parent_id_filter_result)
}

# Filter by owner
data "huaweicloud_dataarts_architecture_process" "filter_by_owner" {
  workspace_id = "%[2]s"
  owner        = huaweicloud_dataarts_architecture_process.test.owner

  depends_on = [huaweicloud_dataarts_architecture_process.test]
}

locals {
  owner_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_process.filter_by_owner.processes[*].owner:
    strcontains(v, huaweicloud_dataarts_architecture_process.test.owner)
  ]
}

output "is_owner_filter_useful" {
  value = length(local.owner_filter_result) > 0 && alltrue(local.owner_filter_result)
}

# Filter by creator
data "huaweicloud_dataarts_architecture_process" "filter_by_creator" {
  workspace_id = "%[2]s"
  create_by    = huaweicloud_dataarts_architecture_process.test.created_by

  depends_on = [huaweicloud_dataarts_architecture_process.test]
}

locals {
  creator_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_process.filter_by_creator.processes[*].create_by:
    v == huaweicloud_dataarts_architecture_process.test.created_by
  ]
}

output "is_creator_filter_useful" {
  value = length(local.creator_filter_result) > 0 && alltrue(local.creator_filter_result)
}

# Filter by time range
locals {
  begin_time = formatdate("YYYY-MM-DD'T'HH:mm:ss'Z'", timeadd(huaweicloud_dataarts_architecture_process.test.created_at, "-24h"))
  end_time   = formatdate("YYYY-MM-DD'T'HH:mm:ss'Z'", timeadd(huaweicloud_dataarts_architecture_process.test.created_at, "24h"))
}

data "huaweicloud_dataarts_architecture_process" "filter_by_time_range" {
  workspace_id = "%[2]s"
  begin_time   = local.begin_time
  end_time     = local.end_time

  depends_on = [huaweicloud_dataarts_architecture_process.test]
}

locals {
  time_range_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_process.filter_by_time_range.processes[*].create_time:
	timecmp(v, local.begin_time) >= 0 &&
	timecmp(v, local.end_time) <= 0
  ]
}

output "is_time_range_filter_useful" {
  value = length(local.time_range_filter_result) > 0 && alltrue(local.time_range_filter_result)
}
`, testAccDataArchitectureProcess_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}

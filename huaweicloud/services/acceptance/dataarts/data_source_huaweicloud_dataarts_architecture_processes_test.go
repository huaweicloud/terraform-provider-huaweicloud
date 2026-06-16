package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataArchitectureProcesses_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dataarts_architecture_processes.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByName   = "data.huaweicloud_dataarts_architecture_processes.filter_by_name"
		dcFilterByName = acceptance.InitDataSourceCheck(filterByName)

		filterByParentId   = "data.huaweicloud_dataarts_architecture_processes.filter_by_parent_id"
		dcFilterByParentId = acceptance.InitDataSourceCheck(filterByParentId)

		filterByCreateBy   = "data.huaweicloud_dataarts_architecture_processes.filter_by_create_by"
		dcFilterByCreateBy = acceptance.InitDataSourceCheck(filterByCreateBy)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataArchitectureProcesses_nonExistentWorkspace(),
				ExpectError: regexp.MustCompile("error querying DataArts Architecture processes"),
			},
			{
				Config: testAccDataSourceArchitectureProcesses_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "processes.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),

					// filter by name
					dcFilterByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(filterByName, "processes.0.id"),
					resource.TestCheckResourceAttrPair(filterByName, "processes.0.name",
						"huaweicloud_dataarts_architecture_process.test", "name"),
					resource.TestCheckResourceAttrPair(filterByName, "processes.0.name_en",
						"huaweicloud_dataarts_architecture_process.test", "name"),
					resource.TestCheckResourceAttrPair(filterByName, "processes.0.description",
						"huaweicloud_dataarts_architecture_process.test", "description"),
					resource.TestCheckResourceAttrPair(filterByName, "processes.0.owner",
						"huaweicloud_dataarts_architecture_process.test", "owner"),
					resource.TestCheckResourceAttrPair(filterByName, "processes.0.parent_id",
						"huaweicloud_dataarts_architecture_process.test", "parent_id"),
					resource.TestCheckResourceAttrPair(filterByName, "processes.0.prev_id",
						"huaweicloud_dataarts_architecture_process.test", "prev_id"),
					resource.TestCheckResourceAttrPair(filterByName, "processes.0.qualified_id",
						"huaweicloud_dataarts_architecture_process.test", "qualified_id"),
					resource.TestCheckResourceAttrPair(filterByName, "processes.0.created_at",
						"huaweicloud_dataarts_architecture_process.test", "created_at"),
					resource.TestCheckResourceAttrPair(filterByName, "processes.0.updated_at",
						"huaweicloud_dataarts_architecture_process.test", "updated_at"),
					resource.TestCheckResourceAttrPair(filterByName, "processes.0.created_by",
						"huaweicloud_dataarts_architecture_process.test", "created_by"),
					resource.TestCheckResourceAttrPair(filterByName, "processes.0.updated_by",
						"huaweicloud_dataarts_architecture_process.test", "updated_by"),

					// filter by parent ID
					dcFilterByParentId.CheckResourceExists(),
					resource.TestCheckOutput("is_parent_id_filter_useful", "true"),

					// filter by create by
					dcFilterByCreateBy.CheckResourceExists(),
					resource.TestCheckOutput("is_create_by_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataArchitectureProcesses_nonExistentWorkspace() string {
	randUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
data "huaweicloud_dataarts_architecture_processes" "test" {
  workspace_id = "%[1]s"
}
`, randUUID.String())
}

func testAccDataSourceArchitectureProcesses_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_process" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  owner        = "owner1"
  description  = "Created by script"
}

data "huaweicloud_dataarts_architecture_processes" "all" {
  depends_on = [
    huaweicloud_dataarts_architecture_process.test,
  ]

  workspace_id = "%[1]s"
}

# Filter by name
locals {
  name = huaweicloud_dataarts_architecture_process.test.name
}

data "huaweicloud_dataarts_architecture_processes" "filter_by_name" {
  depends_on = [
    huaweicloud_dataarts_architecture_process.test,
  ]

  workspace_id = "%[1]s"
  name         = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_processes.filter_by_name.processes : v.name == local.name
  ]
}
  
output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by parent ID
locals {
  parent_id = huaweicloud_dataarts_architecture_process.test.parent_id
}

data "huaweicloud_dataarts_architecture_processes" "filter_by_parent_id" {
  depends_on = [
    huaweicloud_dataarts_architecture_process.test,
  ]

  workspace_id = "%[1]s"
  parent_id    = local.parent_id
}

locals {
  parent_id_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_processes.filter_by_parent_id.processes : v.parent_id == local.parent_id
  ]
}
  
output "is_parent_id_filter_useful" {
  value = length(local.parent_id_filter_result) > 0 && alltrue(local.parent_id_filter_result)
}

# Filter by create by
locals {
  create_by = huaweicloud_dataarts_architecture_process.test.created_by
}

data "huaweicloud_dataarts_architecture_processes" "filter_by_create_by" {
  depends_on = [
    huaweicloud_dataarts_architecture_process.test,
  ]

  workspace_id = "%[1]s"
  create_by    = local.create_by
}

locals {
  create_by_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_processes.filter_by_create_by.processes : v.created_by == local.create_by
  ]
}
  
output "is_create_by_filter_useful" {
  value = length(local.create_by_filter_result) > 0 && alltrue(local.create_by_filter_result)
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

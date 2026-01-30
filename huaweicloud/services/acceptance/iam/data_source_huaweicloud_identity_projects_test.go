package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataProjects_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_identity_projects.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_identity_projects.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byProjectId   = "data.huaweicloud_identity_projects.filter_by_project_id"
		dcByProjectId = acceptance.InitDataSourceCheck(byProjectId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataProjects_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "projects.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "projects.0.id"),
					resource.TestCheckResourceAttrSet(all, "projects.0.name"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByProjectId.CheckResourceExists(),
					resource.TestCheckOutput("is_project_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataProjects_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_project" "test" {
  name        = "%[1]s_%[2]s"
  status      = "suspended"
}
`, acceptance.HW_REGION_NAME, name)
}

func testAccDataProjects_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# All
data "huaweicloud_identity_projects" "all" {
  # Waiting for the project to be created
  depends_on = [huaweicloud_identity_project.test]
}

# Filter by name
locals {
  name = "%[2]s_%[3]s"
}

data "huaweicloud_identity_projects" "filter_by_name" {
  name = local.name

  # Waiting for the project to be created
  depends_on = [huaweicloud_identity_project.test]
}

locals {
  name_filter_result = [for o in data.huaweicloud_identity_projects.filter_by_name.projects : o.name == local.name]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) >= 1 && alltrue(local.name_filter_result)
}

# Filter by project ID
locals {
  project_id = huaweicloud_identity_project.test.id
}

data "huaweicloud_identity_projects" "filter_by_project_id" {
  project_id = local.project_id

  # Waiting for the project to be created
  depends_on = [huaweicloud_identity_project.test]
}

locals {
  project_id_filter_result = [for o in data.huaweicloud_identity_projects.filter_by_project_id.projects : 
    o.id == local.project_id]
}

output "is_project_id_filter_useful" {
  value = length(local.project_id_filter_result) >= 1 && alltrue(local.project_id_filter_result)
}
`, testAccDataProjects_base(name), acceptance.HW_REGION_NAME, name)
}

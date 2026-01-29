package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceProjects_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		dcName = "data.huaweicloud_identity_projects.test"
		dc     = acceptance.InitDataSourceCheck(dcName)

		dcNameByName = "data.huaweicloud_identity_projects.filter_by_name"
		dcByName     = acceptance.InitDataSourceCheck(dcNameByName)

		dcNameByProjectId = "data.huaweicloud_identity_projects.filter_by_project_id"
		dcByProjectId     = acceptance.InitDataSourceCheck(dcNameByProjectId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProjects_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "projects.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dcName, "projects.0.id"),
					resource.TestCheckResourceAttrSet(dcName, "projects.0.name"),
				),
			},
			{
				Config: testAccDataSourceProjects_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					dcByName.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcNameByName, "projects.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
			{
				Config: testAccDataSourceProjects_basic_step3(name),
				Check: resource.ComposeTestCheckFunc(
					dcByProjectId.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcNameByProjectId, "projects.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("is_project_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceProjects_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_project" "test" {
  name        = "%[1]s_%[2]s"
  status      = "suspended"
  description = "Terraform test project"
}
`, acceptance.HW_REGION_NAME, name)
}

func testAccDataSourceProjects_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identity_projects" "test" {}
`, testAccDataSourceProjects_base(name))
}

func testAccDataSourceProjects_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identity_projects" "filter_by_name" {
  name = "%[2]s_%[3]s"
}

locals {
  name_filter_result = [for o in data.huaweicloud_identity_projects.filter_by_name.projects : o.name == "%[2]s_%[3]s"]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) >= 1 && alltrue(local.name_filter_result)
}
`, testAccDataSourceProjects_base(name), acceptance.HW_REGION_NAME, name)
}

func testAccDataSourceProjects_basic_step3(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identity_projects" "filter_by_project_id" {
  project_id = huaweicloud_identity_project.test.id
}

locals {
  project_id_filter_result = [for o in data.huaweicloud_identity_projects.filter_by_project_id.projects : 
    o.id == huaweicloud_identity_project.test.id]
}

output "is_project_id_filter_useful" {
  value = length(local.project_id_filter_result) >= 1 && alltrue(local.project_id_filter_result)
}
`, testAccDataSourceProjects_base(name))
}

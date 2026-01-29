package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityProjects_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_identity_projects.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProjects_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dcName, "projects.#", "1"),
					resource.TestCheckResourceAttr(dcName, "projects.0.name", "MOS"),
					resource.TestCheckResourceAttr(dcName, "projects.0.enabled", "true"),
				),
			},
		},
	})
}

func TestAccIdentityProjects_subProject(t *testing.T) {
	var (
		dcName = "data.huaweicloud_identity_projects.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProjects_subProject,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dcName, "projects.#", "1"),
					resource.TestCheckResourceAttr(dcName, "projects.0.name", "cn-north-4_test"),
					resource.TestCheckResourceAttr(dcName, "projects.0.enabled", "true"),
				),
			},
		},
	})
}

func TestAccIdentityProjects_projectId(t *testing.T) {
	var (
		dcName = "data.huaweicloud_identity_projects.test"
		dc     = acceptance.InitDataSourceCheck(dcName)

		projectName = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProjects_projectId(projectName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dcName, "projects.#", "1"),
					resource.TestCheckResourceAttr(dcName, "projects.0.name",
						fmt.Sprintf("%s_%s", acceptance.HW_REGION_NAME, projectName)),
					resource.TestCheckResourceAttr(dcName, "projects.0.enabled", "true"),
				),
			},
		},
	})
}

const testAccIdentityProjects_basic string = `
data "huaweicloud_identity_projects" "test" {
  name = "MOS"
}
`

const testAccIdentityProjects_subProject string = `
data "huaweicloud_identity_projects" "test" {
  name = "cn-north-4_test"
}
`

func testAccIdentityProjects_projectId(projectName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_project" "test" {
  name        = "%s_%s"
  status      = "suspended"
  description = "Created by acc test"
}

data "huaweicloud_identity_projects" "test" {
  project_id = huaweicloud_identity_project.test.id
}

`, acceptance.HW_REGION_NAME, projectName)
}

package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityProjectsDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identity_projects.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProjectsDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "projects.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "projects.0.name", "MOS"),
					resource.TestCheckResourceAttr(dataSourceName, "projects.0.enabled", "true"),
				),
			},
		},
	})
}

func TestAccIdentityProjectsDataSource_subProject(t *testing.T) {
	dataSourceName := "data.huaweicloud_identity_projects.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProjectsDataSource_subProject,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "projects.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "projects.0.name", "cn-north-4_test"),
					resource.TestCheckResourceAttr(dataSourceName, "projects.0.enabled", "true"),
				),
			},
		},
	})
}

func TestAccIdentityProjectsDataSource_projectId(t *testing.T) {
	dataSourceName := "data.huaweicloud_identity_projects.test"
	projectName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProjectsDataSource_projectId(projectName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "projects.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "projects.0.name",
						fmt.Sprintf("%s_%s", acceptance.HW_REGION_NAME, projectName)),
					resource.TestCheckResourceAttr(dataSourceName, "projects.0.enabled", "true"),
				),
			},
		},
	})
}

const testAccIdentityProjectsDataSource_basic string = `
data "huaweicloud_identity_projects" "test" {
  name = "MOS"
}
`

const testAccIdentityProjectsDataSource_subProject string = `
data "huaweicloud_identity_projects" "test" {
  name = "cn-north-4_test"
}
`

func testAccIdentityProjectsDataSource_projectId(projectName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_project" "project_1" {
  name        = "%s_%s"
  status      = "suspended"
  description = "An updated project"
}

data "huaweicloud_identity_projects" "test" {
  project_id = huaweicloud_identity_project.project_1.id
}

`, acceptance.HW_REGION_NAME, projectName)
}

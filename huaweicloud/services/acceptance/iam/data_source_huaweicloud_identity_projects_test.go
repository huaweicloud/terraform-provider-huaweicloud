package iam

import (
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

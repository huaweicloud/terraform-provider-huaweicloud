package cpts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCptsProjects_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cpts_projects.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCptsProjects_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "projects.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "projects.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "projects.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "projects.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "projects.0.source"),
					resource.TestCheckResourceAttrSet(dataSource, "projects.0.updated_at"),
				),
			},
		},
	})
}

func testDataSourceCptsProjects_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cpts_projects" "test" {
  depends_on = [huaweicloud_cpts_project.test]
}
`, testProject_basic(name))
}

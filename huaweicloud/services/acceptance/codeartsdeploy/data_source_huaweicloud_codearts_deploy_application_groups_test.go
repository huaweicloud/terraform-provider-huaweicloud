package codeartsdeploy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCodeartsDeployApplicationGroups_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_deploy_application_groups.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCodeartsDeployApplicationGroups_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "groups.#", "3"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.path"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.ordinal"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.created_by"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.updated_by"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.application_count"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.children.0"),
				),
			},
		},
	})
}

func testDataSourceCodeartsDeployApplicationGroups_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_deploy_application_groups" "test" {
  depends_on = [huaweicloud_codearts_deploy_application_group.level2]

  project_id = huaweicloud_codearts_project.test.id
}
`, testDeployApplicationGroup_secondLevel(name))
}

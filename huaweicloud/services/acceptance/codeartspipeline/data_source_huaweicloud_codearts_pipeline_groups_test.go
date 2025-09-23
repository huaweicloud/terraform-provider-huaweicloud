package codeartspipeline

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCodeartsPipelineGroups_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_pipeline_groups.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCodeartsPipelineGroups_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "groups.#", "2"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.path_id"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.ordinal"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.creator"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.children.0"),
				),
			},
		},
	})
}

func testDataSourceCodeartsPipelineGroups_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_pipeline_groups" "test" {
  depends_on = [huaweicloud_codearts_pipeline_group.level2]

  project_id = huaweicloud_codearts_project.test.id
}
`, testPipelineGroup_secondLevel(name))
}

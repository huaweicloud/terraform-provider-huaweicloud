package codeartspipeline

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCodeartsPipelineModifyHistories_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_pipeline_modify_histories.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCodeartsPipelineModifyHistories_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "histories.#"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.modify_type"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.creator_name"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.creator_nick_name"),
				),
			},
		},
	})
}

func testDataSourceCodeartsPipelineModifyHistories_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_pipeline_modify_histories" "test" {
  project_id  = huaweicloud_codearts_project.test.id
  pipeline_id = huaweicloud_codearts_pipeline.test.id
}
`, testPipeline_basic(name))
}

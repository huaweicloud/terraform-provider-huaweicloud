package codeartspipeline

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePipelineTags_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_pipeline_tags.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePipelineTags_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.color"),
				),
			},
		},
	})
}

func testDataSourcePipelineTags_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_pipeline_tags" "test" {
  depends_on = [huaweicloud_codearts_pipeline_tag.test]

  project_id  = huaweicloud_codearts_project.test.id
}
`, testPipelineTag_basic(name))
}

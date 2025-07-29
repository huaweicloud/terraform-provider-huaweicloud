package codeartspipeline

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCodeartsPipelinePublishers_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_pipeline_publishers.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCodeartsPipelinePublishers_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "publishers.#"),
					resource.TestCheckResourceAttrSet(dataSource, "publishers.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "publishers.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "publishers.0.en_name"),
					resource.TestCheckResourceAttrSet(dataSource, "publishers.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "publishers.0.website"),
					resource.TestCheckResourceAttrSet(dataSource, "publishers.0.support_url"),
					resource.TestCheckResourceAttrSet(dataSource, "publishers.0.source_url"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCodeartsPipelinePublishers_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_pipeline_publishers" "test" {
  depends_on = [huaweicloud_codearts_pipeline_publisher.test]
}

// filter by name
data "huaweicloud_codearts_pipeline_publishers" "filter_by_name" {
  name = huaweicloud_codearts_pipeline_publisher.test.name
}

locals {
  filter_result_by_name = [for v in data.huaweicloud_codearts_pipeline_publishers.filter_by_name.publishers[*].name :
    v == huaweicloud_codearts_pipeline_publisher.test.name]
}

output "is_name_filter_useful" {
  value = length(local.filter_result_by_name) > 0 && alltrue(local.filter_result_by_name)
}
`, testPipelinePublisher_basic(name))
}

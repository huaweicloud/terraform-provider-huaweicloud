package codeartspipeline

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCodeartsPipelineMicroServices_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_pipeline_micro_services.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCodeartsPipelineMicroServices_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "micro_services.#"),
					resource.TestCheckResourceAttrSet(dataSource, "micro_services.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "micro_services.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "micro_services.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "micro_services.0.creator_id"),
					resource.TestCheckResourceAttrSet(dataSource, "micro_services.0.updater_id"),
					resource.TestCheckResourceAttrSet(dataSource, "micro_services.0.creator_name"),
					resource.TestCheckResourceAttrSet(dataSource, "micro_services.0.updater_name"),
					resource.TestCheckResourceAttrSet(dataSource, "micro_services.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "micro_services.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "micro_services.0.status"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCodeartsPipelineMicroServices_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_pipeline_micro_services" "test" {
  depends_on = [huaweicloud_codearts_pipeline_micro_service.test]

  project_id = huaweicloud_codearts_project.test.id
}

// filter by name
data "huaweicloud_codearts_pipeline_micro_services" "filter_by_name" {
  project_id = huaweicloud_codearts_project.test.id
  name       = huaweicloud_codearts_pipeline_micro_service.test.name
}

locals {
  filter_result_by_name = [for v in data.huaweicloud_codearts_pipeline_micro_services.filter_by_name.micro_services[*].name :
    v == huaweicloud_codearts_pipeline_micro_service.test.name]
}

output "is_name_filter_useful" {
  value = length(local.filter_result_by_name) > 0 && alltrue(local.filter_result_by_name)
}
`, testPipelineMicroService_update(name))
}

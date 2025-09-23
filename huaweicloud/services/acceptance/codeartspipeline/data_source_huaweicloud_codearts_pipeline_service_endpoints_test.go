package codeartspipeline

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePipelineServiceEndpoints_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_pipeline_service_endpoints.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePipelineServiceEndpoints_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "endpoints.#"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoints.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoints.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoints.0.module_id"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoints.0.url"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoints.0.created_by.#"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoints.0.created_by.0.user_id"),
					resource.TestCheckResourceAttrSet(dataSource, "endpoints.0.created_by.0.user_name"),

					resource.TestCheckOutput("is_module_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourcePipelineServiceEndpoints_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_pipeline_service_endpoints" "test" {
  depends_on = [huaweicloud_codearts_pipeline_service_endpoint.test]

  project_id = huaweicloud_codearts_project.test.id
}

// filter by module_id
data "huaweicloud_codearts_pipeline_service_endpoints" "filter_by_module_id" {
  project_id = huaweicloud_codearts_project.test.id
  module_id  = huaweicloud_codearts_pipeline_service_endpoint.test.module_id
}

locals {
  filter_result_by_module_id = [for v in data.huaweicloud_codearts_pipeline_service_endpoints.filter_by_module_id.endpoints[*].module_id :
    v == huaweicloud_codearts_pipeline_service_endpoint.test.module_id]
}

output "is_module_id_filter_useful" {
  value = length(local.filter_result_by_module_id) > 0 && alltrue(local.filter_result_by_module_id)
}
`, testPipelineServiceEndpoint_basic(name))
}

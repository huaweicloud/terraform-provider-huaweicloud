package codeartspipeline

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePipelineUserPermissions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_pipeline_user_permissions.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePipelineUserPermissions_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "users.#"),
					resource.TestCheckResourceAttrSet(dataSource, "users.0.user_id"),
					resource.TestCheckResourceAttrSet(dataSource, "users.0.user_name"),
					resource.TestCheckResourceAttrSet(dataSource, "users.0.operation_query"),
					resource.TestCheckResourceAttrSet(dataSource, "users.0.operation_execute"),
					resource.TestCheckResourceAttrSet(dataSource, "users.0.operation_update"),
					resource.TestCheckResourceAttrSet(dataSource, "users.0.operation_delete"),
					resource.TestCheckResourceAttrSet(dataSource, "users.0.operation_authorize"),
					resource.TestCheckResourceAttrSet(dataSource, "users.0.role_id"),
					resource.TestCheckResourceAttrSet(dataSource, "users.0.role_name"),

					resource.TestCheckOutput("is_user_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourcePipelineUserPermissions_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_pipeline_user_permissions" "test" {
  project_id  = huaweicloud_codearts_project.test.id
  pipeline_id = huaweicloud_codearts_pipeline.test.id
}

// filter by user name
data "huaweicloud_codearts_pipeline_user_permissions" "filter_by_user_name" {
  project_id  = huaweicloud_codearts_project.test.id
  pipeline_id = huaweicloud_codearts_pipeline.test.id
  user_name   = data.huaweicloud_codearts_pipeline_user_permissions.test.users[0].user_name
}

locals {
  filter_result_by_user_name = [for v in data.huaweicloud_codearts_pipeline_user_permissions.filter_by_user_name.users[*].user_name :
    v == data.huaweicloud_codearts_pipeline_user_permissions.test.users[0].user_name]
}

output "is_user_name_filter_useful" {
  value = length(local.filter_result_by_user_name) > 0 && alltrue(local.filter_result_by_user_name)
}
`, testPipeline_basic(name))
}

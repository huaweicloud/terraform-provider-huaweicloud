package codeartsdeploy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCodeartsDeployEnvironments_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_deploy_environments.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCodeartsDeployEnvironments_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "environments.#"),
					resource.TestCheckResourceAttrSet(dataSource, "environments.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "environments.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "environments.0.os_type"),
					resource.TestCheckResourceAttrSet(dataSource, "environments.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "environments.0.created_by.0.user_id"),
					resource.TestCheckResourceAttrSet(dataSource, "environments.0.created_by.0.user_name"),
					resource.TestCheckResourceAttrSet(dataSource, "environments.0.deploy_type"),
					resource.TestCheckResourceAttrSet(dataSource, "environments.0.instance_count"),
					resource.TestCheckResourceAttrSet(dataSource, "environments.0.permission.0.can_delete"),
					resource.TestCheckResourceAttrSet(dataSource, "environments.0.permission.0.can_deploy"),
					resource.TestCheckResourceAttrSet(dataSource, "environments.0.permission.0.can_edit"),
					resource.TestCheckResourceAttrSet(dataSource, "environments.0.permission.0.can_manage"),
					resource.TestCheckResourceAttrSet(dataSource, "environments.0.permission.0.can_view"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCodeartsDeployEnvironments_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_deploy_environments" "test" {
  depends_on = [huaweicloud_codearts_deploy_environment.test]

  project_id     = huaweicloud_codearts_project.test.id
  application_id = huaweicloud_codearts_deploy_application.test.id
}

// filter by name
data "huaweicloud_codearts_deploy_environments" "filter_by_name" {
  project_id     = huaweicloud_codearts_project.test.id
  application_id = huaweicloud_codearts_deploy_application.test.id
  name           = huaweicloud_codearts_deploy_environment.test.name
}

locals {
  filter_result_by_name = [for v in data.huaweicloud_codearts_deploy_environments.filter_by_name.environments[*].name : 
    v == huaweicloud_codearts_deploy_environment.test.name]
}

output "is_name_filter_useful" {
  value = length(local.filter_result_by_name) == 1 && alltrue(local.filter_result_by_name)
}
`, testDeployEnvironment_basic(name))
}

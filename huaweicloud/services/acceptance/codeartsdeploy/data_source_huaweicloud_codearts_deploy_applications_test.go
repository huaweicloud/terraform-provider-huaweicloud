package codeartsdeploy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCodeartsDeployApplications_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_deploy_applications.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCodeartsDeployApplications_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "applications.#"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.project_name"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.release_id"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.can_modify"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.can_manage"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.can_create_env"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.can_execute"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.can_copy"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.can_delete"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.can_view"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.can_disable"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.is_care"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.is_disable"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.create_user_id"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.create_tenant_id"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.arrange_infos.#"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.arrange_infos.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.arrange_infos.0.state"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.arrange_infos.0.deploy_system"),

					resource.TestCheckOutput("is_state_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCodeartsDeployApplications_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_deploy_applications" "test" {
  depends_on = [huaweicloud_codearts_deploy_application.test]

  project_id = huaweicloud_codearts_project.test.id
}

// filter by state
data "huaweicloud_codearts_deploy_applications" "filter_by_state" {
  depends_on = [huaweicloud_codearts_deploy_application.test]
  
  project_id = huaweicloud_codearts_project.test.id
  states     = ["not_executed"]
}

locals {
  filter_result_by_state = [for v in data.huaweicloud_codearts_deploy_applications.filter_by_state.applications[*].arrange_infos[0].state : 
    v == "Draft"]
}

output "is_state_filter_useful" {
  value = length(local.filter_result_by_state) == 1 && alltrue(local.filter_result_by_state)
}
`, testDeployApplication_basic(name))
}

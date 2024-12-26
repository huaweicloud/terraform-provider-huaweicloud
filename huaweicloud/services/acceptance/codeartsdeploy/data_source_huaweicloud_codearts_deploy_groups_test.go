package codeartsdeploy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCodeartsDeployGroups_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_deploy_groups.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCodeartsDeployGroups_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "groups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.env_count"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.host_count"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.is_proxy_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.os_type"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.created_by"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.permission.#"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.permission.0.can_view"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.permission.0.can_edit"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.permission.0.can_delete"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.permission.0.can_add_host"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.permission.0.can_manage"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.permission.0.can_copy"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_os_type_filter_useful", "true"),
					resource.TestCheckOutput("is_proxy_mode_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCodeartsDeployGroups_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_deploy_groups" "test" {
  depends_on = [huaweicloud_codearts_deploy_group.test]

  project_id = huaweicloud_codearts_project.test.id
}

// filter by name
data "huaweicloud_codearts_deploy_groups" "filter_by_name" {
  project_id = huaweicloud_codearts_project.test.id
  name       = huaweicloud_codearts_deploy_group.test.name
}

locals {
  filter_result_by_name = [for v in data.huaweicloud_codearts_deploy_groups.filter_by_name.groups[*].name : 
    v == huaweicloud_codearts_deploy_group.test.name]
}

output "is_name_filter_useful" {
  value = length(local.filter_result_by_name) > 0
}

// filter by os type
data "huaweicloud_codearts_deploy_groups" "filter_by_os_type" {
  project_id = huaweicloud_codearts_project.test.id
  os_type    = huaweicloud_codearts_deploy_group.test.os_type
}

locals {
  filter_result_by_os_type = [for v in data.huaweicloud_codearts_deploy_groups.filter_by_os_type.groups[*].os_type : 
    v == huaweicloud_codearts_deploy_group.test.os_type]
}

output "is_os_type_filter_useful" {
  value = length(local.filter_result_by_os_type) > 0 && alltrue(local.filter_result_by_os_type) 
}

// filter by proxy mode
data "huaweicloud_codearts_deploy_groups" "filter_by_proxy_mode" {
  project_id    = huaweicloud_codearts_project.test.id
  is_proxy_mode = huaweicloud_codearts_deploy_group.test.is_proxy_mode
}

locals {
  filter_result_by_proxy_mode = [for v in data.huaweicloud_codearts_deploy_groups.filter_by_proxy_mode.groups[*].is_proxy_mode : 
    v == huaweicloud_codearts_deploy_group.test.is_proxy_mode]
}

output "is_proxy_mode_filter_useful" {
  value = length(local.filter_result_by_proxy_mode) > 0 && alltrue(local.filter_result_by_proxy_mode) 
}
`, testDeployGroup_basic(name))
}

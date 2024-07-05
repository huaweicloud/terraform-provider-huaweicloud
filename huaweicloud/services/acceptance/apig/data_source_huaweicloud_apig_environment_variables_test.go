package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceApigEnvironmentVariables_basic(t *testing.T) {
	var (
		rName      = acceptance.RandomAccResourceName()
		dataSource = "data.huaweicloud_apig_environment_variables.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byEnvId   = "data.huaweicloud_apig_environment_variables.filter_by_env_id"
		dcByEnvId = acceptance.InitDataSourceCheck(byEnvId)

		byName   = "data.huaweicloud_apig_environment_variables.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNotFoundName   = "data.huaweicloud_apig_environment_variables.not_found"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceApigEnvironmentVariables_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "variables.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByEnvId.CheckResourceExists(),
					resource.TestCheckOutput("is_env_id_filter_useful", "true"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckResourceAttrPair(byName, "variables.0.id", "huaweicloud_apig_environment_variable.test", "id"),
					resource.TestCheckResourceAttrPair(byName, "variables.0.group_id", "huaweicloud_apig_group.test", "id"),
					resource.TestCheckResourceAttrSet(byName, "variables.0.value"),
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_not_found_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceApigEnvironmentVariables_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_apig_environment_variables" "test" {
  depends_on = [
    huaweicloud_apig_environment_variable.test
  ]
  
  instance_id = local.instance_id
  group_id    = huaweicloud_apig_group.test.id
}

locals {
  env_id = huaweicloud_apig_environment_variable.test.env_id
}

data "huaweicloud_apig_environment_variables" "filter_by_env_id" {
  instance_id = local.instance_id
  group_id    = huaweicloud_apig_group.test.id
  env_id      = local.env_id
}

locals {
  env_id_filter_result = [
    for v in data.huaweicloud_apig_environment_variables.filter_by_env_id.variables[*].env_id : v == local.env_id
  ]
}

output "is_env_id_filter_useful" {
  value = length(local.env_id_filter_result) > 0 && alltrue(local.env_id_filter_result)
}

locals {
  variable_name = huaweicloud_apig_environment_variable.test.name
}

data "huaweicloud_apig_environment_variables" "filter_by_name" {
  instance_id = local.instance_id
  group_id    = huaweicloud_apig_group.test.id
  name        = local.variable_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_apig_environment_variables.filter_by_name.variables[*].name : v == local.variable_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

data "huaweicloud_apig_environment_variables" "not_found" {
  instance_id = local.instance_id
  group_id    = huaweicloud_apig_group.test.id
  name        = "not_found"
}

output "is_name_not_found_filter_useful" {
  value = length(data.huaweicloud_apig_environment_variables.not_found.variables) == 0
}
`, testAccEnvironmentVariable_basic(name))
}

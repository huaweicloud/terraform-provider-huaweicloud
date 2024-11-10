package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running the test case, please ensure that there is at least one WAF dedicated instance in the current region.
func TestAccDataSourceDedicatedInstances_basic(t *testing.T) {
	var (
		datasourceName = "data.huaweicloud_waf_dedicated_instances.test"
		dc             = acceptance.InitDataSourceCheck(datasourceName)

		byID   = "data.huaweicloud_waf_dedicated_instances.id_filter"
		dcByID = acceptance.InitDataSourceCheck(byID)

		byName   = "data.huaweicloud_waf_dedicated_instances.name_filter"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byAllParameters   = "data.huaweicloud_waf_dedicated_instances.all_parameters_filter"
		dcByAllParameters = acceptance.InitDataSourceCheck(byAllParameters)

		byNonExistID   = "data.huaweicloud_waf_dedicated_instances.non_exist_id_filter"
		dcByNonExistID = acceptance.InitDataSourceCheck(byNonExistID)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedInstances_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(datasourceName, "instances.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "instances.0.id"),
					resource.TestCheckResourceAttrSet(datasourceName, "instances.0.name"),
					resource.TestCheckResourceAttrSet(datasourceName, "instances.0.available_zone"),
					resource.TestCheckResourceAttrSet(datasourceName, "instances.0.ecs_flavor"),
					resource.TestCheckResourceAttrSet(datasourceName, "instances.0.cpu_architecture"),
					resource.TestCheckResourceAttrSet(datasourceName, "instances.0.security_group.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "instances.0.server_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "instances.0.service_ip"),
					resource.TestCheckResourceAttrSet(datasourceName, "instances.0.subnet_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "instances.0.vpc_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "instances.0.run_status"),
					resource.TestCheckResourceAttrSet(datasourceName, "instances.0.access_status"),
					resource.TestCheckResourceAttrSet(datasourceName, "instances.0.upgradable"),

					dcByID.CheckResourceExists(),
					resource.TestCheckOutput("id_filter_is_useful", "true"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),

					dcByAllParameters.CheckResourceExists(),
					resource.TestCheckOutput("all_parameters_filter_is_useful", "true"),

					dcByNonExistID.CheckResourceExists(),
					resource.TestCheckOutput("non_exist_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccWafDedicatedInstances_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_waf_dedicated_instances" "test" {
  enterprise_project_id = "%[1]s"
}

# Filter by ID
locals {
  id = data.huaweicloud_waf_dedicated_instances.test.instances.0.id
}

data "huaweicloud_waf_dedicated_instances" "id_filter" {
  id                    = local.id
  enterprise_project_id = "%[1]s"
}

locals {
  id_filter_result = [
    for v in data.huaweicloud_waf_dedicated_instances.id_filter.instances[*].id : v == local.id
  ]
}

output "id_filter_is_useful" {
  value = length(local.id_filter_result) > 0 && alltrue(local.id_filter_result)  
}

# Filter by name
locals {
  name = data.huaweicloud_waf_dedicated_instances.test.instances.0.name
}

data "huaweicloud_waf_dedicated_instances" "name_filter" {
  name                  = local.name
  enterprise_project_id = "%[1]s"
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_waf_dedicated_instances.name_filter.instances[*].name : v == local.name
  ]
}

output "name_filter_is_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)  
}

# Filter by all parameters
data "huaweicloud_waf_dedicated_instances" "all_parameters_filter" {
  id                    = local.id
  name                  = local.name
  enterprise_project_id = "%[1]s"
}

locals {
  all_parameters_filter_result = [
    for v in data.huaweicloud_waf_dedicated_instances.all_parameters_filter.instances[*] :
    v.name == local.name && v.id == local.id
  ]
}

output "all_parameters_filter_is_useful" {
  value = length(local.all_parameters_filter_result) > 0 && alltrue(local.all_parameters_filter_result)  
}

# Filter by non-exist ID
data "huaweicloud_waf_dedicated_instances" "non_exist_id_filter" {
  id                    = "non-exist-id"
  enterprise_project_id = "%[1]s"
}

output "non_exist_id_filter_is_useful" {
  value = length(data.huaweicloud_waf_dedicated_instances.non_exist_id_filter.instances) == 0
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

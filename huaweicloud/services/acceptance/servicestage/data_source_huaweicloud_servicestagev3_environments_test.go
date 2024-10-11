package servicestage

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataV3Environments_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_servicestagev3_environments.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byId   = "data.huaweicloud_servicestagev3_environments.filter_by_id"
		dcById = acceptance.InitDataSourceCheck(byId)

		byName   = "data.huaweicloud_servicestagev3_environments.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNotFoundName   = "data.huaweicloud_servicestagev3_environments.filter_by_not_found_name"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)

		byEpsId   = "data.huaweicloud_servicestagev3_environments.filter_by_eps_id"
		dcByEpsId = acceptance.InitDataSourceCheck(byEpsId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataV3Environments_basic_step1(name),
			},
			{
				// Update the environment name and make sure the attribute 'updated_at' not empty.
				Config: testAccDataV3Environments_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "environments.#", regexp.MustCompile(`[1-9]\d*`)),
					dcById.CheckResourceExists(),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckResourceAttr(byId, "environments.#", "1"),
					resource.TestCheckResourceAttrSet(byId, "environments.0.id"),
					resource.TestCheckResourceAttr(byId, "environments.0.name", updateName),
					resource.TestCheckResourceAttr(byId, "environments.0.description", "Created by terraform test"),
					resource.TestCheckResourceAttrPair(byId, "environments.0.vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttr(byId, "environments.0.enterprise_project_id",
						acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(byId, "environments.0.tags.%", "2"),
					resource.TestCheckResourceAttr(byId, "environments.0.tags.foo", "bar"),
					resource.TestCheckResourceAttr(byId, "environments.0.tags.owner", "terraform"),
					resource.TestMatchResourceAttr(byId, "environments.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(byId, "environments.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found_name_filter_useful", "true"),
					dcByEpsId.CheckResourceExists(),
					resource.TestCheckOutput("is_eps_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataV3Environments_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_servicestagev3_environment" "test" {
  name                  = "%[2]s"
  description           = "Created by terraform test"
  vpc_id                = huaweicloud_vpc.test.id
  enterprise_project_id = "%[3]s"

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}
`, common.TestVpc(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccDataV3Environments_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_servicestagev3_environments" "test" {
  depends_on = [
    huaweicloud_servicestagev3_environment.test
  ]
}

# Filter by ID
locals {
  environment_id = huaweicloud_servicestagev3_environment.test.id
}

data "huaweicloud_servicestagev3_environments" "filter_by_id" {
  depends_on = [
    huaweicloud_servicestagev3_environment.test
  ]

  environment_id = local.environment_id
}

locals {
  id_filter_result = [
    for v in data.huaweicloud_servicestagev3_environments.filter_by_id.environments[*].id : v == local.environment_id
  ]
}

output "is_id_filter_useful" {
  value = length(local.id_filter_result) > 0 && alltrue(local.id_filter_result)
}

# Filter by name
locals {
  environment_name = huaweicloud_servicestagev3_environment.test.name
}

data "huaweicloud_servicestagev3_environments" "filter_by_name" {
  depends_on = [
    huaweicloud_servicestagev3_environment.test
  ]

  name = local.environment_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_servicestagev3_environments.filter_by_name.environments[*].name : v == local.environment_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by name (not found)
data "huaweicloud_servicestagev3_environments" "filter_by_not_found_name" {
  depends_on = [
    huaweicloud_servicestagev3_environment.test
  ]

  name = "name_not_found"
}

output "is_not_found_name_filter_useful" {
  value = length(data.huaweicloud_servicestagev3_environments.filter_by_not_found_name.environments) == 0
}

data "huaweicloud_servicestagev3_environments" "filter_by_eps_id" {
  depends_on = [
    huaweicloud_servicestagev3_environment.test
  ]

  enterprise_project_id = "%[2]s"
}

locals {
  eps_id_filter_result = [
    for v in data.huaweicloud_servicestagev3_environments.filter_by_eps_id.environments[*].enterprise_project_id : v == "%[2]s"
  ]
}

output "is_eps_id_filter_useful" {
  value = length(local.eps_id_filter_result) > 0 && alltrue(local.eps_id_filter_result)
}
`, testAccDataV3Environments_basic_step1(name), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

package cae

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceApplications_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		resourceName = "data.huaweicloud_cae_applications.test"
		all          = acceptance.InitDataSourceCheck(resourceName)

		byApplicationId   = "data.huaweicloud_cae_applications.filter_by_application_id"
		dcByApplicationId = acceptance.InitDataSourceCheck(byApplicationId)

		byName   = "data.huaweicloud_cae_applications.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byEnterpriseProjectId   = "data.huaweicloud_cae_applications.filter_by_enterprise_project_id"
		dcByEnterpriseProjectId = acceptance.InitDataSourceCheck(byEnterpriseProjectId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCaeEnvironmentIds(t, 2)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceApplications_basic(name),
				Check: resource.ComposeTestCheckFunc(
					all.CheckResourceExists(),
					resource.TestMatchResourceAttr(resourceName, "applications.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByApplicationId.CheckResourceExists(),
					resource.TestCheckOutput("application_id_filter_is_useful", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "applications.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "applications.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "applications.0.created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "applications.0.updated_at"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					dcByEnterpriseProjectId.CheckResourceExists(),
					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceApplications_base(name string) string {
	return fmt.Sprintf(`
locals {
  env_ids = split(",", "%[1]s")
}

# Query the enterprise project ID of the environment.
data "huaweicloud_cae_environments" "test" {
  environment_id = local.env_ids[1]
}

# The first application belongs to the default enterprise project.
# The second application belongs to the non-default enterprise project.
resource "huaweicloud_cae_application" "test" {
  count = 2

  environment_id        = local.env_ids[count.index]
  name                  = "%[2]s${count.index}"
  enterprise_project_id = count.index == 1 ? try(data.huaweicloud_cae_environments.test.environments[0].annotations.enterprise_project_id,
  null) : null
}
`, acceptance.HW_CAE_ENVIRONMENT_IDS, name)
}

func testAccDatasourceApplications_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cae_applications" "test" {
  environment_id = huaweicloud_cae_application.test[0].environment_id
}

# Filter by application ID.
locals {
  application_id = huaweicloud_cae_application.test[0].id
}

data "huaweicloud_cae_applications" "filter_by_application_id" {
  environment_id = huaweicloud_cae_application.test[0].environment_id
  application_id = local.application_id
}

locals {
  application_id_filter_result = [for v in data.huaweicloud_cae_applications.filter_by_application_id.applications[*].id : v == local.application_id]
}

output "application_id_filter_is_useful" {
  value = length(local.application_id_filter_result) > 0 && alltrue(local.application_id_filter_result)
}

# Filter by application name.
locals {
  application_name = huaweicloud_cae_application.test[0].name
}

data "huaweicloud_cae_applications" "filter_by_name" {
  environment_id = huaweicloud_cae_application.test[0].environment_id
  name           = local.application_name

  depends_on = [huaweicloud_cae_application.test]
}

locals {
  name_filter_result = [for v in data.huaweicloud_cae_applications.filter_by_name.applications[*].name : v == local.application_name]
}

output "name_filter_is_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by enterprise project ID.	
data "huaweicloud_cae_applications" "filter_by_enterprise_project_id" {
  environment_id        = huaweicloud_cae_application.test[1].environment_id
  enterprise_project_id = huaweicloud_cae_application.test[1].enterprise_project_id
}

# Only check the enterprise_project_id parameter is valid.
output "enterprise_project_id_filter_is_useful" {
  value = length(data.huaweicloud_cae_applications.filter_by_enterprise_project_id.applications) > 0
}
`, testAccDatasourceApplications_base(name))
}

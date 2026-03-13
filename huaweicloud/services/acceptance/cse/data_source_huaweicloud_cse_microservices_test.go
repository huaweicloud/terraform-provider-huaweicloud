package cse

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataMicroservices_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_cse_microservices.test"
		dc  = acceptance.InitDataSourceCheck(all)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCSEMicroserviceEngineID(t)
			acceptance.TestAccPreCheckCSEMicroserviceEngineAdminPassword(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataMicroservices_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "microservices.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckOutput("is_microservice_listed", "true"),
					resource.TestCheckOutput("is_microservice_name_set", "true"),
					resource.TestCheckOutput("is_microservice_app_id_set", "true"),
					resource.TestCheckOutput("is_microservice_version_set", "true"),
					resource.TestCheckOutput("is_microservice_environment_set", "true"),
					resource.TestCheckOutput("is_microservice_description_set", "true"),
					resource.TestCheckOutput("is_microservice_level_set", "true"),
					resource.TestCheckOutput("is_created_at_set_and_valid", "true"),
					resource.TestCheckOutput("is_updated_at_set_and_valid", "true"),
				),
			},
		},
	})
}

func testAccDataMicroservices_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_cse_microservice_engines" "test" {}

locals {
  id_filter_result = [for o in data.huaweicloud_cse_microservice_engines.test.engines : o if o.id == "%[1]s"]
}

resource "huaweicloud_cse_microservice" "test" {
  auth_address    = local.id_filter_result[0].service_registry_addresses[0].public
  connect_address = local.id_filter_result[0].service_registry_addresses[0].public

  name        = "%[2]s"
  app_name    = "%[2]s"
  environment = "development"
  version     = "1.0.1"
  description = "Created by terraform test"
  level       = "BACK"

  admin_user = "root"
  admin_pass = "%[3]s"

  lifecycle {
    ignore_changes = [
      admin_pass,
    ]
  }
}

data "huaweicloud_cse_microservices" "test" {
  depends_on = [huaweicloud_cse_microservice.test]

  auth_address    = local.id_filter_result[0].service_registry_addresses[0].public
  connect_address = local.id_filter_result[0].service_registry_addresses[0].public
  admin_user      = "root"
  admin_pass      = "%[3]s"
}

locals {
  filter_result = try([
    for o in data.huaweicloud_cse_microservices.test.microservices : o if o.id == huaweicloud_cse_microservice.test.id
  ][0], null)
}

output "is_microservice_listed" {
  value = try(length(data.huaweicloud_cse_microservices.test.microservices) > 0, false)
}

output "is_microservice_name_set" {
  value = try(local.filter_result.name == huaweicloud_cse_microservice.test.name, false)
}

output "is_microservice_app_id_set" {
  value = try(local.filter_result.app_id == huaweicloud_cse_microservice.test.app_name, false)
}

output "is_microservice_version_set" {
  value = try(local.filter_result.version == huaweicloud_cse_microservice.test.version, false)
}

output "is_microservice_environment_set" {
  value = try(local.filter_result.environment == huaweicloud_cse_microservice.test.environment, false)
}

output "is_microservice_description_set" {
  value = try(local.filter_result.description == huaweicloud_cse_microservice.test.description, false)
}

output "is_microservice_level_set" {
  value = try(local.filter_result.level == huaweicloud_cse_microservice.test.level, false)
}

output "is_created_at_set_and_valid" {
  value = try(length(regexall("^\\d{4}\\-\\d{2}\\-\\d{2}T\\d{2}:\\d{2}:\\d{2}(?:Z|[+-]\\d{2}:\\d{2})$",
    local.filter_result.created_at)) > 0, false)
}

output "is_updated_at_set_and_valid" {
  value = try(length(regexall("^\\d{4}\\-\\d{2}\\-\\d{2}T\\d{2}:\\d{2}:\\d{2}(?:Z|[+-]\\d{2}:\\d{2})$",
    local.filter_result.updated_at)) > 0, false)
}
`, acceptance.HW_CSE_MICROSERVICE_ENGINE_ID, name, acceptance.HW_CSE_MICROSERVICE_ENGINE_ADMIN_PASSWORD)
}

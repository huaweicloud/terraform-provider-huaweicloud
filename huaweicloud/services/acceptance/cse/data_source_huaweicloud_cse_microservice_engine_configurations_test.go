package cse

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataMicroserviceEngineConfigurations_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_cse_microservice_engine_configurations.test"
		dc  = acceptance.InitDataSourceCheck(all)
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
				Config: testAccDataMicroserviceEngineConfigurations_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_configuration_queried", "true"),
				),
			},
		},
	})
}

func testAccDataMicroserviceEngineConfigurations_basic_base(name string) string {
	return fmt.Sprintf(`
variable "enterprise_project_id" {
  type    = string
  default = "%[1]s"
}

data "huaweicloud_cse_microservice_engines" "test" {}

locals {
  id_filter_result = [
    for o in data.huaweicloud_cse_microservice_engines.test.engines : o if o.id == "%[2]s"
  ]
}

resource "huaweicloud_cse_microservice_engine_configuration" "test" {
  auth_address    = try(local.id_filter_result[0].service_registry_addresses[0].public, null)
  connect_address = try(local.id_filter_result[0].config_center_addresses[0].public, null)
  admin_user      = "root"
  admin_pass      = "%[3]s"

  key        = "%[4]s"
  value_type = "json"
  value      = jsonencode({
    "foo": "baar"
  })

  status                = "disabled"
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  tags = {
    owner = "terraform"
  }

  lifecycle {
    ignore_changes = [
      admin_pass,
    ]
  }
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST,
		acceptance.HW_CSE_MICROSERVICE_ENGINE_ID,
		acceptance.HW_CSE_MICROSERVICE_ENGINE_ADMIN_PASSWORD,
		name)
}

func testAccDataMicroserviceEngineConfigurations_basic() string {
	name := acceptance.RandomAccResourceNameWithDash()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cse_microservice_engine_configurations" "test" {
  depends_on = [huaweicloud_cse_microservice_engine_configuration.test]

  auth_address    = try(local.id_filter_result[0].service_registry_addresses[0].public, null)
  connect_address = try(local.id_filter_result[0].config_center_addresses[0].public, null)
  admin_user      = "root"
  admin_pass      = "%[2]s"

  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

output "is_configuration_queried" {
  value = contains(data.huaweicloud_cse_microservice_engine_configurations.test.configurations[*].id,
    huaweicloud_cse_microservice_engine_configuration.test.id)
}
`, testAccDataMicroserviceEngineConfigurations_basic_base(name),
		acceptance.HW_CSE_MICROSERVICE_ENGINE_ADMIN_PASSWORD)
}

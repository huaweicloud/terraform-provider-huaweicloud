package cse

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataMicroserviceInstances_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		dataSource = "data.huaweicloud_cse_microservice_instances.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
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
				Config:      testAccDataMicroserviceInstances_expectErr(),
				ExpectError: regexp.MustCompile(`Micro-service does not exist`),
			},
			{
				Config: testAccDataMicroserviceInstances_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "instances.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckOutput("is_host_name_set", "true"),
					resource.TestCheckOutput("is_endpoints_set", "true"),
					resource.TestCheckOutput("is_version_set", "true"),
					resource.TestCheckOutput("is_properties_set", "true"),
					resource.TestCheckOutput("is_health_check_set", "true"),
					resource.TestCheckOutput("is_data_center_set", "true"),
					resource.TestCheckOutput("is_created_at_set_and_valid", "true"),
					resource.TestCheckOutput("is_updated_at_set_and_valid", "true"),
				),
			},
		},
	})
}

func testAccDataMicroserviceInstances_expectErr() string {
	randUUID, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_cse_microservice_engines" "test" {}

locals {
  id_filter_result = [for o in data.huaweicloud_cse_microservice_engines.test.engines : o if o.id == "%[1]s"]
}

data "huaweicloud_cse_microservice_instances" "test" {
  auth_address    = local.id_filter_result[0].service_registry_addresses[0].public
  connect_address = local.id_filter_result[0].service_registry_addresses[0].public
  admin_user      = "root"
  admin_pass      = "%[2]s"

  microservice_id = "%[3]s"
}
`, acceptance.HW_CSE_MICROSERVICE_ENGINE_ID, acceptance.HW_CSE_MICROSERVICE_ENGINE_ADMIN_PASSWORD, randUUID)
}

func testAccDataMicroserviceInstances_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cse_microservice_instance" "test" {
  auth_address    = local.id_filter_result[0].service_registry_addresses[0].public
  connect_address = local.id_filter_result[0].service_registry_addresses[0].public
  microservice_id = huaweicloud_cse_microservice.test.id
  host_name       = "localhost_with_auth_address"
  endpoints       = ["grpc://127.0.1.132:9980", "rest://127.0.0.111:8081"]
  version         = "1.0.1"

  properties = {
    "nodeIP" = "127.0.0.1"
  }

  health_check {
    mode        = "push"
    interval    = 30
    max_retries = 3
    port        = 8080
  }

  data_center {
    region            = "%[2]s"
    name              = "dc1"
    availability_zone = data.huaweicloud_availability_zones.test.names[0]
  }

  admin_user = "root"
  admin_pass = "%[3]s"
}

data "huaweicloud_cse_microservice_instances" "test" {
  depends_on = [huaweicloud_cse_microservice_instance.test]

  auth_address    = local.id_filter_result[0].service_registry_addresses[0].public
  connect_address = local.id_filter_result[0].service_registry_addresses[0].public
  admin_user      = "root"
  admin_pass      = "%[3]s"

  microservice_id = huaweicloud_cse_microservice.test.id
}

locals {
  filter_result = try([
  for v in data.huaweicloud_cse_microservice_instances.test.instances : v if v.id == huaweicloud_cse_microservice_instance.test.id][0], null)
}

output "is_host_name_set" {
  value = try(local.filter_result.host_name == huaweicloud_cse_microservice_instance.test.host_name, false)
}

output "is_endpoints_set" {
  value = try(length(local.filter_result.endpoints) == length(huaweicloud_cse_microservice_instance.test.endpoints), false)
}

output "is_version_set" {
  value = try(local.filter_result.version == huaweicloud_cse_microservice_instance.test.version, false)
}

output "is_properties_set" {
  value = try(length(local.filter_result.properties) > 0, false)
}

output "is_health_check_set" {
  value = try(alltrue([
    length(local.filter_result.health_check) > 0,
    local.filter_result.health_check[0].interval == huaweicloud_cse_microservice_instance.test.health_check[0].interval,
    local.filter_result.health_check[0].max_retries == huaweicloud_cse_microservice_instance.test.health_check[0].max_retries,
    local.filter_result.health_check[0].mode == huaweicloud_cse_microservice_instance.test.health_check[0].mode,
    local.filter_result.health_check[0].port == huaweicloud_cse_microservice_instance.test.health_check[0].port,
  ]), false)
}

output "is_data_center_set" {
  value = try(alltrue([
    length(local.filter_result.data_center) > 0,
    local.filter_result.data_center[0].name == huaweicloud_cse_microservice_instance.test.data_center[0].name,
    local.filter_result.data_center[0].region == huaweicloud_cse_microservice_instance.test.data_center[0].region,
    local.filter_result.data_center[0].availability_zone == huaweicloud_cse_microservice_instance.test.data_center[0].availability_zone,
  ]), false)
}

output "is_created_at_set_and_valid" {
  value = try(length(regexall("^\\d{4}\\-\\d{2}\\-\\d{2}T\\d{2}:\\d{2}:\\d{2}(?:Z|[+-]\\d{2}:\\d{2})$",
  local.filter_result.created_at)) > 0, false)
}

output "is_updated_at_set_and_valid" {
  value = try(length(regexall("^\\d{4}\\-\\d{2}\\-\\d{2}T\\d{2}:\\d{2}:\\d{2}(?:Z|[+-]\\d{2}:\\d{2})$",
  local.filter_result.updated_at)) > 0, false)
}
`, testAccMicroserviceInstance_base(name), acceptance.HW_REGION_NAME, acceptance.HW_CSE_MICROSERVICE_ENGINE_ADMIN_PASSWORD)
}

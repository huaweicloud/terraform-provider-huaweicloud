package cse

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataMicroserviceEngineFlavors_basic(t *testing.T) {
	var (
		withDefaultVersion   = "data.huaweicloud_cse_microservice_engine_flavors.version_omitted"
		dcWithDefaultVersion = acceptance.InitDataSourceCheck(withDefaultVersion)

		byVersion   = "data.huaweicloud_cse_microservice_engine_flavors.filter_by_version"
		dcByVersion = acceptance.InitDataSourceCheck(byVersion)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataMicroserviceEngineFlavors_basic,
				Check: resource.ComposeTestCheckFunc(
					dcWithDefaultVersion.CheckResourceExists(),
					resource.TestMatchResourceAttr(withDefaultVersion, "flavors.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckOutput("all_flavor_ids_set", "true"),
					resource.TestCheckOutput("is_available_cpu_memory_set_and_correct", "true"),
					resource.TestCheckOutput("is_linear_set_and_correct", "true"),
					resource.TestCheckOutput("is_available_prefix_set_and_correct", "true"),
					dcByVersion.CheckResourceExists(),
					resource.TestCheckOutput("is_version_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataMicroserviceEngineFlavors_basic string = `
data "huaweicloud_cse_microservice_engine_flavors" "version_omitted" {}

# Check whether flavor ID is set
locals {
  flavor_id_validate_result = [
    for o in data.huaweicloud_cse_microservice_engine_flavors.version_omitted.flavors[*].id : o != ""
  ]
}

output "all_flavor_ids_set" {
  value = length(local.flavor_id_validate_result) > 0 && alltrue(local.flavor_id_validate_result)
}

# Check whether spec.available_cpu_memory is in expected format
locals {
  available_cpu_memory_filter_result = [
    for o in data.huaweicloud_cse_microservice_engine_flavors.version_omitted.flavors : length(regexall("\\d+\\-\\d+",
      o.spec[0].available_cpu_memory)) > 0
  ]
}

output "is_available_cpu_memory_set_and_correct" {
  value = length(local.available_cpu_memory_filter_result) > 0 && alltrue(local.available_cpu_memory_filter_result)
}

# Check whether spec.linear is in expected format
locals {
  linear_filter_result = [
    for o in data.huaweicloud_cse_microservice_engine_flavors.version_omitted.flavors : contains(["true", "false"], o.spec[0].linear)
  ]
}

output "is_linear_set_and_correct" {
  value = length(local.linear_filter_result) > 0 && alltrue(local.linear_filter_result)
}

# Check whether spec.available_prefix is set
locals {
  available_prefix_filter_result = [
    for o in data.huaweicloud_cse_microservice_engine_flavors.version_omitted.flavors : can(regex("[\\w,]+",
      o.spec[0].available_prefix)) && length(split(",", o.spec[0].available_prefix))> 0
  ]
}

output "is_available_prefix_set_and_correct" {
  value = length(local.available_prefix_filter_result) > 0 && alltrue(local.available_prefix_filter_result)
}

data "huaweicloud_cse_microservice_engine_flavors" "filter_by_version" {
  version = "Nacos2"
}

locals {
  version_filter_result = [
    for o in data.huaweicloud_cse_microservice_engine_flavors.filter_by_version.flavors : strcontains(o.id, "nacos2")
  ]
}

output "is_version_filter_useful" {
  value = length(local.version_filter_result) > 0 && alltrue(local.version_filter_result)
}
`

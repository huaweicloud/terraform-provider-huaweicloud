package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running the test, please ensure that you have created a DSC instance and a AES encryption configuration.
func TestAccDataSourceEncryptConfigurations_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dsc_encrypt_configurations.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byConfigurationName   = "data.huaweicloud_dsc_encrypt_configurations.filter_by_configuration_name"
		dcByConfigurationName = acceptance.InitDataSourceCheck(byConfigurationName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDscInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEncryptConfigurations_basic(),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(all, "configurations.#"),
					resource.TestCheckResourceAttrSet(all, "configurations.0.id"),
					resource.TestCheckResourceAttrSet(all, "configurations.0.configuration_name"),
					resource.TestCheckResourceAttrSet(all, "configurations.0.algorithm_name"),
					resource.TestCheckResourceAttrSet(all, "configurations.0.algorithm_type"),
					resource.TestCheckResourceAttrSet(all, "configurations.0.encrypt_mode"),
					resource.TestCheckResourceAttrSet(all, "configurations.0.filling_method"),
					resource.TestCheckResourceAttrSet(all, "configurations.0.kms_context.0.kms_key_id"),
					// Filter by 'configuration_name' parameter.
					dcByConfigurationName.CheckResourceExists(),
					resource.TestCheckOutput("is_configuration_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceEncryptConfigurations_basic() string {
	return `
# Without any filter parameters.
data "huaweicloud_dsc_encrypt_configurations" "test" {
  algorithm_type = "AES"
}

locals {
  configuration_name = try(data.huaweicloud_dsc_encrypt_configurations.test.configurations[0].configuration_name, "NOT_FOUND")
}

# Filter by 'configuration_name' parameter.
data "huaweicloud_dsc_encrypt_configurations" "filter_by_configuration_name" {
  algorithm_type     = "AES"
  configuration_name = local.configuration_name
}

locals {
  configuration_name_filter_result = [
    for v in data.huaweicloud_dsc_encrypt_configurations.filter_by_configuration_name.configurations[*].configuration_name :
    strcontains(v, local.configuration_name)
  ]
}

output "is_configuration_name_filter_useful" {
  value = length(local.configuration_name_filter_result) > 0 && alltrue(local.configuration_name_filter_result)
}
`
}

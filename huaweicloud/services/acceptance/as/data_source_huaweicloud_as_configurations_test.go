package as

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceASConfiguration_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_as_configurations.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byName   = "data.huaweicloud_as_configurations.name_filter"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byImageID   = "data.huaweicloud_as_configurations.image_id_filter"
		dcByImageID = acceptance.InitDataSourceCheck(byImageID)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please configure the scaling group configuration in advance.
			acceptance.TestAccPreCheckASScalingConfigurationID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceASConfiguration_conf,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "configurations.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "configurations.0.instance_config.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "configurations.0.instance_config.0.charging_mode"),
					resource.TestCheckResourceAttrSet(dataSourceName, "configurations.0.instance_config.0.disk.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "configurations.0.instance_config.0.flavor"),
					resource.TestCheckResourceAttrSet(dataSourceName, "configurations.0.instance_config.0.flavor_priority_policy"),
					resource.TestCheckResourceAttrSet(dataSourceName, "configurations.0.instance_config.0.image"),
					resource.TestCheckResourceAttrSet(dataSourceName, "configurations.0.instance_config.0.key_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "configurations.0.scaling_configuration_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "configurations.0.scaling_configuration_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "configurations.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "configurations.0.create_time"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),

					dcByImageID.CheckResourceExists(),
					resource.TestCheckOutput("is_image_id_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceASConfiguration_conf = `
data "huaweicloud_as_configurations" "test" {
}

# Filter by name
locals {
  name = data.huaweicloud_as_configurations.test.configurations[0].scaling_configuration_name
}

data "huaweicloud_as_configurations" "name_filter" {
  name = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_as_configurations.name_filter.configurations[*].scaling_configuration_name : v == local.name
  ]
}

output "is_name_filter_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

# Filter by image_id
locals {
  image_id = data.huaweicloud_as_configurations.test.configurations[0].instance_config[0].image
}

data "huaweicloud_as_configurations" "image_id_filter" {
  image_id = local.image_id
}

locals {
  image_id_filter_result = [
    for v in data.huaweicloud_as_configurations.image_id_filter.configurations[*].instance_config[0].image : v == local.image_id
  ]
}

output "is_image_id_filter_useful" {
  value = alltrue(local.image_id_filter_result) && length(local.image_id_filter_result) > 0
}
`

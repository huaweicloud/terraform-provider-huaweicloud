package cdn

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataStatisticConfiguration_basic(t *testing.T) {
	var (
		rName = "data.huaweicloud_cdn_statistic_configuration.test"
		dc    = acceptance.InitDataSourceCheck(rName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataStatisticConfiguration_invalid,
				ExpectError: regexp.MustCompile(`config type is invalid`),
			},
			{
				Config: testAccDataStatisticConfiguration_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "config_type", "0"),
					resource.TestCheckResourceAttrSet(rName, "configurations.#"),
					resource.TestCheckOutput("is_config_type_set_and_valid", "true"),
					resource.TestCheckOutput("is_resource_type_set_and_valid", "true"),
					resource.TestCheckOutput("is_resource_name_set_and_valid", "true"),
					resource.TestCheckOutput("is_config_info_set_and_valid", "true"),
				),
			},
		},
	})
}

const testAccDataStatisticConfiguration_invalid = `
data "huaweicloud_cdn_statistic_configuration" "test" {
  config_type = 5
}
`

func testAccDataStatisticConfiguration_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_statistic_configuration" "test" {
  resource_type = "domain"
  resource_name = "%[1]s"
  config_type   = 0

  config_info {
    url {
      enable = true
    }
    ua {
      enable = true
    }
  }
}

data "huaweicloud_cdn_statistic_configuration" "test" {
  depends_on = [huaweicloud_cdn_statistic_configuration.test]

  config_type = 0
}

locals {
  configurations = data.huaweicloud_cdn_statistic_configuration.test.configurations
  first_config   = length(local.configurations) > 0 ? local.configurations[0] : null
}

output "is_config_type_set_and_valid" {
  value = local.first_config != null && local.first_config.config_type != null
}

output "is_resource_type_set_and_valid" {
  value = local.first_config != null && local.first_config.resource_type != null
}

output "is_resource_name_set_and_valid" {
  value = local.first_config != null && local.first_config.resource_name != null
}

output "is_config_info_set_and_valid" {
  value = local.first_config != null && length(local.first_config.config_info) > 0
}
`, acceptance.HW_CDN_DOMAIN_NAME)
}

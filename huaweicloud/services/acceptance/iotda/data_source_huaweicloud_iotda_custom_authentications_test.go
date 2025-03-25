package iotda

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCustomAuthentications_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_iotda_custom_authentications.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
		name           = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCustomAuthentications_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "authorizers.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "authorizers.0.authorizer_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "authorizers.0.authorizer_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "authorizers.0.func_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "authorizers.0.func_urn"),
					resource.TestCheckResourceAttrSet(dataSourceName, "authorizers.0.signing_enable"),
					resource.TestCheckResourceAttrSet(dataSourceName, "authorizers.0.default_authorizer"),
					resource.TestCheckResourceAttrSet(dataSourceName, "authorizers.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "authorizers.0.cache_enable"),
					resource.TestCheckResourceAttrSet(dataSourceName, "authorizers.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "authorizers.0.update_time"),

					resource.TestCheckOutput("authorizer_name_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccDataSourceCustomAuthentications_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_iotda_custom_authentications" "test" {
  depends_on = [huaweicloud_iotda_custom_authentication.test]
}

locals {
  authorizer_name = data.huaweicloud_iotda_custom_authentications.test.authorizers[0].authorizer_name
}

data "huaweicloud_iotda_custom_authentications" "authorizer_name_filter" {
  authorizer_name = local.authorizer_name
}

output "authorizer_name_filter_useful" {
  value = length(data.huaweicloud_iotda_custom_authentications.authorizer_name_filter.authorizers) > 0 && alltrue(
    [
      for v in data.huaweicloud_iotda_custom_authentications.authorizer_name_filter.authorizers[*].authorizer_name :
      v == local.authorizer_name
    ]
  )
}

data "huaweicloud_iotda_custom_authentications" "not_found" {
  authorizer_name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_iotda_custom_authentications.not_found.authorizers) == 0
}
`, testAccCustomAuthentication_basic(name))
}

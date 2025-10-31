package nat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceNatPrivateGatewaySpecs_basic(t *testing.T) {
	var (
		natGatewaySpecs = "data.huaweicloud_nat_private_gateway_specs.test"
		dcTags          = acceptance.InitDataSourceCheck(natGatewaySpecs)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceNatPrivateGatewaySpecs_basic,
				Check: resource.ComposeTestCheckFunc(
					dcTags.CheckResourceExists(),

					resource.TestCheckResourceAttrSet(natGatewaySpecs, "specs.#"),
					resource.TestCheckOutput("specs_contains_any", "true"),
				),
			},
		},
	},
	)
}

const testDataSourceNatPrivateGatewaySpecs_basic = `
data "huaweicloud_nat_private_gateway_specs" "test" {
}

output "specs_contains_any" {
  value = length(data.huaweicloud_nat_private_gateway_specs.test.specs) > 0 && anytrue([
    for v in data.huaweicloud_nat_private_gateway_specs.test.specs : (
		(v.name == "Small" && v.code == "1" && v.cbc_code == "privatenat_small") ||
		(v.name == "Medium" && v.code == "2" && v.cbc_code == "privatenat_medium") ||
		(v.name == "Large" && v.code == "3" && v.cbc_code == "privatenat_large") ||
		(v.name == "Extra-large" && v.code == "4" && v.cbc_code == "privatenat_xlarge")
	)
  ])
}

`

package nat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceNatGatewaySpecs_basic(t *testing.T) {
	var (
		natGatewaySpecs = "data.huaweicloud_nat_gateway_specs.test"
		dcTags          = acceptance.InitDataSourceCheck(natGatewaySpecs)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceNatGatewaySpecs_basic,
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

const testDataSourceNatGatewaySpecs_basic = `
data "huaweicloud_nat_gateway_specs" "test" {
}

output "specs_contains_any" {
  value = length(data.huaweicloud_nat_gateway_specs.test.specs) > 0 && anytrue([
    for v in ["1", "2", "3", "4", "5", "6"] : contains(data.huaweicloud_nat_gateway_specs.test.specs, v)
  ])
}

`

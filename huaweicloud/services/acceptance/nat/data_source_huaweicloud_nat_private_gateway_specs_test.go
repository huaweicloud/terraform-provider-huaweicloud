package nat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePrivateGatewaySpecs_basic(t *testing.T) {
	var (
		datasourceName = "data.huaweicloud_nat_private_gateway_specs.test"
		dc             = acceptance.InitDataSourceCheck(datasourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePrivateGatewaySpecs_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),

					resource.TestCheckResourceAttrSet(datasourceName, "specs.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "specs.0.name"),
					resource.TestCheckResourceAttrSet(datasourceName, "specs.0.code"),
					resource.TestCheckResourceAttrSet(datasourceName, "specs.0.cbc_code"),
					resource.TestCheckResourceAttrSet(datasourceName, "specs.0.rule_max"),
					resource.TestCheckResourceAttrSet(datasourceName, "specs.0.sess_max"),
					resource.TestCheckResourceAttrSet(datasourceName, "specs.0.bps_max"),
					resource.TestCheckResourceAttrSet(datasourceName, "specs.0.pps_max"),
					resource.TestCheckResourceAttrSet(datasourceName, "specs.0.qps_max"),
				),
			},
		},
	},
	)
}

const testDataSourcePrivateGatewaySpecs_basic = `data "huaweicloud_nat_private_gateway_specs" "test" {}`

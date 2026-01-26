package nat

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGatewaySpecs_basic(t *testing.T) {
	var (
		datasourceName = "data.huaweicloud_nat_gateway_specs.test"
		dc             = acceptance.InitDataSourceCheck(datasourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGatewaySpecs_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(datasourceName, "specs.#"),
				),
			},
		},
	},
	)
}

const testDataSourceGatewaySpecs_basic = `data "huaweicloud_nat_gateway_specs" "test" {}`

package vpn

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceVpnGatewayFlavors_basic(t *testing.T) {
	rName1 := "data.huaweicloud_vpn_gateway_flavors.test"
	rName2 := "data.huaweicloud_vpn_gateway_flavors.name"
	rName3 := "data.huaweicloud_vpn_gateway_flavors.attachment_type"
	dc := acceptance.InitDataSourceCheck(rName1)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceVpnGatewayFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName1, "flavors.#"),
					resource.TestCheckResourceAttr(rName2, "flavors.0", "basic"),
					resource.TestCheckResourceAttrSet(rName3, "flavors.#"),
				),
			},
		},
	})
}

func testAccDatasourceVpnGatewayFlavors_basic() string {
	return `
data "huaweicloud_vpn_gateway_flavors" "test" {
  availability_zone = "cn-north-4a"
}

data "huaweicloud_vpn_gateway_flavors" "name" {
  availability_zone = "cn-north-4a"
  name              = "basic"
}

data "huaweicloud_vpn_gateway_flavors" "attachment_type" {
  availability_zone = "cn-north-4a"
  attachment_type   = "er"
}
`
}

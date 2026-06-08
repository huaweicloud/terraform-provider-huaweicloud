package vpn

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpnPeerConfigurationSupportedDevices_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpn_peer_configuration_supported_devices.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVpnPeerConfigurationSupportedDevices_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "supported_devices.#"),
					resource.TestCheckResourceAttrSet(dataSource, "supported_devices.0.vendor"),
					resource.TestCheckResourceAttrSet(dataSource, "supported_devices.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "supported_devices.0.model"),
					resource.TestCheckResourceAttrSet(dataSource, "supported_devices.0.version"),
				),
			},
		},
	})
}

func testDataSourceVpnPeerConfigurationSupportedDevices_basic() string {
	return `
data "huaweicloud_vpn_peer_configuration_supported_devices" "test" {}
`
}

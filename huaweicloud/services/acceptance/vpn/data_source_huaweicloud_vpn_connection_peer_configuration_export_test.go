package vpn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpnConnectionPeerConfigurationExport_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpn_connection_peer_configuration_export.test"
	dc := acceptance.InitDataSourceCheck(dataSource)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVpnConnectionPeerConfigurationExport_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "peer_config"),
				),
			},
		},
	})
}

func testDataSourceVpnConnectionPeerConfigurationExport_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpn_peer_configuration_supported_devices" "test" {}

data "huaweicloud_vpn_connection_peer_configuration_export" "test" {
  vpn_connection_id = huaweicloud_vpn_connection.test.id
  vendor            = data.huaweicloud_vpn_peer_configuration_supported_devices.test.supported_devices[0].vendor
  type              = data.huaweicloud_vpn_peer_configuration_supported_devices.test.supported_devices[0].type
  model             = data.huaweicloud_vpn_peer_configuration_supported_devices.test.supported_devices[0].model
  version           = data.huaweicloud_vpn_peer_configuration_supported_devices.test.supported_devices[0].version
}
`, testConnection_basic(name, "172.16.1.4"))
}

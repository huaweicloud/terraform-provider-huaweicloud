package vpn

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpnServers_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpn_servers.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckVPNP2cServer(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVpnServers_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "vpn_servers.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "vpn_servers.0.p2c_vgw_id"),
					resource.TestCheckResourceAttrSet(dataSource, "vpn_servers.0.client_cidr"),
					resource.TestCheckResourceAttrSet(dataSource, "vpn_servers.0.local_subnets.#"),
					resource.TestCheckResourceAttrSet(dataSource, "vpn_servers.0.client_auth_type"),
					resource.TestCheckResourceAttrSet(dataSource, "vpn_servers.0.tunnel_protocol"),
					resource.TestCheckResourceAttrSet(dataSource, "vpn_servers.0.server_certificate.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "vpn_servers.0.server_certificate.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "vpn_servers.0.server_certificate.0.issuer"),
					resource.TestCheckResourceAttrSet(dataSource, "vpn_servers.0.server_certificate.0.subject"),
					resource.TestCheckResourceAttrSet(dataSource, "vpn_servers.0.server_certificate.0.serial_number"),
					resource.TestCheckResourceAttrSet(dataSource, "vpn_servers.0.server_certificate.0.expiration_time"),
					resource.TestCheckResourceAttrSet(dataSource, "vpn_servers.0.server_certificate.0.signature_algorithm"),
					resource.TestCheckResourceAttrSet(dataSource, "vpn_servers.0.ssl_options.0.protocol"),
					resource.TestCheckResourceAttrSet(dataSource, "vpn_servers.0.ssl_options.0.port"),
					resource.TestCheckResourceAttrSet(dataSource, "vpn_servers.0.ssl_options.0.encryption_algorithm"),
					resource.TestCheckResourceAttrSet(dataSource, "vpn_servers.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "vpn_servers.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "vpn_servers.0.updated_at"),
				),
			},
		},
	})
}

const testDataSourceVpnServers_basic = `data "huaweicloud_vpn_servers" "test" {}`

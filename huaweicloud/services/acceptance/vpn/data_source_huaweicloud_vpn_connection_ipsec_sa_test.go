package vpn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpnConnectionIpsecSa_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpn_connection_ipsec_sa.test"
	dc := acceptance.InitDataSourceCheck(dataSource)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVpnConnectionIpsecSa_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "sa_infos.#"),
					resource.TestCheckResourceAttrSet(dataSource, "sa_infos.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "sa_infos.0.source_ip_cidr"),
					resource.TestCheckResourceAttrSet(dataSource, "sa_infos.0.dest_ip_cidr"),
					resource.TestCheckResourceAttrSet(dataSource, "sa_infos.0.packets_sent"),
					resource.TestCheckResourceAttrSet(dataSource, "sa_infos.0.packets_recv"),
					resource.TestCheckResourceAttrSet(dataSource, "sa_infos.0.traffic_sent"),
					resource.TestCheckResourceAttrSet(dataSource, "sa_infos.0.traffic_recv"),
					resource.TestCheckResourceAttrSet(dataSource, "sa_infos.0.collected_at"),
				),
			},
		},
	})
}

func testDataSourceVpnConnectionIpsecSa_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpn_connection_ipsec_sa" "test" {
  vpn_connection_id = huaweicloud_vpn_connection.test.id
}
`, testConnection_basic(name, "172.16.1.4"))
}

package vpn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVPNConnectionHealthChecksDataSource_Basic(t *testing.T) {
	resourceName := "data.huaweicloud_vpn_connection_health_checks.services"
	dc := acceptance.InitDataSourceCheck(resourceName)
	rName := acceptance.RandomAccResourceName()
	ipAddress := "172.16.1.4"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceConnectionHealthChecks_basic(rName, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "connection_id",
						"huaweicloud_vpn_connection.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_ip",
						"huaweicloud_vpn_gateway.test", "master_eip.0.id"),
					resource.TestCheckResourceAttrPair(resourceName, "source_ip",
						"huaweicloud_vpn_customer_gateway.test", "ip"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
		},
	})
}

func testDataSourceConnectionHealthChecks_basic(rName string, ipAddress string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpn_connection_health_checks" "services" {
    status         = "ACTIVE"
    connection_id  = huaweicloud_vpn_connection.test.id
    source_ip      = huaweicloud_vpn_customer_gateway.test.ip
    destination_ip = huaweicloud_vpn_gateway.test.master_eip[0].id
}`, testConnectionHealthCheck_basic(rName, ipAddress))
}

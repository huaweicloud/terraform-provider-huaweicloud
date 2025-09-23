package vpn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVPNConnectionHealthChecksDataSource_Basic(t *testing.T) {
	resourceName := "data.huaweicloud_vpn_connection_health_checks.filter_by_connection_id"
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
					resource.TestCheckResourceAttrSet(resourceName, "connection_health_checks.#"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_health_checks.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_health_checks.0.connection_id"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_health_checks.0.destination_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_health_checks.0.source_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_health_checks.0.status"),
					resource.TestCheckOutput("source_ip_filter_is_useful", "true"),
					resource.TestCheckOutput("connection_id_filter_is_useful", "true"),
					resource.TestCheckOutput("destination_ip_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceConnectionHealthChecks_basic(rName string, ipAddress string) string {
	return fmt.Sprintf(`
%s

locals {
  source_ip      = huaweicloud_vpn_gateway.test.master_eip[0].ip_address
  connection_id  = huaweicloud_vpn_connection.test.id
  destination_ip = huaweicloud_vpn_customer_gateway.test.id_value
}

data "huaweicloud_vpn_connection_health_checks" "filter_by_source_ip" {
  source_ip = local.source_ip
  
  depends_on = [
    huaweicloud_vpn_connection_health_check.test,
  ]
}

data "huaweicloud_vpn_connection_health_checks" "filter_by_connection_id" {
  connection_id = local.connection_id

  depends_on = [
    huaweicloud_vpn_connection_health_check.test,
  ]
}

data "huaweicloud_vpn_connection_health_checks" "filter_by_destination_ip" {
  destination_ip = local.destination_ip

  depends_on = [
    huaweicloud_vpn_connection_health_check.test,
  ]
}

output "source_ip_filter_is_useful" {
  value = length(data.huaweicloud_vpn_connection_health_checks.filter_by_source_ip.connection_health_checks) > 0 && alltrue(
    [for v in data.huaweicloud_vpn_connection_health_checks.filter_by_source_ip.connection_health_checks[*].source_ip : v == local.source_ip]
  )
}

output "connection_id_filter_is_useful" {
  value = length(data.huaweicloud_vpn_connection_health_checks.filter_by_connection_id.connection_health_checks) > 0 && alltrue([
    for v in data.huaweicloud_vpn_connection_health_checks.filter_by_connection_id.connection_health_checks[*].connection_id : 
      v == local.connection_id
  ])
}

output "destination_ip_filter_is_useful" {
  value = length(data.huaweicloud_vpn_connection_health_checks.filter_by_destination_ip.connection_health_checks) > 0 && alltrue([
    for v in data.huaweicloud_vpn_connection_health_checks.filter_by_destination_ip.connection_health_checks[*].destination_ip : 
      v == local.destination_ip
  ])
}
`, testConnectionHealthCheck_basic(rName, ipAddress))
}

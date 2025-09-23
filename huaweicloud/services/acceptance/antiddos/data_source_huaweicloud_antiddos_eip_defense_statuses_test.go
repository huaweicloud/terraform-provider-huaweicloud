package antiddos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEipDefenseStatuses_basic(t *testing.T) {
	dataSource := "data.huaweicloud_antiddos_eip_defense_statuses.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEipDefenseStatuses_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "ddos_status.#"),
					resource.TestCheckResourceAttrSet(dataSource, "ddos_status.0.blackhole_endtime"),
					resource.TestCheckResourceAttrSet(dataSource, "ddos_status.0.protect_type"),
					resource.TestCheckResourceAttrSet(dataSource, "ddos_status.0.traffic_threshold"),
					resource.TestCheckResourceAttrSet(dataSource, "ddos_status.0.http_threshold"),
					resource.TestCheckResourceAttrSet(dataSource, "ddos_status.0.eip_id"),
					resource.TestCheckResourceAttrSet(dataSource, "ddos_status.0.public_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "ddos_status.0.network_type"),
					resource.TestCheckResourceAttrSet(dataSource, "ddos_status.0.status"),

					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_public_ip_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceEipDefenseStatuses_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip" "test" {
  name = "%[1]s"

  publicip {
    type       = "5_bgp"
    ip_version = 4
  }

  bandwidth {
    share_type  = "PER"
    name        = "%[1]s"
    size        = 5
    charge_mode = "traffic"
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name)
}

func testDataSourceEipDefenseStatuses_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_antiddos_eip_defense_statuses" "test" {
  depends_on = [huaweicloud_vpc_eip.test]
}

locals {
  status    = data.huaweicloud_antiddos_eip_defense_statuses.test.ddos_status[0].status
  public_ip = data.huaweicloud_antiddos_eip_defense_statuses.test.ddos_status[0].public_ip
}

data "huaweicloud_antiddos_eip_defense_statuses" "filter_by_status" {
  status = local.status
}

data "huaweicloud_antiddos_eip_defense_statuses" "filter_by_public_ip" {
  ip = local.public_ip
}

locals {
  list_by_status    = data.huaweicloud_antiddos_eip_defense_statuses.filter_by_status.ddos_status
  list_by_public_ip = data.huaweicloud_antiddos_eip_defense_statuses.filter_by_public_ip.ddos_status
}

output "is_status_filter_useful" {
  value = length(local.list_by_status) > 0 && alltrue(
    [for v in local.list_by_status[*].status : v == local.status]
  )
}

output "is_public_ip_filter_useful" {
  value = length(local.list_by_public_ip) > 0 && alltrue(
    [for v in local.list_by_public_ip[*].public_ip : v == local.public_ip]
  )
}
`, testDataSourceEipDefenseStatuses_base(name))
}

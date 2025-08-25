package antiddos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEipDefenseStatusesV2_basic(t *testing.T) {
	dataSource := "data.huaweicloud_antiddosv2_eip_defense_statuses.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEipDefenseStatusesV2_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "ddos_status.#"),
					resource.TestCheckResourceAttrSet(dataSource, "ddos_status.0.floating_ip_id"),
					resource.TestCheckResourceAttrSet(dataSource, "ddos_status.0.floating_ip_address"),
					resource.TestCheckResourceAttrSet(dataSource, "ddos_status.0.product_type"),
					resource.TestCheckResourceAttrSet(dataSource, "ddos_status.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "ddos_status.0.clean_threshold"),
					resource.TestCheckResourceAttrSet(dataSource, "ddos_status.0.block_threshold"),

					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_ip_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceEipDefenseStatusesV2_base(name string) string {
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
}
`, name)
}

func testDataSourceEipDefenseStatusesV2_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_antiddosv2_eip_defense_statuses" "test" {
  depends_on = [huaweicloud_vpc_eip.test]
}

locals {
  status              = data.huaweicloud_antiddosv2_eip_defense_statuses.test.ddos_status[0].status
  floating_ip_address = data.huaweicloud_antiddosv2_eip_defense_statuses.test.ddos_status[0].floating_ip_address
}

data "huaweicloud_antiddosv2_eip_defense_statuses" "filter_by_status" {
  status = local.status
}

data "huaweicloud_antiddosv2_eip_defense_statuses" "filter_by_ip" {
  ips = local.floating_ip_address
}

locals {
  list_by_status = data.huaweicloud_antiddosv2_eip_defense_statuses.filter_by_status.ddos_status
  list_by_ip     = data.huaweicloud_antiddosv2_eip_defense_statuses.filter_by_ip.ddos_status
}

output "is_status_filter_useful" {
  value = length(local.list_by_status) > 0 && alltrue(
    [for v in local.list_by_status[*].status : v == local.status]
  )
}

output "is_ip_filter_useful" {
  value = length(local.list_by_ip) > 0 && alltrue(
    [for v in local.list_by_ip[*].floating_ip_address : v == local.floating_ip_address]
  )
}
`, testDataSourceEipDefenseStatusesV2_base(name))
}

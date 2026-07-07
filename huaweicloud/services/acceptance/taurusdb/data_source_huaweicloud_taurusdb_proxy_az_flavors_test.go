package taurusdb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBProxyAzFlavors_basic(t *testing.T) {
	dataSource := "data.huaweicloud_taurusdb_proxy_az_flavors.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTaurusDBProxyAzFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_flavor_groups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_flavor_groups.0.group_type"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_flavor_groups.0.proxy_flavors.#"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_flavor_groups.0.proxy_flavors.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_flavor_groups.0.proxy_flavors.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_flavor_groups.0.proxy_flavors.0.vcpus"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_flavor_groups.0.proxy_flavors.0.ram"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_flavor_groups.0.proxy_flavors.0.db_type"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_flavor_groups.0.proxy_flavors.0.az_status.%"),
				),
			},
		},
	})
}

func testAccDataSourceTaurusDBProxyAzFlavors_basic() string {
	return `
data "huaweicloud_availability_zones" "test" {}

locals {
  az_codes = join(",", data.huaweicloud_availability_zones.test.names)
}

data "huaweicloud_taurusdb_proxy_az_flavors" "test" {
  az_codes          = local.az_codes
  proxy_engine_name = "taurusproxy"
}
`
}

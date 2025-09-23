package eip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBandwidthAddonPackages_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpc_bandwidth_addon_packages.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckVpcEipBandwidthAddOnPackageEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceBandwidthAddonPackages_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_pkgs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_pkgs.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_pkgs.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_pkgs.0.bandwidth_id"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_pkgs.0.pkg_size"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_pkgs.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_pkgs.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_pkgs.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_pkgs.0.processed_time"),
				),
			},
		},
	})
}

const testDataSourceBandwidthAddonPackages_basic = `data "huaweicloud_vpc_bandwidth_addon_packages" "test" {}`

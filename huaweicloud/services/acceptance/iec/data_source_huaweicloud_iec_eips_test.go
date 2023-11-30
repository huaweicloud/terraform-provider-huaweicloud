package iec

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEIPsDataSource_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_iec_eips.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPsDataSource_config,
			},
			{
				Config: testAccEIPsDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "eips.0.ip_version", "4"),
					resource.TestCheckResourceAttr(dataSourceName, "eips.0.bandwidth_share_type", "WHOLE"),
					resource.TestCheckResourceAttrSet(dataSourceName, "site_info"),
					resource.TestCheckResourceAttrSet(dataSourceName, "eips.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "eips.1.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "eips.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "eips.0.public_ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "eips.0.bandwidth_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "eips.0.bandwidth_size"),
				),
			},
		},
	})
}

var testAccEIPsDataSource_config = `
data "huaweicloud_iec_sites" "sites_test" {}

resource "huaweicloud_iec_eip" "eip_test1" {
  site_id = data.huaweicloud_iec_sites.sites_test.sites[0].id
}

resource "huaweicloud_iec_eip" "eip_test2" {
  site_id = data.huaweicloud_iec_sites.sites_test.sites[0].id
}
`

func testAccEIPsDataSource_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_iec_eips" "test" {
  depends_on = [
    huaweicloud_iec_eip.eip_test1,
    huaweicloud_iec_eip.eip_test2,
  ]

  site_id = data.huaweicloud_iec_sites.sites_test.sites[0].id
}
`, testAccEIPsDataSource_config)
}

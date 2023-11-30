package iec

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBandWidthsDataSource_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_iec_bandwidths.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBWsDataSource_config,
			},
			{
				Config: testAccBWsDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "bandwidths.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "site_info"),
					resource.TestCheckResourceAttrSet(dataSourceName, "bandwidths.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "bandwidths.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "bandwidths.0.size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "bandwidths.0.line"),
					resource.TestCheckResourceAttrSet(dataSourceName, "bandwidths.0.status"),
				),
			},
		},
	})
}

var testAccBWsDataSource_config = `
data "huaweicloud_iec_sites" "sites_test" {}

resource "huaweicloud_iec_eip" "eip_test1" {
  site_id = data.huaweicloud_iec_sites.sites_test.sites[0].id
}

resource "huaweicloud_iec_eip" "eip_test2" {
  site_id = data.huaweicloud_iec_sites.sites_test.sites[0].id
  line_id = data.huaweicloud_iec_sites.sites_test.sites[0].lines[0].id
}
`

func testAccBWsDataSource_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_iec_bandwidths" "test" {
  depends_on = [
    huaweicloud_iec_eip.eip_test1,
    huaweicloud_iec_eip.eip_test2,
  ]

  site_id = data.huaweicloud_iec_sites.sites_test.sites[0].id
}
`, testAccBWsDataSource_config)
}

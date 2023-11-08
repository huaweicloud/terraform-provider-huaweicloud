package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIECBandWidthsDataSource_basic(t *testing.T) {
	resourceName := "data.huaweicloud_iec_bandwidths.bandwidths"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		CheckDestroy:              testAccCheckIecEIPDestroy,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccIECBWsDataSource_config,
			},
			{
				Config: testAccIECBWsDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccIECBandWidthsDataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "bandwidths.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "site_info"),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidths.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidths.1.id"),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidths.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidths.0.size"),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidths.0.line"),
					resource.TestCheckResourceAttrSet(resourceName, "bandwidths.0.status"),
				),
			},
		},
	})
}

func testAccIECBandWidthsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find IEC public IPs data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("IEC public IPs data source ID not set")
		}

		return nil
	}
}

var testAccIECBWsDataSource_config string = `
data "huaweicloud_iec_sites" "sites_test" {}

resource "huaweicloud_iec_eip" "eip_test1" {
  site_id = data.huaweicloud_iec_sites.sites_test.sites[0].id
}

resource "huaweicloud_iec_eip" "eip_test2" {
  site_id = data.huaweicloud_iec_sites.sites_test.sites[0].id
  line_id = data.huaweicloud_iec_sites.sites_test.sites[0].lines[1].id
}
`

func testAccIECBWsDataSource_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_iec_bandwidths" "bandwidths" {
  site_id = data.huaweicloud_iec_sites.sites_test.sites[0].id
}
`, testAccIECBWsDataSource_config)
}

package dns

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDNSPublicZoneLines_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dns_public_zone_lines.test"
	name := fmt.Sprintf("acpttest-zone-%s.com", acctest.RandString(5))
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDNSPublicZoneLines_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "lines.#"),
					resource.TestCheckResourceAttrSet(dataSource, "lines.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "lines.0.line"),
					resource.TestCheckResourceAttrSet(dataSource, "lines.0.create_time"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDNSPublicZoneLines_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dns_public_zone_lines" "test" {
  zone_id = huaweicloud_dns_zone.test.id
}
`, testAccDNSZone_basic(name))
}

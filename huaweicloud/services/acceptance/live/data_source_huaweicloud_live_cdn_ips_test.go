package live

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCdnIps_basic(t *testing.T) {
	dataSource := "data.huaweicloud_live_cdn_ips.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCdnIps_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "cdn_ips.#", "2"),
				),
			},
		},
	})
}

const testDataSourceCdnIps_basic = `
data "huaweicloud_live_cdn_ips" "test" {
  ip_addresses = ["192.168.0.1", "192.168.0.2"]
}
`

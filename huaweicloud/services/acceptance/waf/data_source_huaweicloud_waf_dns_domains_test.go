package waf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDNSDomains_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_waf_dns_domains.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDNSDomains_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.domain"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.servers.#"),
				),
			},
		},
	})
}

const testAccDataSourceDNSDomains_basic = `
data "huaweicloud_waf_dns_domains" "test" {}
`

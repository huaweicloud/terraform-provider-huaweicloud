package waf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSourceIps_basic(t *testing.T) {
	dataSource := "data.huaweicloud_waf_source_ips.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceWafSourceIps_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "source_ips.#"),
					resource.TestCheckResourceAttrSet(dataSource, "source_ips.0.ips.#"),
					resource.TestCheckResourceAttrSet(dataSource, "source_ips.0.update_time"),
				),
			},
		},
	})
}

const testDataSourceDataSourceWafSourceIps_basic = `
data "huaweicloud_waf_source_ips" "test" {
}
`

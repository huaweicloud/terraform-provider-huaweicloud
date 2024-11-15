package cc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcCentralNetworkQuotas_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_central_network_quotas.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcCentralNetworkQuotas_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.quota_key"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.quota_limit"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.unit"),
				),
			},
			{
				Config: testDataSourceCcCentralNetworkQuotas_quotaType(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "quotas.0.quota_key", "central_networks_per_account"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.quota_limit"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.unit"),
				),
			},
		},
	})
}

const testDataSourceCcCentralNetworkQuotas_basic = `data "huaweicloud_cc_central_network_quotas" "test" {}`

func testDataSourceCcCentralNetworkQuotas_quotaType() string {
	return `
data "huaweicloud_cc_central_network_quotas" "test" {
  quota_type = ["central_networks_per_account"]
}`
}

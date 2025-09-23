package waf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcAttackEvents_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_waf_events.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			// Because there is no available data for testing, the test case is only
			// used to verify that the API can be invoked.
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcAttackEvents_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.#"),

					resource.TestCheckOutput("test_verify", "true"),
				),
			},
		},
	})
}

const testAccDataSourcAttackEvents_basic = `
data "huaweicloud_waf_events" "test" {
  recent = "3days"
}

output "test_verify" {
  value = length(data.huaweicloud_waf_events.test.items) == 0
}
`

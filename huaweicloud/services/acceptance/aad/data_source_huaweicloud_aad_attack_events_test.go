package antiddos

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Due to testing environment limitations, this test case can only test the scenario with empty `events`.
func TestAccDataSourceAadAttackEvents_basic(t *testing.T) {
	dataSource := "data.huaweicloud_aad_attack_events.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAadAttackEvents_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "events.#"),
				),
			},
		},
	})
}

const testDataSourceAadAttackEvents_basic = `
data "huaweicloud_aad_attack_events" "test" {
  start_time    = "1755734400"
  end_time      = "1755820800"
  recent        = "today"
  overseas_type = "0"
}
`

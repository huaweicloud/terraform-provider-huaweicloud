package antiddos

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Due to testing environment limitations, this test case can only test the scenario with empty `curve`.
func TestAccDataSourceBandwidthCurve_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_aad_bandwidth_curve.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceBandwidthCurve_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "curve.#"),
				),
			},
		},
	})
}

func testDataSourceBandwidthCurve_basic() string {
	return `
data "huaweicloud_aad_bandwidth_curve" "test" {
  value_type    = "mean"
  recent        = "1month"
  overseas_type = "0"
}
`
}

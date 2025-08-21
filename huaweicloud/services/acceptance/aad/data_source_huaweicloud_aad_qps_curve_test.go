package antiddos

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Due to testing environment limitations, this test case can only test the scenario with empty `curve`.
func TestAccDataSourceQPSCurve_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_aad_qps_curve.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceQPSCurve_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "curve.#"),
				),
			},
		},
	})
}

const testDataSourceQPSCurve_basic = `
data "huaweicloud_aad_qps_curve" "test" {
  value_type = "mean"
  recent     = "today"
}
`

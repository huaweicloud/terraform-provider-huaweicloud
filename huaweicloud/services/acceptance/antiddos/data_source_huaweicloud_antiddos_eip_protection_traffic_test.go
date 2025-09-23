package antiddos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Please prepare the EIP ID and IP address in advance and ensure that the EIP has access traffic.
func TestAccDataSourceEipProtectionTraffic_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_antiddos_eip_protection_traffic.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckEipIDAndIP(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEipProtectionTraffic_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.period_start"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.bps_in"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.bps_attack"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.total_bps"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.pps_in"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.pps_attack"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.total_pps"),
				),
			},
		},
	})
}

func testDataSourceEipProtectionTraffic_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_antiddos_eip_protection_traffic" "test" {
  floating_ip_id = "%s"
  ip             = "%s"
}
`, acceptance.HW_EIP_ID, acceptance.HW_EIP_ADDRESS)
}

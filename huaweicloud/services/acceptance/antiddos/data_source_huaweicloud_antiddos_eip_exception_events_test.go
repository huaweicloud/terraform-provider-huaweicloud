package antiddos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Please prepare the EIP ID and IP address in advance and ensure that the EIP has been bound to Anti-DDoS.
// Due to testing environment limitations, this test case can only test the scenario with empty `logs`.
func TestAccDataSourceEipExceptionEvents_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_antiddos_eip_exception_events.test"
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
				Config: testDataSourceEipExceptionEvents_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "logs.#"),
				),
			},
		},
	})
}

func testDataSourceEipExceptionEvents_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_antiddos_eip_exception_events" "test" {
  floating_ip_id = "%s"
  ip             = "%s"
  sort_dir       = "asc"
}
`, acceptance.HW_EIP_ID, acceptance.HW_EIP_ADDRESS)
}

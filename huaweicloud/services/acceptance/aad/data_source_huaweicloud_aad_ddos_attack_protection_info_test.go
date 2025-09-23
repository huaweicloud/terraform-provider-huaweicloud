package antiddos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Note: Due to limited test conditions, this test case cannot be executed successfully.
func TestAccDataSourceDdosAttackProtectionInfo_basic(t *testing.T) {
	var (
		dataSourceBps = "data.huaweicloud_aad_ddos_attack_protection_info.test_bps"
		dcBps         = acceptance.InitDataSourceCheck(dataSourceBps)

		dataSourcePps = "data.huaweicloud_aad_ddos_attack_protection_info.test_pps"
		dcPps         = acceptance.InitDataSourceCheck(dataSourcePps)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAadInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDdosAttackProtectionInfo_basic(),
				Check: resource.ComposeTestCheckFunc(
					dcBps.CheckResourceExists(),
					dcPps.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceBps, "flow_bps.#"),
					resource.TestCheckResourceAttrSet(dataSourcePps, "flow_pps.#"),
				),
			},
		},
	})
}

// Parameter `ip` are mock data.

func testDataSourceDdosAttackProtectionInfo_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_aad_ddos_attack_protection_info" "test_bps" {
  instance_id = "%[1]s"
  ip          = "12.1.2.117"
  type        = "bps"
  start_time  = "1755734400"
  end_time    = "1755820800"
}

data "huaweicloud_aad_ddos_attack_protection_info" "test_pps" {
  instance_id = "%[1]s"
  ip          = "12.1.2.117"
  type        = "pps"
  start_time  = "1755734400"
  end_time    = "1755820800"
}
`, acceptance.HW_AAD_INSTANCE_ID)
}

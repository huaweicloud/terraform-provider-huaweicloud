package antiddos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Note: Due to limited test conditions, this test case cannot be executed successfully.
func TestAccIpDdosStatisticsDataSource_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_aad_ip_ddos_statistics.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAadInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testIpDdosStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "attack_kbps_peak"),
					resource.TestCheckResourceAttrSet(dataSourceName, "in_kbps_peak"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ddos_count"),
					resource.TestCheckResourceAttrSet(dataSourceName, "timestamp"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vip"),
				),
			},
		},
	})
}

// Parameter `ip` are mock data.
func testIpDdosStatistics_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_aad_ip_ddos_statistics" "test" {
  instance_id = "%[1]s"
  ip          = "12.1.2.117"
  start_time  = "1755734400"
  end_time    = "1755820800"
}
`, acceptance.HW_AAD_INSTANCE_ID)
}

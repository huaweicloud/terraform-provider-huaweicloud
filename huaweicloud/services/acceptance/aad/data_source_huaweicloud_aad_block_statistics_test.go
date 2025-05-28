package antiddos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAadBlockStatisticsDataSource_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_aad_block_statistics.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare domain_id before running this test cases.
			acceptance.TestAccPrecheckDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAadBlockStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_unblocking_times"),
					resource.TestCheckResourceAttrSet(dataSourceName, "manual_unblocking_times"),
					resource.TestCheckResourceAttrSet(dataSourceName, "automatic_unblocking_times"),
					resource.TestCheckResourceAttrSet(dataSourceName, "current_blocked_ip_numbers"),
				),
			},
		},
	})
}

func testAadBlockStatistics_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_aad_block_statistics" "test" {
  domain_id = "%[1]s"
}
`, acceptance.HW_DOMAIN_ID)
}

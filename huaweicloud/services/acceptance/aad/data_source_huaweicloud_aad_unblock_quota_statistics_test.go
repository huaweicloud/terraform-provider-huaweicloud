package antiddos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccUnblockQuotaStatisticsDataSource_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_aad_unblock_quota_statistics.test"
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
				Config: testUnblockQuotaStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_unblocking_quota"),
					resource.TestCheckResourceAttrSet(dataSourceName, "remaining_unblocking_quota"),
					resource.TestCheckResourceAttrSet(dataSourceName, "unblocking_quota_today"),
					resource.TestCheckResourceAttrSet(dataSourceName, "remaining_unblocking_quota_today"),
				),
			},
		},
	})
}

func testUnblockQuotaStatistics_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_aad_unblock_quota_statistics" "test" {
  domain_id = "%[1]s"
}
`, acceptance.HW_DOMAIN_ID)
}

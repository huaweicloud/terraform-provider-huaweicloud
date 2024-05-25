package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRmsAggregationPendingRequests_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_rms_resource_aggregation_pending_requests.basic"
	dataSource2 := "data.huaweicloud_rms_resource_aggregation_pending_requests.filter_by_account"
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)

	accountID := acctest.RandStringFromCharSet(32, randomCharSet)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRmsAggregationPendingRequests_basic(accountID),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
				),
			},
		},
	})
}

func testDataSourceDataSourceRmsAggregationPendingRequests_basic(accountID string) string {
	return fmt.Sprintf(`
data "huaweicloud_rms_resource_aggregation_pending_requests" "basic" {}

data "huaweicloud_rms_resource_aggregation_pending_requests" "filter_by_account" {
  account_id = "%s"
}
`, accountID)
}

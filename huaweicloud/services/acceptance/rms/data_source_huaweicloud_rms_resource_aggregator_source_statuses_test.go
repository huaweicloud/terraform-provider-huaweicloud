package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRmsAggregatorSourceStatuses_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_rms_resource_aggregator_source_statuses.basic"
	dataSource2 := "data.huaweicloud_rms_resource_aggregator_source_statuses.filter_by_status"
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)

	rName := acceptance.RandomAccResourceName()
	account1 := acctest.RandStringFromCharSet(32, randomCharSet)
	account2 := acctest.RandStringFromCharSet(32, randomCharSet)
	accounts := []string{account1, account2}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRmsAggregatorSourceStatuses_basic(rName, accounts),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceRmsAggregatorSourceStatuses_basic(name string, accounts []string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rms_resource_aggregator_source_statuses" "basic" {
  aggregator_id = huaweicloud_rms_resource_aggregator.test.id

  depends_on = [huaweicloud_rms_resource_aggregator.test]
}

data "huaweicloud_rms_resource_aggregator_source_statuses" "filter_by_status" {
  aggregator_id = huaweicloud_rms_resource_aggregator.test.id
  status        = "FAILED"

  depends_on = [huaweicloud_rms_resource_aggregator.test]
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_rms_resource_aggregator_source_statuses.filter_by_status.aggregated_source_statuses[*].last_update_status :
  v == "FAILED"]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_rms_resource_aggregator_source_statuses.basic.aggregated_source_statuses) > 0
}

output "is_status_filter_useful" {
  value = alltrue(local.type_filter_result) && length(local.type_filter_result) > 0
}
`, testAggregator_config(name, accounts))
}

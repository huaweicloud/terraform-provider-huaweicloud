package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRmsAggregators_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_rms_resource_aggregators.basic"
	dataSource2 := "data.huaweicloud_rms_resource_aggregators.filter_by_name"
	dataSource3 := "data.huaweicloud_rms_resource_aggregators.filter_by_type"
	dataSource4 := "data.huaweicloud_rms_resource_aggregators.filter_by_id"
	rName := acceptance.RandomAccResourceName()
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)
	dc3 := acceptance.InitDataSourceCheck(dataSource3)
	dc4 := acceptance.InitDataSourceCheck(dataSource4)

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
				Config: testDataSourceDataSourceRmsAggregators_basic(rName, accounts),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					dc3.CheckResourceExists(),
					dc4.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceRmsAggregators_basic(name string, accounts []string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rms_resource_aggregators" "basic" {
  depends_on = [huaweicloud_rms_resource_aggregator.test]
}

data "huaweicloud_rms_resource_aggregators" "filter_by_name" {
  name = "%[2]s"

  depends_on = [huaweicloud_rms_resource_aggregator.test]
}

data "huaweicloud_rms_resource_aggregators" "filter_by_type" {
  type = "ACCOUNT"

  depends_on = [huaweicloud_rms_resource_aggregator.test]
}

data "huaweicloud_rms_resource_aggregators" "filter_by_id" {
  aggregator_id = huaweicloud_rms_resource_aggregator.test.id

  depends_on = [huaweicloud_rms_resource_aggregator.test]
}

locals {
  name_filter_result = [for v in data.huaweicloud_rms_resource_aggregators.filter_by_name.aggregators[*].name : v == "%[2]s"]
  type_filter_result = [for v in data.huaweicloud_rms_resource_aggregators.filter_by_type.aggregators[*].type : v == "ACCOUNT"]
  id_filter_result   = [
    for v in data.huaweicloud_rms_resource_aggregators.filter_by_id.aggregators[*].id : v == huaweicloud_rms_resource_aggregator.test.id]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_rms_resource_aggregators.basic.aggregators) > 0
}

output "is_name_filter_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

output "is_type_filter_useful" {
  value = alltrue(local.type_filter_result) && length(local.type_filter_result) > 0
}

output "is_id_filter_useful" {
  value = alltrue(local.id_filter_result) && length(local.id_filter_result) > 0
}
`, testAggregator_config(name, accounts), name)
}

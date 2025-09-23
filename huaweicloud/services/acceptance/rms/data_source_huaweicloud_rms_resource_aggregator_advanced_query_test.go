package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAggregatorAdvancedQuery_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rms_resource_aggregator_advanced_query.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAggregatorAdvancedQuery_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_name_set", "true"),
					resource.TestCheckOutput("is_id_set", "true"),
					resource.TestCheckOutput("is_query_info_correct", "true"),
				),
			},
		},
	})
}

func testDataSourceAggregatorAdvancedQuery_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rms_resource_aggregator" "test" {
  name        = "%[1]s"
  type        = "ACCOUNT"
  account_ids = ["%[2]s"]

  # wait 30 seconds to let the policies evaluate
  provisioner "local-exec" {
    command = "sleep 30"
  }
}
`, name, acceptance.HW_DOMAIN_ID)
}

func testDataSourceAggregatorAdvancedQuery_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rms_resource_aggregator_advanced_query" "test" {
  aggregator_id = huaweicloud_rms_resource_aggregator.test.id
  expression    = "select name, id from aggregator_resources where provider = 'ecs' and type = 'cloudservers'"
}

locals {
  name_set = [
    for v in data.huaweicloud_rms_resource_aggregator_advanced_query.test.results[*].name : v != ""
  ]
  id_set = [
    for v in data.huaweicloud_rms_resource_aggregator_advanced_query.test.results[*].id : v != ""
  ]
  query_info_correct = [
    for v in data.huaweicloud_rms_resource_aggregator_advanced_query.test.query_info[*].select_fields :
	length(setsubtract(v, ["name", "id"])) == 0
  ]
}

output "is_name_set" {
  value = alltrue(local.name_set) && length(local.name_set) > 0
}

output "is_id_set" {
  value = alltrue(local.id_set) && length(local.id_set) > 0
}

output "is_query_info_correct" {
  value = alltrue(local.query_info_correct) && length(local.query_info_correct) > 0
}
`, testDataSourceAggregatorAdvancedQuery_base(name), name)
}

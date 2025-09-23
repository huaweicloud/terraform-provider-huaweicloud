package rms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAdvancedQuery_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rms_advanced_query.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAdvancedQuery_basic,
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

const testDataSourceAdvancedQuery_basic = `
data "huaweicloud_rms_advanced_query" "test" {
  expression = "select name, id from tracked_resources where provider = 'ecs' and type = 'cloudservers'"
}

locals {
  name_set = [
    for v in data.huaweicloud_rms_advanced_query.test.results[*].name : v != ""
  ]
  id_set = [
    for v in data.huaweicloud_rms_advanced_query.test.results[*].id : v != ""
  ]
  query_info_correct = [
    for v in data.huaweicloud_rms_advanced_query.test.query_info[*].select_fields :
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
`

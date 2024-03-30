package rms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRmsAdvancedQuerySchemas_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_rms_advanced_query_schemas.basic"
	dataSource2 := "data.huaweicloud_rms_advanced_query_schemas.filter_by_type"
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRmsAdvancedQuerySchemas_basic,
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
				),
			},
		},
	})
}

const testDataSourceDataSourceRmsAdvancedQuerySchemas_basic = `
data "huaweicloud_rms_advanced_query_schemas" "basic" {}

data "huaweicloud_rms_advanced_query_schemas" "filter_by_type" {
  type = "ecs.cloudservers"
}

locals {
  type_filter_result = [for v in data.huaweicloud_rms_advanced_query_schemas.filter_by_type.schemas[*].type : v == "ecs.cloudservers"]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_rms_advanced_query_schemas.basic.schemas) > 0
}

output "is_type_filter_useful" {
  value = alltrue(local.type_filter_result) && length(local.type_filter_result) > 0
}
`

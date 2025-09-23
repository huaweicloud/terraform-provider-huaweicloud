package codeartspipeline

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePipelineRules_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_pipeline_rules.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePipelineRules_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "rules.#"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.operator"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.operate_time"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourcePipelineRules_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_pipeline_rules" "test" {}

// filter by name
data "huaweicloud_codearts_pipeline_rules" "filter_by_name" {
  name = huaweicloud_codearts_pipeline_rule.test.name
}

locals {
  filter_result_by_name = [for v in data.huaweicloud_codearts_pipeline_rules.filter_by_name.rules[*].name :
    v == huaweicloud_codearts_pipeline_rule.test.name]
}

output "is_name_filter_useful" {
  value = length(local.filter_result_by_name) > 0 && alltrue(local.filter_result_by_name)
}

// filter by type
data "huaweicloud_codearts_pipeline_rules" "filter_by_type" {
  type = huaweicloud_codearts_pipeline_rule.test.type
}

locals {
  filter_result_by_type = [for v in data.huaweicloud_codearts_pipeline_rules.filter_by_type.rules[*].type :
    v == huaweicloud_codearts_pipeline_rule.test.type]
}

output "is_type_filter_useful" {
  value = length(local.filter_result_by_type) > 0 && alltrue(local.filter_result_by_type)
}
`, testPipelineRule_basic("0.0.3", name))
}

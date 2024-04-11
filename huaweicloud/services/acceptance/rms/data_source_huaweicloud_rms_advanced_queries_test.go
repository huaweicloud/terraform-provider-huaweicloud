package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRmsAdvancedQueries_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_rms_advanced_queries.basic"
	dataSource2 := "data.huaweicloud_rms_advanced_queries.filter_by_name"
	rName := acceptance.RandomAccResourceName()
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRmsAdvancedQueries_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceRmsAdvancedQueries_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rms_advanced_queries" "basic" {
  depends_on = [huaweicloud_rms_advanced_query.test]
}

data "huaweicloud_rms_advanced_queries" "filter_by_name" {
  name = "%[2]s"

  depends_on = [huaweicloud_rms_advanced_query.test]
}

locals {
  name_filter_result = [for v in data.huaweicloud_rms_advanced_queries.filter_by_name.queries[*].name : v == "%[2]s"]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_rms_advanced_queries.basic.queries) > 0
}

output "is_name_filter_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}
`, testAdvancedQuery_basic(name), name)
}

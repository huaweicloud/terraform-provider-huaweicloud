package accessanalyzer

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAccessAnalyzers_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_access_analyzers.basic"
	dataSource2 := "data.huaweicloud_access_analyzers.filter_by_type"
	rName := acceptance.RandomAccResourceName()
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccessAnalyzers_basic(rName),
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

func testAccDataSourceAccessAnalyzers_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_access_analyzers" "basic" {
  depends_on = [huaweicloud_access_analyzer.test]
}

data "huaweicloud_access_analyzers" "filter_by_type" {
  type = "account"

  depends_on = [huaweicloud_access_analyzer.test]
}


locals {
  type_filter_result = [for v in data.huaweicloud_access_analyzers.filter_by_type.analyzers[*].type : v == "account"]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_access_analyzers.basic.analyzers) > 0
}

output "is_type_filter_useful" {
  value = alltrue(local.type_filter_result) && length(local.type_filter_result) > 0
}
`, testAccAnalyzer_basic(name))
}

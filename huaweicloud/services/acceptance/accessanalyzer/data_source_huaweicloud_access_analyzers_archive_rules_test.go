package accessanalyzer

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAccessAnalyzerArchiveRules_basic(t *testing.T) {
	dataSource := "data.huaweicloud_access_analyzer_archive_rules.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccessAnalyzerArchiveRules_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAccessAnalyzerArchiveRules_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_access_analyzer_archive_rules" "test" {
  analyzer_id = huaweicloud_access_analyzer.test.id

  depends_on = [huaweicloud_access_analyzer_archive_rule.test]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_access_analyzer_archive_rules.test.archive_rules) > 0
}
`, testAccArchiveRule_basic(name))
}

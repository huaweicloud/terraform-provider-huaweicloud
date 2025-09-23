package aom

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAlarmRulesTemplates_basic(t *testing.T) {
	dataSource := "data.huaweicloud_aom_alarm_rules_templates.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAlarmRulesTemplates_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "templates.#"),

					resource.TestCheckOutput("id_filter_validation", "true"),
				),
			},
		},
	})
}

func testDataSourceAlarmRulesTemplates_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_aom_alarm_rules_templates" "test" {
  depends_on = [huaweicloud_aom_alarm_rules_template.test]
}

data "huaweicloud_aom_alarm_rules_templates" "id_filter" {
  template_id = huaweicloud_aom_alarm_rules_template.test.id
}

locals {
  test_id_filter_results = data.huaweicloud_aom_alarm_rules_templates.id_filter
}

output "id_filter_validation" {
  value = alltrue([for v in local.test_id_filter_results.templates[*].id : v == huaweicloud_aom_alarm_rules_template.test.id])
}
`, testAlarmRulesTemplate_basic(name))
}

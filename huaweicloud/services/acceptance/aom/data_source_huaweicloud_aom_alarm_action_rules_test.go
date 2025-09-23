package aom

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAomAlarmActionRules_basic(t *testing.T) {
	dataSource := "data.huaweicloud_aom_alarm_action_rules.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceAomAlarmActionRules_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "action_rules.#"),
					resource.TestCheckResourceAttrSet(dataSource, "action_rules.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "action_rules.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "action_rules.0.notification_template"),
					resource.TestCheckResourceAttrSet(dataSource, "action_rules.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "action_rules.0.smn_topics.#"),
					resource.TestMatchResourceAttr(dataSource, "action_rules.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testDataSourceDataSourceAomAlarmActionRules_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_aom_alarm_action_rules" "test" {
  depends_on = [huaweicloud_aom_alarm_action_rule.test]
}
`, testAlarmActionRule_basic(name))
}

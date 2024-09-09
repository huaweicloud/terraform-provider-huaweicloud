package aom

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAlarmGroupRules_basic(t *testing.T) {
	dataSource := "data.huaweicloud_aom_alarm_group_rules.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAlarmGroupRules_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "rules.#"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.detail.#"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.group_by.#"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.group_wait"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.group_interval"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.group_repeat_waiting"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.enterprise_project_id"),
				),
			},
		},
	})
}

func testDataSourceAlarmGroupRules_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_aom_alarm_group_rules" "test" {
  depends_on = [huaweicloud_aom_alarm_group_rule.test]

  enterprise_project_id = huaweicloud_aom_alarm_group_rule.test.enterprise_project_id
}
`, testAlarmGroupRule_basic(name))
}

package aom

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAlarmSilenceRules_basic(t *testing.T) {
	dataSource := "data.huaweicloud_aom_alarm_silence_rules.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAlarmSilenceRules_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "rules.#"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.time_zone"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.silence_time.#"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.silence_conditions.#"),
					resource.TestMatchResourceAttr(dataSource,
						"rules.0.created_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(dataSource,
						"rules.0.updated_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testDataSourceAlarmSilenceRules_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_aom_alarm_silence_rules" "test" {
  depends_on = [huaweicloud_aom_alarm_silence_rule.test]
}
`, testAlarmSilenceRule_basic(name))
}

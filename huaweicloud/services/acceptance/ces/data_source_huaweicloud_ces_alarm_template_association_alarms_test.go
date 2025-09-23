package ces

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCesAlarmTemplateAssociationAlarms_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ces_alarm_template_association_alarms.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCesAlarmTemplateAssociatedWithAlarmRules(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCesAlarmTemplateAssociationAlarms_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.alarm_id"),
					resource.TestCheckResourceAttrSet(dataSource, "alarms.0.name"),
				),
			},
		},
	})
}

func testDataSourceCesAlarmTemplateAssociationAlarms_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_ces_alarm_template_association_alarms" "test" {
  template_id = "%[1]s" 
}
`, acceptance.HW_CES_ALARM_TEMPLATE_ID)
}

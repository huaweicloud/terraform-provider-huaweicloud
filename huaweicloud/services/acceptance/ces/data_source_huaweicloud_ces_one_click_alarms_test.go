package ces

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCesOneClickAlarms_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ces_one_click_alarms.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCesOneClickAlarms_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "one_click_alarms.0.one_click_alarm_id"),
					resource.TestCheckResourceAttrSet(dataSource, "one_click_alarms.0.namespace"),
					resource.TestCheckResourceAttrSet(dataSource, "one_click_alarms.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "one_click_alarms.0.enabled"),
				),
			},
		},
	})
}

func testDataSourceCesOneClickAlarms_basic() string {
	return `data "huaweicloud_ces_one_click_alarms" "test" {}`
}

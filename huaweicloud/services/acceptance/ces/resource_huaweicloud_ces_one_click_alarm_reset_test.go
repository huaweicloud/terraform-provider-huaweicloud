package ces

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceCesOneClickAlarmReset_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCesOneClickAlarmReset_basic(),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testCesOneClickAlarmReset_basic() string {
	return `
resource "huaweicloud_ces_one_click_alarm_reset" "test" {
  one_click_alarm_id = "OBSSystemOneClickAlarm"
}
`
}

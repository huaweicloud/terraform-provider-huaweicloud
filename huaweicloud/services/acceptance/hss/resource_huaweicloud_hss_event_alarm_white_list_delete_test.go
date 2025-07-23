package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEventAlarmWhiteListDelete_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires the preparation of variable values related to the alarm white list in advance.
			acceptance.TestAccPreCheckHSSEventAlarmWhiteListDeleteEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testEventAlarmWhiteListDelete_basic(),
			},
		},
	})
}

func testEventAlarmWhiteListDelete_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_event_alarm_white_list_delete" "test" {
  enterprise_project_id = "all_granted_eps"

  data_list {
    event_type  = %[1]s
    hash        = "%[2]s"
    description = "%[3]s"
  }
}
`, acceptance.HW_HSS_EVENT_ALARM_WHITE_LIST_EVENT_TYPE, acceptance.HW_HSS_EVENT_ALARM_WHITE_LIST_HASH,
		acceptance.HW_HSS_EVENT_ALARM_WHITE_LIST_DESCRIPTION)
}

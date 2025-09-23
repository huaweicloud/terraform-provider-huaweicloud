package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceAlarmLinkedIncident_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckUserId(t)
			acceptance.TestAccPreCheckCocAlarmID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAlarmLinkedIncident_basic(rName),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testAlarmLinkedIncident_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_coc_alarm_linked_incident" "test" {
  alarm_ids                = "%[2]s"
  current_cloud_service_id = huaweicloud_coc_application.test.id
  description              = "alarm to incident"
  is_service_interrupt     = false
  level_id                 = "level_50"
  mtm_type                 = "inc_type_p_change_issues"
  title                    = "%[3]s"
  enterprise_project_id    = "0"
  assignee                 = "%[4]s"
  is_change_event          = false
  source_id                = "incident_source_alarm"
}
`, testAccApplication_basic(rName), acceptance.HW_COC_ALARM_ID, rName, acceptance.HW_USER_ID)
}

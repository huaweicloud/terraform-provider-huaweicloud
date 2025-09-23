package coc

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceCocTicketAdd_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckUserId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCocTicketAdd_basic(rName),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testCocTicketAdd_basic(name string) string {
	currentTime := time.Now()
	tenMinutesAgo := currentTime.Add(-10*time.Minute).Unix() * 1e3
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_coc_ticket_add" "test" {
  ticket_type              = "issues_mgmt"
  title                    = "%[2]s"
  description              = "this is description"
  enterprise_project_id    = "0"
  issue_ticket_type        = "issues_type_1000"
  virtual_schedule_type    = "issues_mgmt_virtual_schedule_type_2000"
  regions                  = "%[3]s"
  level                    = "issues_level_4000"
  root_cause_cloud_service = huaweicloud_coc_application.test.id
  source                   = "issues_mgmt_associated_type_1000"
  source_id                = huaweicloud_coc_incident.test.id
  found_time               = %[4]v
  issue_contact_person     = "%[5]s"
}
`, testIncident_basic(name), name, acceptance.HW_REGION_NAME, tenMinutesAgo, acceptance.HW_USER_ID)
}

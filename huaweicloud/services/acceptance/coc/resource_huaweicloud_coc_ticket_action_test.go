package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceTicketAction_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckUserId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.12.1",
			},
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testTicketAction_basic(rName),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testTicketAction_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "time_sleep" "wait_10_seconds" {
  depends_on = [huaweicloud_coc_issue.test]

  create_duration = "10s"
}

resource "huaweicloud_coc_ticket_action" "test" {
  depends_on  = [time_sleep.wait_10_seconds]
  ticket_type = "issues_mgmt"
  ticket_id   = huaweicloud_coc_issue.test.id
  action      = "gocm_issues_accepte"
}
`, tesIssue_basic(rName))
}

package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceChangeUpdate_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocTicketID(t)
			acceptance.TestAccPreCheckCocSubTicketID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testChangeUpdate_basic(),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testChangeUpdate_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_coc_change_update" "test" {
  ticket_id = "%[1]s"
  action    = "change_start_change_success"
  sub_tickets {
    ticket_id = "%[2]s"
  }
}
`, acceptance.HW_COC_TICKET_ID, acceptance.HW_COC_SUB_TICKET_ID)
}

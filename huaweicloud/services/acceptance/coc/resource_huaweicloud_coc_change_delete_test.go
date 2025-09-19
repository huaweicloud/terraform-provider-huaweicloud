package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceChangeDelete_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocTicketID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testChangeDelete_basic(),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testChangeDelete_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_coc_change_delete" "test" {
  ticket_type = "change"
  ticket_id   = "%s"
}
`, acceptance.HW_COC_TICKET_ID)
}

package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocChangeSubTickets_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_change_sub_tickets.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocTicketID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocChangeSubTickets_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.is_deleted"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.main_ticket_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.parent_ticket_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.ticket_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.real_ticket_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.ticket_path"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.target_value"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.target_type"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.creator"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.operator"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocChangeSubTickets_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_coc_change_sub_tickets" "test" {
  ticket_type = "change"
  ticket_id   = "%s"
  type        = "child_ticket"
}
`, acceptance.HW_COC_TICKET_ID)
}

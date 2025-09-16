package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocTicketOperationHistories_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_ticket_operation_histories.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocTicketOperationHistories_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.action_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.action"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.sub_action"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.operator"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.ticket_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.is_deleted"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.action_name_zh"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.action_name_en"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.final_sub_action"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.enum_data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.enum_data_list.0.is_deleted"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.enum_data_list.0.match_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.enum_data_list.0.ticket_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.enum_data_list.0.real_ticket_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.enum_data_list.0.name_zh"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.enum_data_list.0.name_en"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.enum_data_list.0.user_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.enum_data_list.0.biz_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.enum_data_list.0.prop_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.enum_data_list.0.model_id"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocTicketOperationHistories_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_coc_ticket_operation_histories" "test" {
  ticket_type = "incident"
  string_filters {
    operator = "="
    field    = "ticket_id"
    values   = [huaweicloud_coc_incident.test.id]
  }
  sort_filter {
    operator = "desc"
    field    = "start_time"
    values   = ["start_time"]
  }
}
`, testIncident_basic(name))
}

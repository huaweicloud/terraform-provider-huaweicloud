package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocIssueTickets_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_issue_tickets.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckUserId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocIssueTickets_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.issue_correlation_sla"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.level"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.root_cause_cloud_service"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.current_cloud_service_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.source"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.source_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.virtual_schedule_type"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.title"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.regions"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.regions_search"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.handle_time"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.found_time"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.is_common_issue"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.is_need_change"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.is_enable_suspension"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.is_start_process_async"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.is_update_null"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.creator"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.operator"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.is_return_full_info"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.is_start_process"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.ticket_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.real_ticket_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.assignee"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.participator"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.work_flow_status"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.baseline_status"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.ticket_type"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.phase"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.enum_data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.enum_data_list.0.is_deleted"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.enum_data_list.0.match_type"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.enum_data_list.0.ticket_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.enum_data_list.0.real_ticket_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.enum_data_list.0.name_zh"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.enum_data_list.0.name_en"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.enum_data_list.0.biz_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.enum_data_list.0.prop_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.enum_data_list.0.model_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.meta_data_version"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.is_deleted"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.ticket_type_id"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocIssueTickets_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_coc_issue_tickets" "test" {
  ticket_type = "issues_mgmt"
  string_filters {
      operator = "="
      field    = "ticket_id"
      values   = [huaweicloud_coc_issue.test.ticket_id]
  }
  sort_filter {
    operator = "desc"
    field    = "fount_time"
    values   = ["fount_time"]
  }
  contain_sub_ticket = false
}
`, tesIssue_basic(name))
}

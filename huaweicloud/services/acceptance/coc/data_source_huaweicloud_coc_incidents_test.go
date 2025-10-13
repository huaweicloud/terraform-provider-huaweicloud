package coc

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocIncidents_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_incidents.test"
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
				Config: testDataSourceDataSourceCocIncidents_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.current_cloud_service_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.level_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.source_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.mtm_type"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.title"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.ticket_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.is_service_interrupt"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.work_flow_status"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.assignee"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.creator"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.operator"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.enum_data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.enum_data_list.0.prop_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.enum_data_list.0.biz_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.enum_data_list.0.name_zh"),
					resource.TestCheckResourceAttrSet(dataSource, "tickets.0.enum_data_list.0.name_en"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocIncidents_basic(name string) string {
	currentTime := time.Now()
	tenMinutesAgo := fmt.Sprintf("%d", currentTime.Add(-10*time.Minute).Unix()*1e3)
	tenMinutesLater := fmt.Sprintf("%d", currentTime.Add(10*time.Minute).Unix()*1e3)
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_coc_incidents" "test" {
  depends_on = [huaweicloud_coc_incident.test]

  contain_sub_ticket = true
  string_filters {
    field    = "create_time"
    operator = ">="
    values   = ["%[2]s"]
    name     = "create_time_start"
  }
  string_filters {
    field    = "create_time"
    operator = "<="
    values   = ["%[3]s"]
    name     = "create_time_end"
  }
  sort_filter {
    field    = "create_time"
    operator = "desc"
    values   = ["create_time"]
    name     = "create_time"
  }
}
`, testIncident_basic(name), tenMinutesAgo, tenMinutesLater)
}

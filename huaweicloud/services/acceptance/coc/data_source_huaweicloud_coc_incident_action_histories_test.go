package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocIncidentActionHistories_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_incident_action_histories.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.12.1",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocIncidentActionHistories_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.action"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.action_name_zh"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.action_name_en"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.operator"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.enum_data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.enum_data_list.0.prop_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.enum_data_list.0.biz_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.enum_data_list.0.name_zh"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.enum_data_list.0.name_en"),
					resource.TestCheckOutput("sort_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocIncidentActionHistories_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "time_sleep" "wait_10_seconds" {
  depends_on = [huaweicloud_coc_incident.test]

  create_duration = "10s"
}

data "huaweicloud_coc_incident_action_histories" "test" {
  depends_on  = [time_sleep.wait_10_seconds]
  incident_id = huaweicloud_coc_incident.test.id

  sort_filter {
    operator = "desc"
    field    = "start_time"
    name     = "start_time"
    values   = ["start_time"]
  }
}

output "sort_filter_is_useful" {
  value = length(data.huaweicloud_coc_incident_action_histories.test.data) > 0
}
`, testIncident_basic(name))
}

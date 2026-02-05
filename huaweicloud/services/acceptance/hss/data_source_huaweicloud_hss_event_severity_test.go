package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEventSeverity_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_event_severity.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEventSeverity_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "low_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "medium_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "high_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "critical_num"),
				),
			},
		},
	})
}

func testDataSourceEventSeverity_basic() string {
	return `
data "huaweicloud_hss_event_severity" "test" {
  category              = "host"
  enterprise_project_id = "all_granted_eps"
  begin_time            = "1768287600000"
  end_time              = "1768298400000"
  event_type            = 1001
  handle_status         = "unhandled"
  severity              = "Security"
  severity_list         = ["Security", "Low"]
  attack_tag            = "attack_success"
  asset_value           = "common"
  tag_list              = ["热点事件"]
  att_ck                = "Reconnaissance"
}
`
}

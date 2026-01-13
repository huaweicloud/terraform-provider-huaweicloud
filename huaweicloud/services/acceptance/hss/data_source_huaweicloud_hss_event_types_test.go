package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEventTypes_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_event_types.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEventTypes_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.event_type_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.event_type_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.event_type_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.event_type_list.0.event_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.event_type_list.0.event_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.event_type_list.0.status"),
				),
			},
		},
	})
}

func testDataSourceEventTypes_basic() string {
	return `
data "huaweicloud_hss_event_types" "test" {
  category              = "host"
  enterprise_project_id = "all_granted_eps"
  begin_time            = "1768287600000"
  end_time              = "1768298400000"
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

package antiddos

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceConfigRanges_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_antiddos_config_ranges.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceConfigRanges_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "connection_limited_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "http_limited_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "traffic_limited_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "connection_limited_list.0.cleaning_access_pos_id"),
					resource.TestCheckResourceAttrSet(dataSource, "connection_limited_list.0.new_connection_limited"),
					resource.TestCheckResourceAttrSet(dataSource, "connection_limited_list.0.total_connection_limited"),
					resource.TestCheckResourceAttrSet(dataSource, "http_limited_list.0.http_packet_per_second"),
					resource.TestCheckResourceAttrSet(dataSource, "http_limited_list.0.http_request_pos_id"),
					resource.TestCheckResourceAttrSet(dataSource, "traffic_limited_list.0.packet_per_second"),
					resource.TestCheckResourceAttrSet(dataSource, "traffic_limited_list.0.traffic_per_second"),
					resource.TestCheckResourceAttrSet(dataSource, "traffic_limited_list.0.traffic_pos_id"),
				),
			},
		},
	})
}

const testDataSourceConfigRanges_basic = `
data "huaweicloud_antiddos_config_ranges" "test" {
}
`

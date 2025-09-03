package lts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceStreamCharts_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_lts_stream_charts.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLTSLogConvergeMappingConfig(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceStreamCharts_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "region"),
					resource.TestMatchResourceAttr(dataSourceName, "charts.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "charts.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "charts.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "charts.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "charts.0.log_group_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "charts.0.log_group_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "charts.0.log_stream_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "charts.0.log_stream_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "charts.0.sql"),
					resource.TestCheckResourceAttrSet(dataSourceName, "charts.0.config.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "charts.0.config.0.page_size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "charts.0.config.0.can_sort"),
					resource.TestCheckResourceAttrSet(dataSourceName, "charts.0.config.0.can_search"),

					resource.TestCheckOutput("is_charts_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceStreamCharts_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_lts_stream_charts" "test" {
  log_group_id  = "%[1]s"
  log_stream_id = "%[2]s"
}

locals {
  charts_count = length(data.huaweicloud_lts_stream_charts.test.charts)
}

output "is_charts_filter_useful" {
  value = local.charts_count >= 0
}
`, acceptance.HW_LTS_LOG_CONVERGE_SOURCE_LOG_GROUP_ID, acceptance.HW_LTS_LOG_CONVERGE_SOURCE_LOG_STREAM_ID)
}

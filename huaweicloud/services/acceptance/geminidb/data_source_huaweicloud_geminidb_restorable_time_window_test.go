package geminidb

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGeminiDBRestorableTimeWindow_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_geminidb_restorable_time_window.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGeminiDBRestorableTimeWindow_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "restorable_time_periods.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "restorable_time_periods.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "restorable_time_periods.0.end_time"),

					resource.TestCheckOutput("filter_by_time_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceGeminiDBRestorableTimeWindow_basic() string {
	beginTime := time.Now().UTC()
	beginTimeString := beginTime.Format("2006-01-02T15:04:05+0800")
	endTime := time.Now().UTC().Add(8 * time.Hour)
	endTimeString := endTime.Format("2006-01-02T15:04:05+0800")
	return fmt.Sprintf(`
data "huaweicloud_geminidb_restorable_time_window" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_geminidb_restorable_time_window" "time_filter_is_useful" {
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
}

output "filter_by_time_is_useful" {
  value = length(data.huaweicloud_geminidb_restorable_time_window.time_filter_is_useful.restorable_time_periods) > 0
}
`, acceptance.HW_GEMINIDB_INSATNCE_ID, beginTimeString, endTimeString)
}

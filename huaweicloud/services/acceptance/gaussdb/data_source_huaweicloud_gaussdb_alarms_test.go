package gaussdb

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAlarms_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_gaussdb_alarms.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAlarms_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "history_records.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttrSet(all, "history_records.0.alarm_id"),
					resource.TestCheckResourceAttrSet(all, "history_records.0.name"),
					resource.TestCheckResourceAttrSet(all, "history_records.0.status"),
					resource.TestCheckResourceAttrSet(all, "history_records.0.alarm_type"),
					resource.TestCheckResourceAttrSet(all, "history_records.0.level"),
					resource.TestCheckResourceAttrSet(all, "history_records.0.instance_id"),
					resource.TestCheckResourceAttrSet(all, "history_records.0.instance_name"),
					resource.TestCheckResourceAttrSet(all, "history_records.0.begin_time"),
					resource.TestCheckResourceAttrSet(all, "history_records.0.update_time"),
					resource.TestCheckOutput("level_filter_is_useful", "true"),
				),
			},
		},
	},
	)
}

func testDataSourceAlarms_basic() string {
	startTime := time.Now().UTC().Add(-30 * 24 * time.Hour).Format("2006-01-02T15:04:05+0000")

	return fmt.Sprintf(`
data "huaweicloud_gaussdb_alarms" "all" {
  start_time = "%[1]s"
}

# Filter by 'level' parameter
data "huaweicloud_gaussdb_alarms" "level_filter" {
  start_time = "%[1]s"
  level      = 1
}

output "level_filter_is_useful" {
 value = length(data.huaweicloud_gaussdb_alarms.level_filter.history_records) >= 0
}
`, startTime,
	)
}

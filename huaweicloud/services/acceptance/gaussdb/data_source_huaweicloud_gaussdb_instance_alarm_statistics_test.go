package gaussdb

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceInstanceAlarmStatistics_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_gaussdb_instance_alarm_statistics.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(
		t, resource.TestCase{
			PreCheck:          func() { acceptance.TestAccPreCheck(t) },
			ProviderFactories: acceptance.TestAccProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testDataSourceInstanceAlarmStatistics_basic(),
					Check: resource.ComposeTestCheckFunc(
						dc.CheckResourceExists(),
						resource.TestCheckResourceAttrSet(all, "ring_percentage"),
						resource.TestMatchResourceAttr(all, "instance_alarm_level_statistics.#", regexp.MustCompile(`^[0-9]+$`)),
						resource.TestCheckResourceAttrSet(all, "instance_alarm_level_statistics.0.instance_id"),
						resource.TestCheckResourceAttrSet(all, "instance_alarm_level_statistics.0.instance_name"),
						resource.TestCheckResourceAttrSet(all, "instance_alarm_level_statistics.0.total_count"),
						resource.TestMatchResourceAttr(all, "instance_alarm_level_statistics.0.alarm_level_statistics.#",
							regexp.MustCompile(`^[0-9]+$`)),
						resource.TestCheckResourceAttrSet(all, "instance_alarm_level_statistics.0.alarm_level_statistics.0.count"),
						resource.TestCheckResourceAttrSet(all, "instance_alarm_level_statistics.0.alarm_level_statistics.0.level_name"),
						resource.TestMatchResourceAttr(all, "total_alarm_level_statistics.#", regexp.MustCompile(`^[0-9]+$`)),
						resource.TestCheckResourceAttrSet(all, "total_alarm_level_statistics.0.count"),
						resource.TestCheckResourceAttrSet(all, "total_alarm_level_statistics.0.level_name"),
					),
				},
			},
		},
	)
}

func testDataSourceInstanceAlarmStatistics_basic() string {
	startTime := time.Now().UTC().Add(-30 * 24 * time.Hour).Format("2006-01-02T15:04:05+0000")

	return fmt.Sprintf(`
data "huaweicloud_gaussdb_instance_alarm_statistics" "test" {
  start_time = "%s"
  top_num    = 10
}
`, startTime)
}

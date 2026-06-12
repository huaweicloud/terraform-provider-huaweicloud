package dataarts

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// This test requires alarm records to be created manually before running.
// Please ensure that alarm records exist in the specified time range.
func TestAccDataFactoryAlarmRecords_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_dataarts_factory_alarm_records.test"
		now            = time.Now()
		startTime      = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Format(time.RFC3339)
		endTime        = time.Date(now.Year(), now.Month(), now.Day(), 24, 0, 0, 0, now.Location()).Format(time.RFC3339)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataFactoryAlarmRecords_basic_step1(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceName, "records.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "region"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.alarm_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.job_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.schedule_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.send_msg"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.plan_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.remind_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.send_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "records.0.job_id"),
				),
			},
			{
				Config: testAccDataFactoryAlarmRecords_basic_step2(startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceName, "records.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					checkAlarmTimeInRange(dataSourceName, startTime, endTime),
				),
			},
		},
	})
}

func testAccDataFactoryAlarmRecords_basic_step1() string {
	return fmt.Sprintf(`
data "huaweicloud_dataarts_factory_alarm_records" "test" {
  workspace = "%[1]s"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID)
}

func testAccDataFactoryAlarmRecords_basic_step2(startTime, endTime string) string {
	return fmt.Sprintf(`
data "huaweicloud_dataarts_factory_alarm_records" "test" {
  workspace  = "%[1]s"
  start_time = "%[2]s"
  end_time   = "%[3]s"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, startTime, endTime)
}

func checkAlarmTimeInRange(dataSourceName, startTime, endTime string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("can not find the data source: %s", dataSourceName)
		}

		start, err := time.Parse(time.RFC3339, startTime)
		if err != nil {
			return fmt.Errorf("error parsing start_time %s: %s", startTime, err)
		}

		end, err := time.Parse(time.RFC3339, endTime)
		if err != nil {
			return fmt.Errorf("error parsing end_time %s: %s", endTime, err)
		}

		records := rs.Primary.Attributes["records.#"]
		if records == "0" {
			return nil
		}

		// Check the first record's alarm_time
		alarmTimeStr := rs.Primary.Attributes["records.0.alarm_time"]
		alarmTime, err := time.Parse(time.RFC3339, alarmTimeStr)
		if err != nil {
			return fmt.Errorf("error parsing alarm_time %s: %s", alarmTimeStr, err)
		}

		if alarmTime.Before(start) || alarmTime.After(end) {
			return fmt.Errorf("alarm_time %s is not in the range [%s, %s]", alarmTimeStr, startTime, endTime)
		}

		return nil
	}
}

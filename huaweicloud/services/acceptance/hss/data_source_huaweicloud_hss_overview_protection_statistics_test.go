package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOverviewProtectionStatistics_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_overview_protection_statistics.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceOverviewProtectionStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "vul_library_update_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "protect_days"),
					resource.TestCheckResourceAttrSet(dataSourceName, "threat_library_update_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "vul_detected_total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "baseline_detected_total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "finger_scan_total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "alarm_detected_total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "ransomware_alarm_detected_total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "file_alarm_detected_total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rasp_alarm_detected_total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "wtp_alarm_detected_total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "image_risk_total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "container_alarm_total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "container_firewall_policy_total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "auto_kill_virus_status"),
				),
			},
		},
	})
}

func testDataSourceOverviewProtectionStatistics_basic() string {
	return `
data "huaweicloud_hss_overview_protection_statistics" "test" {
  enterprise_project_id = "0"
}
`
}

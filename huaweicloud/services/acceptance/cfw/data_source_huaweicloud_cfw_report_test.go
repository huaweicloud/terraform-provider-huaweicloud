package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceReport_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cfw_report.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
			acceptance.TestAccPreCheckCfwReportProfile(t)
			acceptance.TestAccPreCheckCfwReport(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceReport_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "attack_info.0.dst_ip.#"),
					resource.TestCheckResourceAttrSet(dataSource, "attack_info.0.ips_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "attack_info.0.level.#"),
					resource.TestCheckResourceAttrSet(dataSource, "attack_info.0.rule.#"),
					resource.TestCheckResourceAttrSet(dataSource, "attack_info.0.src_ip.#"),
					resource.TestCheckResourceAttrSet(dataSource, "attack_info.0.trend.#"),
					resource.TestCheckResourceAttrSet(dataSource, "attack_info.0.type.#"),
					resource.TestCheckResourceAttrSet(dataSource, "category"),
					resource.TestCheckResourceAttrSet(dataSource, "internet_firewall.0.eip.#"),
					resource.TestCheckResourceAttrSet(dataSource, "internet_firewall.0.in2out.#"),
					resource.TestCheckResourceAttrSet(dataSource, "internet_firewall.0.out2in.#"),
					resource.TestCheckResourceAttrSet(dataSource, "internet_firewall.0.overview.#"),
					resource.TestCheckResourceAttrSet(dataSource, "internet_firewall.0.traffic_trend.#"),
					resource.TestCheckResourceAttrSet(dataSource, "send_time"),
					resource.TestCheckResourceAttrSet(dataSource, "statistic_period.#"),
				),
			},
		},
	})
}

func testDataSourceReport_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_report" "test" {
  fw_instance_id    = "%[1]s"
  report_profile_id = "%[2]s"
  report_id         = "%[3]s"
}
`, acceptance.HW_CFW_INSTANCE_ID, acceptance.HW_CFW_REPORT_PROFILE_ID, acceptance.HW_CFW_REPORT_ID)
}

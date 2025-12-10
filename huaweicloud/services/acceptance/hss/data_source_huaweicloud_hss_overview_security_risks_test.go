package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOverviewSecurityRisks_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_overview_security_risks.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOverviewSecurityRisks_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_risk.#"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_risk.0.risk_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_risk.0.risk_list.0.severity"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_risk.0.risk_list.0.risk_num"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_risk.0.deduct_score"),
					resource.TestCheckResourceAttrSet(dataSource, "baseline_risk.#"),
					resource.TestCheckResourceAttrSet(dataSource, "baseline_risk.0.existed_pwd_host_num"),
					resource.TestCheckResourceAttrSet(dataSource, "asset_risk.#"),
					resource.TestCheckResourceAttrSet(dataSource, "asset_risk.0.existed_danger_port_host_num"),
					resource.TestCheckResourceAttrSet(dataSource, "security_protect_risk.#"),
					resource.TestCheckResourceAttrSet(dataSource, "security_protect_risk.0.un_open_protection_host_num"),
					resource.TestCheckResourceAttrSet(dataSource, "vul_risk.#"),
					resource.TestCheckResourceAttrSet(dataSource, "vul_risk.0.un_scanned_host_num"),
					resource.TestCheckResourceAttrSet(dataSource, "image_risk.#"),
					resource.TestCheckResourceAttrSet(dataSource, "image_risk.0.risk_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "image_risk.0.risk_list.0.severity"),
					resource.TestCheckResourceAttrSet(dataSource, "image_risk.0.risk_list.0.image_num"),
				),
			},
		},
	})
}

const testAccDataSourceOverviewSecurityRisks_basic = `
data "huaweicloud_hss_overview_security_risks" "test" {
  enterprise_project_id = "0"
}
`

package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAntivirusStatistic_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_antivirus_statistic.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID with container version host protection enabled.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAntivirusStatistic_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "total_malware_num"),
					resource.TestCheckResourceAttrSet(dataSource, "malware_host_num"),
					resource.TestCheckResourceAttrSet(dataSource, "total_task_num"),
					resource.TestCheckResourceAttrSet(dataSource, "scanning_task_num"),
					resource.TestCheckResourceAttrSet(dataSource, "latest_scan_time"),
					resource.TestCheckResourceAttrSet(dataSource, "scan_type"),
				),
			},
		},
	})
}

const testAccDataSourceAntivirusStatistic_basic string = `
data "huaweicloud_hss_antivirus_statistic" "test" {
  enterprise_project_id = "all_granted_eps"
}
`

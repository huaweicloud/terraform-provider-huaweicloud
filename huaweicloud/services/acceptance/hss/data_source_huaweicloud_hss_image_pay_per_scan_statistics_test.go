package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceImagePayPerScanStatistics_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_image_pay_per_scan_statistics.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceImagePayPerScanStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "repository_scan_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cicd_scan_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "free_quota_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "already_given"),
					resource.TestCheckResourceAttrSet(dataSourceName, "image_free_quota"),
					resource.TestCheckResourceAttrSet(dataSourceName, "high_risk.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "high_risk.0.local"),
					resource.TestCheckResourceAttrSet(dataSourceName, "high_risk.0.registriy"),
					resource.TestCheckResourceAttrSet(dataSourceName, "high_risk.0.cicd"),
					resource.TestCheckResourceAttrSet(dataSourceName, "has_risk.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "has_risk.0.local"),
					resource.TestCheckResourceAttrSet(dataSourceName, "has_risk.0.registriy"),
					resource.TestCheckResourceAttrSet(dataSourceName, "has_risk.0.cicd"),
					resource.TestCheckResourceAttrSet(dataSourceName, "total.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "total.0.local"),
					resource.TestCheckResourceAttrSet(dataSourceName, "total.0.registriy"),
					resource.TestCheckResourceAttrSet(dataSourceName, "total.0.cicd"),
					resource.TestCheckResourceAttrSet(dataSourceName, "unscan.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "unscan.0.local"),
					resource.TestCheckResourceAttrSet(dataSourceName, "unscan.0.registriy"),
					resource.TestCheckResourceAttrSet(dataSourceName, "unscan.0.cicd"),
				),
			},
		},
	})
}

func testAccDataSourceImagePayPerScanStatistics_basic() string {
	return `
data "huaweicloud_hss_image_pay_per_scan_statistics" "test" {
  enterprise_project_id = "0"
}
`
}

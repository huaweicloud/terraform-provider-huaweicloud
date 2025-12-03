package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOperationalReportWelfare_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_operational_report_welfare.test"
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
				Config: testAccDataSourceOperationalReportWelfare_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "hot_info.#"),
					resource.TestCheckResourceAttrSet(dataSource, "hot_info.0.title"),
					resource.TestCheckResourceAttrSet(dataSource, "hot_info.0.url_json"),
					resource.TestCheckResourceAttrSet(dataSource, "version_update_info.#"),
					resource.TestCheckResourceAttrSet(dataSource, "version_update_info.0.title"),
					resource.TestCheckResourceAttrSet(dataSource, "version_update_info.0.url_json"),
					resource.TestCheckResourceAttrSet(dataSource, "activities_info.#"),
					resource.TestCheckResourceAttrSet(dataSource, "activities_info.0.title"),
					resource.TestCheckResourceAttrSet(dataSource, "activities_info.0.url_json"),
				),
			},
		},
	})
}

const testAccDataSourceOperationalReportWelfare_basic = `data "huaweicloud_hss_operational_report_welfare" "test" {}`

package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWebtamperProtectionStatistics_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_webtamper_protection_statistics.test"
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
				Config: testAccDataSourceWebtamperProtectionStatistics_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "protect_host_num"),
					resource.TestCheckResourceAttrSet(dataSource, "protect_success_host_num"),
				),
			},
		},
	})
}

const testAccDataSourceWebtamperProtectionStatistics_basic = `
data "huaweicloud_hss_webtamper_protection_statistics" "test" {
  enterprise_project_id = "0"
}
`

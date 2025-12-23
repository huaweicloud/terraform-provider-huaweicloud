package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBaselineStatistic_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_baseline_statistic.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBaselineStatistic_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "host_weak_pwd"),
					resource.TestCheckResourceAttrSet(dataSource, "pwd_policy"),
					resource.TestCheckResourceAttrSet(dataSource, "security_check"),
				),
			},
		},
	})
}

const testAccDataSourceBaselineStatistic_basic = `
data "huaweicloud_hss_baseline_statistic" "test" {
  enterprise_project_id = "0"
}
`

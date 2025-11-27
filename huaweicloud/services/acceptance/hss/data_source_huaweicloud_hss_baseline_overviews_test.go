package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBaselineOverviews_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_baseline_overviews.test"
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
				Config: testAccDataSourceBaselineOverviews_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "host_num"),
					resource.TestCheckResourceAttrSet(dataSource, "check_type_num"),
					resource.TestCheckResourceAttrSet(dataSource, "weak_pwd_total_host"),
					resource.TestCheckResourceAttrSet(dataSource, "weak_pwd_risk"),
				),
			},
		},
	})
}

const testAccDataSourceBaselineOverviews_basic = `
data "huaweicloud_hss_baseline_overviews" "test" {
  enterprise_project_id = "0"
}
`

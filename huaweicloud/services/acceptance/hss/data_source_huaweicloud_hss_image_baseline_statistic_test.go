package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceImageBaselineStatistic_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_image_baseline_statistic.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceImageBaselineStatistic_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "image_weak_pwd"),
					resource.TestCheckResourceAttrSet(dataSourceName, "pwd_policy"),
					resource.TestCheckResourceAttrSet(dataSourceName, "security_check"),
				),
			},
		},
	})
}

func testDataSourceImageBaselineStatistic_basic() string {
	return `
data "huaweicloud_hss_image_baseline_statistic" "test" {
  image_type            = "private_image"
  enterprise_project_id = "all_granted_eps"
}
`
}

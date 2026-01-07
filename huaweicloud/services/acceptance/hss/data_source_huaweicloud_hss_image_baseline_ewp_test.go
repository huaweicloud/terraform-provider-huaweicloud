package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceImageBaselineEWP_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_image_baseline_ewp.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceImageBaselineEWP_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "extended_weak_password.#"),
				),
			},
		},
	})
}

func testDataSourceImageBaselineEWP_basic() string {
	return `
data "huaweicloud_hss_image_baseline_ewp" "test" {
  enterprise_project_id = "all_granted_eps"
}
`
}

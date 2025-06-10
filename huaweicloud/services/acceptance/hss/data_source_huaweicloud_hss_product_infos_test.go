package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceProductInfos_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_product_infos.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProductInfos_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.charging_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.is_auto_renew"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.version_info.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.version_info.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.version_info.0.periods.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.version_info.0.periods.0.period_vals"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.version_info.0.periods.0.period_unit"),
				),
			},
		},
	})
}

func testAccDataSourceProductInfos_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_hss_product_infos" "test" {
  site_code             = "HWC_CN"
  enterprise_project_id = "%s"
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

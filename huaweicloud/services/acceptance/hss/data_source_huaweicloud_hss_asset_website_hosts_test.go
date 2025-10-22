package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssetWebsiteHosts_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_asset_website_hosts.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			// Because there is no available data for testing, the test case is only
			// used to verify that the API can be invoked.
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSDomain(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAssetWebsiteHosts_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
				),
			},
		},
	})
}

func testAccDataSourceAssetWebsiteHosts_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_hss_asset_website_hosts" "test" {
  category              = "0"
  domain                = "%s"
  enterprise_project_id = "0"
}
`, acceptance.HW_HSS_DOMAIN)
}

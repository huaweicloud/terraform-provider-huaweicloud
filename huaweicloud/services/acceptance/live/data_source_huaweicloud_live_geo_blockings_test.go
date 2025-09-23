package live

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGeoBlockings_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_live_geo_blockings.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLiveStreamingDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGeoBlockings_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "apps.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "apps.0.app_name", "live"),
					resource.TestCheckResourceAttr(dataSourceName, "apps.0.area_whitelist.#", "5"),
				),
			},
		},
	})
}

func testDataSourceGeoBlockings_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_live_geo_blockings" "test" {
  domain_name = "%[2]s"

  depends_on = [huaweicloud_live_geo_blocking.test]
}
`, testResourceGeoBlocking_basic(), acceptance.HW_LIVE_STREAMING_DOMAIN_NAME)
}

package waf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGeolocationDetail_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_waf_geolocation_detail.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGeolocationDetail_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "continent"),
					resource.TestCheckResourceAttrSet(dataSourceName, "geomap"),
					resource.TestCheckResourceAttrSet(dataSourceName, "locale.CN"),
				),
			},
		},
	})
}

const testAccDataSourceGeolocationDetail_basic = `
data "huaweicloud_waf_geolocation_detail" "test" {
  lang = "en"
}
`

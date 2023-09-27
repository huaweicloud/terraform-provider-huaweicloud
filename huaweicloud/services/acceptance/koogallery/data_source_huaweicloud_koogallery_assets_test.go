package koogallery

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKooGalleryAssetsDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_koogallery_assets.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckKooGallery(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKooGalleryAssetsDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "assets.0.asset_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "assets.0.software_pkg_deployed_object.0.package_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "assets.0.software_pkg_deployed_object.0.internal_path"),
					resource.TestCheckResourceAttrSet(dataSourceName, "assets.0.software_pkg_deployed_object.0.checksum"),
				),
			},
		},
	})
}

var testAccKooGalleryAssetsDataSource_basic = `
data "huaweicloud_koogallery_assets" "test" {
  asset_id      = "5848036e33fb4d4fa408c63fc3a9c8ab"
  deployed_type = "software_package"
  asset_version = "V1.0"
}
`

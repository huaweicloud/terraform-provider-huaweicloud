package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscAssetQuota_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dsc_asset_quota.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscAssetQuota_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "use_size"),
				),
			},
		},
	})
}

const testAccDataSourceDscAssetQuota_basic = `
data "huaweicloud_dsc_asset_quota" "test" {}
`

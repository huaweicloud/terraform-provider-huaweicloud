package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssetPortInfo_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_asset_port_info.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAssetPortInfo_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "type"),
					resource.TestCheckResourceAttrSet(dataSource, "status"),
					resource.TestCheckResourceAttrSet(dataSource, "description"),
					resource.TestCheckResourceAttrSet(dataSource, "description_en"),
				),
			},
		},
	})
}

const testAccDataSourceAssetPortInfo_basic = `
data "huaweicloud_hss_asset_port_info" "test" {
  port                  = 8080
  category              = "0"
  enterprise_project_id = "all_granted_eps"
}
`

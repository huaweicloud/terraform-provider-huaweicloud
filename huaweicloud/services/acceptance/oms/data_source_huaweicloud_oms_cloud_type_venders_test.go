package oms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCloudTypeVenders_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_oms_cloud_type_venders.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCloudTypeVenders_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "venders.#"),
				),
			},
		},
	})
}

const testAccDataSourceCloudTypeVenders_basic = `
data "huaweicloud_oms_cloud_type_venders" "test" {
  type = "src"
}
`

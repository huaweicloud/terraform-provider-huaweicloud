package oms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceObjectstorageCloudType_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_oms_objectstorage_cloud_type.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceObjectstorageCloudType_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "vendors.#"),
				),
			},
		},
	})
}

const testAccDataSourceObjectstorageCloudType_basic = `
data "huaweicloud_oms_objectstorage_cloud_type" "test" {
  type = "src"
}
`

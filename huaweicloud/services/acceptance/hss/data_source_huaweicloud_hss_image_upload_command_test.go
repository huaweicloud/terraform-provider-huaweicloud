package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHssImageUploadCommand_basic(t *testing.T) {
	dataSource := "data.huaweicloud_hss_image_upload_command.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceHssImageUploadCommand_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "image_command"),
					resource.TestCheckResourceAttrSet(dataSource, "images_download_url"),
				),
			},
		},
	})
}

const testDataSourceDataSourceHssImageUploadCommand_basic = `
data "huaweicloud_hss_image_upload_command" "test" {
  registry_addr = "hub.docker.com"
  namespace     = "test"
  username      = "xxx"
  password      = "xxxxxx"
}
`

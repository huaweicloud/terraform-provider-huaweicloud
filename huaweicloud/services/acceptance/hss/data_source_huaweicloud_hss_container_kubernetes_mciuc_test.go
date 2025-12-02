package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHssContainerKubernetesMciuc_basic(t *testing.T) {
	dataSource := "data.huaweicloud_hss_container_kubernetes_mciuc.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceHssContainerKubernetesMciuc_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "image_command"),
					resource.TestCheckResourceAttrSet(dataSource, "secret_command"),
					resource.TestCheckResourceAttrSet(dataSource, "images_download_url"),
				),
			},
		},
	})
}

// The test case uses fake test data.
const testDataSourceDataSourceHssContainerKubernetesMciuc_basic = `
data "huaweicloud_hss_container_kubernetes_mciuc" "test" {
  image_repo   = "hub.docker.com"
  organization = "ywk-test"
  username     = "ywkkkkkkk"
  password     = "agent@2023"
  plug_type    = "docker"
}
`

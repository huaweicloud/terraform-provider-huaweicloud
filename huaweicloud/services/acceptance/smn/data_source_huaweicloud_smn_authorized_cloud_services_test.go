package smn

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSmnAuthorizedCloudServices_basic(t *testing.T) {
	dataSource := "data.huaweicloud_smn_authorized_cloud_services.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSmnAuthorizedCloudServices_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "cloud_services.#"),
					resource.TestCheckResourceAttrSet(dataSource, "cloud_services.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "cloud_services.0.show_name"),
				),
			},
		},
	})
}

func testDataSourceSmnAuthorizedCloudServices_basic() string {
	return `
data "huaweicloud_smn_authorized_cloud_services" "test" {}
`
}

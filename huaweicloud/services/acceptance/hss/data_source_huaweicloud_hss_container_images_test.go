package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceContainerImages_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_container_images.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceContainerImages_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.image_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.image_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.image_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.create_time"),
				),
			},
		},
	})
}

func testDataSourceContainerImages_basic() string {
	return `
data "huaweicloud_hss_container_images" "test" {
  enterprise_project_id = "0"
}
`
}

package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceImageSensitiveInformation_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_image_sensitive_information.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceImageSensitiveInformation_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.sensitive_info_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.severity"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.position"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.file_path"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.content"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.latest_scan_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.handle_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.operate_accept"),
				),
			},
		},
	})
}

func testDataSourceImageSensitiveInformation_basic() string {
	return `
data "huaweicloud_hss_image_sensitive_information" "test" {
  enterprise_project_id = "0"
  image_type            = "registry"
}
`
}

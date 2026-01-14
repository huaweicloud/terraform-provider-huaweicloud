package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceImageNonCompliantApp_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_image_non_compliant_app.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSImageId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceImageNonCompliantApp_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.app_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.app_path"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.app_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.layer_digest"),
				),
			},
		},
	})
}

func testDataSourceImageNonCompliantApp_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_hss_image_non_compliant_app" "test" {
  image_id              = "%s"
  image_type            = "private_image"
  enterprise_project_id = "all_granted_eps"
}
`, acceptance.HW_HSS_IMAGE_ID)
}

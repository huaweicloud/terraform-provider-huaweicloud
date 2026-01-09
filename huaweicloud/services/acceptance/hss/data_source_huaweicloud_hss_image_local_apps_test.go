package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceImageLocalApps_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_image_local_apps.test"
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
				Config: testDataSourceImageLocalApps_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
				),
			},
		},
	})
}

func testDataSourceImageLocalApps_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_hss_image_local_apps" "test" {
  image_id              = "%s"
  enterprise_project_id = "all_granted_eps"
}
`, acceptance.HW_HSS_IMAGE_ID)
}

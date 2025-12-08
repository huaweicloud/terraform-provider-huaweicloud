package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceImageFilesStatistics_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_image_files_statistics.test"
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
				Config: testDataSourceImageFilesStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_files_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_files_size"),
				),
			},
		},
	})
}

func testDataSourceImageFilesStatistics_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_hss_image_files_statistics" "test" {
  image_id   = "%s"
  image_type = "private_image"
}
`, acceptance.HW_HSS_IMAGE_ID)
}

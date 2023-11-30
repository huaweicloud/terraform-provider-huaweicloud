package iec

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccImagesDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_iec_images.images_test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesConfig(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "images.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckResourceAttr(dataSourceName, "region", acceptance.HW_REGION_NAME),
				),
			},
		},
	})
}

func testAccImagesConfig() string {
	return fmt.Sprintf(`
data "huaweicloud_iec_images" "images_test" {
  region = "%s"
}
	`, acceptance.HW_REGION_NAME)
}

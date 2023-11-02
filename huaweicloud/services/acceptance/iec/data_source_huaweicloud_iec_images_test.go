package iec

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIECImagesDataSource_basic(t *testing.T) {
	resourceName := "data.huaweicloud_iec_images.images_test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIECImagesConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIECImagesDataSourceID(resourceName),
					resource.TestMatchResourceAttr(resourceName, "images.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckResourceAttr(resourceName, "region", acceptance.HW_REGION_NAME),
				),
			},
		},
	})
}

func testAccCheckIECImagesDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("IEC images data source ID not set")
		}
		return nil
	}
}

func testAccIECImagesConfig() string {
	return fmt.Sprintf(`
data "huaweicloud_iec_images" "images_test" {
  region = "%s"
}
	`, acceptance.HW_REGION_NAME)
}

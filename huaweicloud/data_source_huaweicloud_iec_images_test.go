package huaweicloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIECImagesDataSource_basic(t *testing.T) {
	resourceName := "data.huaweicloud_iec_images.images_test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIECImagesConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIECImagesDataSourceID(resourceName),
					resource.TestMatchResourceAttr(resourceName, "images.#", regexp.MustCompile("[1-9]\\d*")),
					resource.TestCheckResourceAttr(resourceName, "region", HW_REGION_NAME),
				),
			},
		},
	})
}

func testAccCheckIECImagesDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Root module has no resource called %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("IEC images data source ID not set")
		}
		return nil
	}
}

func testAccIECImagesConfig() string {
	return fmt.Sprintf(`
data "huaweicloud_iec_images" "images_test" {
  region = "%s"
}
	`, HW_REGION_NAME)
}

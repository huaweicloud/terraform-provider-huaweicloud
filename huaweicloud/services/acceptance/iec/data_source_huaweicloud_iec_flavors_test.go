package huaweicloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIECFlavorsDataSource_basic(t *testing.T) {
	resourceName := "data.huaweicloud_iec_flavors.flavors_test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIECFlavorsConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIECFlavorsDataSourceID(resourceName),
					resource.TestMatchResourceAttr(resourceName, "flavors.#", regexp.MustCompile("[1-9]\\d*")),
					resource.TestCheckResourceAttr(resourceName, "region", HW_REGION_NAME),
				),
			},
		},
	})
}

func TestAccIECFlavorsDataSource_FilterName(t *testing.T) {
	resourceName := "data.huaweicloud_iec_flavors.flavors_test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIECFlavorsWithName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIECFlavorsDataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "flavors.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "region", HW_REGION_NAME),
				),
			},
		},
	})
}

func testAccCheckIECFlavorsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Root module has no resource called %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("IEC flavors data source ID not set")
		}
		return nil
	}
}

func testAccIECFlavorsConfig() string {
	return fmt.Sprintf(`
data "huaweicloud_iec_flavors" "flavors_test" {
  region = "%s"
}
	`, HW_REGION_NAME)
}

func testAccIECFlavorsWithName() string {
	return fmt.Sprintf(`
data "huaweicloud_iec_flavors" "flavors_test" {
  region = "%s"
  name   = "c6.large.2"
}
	`, HW_REGION_NAME)
}

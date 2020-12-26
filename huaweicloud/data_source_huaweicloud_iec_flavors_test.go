package huaweicloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccIECFlavorsDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("terraform_test_iec_flavors-%s", acctest.RandString(5))
	resourceName := "data.huaweicloud_iec_flavors.flavors_test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIECFlavorsConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIECFlavorsDataSourceID(resourceName, rName),
					resource.TestMatchResourceAttr(resourceName, "flavors.#", regexp.MustCompile("[1-9]\\d*")),
					resource.TestCheckResourceAttr(resourceName, "region", "cn-north-4"),
				),
			},
		},
	})
}

func testAccCheckIECFlavorsDataSourceID(n, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Root module has no resource called %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("IEC flavors data source ID not set")
		}
		return nil
	}
}

func testAccIECFlavorsConfig(val string) string {
	return fmt.Sprintf(`
data "huaweicloud_iec_flavors" "flavors_test" {
  region = "cn-north-4"
  name   = "%s"
}
	`, val)
}

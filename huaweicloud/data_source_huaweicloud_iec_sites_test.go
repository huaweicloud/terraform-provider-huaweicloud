package huaweicloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccIECSitesDataSource_basic(t *testing.T) {
	resourceName := "data.huaweicloud_iec_sites.sites_test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIECSitesConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIECSitesDataSourceID(resourceName),
					resource.TestMatchResourceAttr(resourceName, "sites.#", regexp.MustCompile("[1-9]\\d*")),
					resource.TestCheckResourceAttr(resourceName, "region", HW_REGION_NAME),
					resource.TestCheckResourceAttr(resourceName, "area", "east"),
					resource.TestCheckResourceAttr(resourceName, "city", "hangzhou"),
				),
			},
		},
	})
}

func testAccCheckIECSitesDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Root module has no resource called %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("IEC sites data source ID not set")
		}
		return nil
	}
}

func testAccIECSitesConfig() string {
	return fmt.Sprintf(`
data "huaweicloud_iec_sites" "sites_test" {
  region = "%s"
  area   = "east"
  city   = "hangzhou"
}
	`, HW_REGION_NAME)
}

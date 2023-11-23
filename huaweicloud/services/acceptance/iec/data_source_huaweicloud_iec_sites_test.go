package iec

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIECSitesDataSource_basic(t *testing.T) {
	resourceName := "data.huaweicloud_iec_sites.sites_test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIecServerDestory,
		Steps: []resource.TestStep{
			{
				Config: testAccIECSitesConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIECSitesDataSourceID(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "sites.#"),
					resource.TestCheckResourceAttr(resourceName, "sites.0.area", "east"),
					resource.TestCheckResourceAttrSet(resourceName, "sites.0.lines.#"),
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

func testAccIECSitesConfig_basic() string {
	return `
data "huaweicloud_iec_sites" "sites_test" {
  area = "east"
}
`
}

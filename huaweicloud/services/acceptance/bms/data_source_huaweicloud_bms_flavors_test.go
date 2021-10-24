package bms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBmsFlavorsDataSource_basic(t *testing.T) {
	resourceName := "data.huaweicloud_bms_flavors.demo"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBmsFlavorsDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBmsFlavorDataSourceID(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.#"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.id"),
					resource.TestCheckResourceAttr(resourceName, "flavors.0.vcpus", "48"),
					resource.TestCheckResourceAttr(resourceName, "flavors.0.cpu_arch", "x86_64"),
				),
			},
		},
	})
}

func testAccCheckBmsFlavorDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find bms flavors data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("BMS flavors data source ID not set")
		}

		return nil
	}
}

const testAccBmsFlavorsDataSource_basic = `
data "huaweicloud_bms_flavors" "demo" {
  vcpus = 48
}
`

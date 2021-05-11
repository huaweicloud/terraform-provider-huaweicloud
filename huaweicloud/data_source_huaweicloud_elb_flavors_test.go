package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccElbFlavorsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccElbFlavorsDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElbFlavorDataSourceID("data.huaweicloud_elb_flavors.this"),
				),
			},
		},
	})
}

func testAccCheckElbFlavorDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find elb flavors data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Elb Flavors data source ID not set")
		}

		return nil
	}
}

const testAccElbFlavorsDataSource_basic = `
data "huaweicloud_elb_flavors" "this" {
  type = "L7"
}
`

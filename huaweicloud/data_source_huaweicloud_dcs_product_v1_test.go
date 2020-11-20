package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDcsProductV1DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsProductV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsProductV1DataSourceID("data.huaweicloud_dcs_product_v1.product1"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_dcs_product_v1.product1", "spec_code", "dcs.single_node"),
				),
			},
		},
	})
}

func testAccCheckDcsProductV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Dcs product data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Dcs product data source ID not set")
		}

		return nil
	}
}

var testAccDcsProductV1DataSource_basic = fmt.Sprintf(`
data "huaweicloud_dcs_product_v1" "product1" {
spec_code = "dcs.single_node"
}
`)

package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDatasourceAccount_basic(t *testing.T) {
	rName := "data.huaweicloud_account.current"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceAccount_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAccountDataSourceID(rName),
					resource.TestCheckResourceAttrSet(rName, "name"),
				),
			},
		},
	})
}

func testAccCheckAccountDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find the account data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("the account data source ID not set")
		}

		return nil
	}
}

const testAccDatasourceAccount_basic = `
data "huaweicloud_account" "current" {}
`

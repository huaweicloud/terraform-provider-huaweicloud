package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGaussDBMysqlDehResourceDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBMysqlDehResourceDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaussDBMysqlDehResourceDataSourceID("data.huaweicloud_gaussdb_mysql_dedicated_resource.test"),
				),
			},
		},
	})
}

func testAccCheckGaussDBMysqlDehResourceDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find GaussDB mysql dedicated resource data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("the GaussDB mysql dedicated resource data source ID not set ")
		}

		return nil
	}
}

const testAccGaussDBMysqlDehResourceDataSource_basic = `
data "huaweicloud_gaussdb_mysql_dedicated_resource" "test" {
}
`

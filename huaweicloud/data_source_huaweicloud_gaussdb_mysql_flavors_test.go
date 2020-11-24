package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccHuaweiCloudGaussdbMysqlFlavorsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHuaweiCloudGaussdbMysqlFlavorsDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaussdbMysqlFlavorsDataSourceID("data.huaweicloud_gaussdb_mysql_flavors.flavor"),
				),
			},
		},
	})
}

func testAccCheckGaussdbMysqlFlavorsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find GaussDB mysql data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("GaussDB mysql data source ID not set ")
		}

		return nil
	}
}

var testAccHuaweiCloudGaussdbMysqlFlavorsDataSource_basic = `
data "huaweicloud_gaussdb_mysql_flavors" "flavor" {
  engine = "gaussdb-mysql"
  version = "8.0"
  availability_zone_mode = "multi"
}
`

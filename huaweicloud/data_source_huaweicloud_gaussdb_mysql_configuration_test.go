package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccGaussdbMysqlConfigurationDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussdbMysqlConfigurationDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.huaweicloud_gaussdb_mysql_configuration.test", "name", "Default-GaussDB-for-MySQL 8.0"),
				),
			},
		},
	})
}

const testAccGaussdbMysqlConfigurationDataSource_basic = `
data "huaweicloud_gaussdb_mysql_configuration" "test" {
  name = "Default-GaussDB-for-MySQL 8.0"
}
`

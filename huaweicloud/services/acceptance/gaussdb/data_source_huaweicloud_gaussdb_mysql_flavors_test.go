package gaussdb

import (
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccHuaweiCloudGaussdbMysqlFlavorsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
			return fmtp.Errorf("Can't find GaussDB mysql data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("GaussDB mysql data source ID not set ")
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

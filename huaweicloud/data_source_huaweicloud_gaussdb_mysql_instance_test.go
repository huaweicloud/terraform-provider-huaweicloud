package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccGaussdbMysqlInstanceDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussdbMysqlInstanceDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaussdbMysqlInstanceDataSourceID("data.huaweicloud_gaussdb_mysql_instance.test"),
				),
			},
		},
	})
}

func testAccCheckGaussdbMysqlInstanceDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find GaussDB mysql instance data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("GaussDB mysql instance data source ID not set ")
		}

		return nil
	}
}

func testAccGaussdbMysqlInstanceDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

resource "huaweicloud_gaussdb_mysql_instance" "test" {
  name        = "%s"
  password    = "Test@123"
  flavor      = "gaussdb.mysql.4xlarge.x86.4"
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id

  security_group_id = data.huaweicloud_networking_secgroup.test.id
}

data "huaweicloud_gaussdb_mysql_instance" "test" {
  name = huaweicloud_gaussdb_mysql_instance.test.name
  depends_on = [
    huaweicloud_gaussdb_mysql_instance.test,
  ]
}
`, testAccVpcConfig_Base(rName), rName)
}

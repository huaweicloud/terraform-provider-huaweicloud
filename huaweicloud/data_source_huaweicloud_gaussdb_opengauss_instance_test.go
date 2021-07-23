package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOpenGaussInstanceDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenGaussInstanceDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOpenGaussInstanceDataSourceID("data.huaweicloud_gaussdb_opengauss_instance.test"),
				),
			},
		},
	})
}

func testAccCheckOpenGaussInstanceDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find GaussDB opengauss instance data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("GaussDB opengauss data source ID not set ")
		}

		return nil
	}
}

func testAccOpenGaussInstanceDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

resource "huaweicloud_gaussdb_opengauss_instance" "test" {
  name        = "%s"
  password    = "Test@123"
  flavor      = "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id

  availability_zone = "cn-north-4a,cn-north-4a,cn-north-4a"
  security_group_id = data.huaweicloud_networking_secgroup.test.id

  ha {
    mode             = "enterprise"
    replication_mode = "sync"
    consistency      = "strong"
  }
  volume {
    type = "ULTRAHIGH"
    size = 40
  }

  sharding_num = 1
  coordinator_num = 2
}

data "huaweicloud_gaussdb_opengauss_instance" "test" {
  name = huaweicloud_gaussdb_opengauss_instance.test.name
}
`, testAccVpcConfig_Base(rName), rName)
}

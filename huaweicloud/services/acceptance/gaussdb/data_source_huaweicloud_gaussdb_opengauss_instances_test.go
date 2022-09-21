package gaussdb

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOpenGaussInstancesDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenGaussInstancesDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOpenGaussInstancesDataSourceID("data.huaweicloud_gaussdb_opengauss_instances.test"),
					resource.TestCheckResourceAttr("data.huaweicloud_gaussdb_opengauss_instances.test", "instances.#", "1"),
				),
			},
		},
	})
}

func testAccCheckOpenGaussInstancesDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find GaussDB opengauss instances data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("GaussDB opengauss data source ID not set ")
		}

		return nil
	}
}

func testAccOpenGaussInstancesDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup" "test" {
  name        = "%s"
  description = "terraform security group rule acceptance test"
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_gaussdb_opengauss_instance" "test" {
  name      = "%s"
  password  = "Test@12345678"
  flavor    = "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"
  vpc_id    = huaweicloud_vpc.test.id
  subnet_id = huaweicloud_vpc_subnet.test.id

  availability_zone = "cn-north-4a,cn-north-4a,cn-north-4a"
  security_group_id = huaweicloud_networking_secgroup.test.id

  ha {
    mode             = "enterprise"
    replication_mode = "sync"
    consistency      = "strong"
  }
  volume {
    type = "ULTRAHIGH"
    size = 40
  }

  sharding_num    = 1
  coordinator_num = 2
}

data "huaweicloud_gaussdb_opengauss_instances" "test" {
  name = huaweicloud_gaussdb_opengauss_instance.test.name
  depends_on = [
    huaweicloud_gaussdb_opengauss_instance.test,
  ]
}
`, testAccVpcConfig_Base(rName), rName, rName)
}

package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/taurusdb/v3/instances"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getResourceProxy(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.GaussdbV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error creating HuaweiCloud GaussDB client: %s", err)
	}

	return instances.GetProxy(client, state.Primary.ID).Extract()
}

func TestAccGaussDBProxy_basic(t *testing.T) {
	var proxy instances.Proxy
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_gaussdb_mysql_proxy.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&proxy,
		getResourceProxy,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlProxy_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "flavor", "gaussdb.proxy.xlarge.arm.2"),
					resource.TestCheckResourceAttr(resourceName, "node_num", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"instance_id",
					"flavor",
					"node_num",
				},
			},
		},
	})
}

func testAccMysqlProxy_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  vpc_id     = huaweicloud_vpc.test.id
  gateway_ip = "192.168.0.1"
  cidr       = "192.168.0.0/24"

  timeouts {
    delete = "20m"
  }
}

resource "huaweicloud_gaussdb_mysql_instance" "test" {
  name        = "%s"
  password    = "Test@12345678"
  flavor      = "gaussdb.mysql.2xlarge.x86.4"
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id

  security_group_id = data.huaweicloud_networking_secgroup.test.id

  enterprise_project_id = "0"
}

resource "huaweicloud_gaussdb_mysql_proxy" "test" {
  instance_id = huaweicloud_gaussdb_mysql_instance.test.id
  flavor      = "gaussdb.proxy.xlarge.arm.2"
  node_num    = 3
}
`, rName, rName, rName)
}

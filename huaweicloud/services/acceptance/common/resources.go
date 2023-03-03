package common

import (
	"fmt"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// TestSecGroup can be referred as `huaweicloud_networking_secgroup.test`
func TestSecGroup(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "test" {
  name                 = "%s"
  delete_default_rules = true
}
`, name)
}

// TestVpc can be referred as `huaweicloud_vpc.test` and `huaweicloud_vpc_subnet.test`
func TestVpc(name string) string {
	randCidr, randGatewayIp := acceptance.RandomCidrAndGatewayIp()
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "%[2]s"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%[1]s"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = "%[2]s"
  gateway_ip = "%[3]s"
}
`, name, randCidr, randGatewayIp)
}

// TestBaseNetwork vpc, subnet, security group
func TestBaseNetwork(name string) string {
	return fmt.Sprintf(`
# base security group without default rules
%s

# base vpc and subnet
%s
`, TestSecGroup(name), TestVpc(name))
}

// TestBaseComputeResources vpc, subnet, security group, availability zone, keypair, image, flavor
func TestBaseComputeResources(name string) string {
	return fmt.Sprintf(`
# base test resources
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}
`, TestBaseNetwork(name))
}

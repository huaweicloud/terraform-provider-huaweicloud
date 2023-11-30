package dns

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dns/v2/endpoints"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getDNSEndpointResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.DNSV21Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating dns client: %s", err)
	}

	return endpoints.Get(client, state.Primary.ID).Extract()
}

func TestAccDNSEndpoint_basic(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_dns_endpoint.test"
	)
	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDNSEndpointResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDNSEndpoint_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "direction", "inbound"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(rName, "resolver_rule_count", "0"),
					resource.TestCheckResourceAttr(rName, "ip_addresses.#", "2"),
					resource.TestCheckResourceAttr(rName, "ip_addresses.0.status", "ACTIVE"),
					resource.TestCheckResourceAttr(rName, "ip_addresses.1.status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(rName, "ip_addresses.0.subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "ip_addresses.1.subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "ip_addresses.0.ip"),
					resource.TestCheckResourceAttrSet(rName, "ip_addresses.1.ip"),
					resource.TestCheckResourceAttrSet(rName, "ip_addresses.0.ip_address_id"),
					resource.TestCheckResourceAttrSet(rName, "ip_addresses.1.ip_address_id"),
					resource.TestCheckResourceAttrSet(rName, "ip_addresses.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "ip_addresses.1.created_at"),
					resource.TestCheckResourceAttrSet(rName, "ip_addresses.0.updated_at"),
					resource.TestCheckResourceAttrSet(rName, "ip_addresses.1.updated_at"),
					resource.TestCheckResourceAttrSet(rName, "vpc_id"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testDNSEndpoint_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "direction", "inbound"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(rName, "resolver_rule_count", "0"),
					resource.TestCheckResourceAttr(rName, "ip_addresses.#", "2"),
					resource.TestCheckResourceAttr(rName, "ip_addresses.0.status", "ACTIVE"),
					resource.TestCheckResourceAttr(rName, "ip_addresses.1.status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(rName, "ip_addresses.0.subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "ip_addresses.1.subnet_id",
						"huaweicloud_vpc_subnet.test_update", "id"),
					resource.TestCheckResourceAttrSet(rName, "ip_addresses.0.ip"),
					resource.TestCheckResourceAttrSet(rName, "ip_addresses.1.ip"),
					resource.TestCheckResourceAttrSet(rName, "ip_addresses.0.ip_address_id"),
					resource.TestCheckResourceAttrSet(rName, "ip_addresses.1.ip_address_id"),
					resource.TestCheckResourceAttrSet(rName, "ip_addresses.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "ip_addresses.1.created_at"),
					resource.TestCheckResourceAttrSet(rName, "ip_addresses.0.updated_at"),
					resource.TestCheckResourceAttrSet(rName, "ip_addresses.1.updated_at"),
					resource.TestCheckResourceAttrSet(rName, "vpc_id"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testVpcSubnet_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name              = "%[1]s"
  cidr              = "192.168.0.0/24"
  gateway_ip        = "192.168.0.1"
  vpc_id            = huaweicloud_vpc.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}`, rName)
}

func testDNSEndpoint_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_endpoint" "test" {
  name      = "%s"
  direction = "inbound"
  ip_addresses {
    subnet_id = huaweicloud_vpc_subnet.test.id
  }
  ip_addresses {
    subnet_id = huaweicloud_vpc_subnet.test.id
  }
}`, testVpcSubnet_basic(rName), rName)
}

func testDNSEndpoint_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_subnet" "test_update" {
  name              = "%[2]s_update"
  cidr              = "192.168.100.0/24"
  gateway_ip        = "192.168.100.1"
  vpc_id            = huaweicloud_vpc.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_dns_endpoint" "test" {
  name      = "%[2]s_update"
  direction = "inbound"
  ip_addresses {
    subnet_id = huaweicloud_vpc_subnet.test.id
  }
  ip_addresses {
    subnet_id = huaweicloud_vpc_subnet.test_update.id
  }
}`, testVpcSubnet_basic(rName), rName)
}

package dns

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dns/v2/resolverrule"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getDNSResolverRuleFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.DNSV21Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating dns client: %s", err)
	}
	body, err := resolverrule.Get(client, state.Primary.ID).Extract()
	if err == nil && body.Status == "DELETED" {
		return nil, fmt.Errorf("DNS resolver rule does not found")
	}
	return body, err
}

func TestAccDNSResolverRule_basic(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_dns_resolver_rule.test"
	)
	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDNSResolverRuleFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDNSResolverRule_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "domain_name", "terraform.test.com."),
					resource.TestCheckResourceAttr(rName, "ip_addresses.#", "1"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet(rName, "rule_type"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrPair(rName, "endpoint_id",
						"huaweicloud_dns_endpoint.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "ip_addresses.0.ip",
						"huaweicloud_dns_endpoint.test", "ip_addresses.0.ip"),
				),
			},
			{
				Config: testDNSResolverRule_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "domain_name", "terraform.test.com."),
					resource.TestCheckResourceAttr(rName, "ip_addresses.#", "1"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet(rName, "rule_type"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrPair(rName, "endpoint_id",
						"huaweicloud_dns_endpoint.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "ip_addresses.0.ip",
						"huaweicloud_dns_endpoint.test", "ip_addresses.1.ip"),
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

func testDNSEndpoint(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%[1]s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.test.id
}

resource "huaweicloud_dns_endpoint" "test" {
  name      = "%[1]s"
  direction = "inbound"
  ip_addresses {
    subnet_id = huaweicloud_vpc_subnet.test.id
  }
  ip_addresses {
    subnet_id = huaweicloud_vpc_subnet.test.id
  }
}`, rName)
}

func testDNSResolverRule_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_resolver_rule" "test" {
  name        = "%s"
  domain_name = "terraform.test.com."
  endpoint_id = huaweicloud_dns_endpoint.test.id
  ip_addresses {
    ip = huaweicloud_dns_endpoint.test.ip_addresses[0].ip
  }
}`, testDNSEndpoint(rName), rName)
}

func testDNSResolverRule_basic_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_resolver_rule" "test" {
  name        = "%s_update"
  domain_name = "terraform.test.com."
  endpoint_id = huaweicloud_dns_endpoint.test.id
  ip_addresses {
    ip = huaweicloud_dns_endpoint.test.ip_addresses[1].ip
  }
}`, testDNSEndpoint(rName), rName)
}

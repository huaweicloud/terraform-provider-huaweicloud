package dns

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dns/v2/resolverrule"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getDNSResolverRuleAssociateResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.DNSV21Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating dns client: %s", err)
	}

	arr := strings.Split(state.Primary.ID, "/")
	if len(arr) != 2 {
		return nil, fmt.Errorf("error getting resolver ruler ID and VPC ID, resource ID is: %s", state.Primary.ID)
	}
	ruleID := arr[0]
	vpcID := arr[1]

	rule, err := resolverrule.Get(client, ruleID).Extract()
	if err != nil {
		return nil, err
	}

	for _, router := range rule.Routers {
		if router.RouterID == vpcID {
			return router, nil
		}
	}
	return nil, fmt.Errorf("the resolver rule associate does not exist")
}

func TestAccDNSResolverRuleAssociate_basic(t *testing.T) {
	var (
		obj         interface{}
		name        = acceptance.RandomAccResourceName()
		rName       = "huaweicloud_dns_resolver_rule_associate.test"
		randUUID, _ = uuid.GenerateUUID()
	)
	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDNSResolverRuleAssociateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testDNSResolverRuleAssociate_ruleNotExist(randUUID),
				ExpectError: regexp.MustCompile("Resolver rule not exist"),
			},
			{
				Config: testDNSResolverRuleAssociate_VPCNotExist(name, randUUID),
				// DNS.0711: The associated VPC does not exist.
				ExpectError: regexp.MustCompile("DNS.0711"),
			},
			{
				Config: testDNSResolverRuleAssociate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrPair(rName, "resolver_rule_id",
						"huaweicloud_dns_resolver_rule.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "vpc_id",
						"huaweicloud_vpc.test", "id"),
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

func testDNSResolverRuleAssociate_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name              = "%[1]s"
  cidr              = "192.168.0.0/24"
  gateway_ip        = "192.168.0.1"
  vpc_id            = huaweicloud_vpc.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
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
}

resource "huaweicloud_dns_resolver_rule" "test" {
  name        = "%[1]s"
  domain_name = "terraform.test.com."
  endpoint_id = huaweicloud_dns_endpoint.test.id
  ip_addresses {
    ip = huaweicloud_dns_endpoint.test.ip_addresses[0].ip
  }
}
`, rName)
}

func testDNSResolverRuleAssociate_ruleNotExist(randUUID string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_resolver_rule_associate" "test" {
  resolver_rule_id = "%[1]s"
  vpc_id           = "%[1]s"
}
`, randUUID)
}

func testDNSResolverRuleAssociate_VPCNotExist(name, randUUID string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_resolver_rule_associate" "test" {
  resolver_rule_id = huaweicloud_dns_resolver_rule.test.id
  vpc_id           = "%s"
}
`, testDNSResolverRuleAssociate_base(name), randUUID)
}

func testDNSResolverRuleAssociate_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_resolver_rule_associate" "test" {
  resolver_rule_id = huaweicloud_dns_resolver_rule.test.id
  vpc_id           = huaweicloud_vpc.test.id
}
`, testDNSResolverRuleAssociate_base(rName))
}

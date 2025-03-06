package dns

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dns"
)

func getResolverRuleAssociatedVpc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.DNSV21Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating dns client: %s", err)
	}

	return dns.GetAssociatedVpcById(client, state.Primary.Attributes["resolver_rule_id"], state.Primary.Attributes["vpc_id"])
}

func TestAccResolverRuleAssociate_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		associatedVpc interface{}
		rName         = "huaweicloud_dns_resolver_rule_associate.test"
		rc            = acceptance.InitResourceCheck(rName, &associatedVpc, getResolverRuleAssociatedVpc)

		randomId, _ = uuid.GenerateUUID()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testResolverRuleAssociate_ruleNotExist(randomId),
				ExpectError: regexp.MustCompile("Resolver rule not exist"),
			},
			{
				Config: testResolverRuleAssociate_VPCNotExist(name, randomId),
				// DNS.0711: The associated VPC does not exist.
				ExpectError: regexp.MustCompile("DNS.0711"),
			},
			{
				Config: testResolverRuleAssociate_basic(name),
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

func testResolverRuleAssociate_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dns_endpoint" "test" {
  name      = "%[2]s"
  direction = "inbound"

  ip_addresses {
    subnet_id = huaweicloud_vpc_subnet.test.id
  }
  ip_addresses {
    subnet_id = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_dns_resolver_rule" "test" {
  name        = "%[2]s"
  domain_name = "terraform.test.com."
  endpoint_id = huaweicloud_dns_endpoint.test.id

  ip_addresses {
    ip = huaweicloud_dns_endpoint.test.ip_addresses[0].ip
  }
}
`, common.TestVpc(rName), rName)
}

func testResolverRuleAssociate_ruleNotExist(randomId string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_resolver_rule_associate" "test" {
  resolver_rule_id = "%[1]s"
  vpc_id           = "%[1]s"
}
`, randomId)
}

func testResolverRuleAssociate_VPCNotExist(name, randomId string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_resolver_rule_associate" "test" {
  resolver_rule_id = huaweicloud_dns_resolver_rule.test.id
  vpc_id           = "%s"
}
`, testResolverRuleAssociate_base(name), randomId)
}

func testResolverRuleAssociate_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dns_resolver_rule_associate" "test" {
  resolver_rule_id = huaweicloud_dns_resolver_rule.test.id
  vpc_id           = huaweicloud_vpc.test.id
}
`, testResolverRuleAssociate_base(rName))
}

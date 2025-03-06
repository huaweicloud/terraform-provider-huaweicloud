package dns

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dns"
)

func getResolverRule(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.DNSV21Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating dns client: %s", err)
	}

	return dns.GetResolverRuleById(client, state.Primary.ID)
}

func TestAccResolverRule_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceNameWithDash()
		domainName = fmt.Sprintf("%s.", name)

		resolverRule interface{}
		rName        = "huaweicloud_dns_resolver_rule.test"
		rc           = acceptance.InitResourceCheck(rName, &resolverRule, getResolverRule)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResolverRule_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "domain_name", domainName),
					resource.TestCheckResourceAttr(rName, "ip_addresses.#", "1"),
					resource.TestCheckResourceAttrPair(rName, "endpoint_id",
						"huaweicloud_dns_endpoint.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "ip_addresses.0.ip",
						"huaweicloud_dns_endpoint.test", "ip_addresses.0.ip"),
					// Check attributtes.
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet(rName, "rule_type"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testAccResolverRule_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "domain_name", domainName),
					resource.TestCheckResourceAttr(rName, "ip_addresses.#", "2"),
					resource.TestCheckResourceAttrSet(rName, "ip_addresses.0.ip"),
					resource.TestCheckResourceAttrSet(rName, "ip_addresses.1.ip"),
					resource.TestCheckResourceAttrPair(rName, "endpoint_id",
						"huaweicloud_dns_endpoint.test", "id"),
					// Check attributtes.
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet(rName, "rule_type"),
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

func testAccResolverRule_base(rName string) string {
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
}`, common.TestVpc(rName), rName)
}

func testAccResolverRule_basic_step1(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dns_resolver_rule" "test" {
  name        = "%[2]s"
  domain_name = "%[2]s"
  endpoint_id = huaweicloud_dns_endpoint.test.id

  ip_addresses {
    ip = huaweicloud_dns_endpoint.test.ip_addresses[0].ip
  }
}`, testAccResolverRule_base(rName), rName)
}

func testAccResolverRule_basic_step2(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dns_resolver_rule" "test" {
  name        = "%[2]s_update"
  domain_name = "%[2]s"
  endpoint_id = huaweicloud_dns_endpoint.test.id

  dynamic "ip_addresses" {
    for_each = huaweicloud_dns_endpoint.test.ip_addresses[*].ip
    content {
      ip = ip_addresses.value
    }
  }
}`, testAccResolverRule_base(rName), rName)
}

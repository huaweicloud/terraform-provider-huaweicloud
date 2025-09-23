package dns

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDNSResolverRules_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dns_resolver_rules.test"
	rName := acceptance.RandomAccResourceNameWithDash()

	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDNSResolverRules_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resolver_rules.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resolver_rules.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "resolver_rules.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "resolver_rules.0.endpoint_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resolver_rules.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "resolver_rules.0.rule_type"),
					resource.TestCheckResourceAttrSet(dataSource, "resolver_rules.0.routers.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resolver_rules.0.routers.0.router_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resolver_rules.0.routers.0.router_region"),
					resource.TestCheckResourceAttrSet(dataSource, "resolver_rules.0.routers.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "resolver_rules.0.domain_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resolver_rules.0.ipaddress_count"),
					resource.TestCheckResourceAttrSet(dataSource, "resolver_rules.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "resolver_rules.0.update_time"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_domain_name_filter_useful", "true"),
					resource.TestCheckOutput("is_endpoint_id_filter_useful", "true"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDNSResolverRules_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dns_resolver_rules" "test" {
  depends_on = [huaweicloud_dns_resolver_rule_associate.test]
}

// filter by name
data "huaweicloud_dns_resolver_rules" "filter_by_name" {
  name = huaweicloud_dns_resolver_rule.test.name
}

locals {
  filter_result_by_name = [for v in data.huaweicloud_dns_resolver_rules.filter_by_name.resolver_rules[*].name :
    v == huaweicloud_dns_resolver_rule.test.name]
}

output "is_name_filter_useful" {
  value = length(local.filter_result_by_name) > 0 && alltrue(local.filter_result_by_name)
}

// filter by domain_name
data "huaweicloud_dns_resolver_rules" "filter_by_domain_name" {
  domain_name = huaweicloud_dns_resolver_rule.test.domain_name
}

locals {
  filter_result_by_domain_name = [for v in data.huaweicloud_dns_resolver_rules.filter_by_domain_name.resolver_rules[*].domain_name :
    v == huaweicloud_dns_resolver_rule.test.domain_name]
}

output "is_domain_name_filter_useful" {
  value = length(local.filter_result_by_domain_name) > 0 && alltrue(local.filter_result_by_domain_name)
}

// filter by endpoint_id
data "huaweicloud_dns_resolver_rules" "filter_by_endpoint_id" {
  endpoint_id = huaweicloud_dns_resolver_rule.test.endpoint_id
}

locals {
  filter_result_by_endpoint_id = [for v in data.huaweicloud_dns_resolver_rules.filter_by_endpoint_id.resolver_rules[*].endpoint_id :
    v == huaweicloud_dns_resolver_rule.test.endpoint_id]
}

output "is_endpoint_id_filter_useful" {
  value = length(local.filter_result_by_endpoint_id) > 0 && alltrue(local.filter_result_by_endpoint_id)
}

// filter by id
data "huaweicloud_dns_resolver_rules" "filter_by_id" {
  resolver_rule_id = huaweicloud_dns_resolver_rule.test.id
}

locals {
  filter_result_by_id = [for v in data.huaweicloud_dns_resolver_rules.filter_by_id.resolver_rules[*].id :
    v == huaweicloud_dns_resolver_rule.test.id]
}

output "is_id_filter_useful" {
  value = length(local.filter_result_by_id) > 0 && alltrue(local.filter_result_by_id)
}
`, testResolverRuleAssociate_basic(name))
}

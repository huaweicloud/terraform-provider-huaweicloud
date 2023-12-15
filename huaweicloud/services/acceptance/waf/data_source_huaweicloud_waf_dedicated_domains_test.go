package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceWAFDedicatedDomains_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_waf_dedicated_domains.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDedicatedDomains_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "domains.0.domain"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.id"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.pci_3ds"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.pci_dds"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.is_dual_az"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.ipv6_enable"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.description"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.policy_id"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.protect_status"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.access_status"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.website_name"),

					resource.TestCheckOutput("domain_filter_is_useful", "true"),

					resource.TestCheckOutput("protect_status_filter_is_useful", "true"),

					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceDedicatedDomains_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_waf_dedicated_domains" "test" {
  depends_on = [huaweicloud_waf_dedicated_domain.domain_1]
}

data "huaweicloud_waf_dedicated_domains" "domain_filter" {
  domain = data.huaweicloud_waf_dedicated_domains.test.domains.0.domain
}

locals {
  domain = data.huaweicloud_waf_dedicated_domains.test.domains.0.domain
}

output "domain_filter_is_useful" {
  value = length(data.huaweicloud_waf_dedicated_domains.domain_filter.domains) > 0 && alltrue(
    [for v in data.huaweicloud_waf_dedicated_domains.domain_filter.domains[*].domain : v == local.domain]
  )  
}

data "huaweicloud_waf_dedicated_domains" "protect_status_filter" {
  protect_status = data.huaweicloud_waf_dedicated_domains.test.domains.0.protect_status
}
  
locals {
  protect_status = data.huaweicloud_waf_dedicated_domains.test.domains.0.protect_status
}
  
output "protect_status_filter_is_useful" {
  value = length(data.huaweicloud_waf_dedicated_domains.protect_status_filter.domains) > 0 && alltrue(
    [for v in data.huaweicloud_waf_dedicated_domains.protect_status_filter.domains[*].protect_status : v == local.protect_status]
  )  
}

data "huaweicloud_waf_dedicated_domains" "enterprise_project_id_filter" {
  enterprise_project_id = data.huaweicloud_waf_dedicated_domains.test.domains.0.enterprise_project_id
}
	
locals {
  enterprise_project_id = data.huaweicloud_waf_dedicated_domains.test.domains.0.enterprise_project_id
}

output "enterprise_project_id_filter_is_useful" {
  value = length(data.huaweicloud_waf_dedicated_domains.enterprise_project_id_filter.domains) > 0 && alltrue(
    [for v in data.huaweicloud_waf_dedicated_domains.enterprise_project_id_filter.domains[*].enterprise_project_id : v == local.enterprise_project_id]
  )  
}
`, testAccWafDedicatedDomainV1_basic(name))
}

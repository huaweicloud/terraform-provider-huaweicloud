package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceWAFDomains_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_waf_domains.test"
	policyNameTestName := "data.huaweicloud_waf_domains.policy_name_filter"
	domainName := fmt.Sprintf("%s.huawei.com", name)
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDomains_basic(name, domainName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "domains.0.id"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.proxy"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.domain"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.policy_id"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.protect_status"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.pci_3ds"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.pci_dss"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.access_status"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.access_code"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(policyNameTestName, "domains.#"),

					resource.TestCheckOutput("domain_filter_is_useful", "true"),

					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceDomains_basic(name string, domainName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_waf_domains" "test" {
  depends_on = [huaweicloud_waf_domain.domain_1]
}

data "huaweicloud_waf_domains" "domain_filter" {
  domain = data.huaweicloud_waf_domains.test.domains.0.domain
}

locals {
  domain = data.huaweicloud_waf_domains.test.domains.0.domain
}

output "domain_filter_is_useful" {
  value = length(data.huaweicloud_waf_domains.domain_filter.domains) > 0 && alltrue(
    [for v in data.huaweicloud_waf_domains.domain_filter.domains[*].domain : v == local.domain]
  )  
}

data "huaweicloud_waf_domains" "policy_name_filter" {
  policy_name = huaweicloud_waf_policy.policy_1.name
}

data "huaweicloud_waf_domains" "enterprise_project_id_filter" {
  enterprise_project_id = data.huaweicloud_waf_domains.test.domains.0.enterprise_project_id
}
	
locals {
  enterprise_project_id = data.huaweicloud_waf_domains.test.domains.0.enterprise_project_id
}

output "enterprise_project_id_filter_is_useful" {
  value = length(data.huaweicloud_waf_domains.enterprise_project_id_filter.domains) > 0 && alltrue(
    [for v in data.huaweicloud_waf_domains.enterprise_project_id_filter.domains[*].enterprise_project_id : v == local.enterprise_project_id]
  )  
}
`, testAccWafDomainV1_policy(name, domainName))
}

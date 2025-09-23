package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running the test case, please ensure that there is at least one WAF dedicated instance in the current region.
func TestAccDataSourceDedicatedDomains_basic(t *testing.T) {
	var (
		name            = acceptance.RandomAccResourceName()
		certificateBody = generateCertificateBody()

		rName = "data.huaweicloud_waf_dedicated_domains.test"
		dc    = acceptance.InitDataSourceCheck(rName)

		byDomain   = "data.huaweicloud_waf_dedicated_domains.domain_filter"
		dcByDomain = acceptance.InitDataSourceCheck(byDomain)

		byProtectStatus   = "data.huaweicloud_waf_dedicated_domains.protect_status_filter"
		dcByProtectStatus = acceptance.InitDataSourceCheck(byProtectStatus)

		byAllParameters   = "data.huaweicloud_waf_dedicated_domains.all_parameters_filter"
		dcByAllParameters = acceptance.InitDataSourceCheck(byAllParameters)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDedicatedDomains_basic(name, certificateBody),
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

					dcByDomain.CheckResourceExists(),
					resource.TestCheckOutput("domain_filter_is_useful", "true"),

					dcByProtectStatus.CheckResourceExists(),
					resource.TestCheckOutput("protect_status_filter_is_useful", "true"),

					dcByAllParameters.CheckResourceExists(),
					resource.TestCheckOutput("all_parameters_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceDedicatedDomains_basic(name, certificateBody string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_waf_dedicated_domains" "test" {
  enterprise_project_id = "%[2]s"

  depends_on = [huaweicloud_waf_dedicated_domain.test]
}

# Filter by domain
locals {
  domain = data.huaweicloud_waf_dedicated_domains.test.domains.0.domain
}

data "huaweicloud_waf_dedicated_domains" "domain_filter" {
  domain                = local.domain
  enterprise_project_id = "%[2]s"
}

locals {
  domain_filter_result = [
    for v in data.huaweicloud_waf_dedicated_domains.domain_filter.domains[*].domain : v == local.domain
  ]
}

output "domain_filter_is_useful" {
  value = length(local.domain_filter_result) > 0 && alltrue(local.domain_filter_result)  
}

# Filter by protect_status
locals {
  protect_status = tostring(data.huaweicloud_waf_dedicated_domains.test.domains.0.protect_status)
}

data "huaweicloud_waf_dedicated_domains" "protect_status_filter" {
  protect_status        = local.protect_status
  enterprise_project_id = "%[2]s"
}
  
locals {
  protect_status_filter_result = [
    for v in data.huaweicloud_waf_dedicated_domains.protect_status_filter.domains[*].protect_status :
    tostring(v) == local.protect_status
  ]
}

output "protect_status_filter_is_useful" {
  value = length(local.protect_status_filter_result) > 0 && alltrue(local.protect_status_filter_result)  
}

# Filter by all parameters
data "huaweicloud_waf_dedicated_domains" "all_parameters_filter" {
  domain                = local.domain
  protect_status        = local.protect_status
  enterprise_project_id = "%[2]s"
}

locals {
  all_parameters_filter_result = [
    for v in data.huaweicloud_waf_dedicated_domains.all_parameters_filter.domains[*] :
    tostring(v.protect_status) == local.protect_status && v.domain == local.domain
  ]
}

output "all_parameters_filter_is_useful" {
  value = length(local.all_parameters_filter_result) > 0 && alltrue(local.all_parameters_filter_result)  
}
`, testAccWafDedicatedDomain_basic(name, certificateBody), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

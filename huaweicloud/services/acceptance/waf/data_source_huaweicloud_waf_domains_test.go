package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running the test case, please ensure that there is at least one WAF cloud instance in the current region.
func TestAccDataSourceDomains_basic(t *testing.T) {
	var (
		name            = acceptance.RandomAccResourceName()
		domainName      = fmt.Sprintf("%s.huawei.com", name)
		certificateBody = testAccWafCertificate_basic(name, generateCertificateBody())

		rName = "data.huaweicloud_waf_domains.test"
		dc    = acceptance.InitDataSourceCheck(rName)

		byDomain   = "data.huaweicloud_waf_domains.domain_filter"
		dcByDomain = acceptance.InitDataSourceCheck(byDomain)

		byPolicyName   = "data.huaweicloud_waf_domains.policy_name_filter"
		dcByPolicyName = acceptance.InitDataSourceCheck(byPolicyName)
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
				Config: testAccDatasourceDomains_basic(certificateBody, name, domainName),
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

					dcByDomain.CheckResourceExists(),
					resource.TestCheckOutput("domain_filter_is_useful", "true"),

					dcByPolicyName.CheckResourceExists(),
					resource.TestCheckOutput("policy_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceDomains_base(certificateBody, name, domainName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_policy" "test" {
  name                  = "%[2]s"
  enterprise_project_id = "%[4]s"
}

resource "huaweicloud_waf_domain" "test" {
  domain                = "%[3]s"
  certificate_id        = huaweicloud_waf_certificate.test.id
  certificate_name      = huaweicloud_waf_certificate.test.name
  policy_id             = huaweicloud_waf_policy.test.id
  keep_policy           = true
  proxy                 = true
  enterprise_project_id = "%[4]s"

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8080
    type            = "ipv4"
    weight          = 3
  }
}
`, certificateBody, name, domainName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccDatasourceDomains_basic(certificateBody, name, domainName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_waf_domains" "test" {
  enterprise_project_id = "%[2]s"

  depends_on = [
    huaweicloud_waf_domain.test
  ]
}

# Filter by domain
locals {
  domain = data.huaweicloud_waf_domains.test.domains.0.domain
}

data "huaweicloud_waf_domains" "domain_filter" {
  domain                = local.domain
  enterprise_project_id = "%[2]s"
}

locals {
  domain_filter_result = [
    for v in data.huaweicloud_waf_domains.domain_filter.domains[*].domain : v == local.domain
  ]
}

output "domain_filter_is_useful" {
  value = length(local.domain_filter_result) > 0 && alltrue(local.domain_filter_result)  
}

# Filter by policy_name
locals {
  policy_name = huaweicloud_waf_policy.test.name
}

data "huaweicloud_waf_domains" "policy_name_filter" {
  policy_name           = local.policy_name
  enterprise_project_id = "%[2]s"

  depends_on = [
    huaweicloud_waf_domain.test
  ]
}

output "policy_name_filter_is_useful" {
  value = length(data.huaweicloud_waf_domains.policy_name_filter.domains) > 0
}
`, testAccDatasourceDomains_base(certificateBody, name, domainName), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

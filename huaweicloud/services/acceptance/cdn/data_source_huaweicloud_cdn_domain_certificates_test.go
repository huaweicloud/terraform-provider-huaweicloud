package cdn

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before executing this use case, please create several pieces of data first.
func TestAccDataDomainCertificates_basic(t *testing.T) {
	rName := "data.huaweicloud_cdn_domain_certificates.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDomainCertificates_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "domain_certificates.0.domain_id"),
					resource.TestCheckResourceAttrSet(rName, "domain_certificates.0.domain_name"),
					resource.TestCheckResourceAttrSet(rName, "domain_certificates.0.certificate_name"),
					resource.TestCheckResourceAttrSet(rName, "domain_certificates.0.certificate_body"),
					resource.TestCheckResourceAttrSet(rName, "domain_certificates.0.certificate_source"),
					resource.TestCheckResourceAttrSet(rName, "domain_certificates.0.expire_at"),
					resource.TestCheckResourceAttrSet(rName, "domain_certificates.0.https_status"),
					resource.TestCheckResourceAttrSet(rName, "domain_certificates.0.force_redirect_https"),
					resource.TestCheckResourceAttrSet(rName, "domain_certificates.0.http2_enabled"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataDomainCertificates_basic() string {
	return `
data "huaweicloud_cdn_domain_certificates" "test" {}

data "huaweicloud_cdn_domain_certificates" "name_filter" {
  name = data.huaweicloud_cdn_domain_certificates.test.domain_certificates.0.domain_name
}
locals {
  name = data.huaweicloud_cdn_domain_certificates.test.domain_certificates.0.domain_name
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_cdn_domain_certificates.name_filter.domain_certificates) > 0 && alltrue(
    [for v in data.huaweicloud_cdn_domain_certificates.name_filter.domain_certificates[*].domain_name : v == local.name]
  )
}
`
}

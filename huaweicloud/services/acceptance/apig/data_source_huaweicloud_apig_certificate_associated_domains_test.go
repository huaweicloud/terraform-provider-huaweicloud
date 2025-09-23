package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataCertificateAssociatedDomains_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_apig_certificate_associated_domains.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byUrlDomain   = "data.huaweicloud_apig_certificate_associated_domains.filter_by_url_domain"
		dcByUrlDomain = acceptance.InitDataSourceCheck(byUrlDomain)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
			acceptance.TestAccPreCheckCertificateBase(t)
			acceptance.TestAccPreCheckCertificateRootCA(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataCertificateAssociatedDomains_certificateNotFound(),
				ExpectError: regexp.MustCompile(`ssl with certId [a-f0-9]+ not found`),
			},
			{
				Config: testAccDataCertificateAssociatedDomains_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "domains.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByUrlDomain.CheckResourceExists(),
					resource.TestCheckOutput("url_domain_filter_is_useful", "true"),
					resource.TestCheckResourceAttrSet(dataSource, "domains.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "domains.0.url_domain"),
					resource.TestCheckResourceAttrSet(dataSource, "domains.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "domains.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "domains.0.min_ssl_version"),
					resource.TestCheckResourceAttr(dataSource, "domains.0.verified_client_certificate_enabled", "true"),
					resource.TestCheckResourceAttrSet(dataSource, "domains.0.api_group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "domains.0.api_group_name"),
				),
			},
		},
	})
}

func testAccDataCertificateAssociatedDomains_base() string {
	name := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
resource "huaweicloud_apig_certificate" "test" {
  instance_id     = "%[1]s"
  name            = "%[2]s"
  type            = "instance"
  content         = "%[3]s"
  private_key     = "%[4]s"
  trusted_root_ca = "%[5]s"
}

data "huaweicloud_apig_instance_ssl_certificates" "test" {
  instance_id = "%[1]s"
  name        = huaweicloud_apig_certificate.test.name

  depends_on = [huaweicloud_apig_certificate.test]
}

locals {
  domain_name = try([for v in data.huaweicloud_apig_instance_ssl_certificates.test.certificates : v.common_name
  if v.id == huaweicloud_apig_certificate.test.id][0], null)
}

resource "huaweicloud_apig_group" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
}

resource "huaweicloud_apig_group_domain_associate" "test" {
  instance_id     = "%[1]s"
  group_id        = huaweicloud_apig_group.test.id
  url_domain      = local.domain_name
  min_ssl_version = "TLSv1.1"
}

resource "huaweicloud_apig_certificate_batch_domains_associate" "test" {
  instance_id    = "%[1]s"
  certificate_id = huaweicloud_apig_certificate.test.id

  verify_enabled_domain_names = [local.domain_name]

  depends_on = [huaweicloud_apig_group_domain_associate.test]
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID,
		name,
		acceptance.HW_CERTIFICATE_CONTENT,
		acceptance.HW_CERTIFICATE_PRIVATE_KEY,
		acceptance.HW_CERTIFICATE_ROOT_CA,
	)
}

func testAccDataCertificateAssociatedDomains_certificateNotFound() string {
	randomUUID, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_apig_certificate_associated_domains" "test" {
  certificate_id = replace("%[1]s", "-", "")
}
`, randomUUID)
}

func testAccDataCertificateAssociatedDomains_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_apig_certificate_associated_domains" "test" {
  certificate_id = huaweicloud_apig_certificate.test.id

  depends_on = [huaweicloud_apig_certificate_batch_domains_associate.test]
}

data "huaweicloud_apig_certificate_associated_domains" "filter_by_url_domain" {
  certificate_id = huaweicloud_apig_certificate.test.id
  url_domain     = local.domain_name

  depends_on = [huaweicloud_apig_certificate_batch_domains_associate.test]
}

locals {
  url_domain_filter_result = [
    for v in data.huaweicloud_apig_certificate_associated_domains.filter_by_url_domain.domains[*].url_domain : v == local.domain_name
  ]
}

output "url_domain_filter_is_useful" {
  value = length(local.url_domain_filter_result) > 0 && alltrue(local.url_domain_filter_result)
}
`, testAccDataCertificateAssociatedDomains_base())
}

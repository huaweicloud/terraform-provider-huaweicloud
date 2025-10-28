package cdn

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCertificateAssociateDomains_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCDNTargetDomainUrls(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateAssociateDomains_basic_step1(name),
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				Config:      testAccCertificateAssociateDomains_basic_step2(),
				ExpectError: regexp.MustCompile("error associating certificate with domains"),
			},
		},
	})
}

func testAccCertificateAssociateDomains_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_certificate_associate_domains" "test" {
  domain_names = "%[1]s"
  https_switch = 1
  cert_name    = "%[2]s"
  certificate  = file("%[3]s")
  private_key  = file("%[4]s")
}
`, acceptance.HW_CDN_TARGET_DOMAIN_URLS, name,
		acceptance.HW_CCM_CERTIFICATE_CONTENT_PATH,
		acceptance.HW_CCM_PRIVATE_KEY_PATH)
}

func testAccCertificateAssociateDomains_basic_step2() string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_certificate_associate_domains" "test_with_invalid_cert" {
  domain_names = "%[1]s"
  https_switch = 1
  cert_name    = "invalid_cert"
  certificate  = "----------------------invalid content----------------------"
  private_key  = "----------------------invalid private key----------------------"
}
`, acceptance.HW_CDN_TARGET_DOMAIN_URLS)
}

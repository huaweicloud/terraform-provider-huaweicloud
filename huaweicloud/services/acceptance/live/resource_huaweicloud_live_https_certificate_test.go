package live

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceHttpsCertificate_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLiveStreamingDomainName(t)
			acceptance.TestAccPreCheckLiveTLSCert(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testResourceHttpsCertificate_basic(),
				ExpectError: regexp.MustCompile("The certificate verify failed"),
			},
		},
	})
}

func testResourceHttpsCertificate_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_live_https_certificate" "test" {
  domain_name        = "%[1]s"
  certificate_format = "PEM"
  force_redirect     = true

  tls_certificate {
    source = "user"
    certificate = file("%[2]s")
    certificate_key = file("%[3]s")
  }
}
`, acceptance.HW_LIVE_STREAMING_DOMAIN_NAME, acceptance.HW_LIVE_HTTPS_TLS_CERT_BODY_PATH, acceptance.HW_LIVE_HTTPS_TLS_CERT_KEY_PATH)
}

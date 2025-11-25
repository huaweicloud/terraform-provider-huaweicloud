package ccm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceCertificateDeploy_cdn(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Deployment operation will incur costs, please enable this environment variable before executing the test case.
			acceptance.TestAccPreCheckCCMEnableFlag(t)
			// Configure the international certificate. The SM2 certificate does not support deployment operations.
			acceptance.TestAccPreCheckCCMSSLCertificateId(t)
			// Make sure the domain name match the certificate.
			acceptance.TestAccPreCheckCdnDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceResourceCcmCertificateDeploy_cdn(),
			},
			{
				Config:      testResourceResourceCcmCertificateDeploy_cdn_pushFailed(),
				ExpectError: regexp.MustCompile("domain name not find"),
			},
		},
	})
}

func testResourceResourceCcmCertificateDeploy_cdn() string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_certificate_deploy" "test" {
  certificate_id = "%s"
  service_name   = "CDN"

  resources {
    domain_name = "%s"
  }
}
`, acceptance.HW_CCM_SSL_CERTIFICATE_ID, acceptance.HW_CDN_DOMAIN_NAME)
}

func testResourceResourceCcmCertificateDeploy_cdn_pushFailed() string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_certificate_deploy" "test" {
  certificate_id = "%s"
  service_name   = "CDN"

  resources {
    domain_name = "error.test.com"
  }
}
`, acceptance.HW_CCM_SSL_CERTIFICATE_ID)
}

func TestAccResourceCertificateDeploy_waf(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Deployment operation will incur costs, please enable this environment variable before executing the test case.
			acceptance.TestAccPreCheckCCMEnableFlag(t)
			// Configure the international certificate. The SM2 certificate does not support deployment operations.
			acceptance.TestAccPreCheckCCMSSLCertificateId(t)
			// Configure the region where the valid WAF certificate is located.
			acceptance.TestAccPrecheckCustomRegion(t)
			// Configure the WAF certificate ID that includes the HTTPS protocol.
			acceptance.TestAccPreCheckWafCertID(t)
			// Configure the instance type of the WAF certificate: premium or cloud.
			acceptance.TestAccPreCheckWafType(t)
			// Configure the EPS ID to which the WAF certificate belongs.
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceResourceCcmCertificateDeploy_waf(),
			},
		},
	})
}

func testResourceResourceCcmCertificateDeploy_waf() string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_certificate_deploy" "test" {
  certificate_id = "%[1]s"
  project_name   = "%[2]s"
  service_name   = "WAF"

  resources {
    id                    = "%[3]s"
    type                  = "%[4]s"
    enterprise_project_id = "%[5]s"
  }
}
`, acceptance.HW_CCM_SSL_CERTIFICATE_ID,
		acceptance.HW_CUSTOM_REGION_NAME,
		acceptance.HW_WAF_CERTIFICATE_ID,
		acceptance.HW_WAF_TYPE,
		acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func TestAccResourceCertificateDeploy_elb(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Deployment operation will incur costs, please enable this environment variable before executing the test case.
			acceptance.TestAccPreCheckCCMEnableFlag(t)
			// Configure the international certificate. The SM2 certificate does not support deployment operations.
			acceptance.TestAccPreCheckCCMSSLCertificateId(t)
			// Configure the region where the valid WAF certificate is located.
			acceptance.TestAccPrecheckCustomRegion(t)
			// Configure the ELB certificate ID.
			acceptance.TestAccPreCheckElbCertID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceResourceCcmCertificateDeploy_elb(),
			},
		},
	})
}

func testResourceResourceCcmCertificateDeploy_elb() string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_certificate_deploy" "test" {
  certificate_id = "%[1]s"
  project_name   = "%[2]s"
  service_name   = "ELB"

  resources {
    id = "%[3]s"
  }
}
`, acceptance.HW_CCM_SSL_CERTIFICATE_ID,
		acceptance.HW_CUSTOM_REGION_NAME,
		acceptance.HW_ELB_CERT_ID)
}

package ccm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/scm/v3/certificates"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

var (
	certificatePath     = acceptance.HW_CCM_CERTIFICATE_CONTENT_PATH
	chainPath           = acceptance.HW_CCM_CERTIFICATE_CHAIN_PATH
	keyPath             = acceptance.HW_CCM_PRIVATE_KEY_PATH
	encCertificatePath  = acceptance.HW_CCM_ENC_CERTIFICATE_PATH
	encKeyPath          = acceptance.HW_CCM_ENC_PRIVATE_KEY_PATH
	enterpriseProjectID = acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
	projectName         = acceptance.HW_CCM_CERTIFICATE_PROJECT
	projectUpdateName   = acceptance.HW_CCM_CERTIFICATE_PROJECT_UPDATED
)

func getCertificateImportResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.ScmV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCM client: %s", err)
	}
	return certificates.Get(client, state.Primary.ID).Extract()
}

// Using to test importing certificates encrypted with SM series cryptographic algorithms.
// Certificates encrypted with SM series cryptographic algorithms cannot be deployed to other cloud services.
func TestAccCertificateImport_basic(t *testing.T) {
	var (
		certInfo     certificates.CertificateEscrowInfo
		resourceName = "huaweicloud_ccm_certificate_import.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&certInfo,
		getCertificateImportResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCCMBaseCertificateImport(t)
			acceptance.TestAccPreCheckCCMEncCertificateImport(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateImport_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "UPLOAD"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", enterpriseProjectID),
					resource.TestCheckResourceAttrSet(resourceName, "push_support"),
					resource.TestCheckResourceAttrSet(resourceName, "authentifications.#"),
					resource.TestCheckResourceAttrSet(resourceName, "domain"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_count"),
					resource.TestCheckResourceAttrSet(resourceName, "not_before"),
					resource.TestCheckResourceAttrSet(resourceName, "not_after"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"certificate",
					"certificate_chain",
					"private_key",
					"enc_certificate",
					"enc_private_key",
				},
			},
		},
	})
}

func testAccCertificateImport_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_certificate_import" "test" {
  name                  = "%[1]s"
  certificate           = file("%[2]s")
  certificate_chain     = file("%[3]s")
  private_key           = file("%[4]s")
  enc_certificate       = file("%[5]s")
  enc_private_key       = file("%[6]s")
  enterprise_project_id = "%[7]s"
}`, name, certificatePath, chainPath, keyPath, encCertificatePath, encKeyPath, enterpriseProjectID)
}

// Using to test importing international standard certificates and pushing to services.
func TestAccCertificateImport_push(t *testing.T) {
	var (
		certInfo     certificates.CertificateEscrowInfo
		resourceName = "huaweicloud_ccm_certificate_import.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&certInfo,
		getCertificateImportResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCCMBaseCertificateImport(t)
			acceptance.TestAccPreCheckCCMCertificatePush(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateImport_push(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "UPLOAD"),
				),
			},
			{
				Config: testAccCertificateImport_push_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "UPLOAD"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate", "certificate_chain", "private_key", "target"},
			},
		},
	})
}

func testAccCertificateImport_push(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_certificate_import" "test" {
  name              = "%[1]s"
  certificate       = file("%[2]s")
  certificate_chain = file("%[3]s")
  private_key       = file("%[4]s")

  target {
    project = ["%[5]s"]
    service = "ELB"
  }

  target {
    service = "CDN"
  }
}`, name, certificatePath, chainPath, keyPath, projectName)
}

func testAccCertificateImport_push_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_certificate_import" "test" {
  name              = "%[1]s"
  certificate       = file("%[2]s")
  certificate_chain = file("%[3]s")
  private_key       = file("%[4]s")

  target {
    project = ["%[5]s", "%[6]s"]
    service = "ELB"
  }

  target {
    service = "CDN"
  }
}`, name, certificatePath, chainPath, keyPath, projectName, projectUpdateName)
}

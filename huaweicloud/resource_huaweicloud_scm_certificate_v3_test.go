package huaweicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/scm/v3/certificates"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccScmCertificationV3_basic(t *testing.T) {
	var certInfo certificates.CertificateEscrowInfo
	resourceName := "huaweicloud_scm_certificate.certificate_1"

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	rUpdateName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckScm(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckScmV3CertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccScmCertificateV3_config(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScmV3CertificateExists(resourceName, &certInfo),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "UPLOAD"),
				),
			},
			{
				Config: testAccScmCertificateV3_config(rUpdateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScmV3CertificateExists(resourceName, &certInfo),
					resource.TestCheckResourceAttr(resourceName, "name", rUpdateName),
					resource.TestCheckResourceAttr(resourceName, "status", "UPLOAD"),
				),
			},
		},
	})
}

func testAccCheckScmV3CertificateExists(n string, c *certificates.CertificateEscrowInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		scmClient, err := config.ScmV3Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud scm client: %s", err)
		}

		found, err := certificates.Get(scmClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Certificate not found")
		}

		*c = *found

		return nil
	}
}

func testAccCheckScmV3CertificateDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	scmClient, err := config.ScmV3Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud scm client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_scm_certificate" {
			continue
		}

		_, err := certificates.Get(scmClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("SSL Certificate still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccScmCertificateV3_config(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_scm_certificate" "certificate_1" {
  name              = "%s"
  certificate       = file("%s")
  certificate_chain = file("%s")
  private_key       = file("%s")
}`, name, HW_CERTIFICATE_KEY_PATH, HW_CERTIFICATE_CHAIN_PATH, HW_CERTIFICATE_PRIVATE_KEY_PATH)
}

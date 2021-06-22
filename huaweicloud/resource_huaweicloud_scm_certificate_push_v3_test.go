package huaweicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/scm/v3/certificates"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccScmCertificationPushV3_basic(t *testing.T) {
	certResourceName := "huaweicloud_scm_certificate.certificate_1"
	pushResourceName := "huaweicloud_scm_certificate_push.push_1"
	targetService := "Enhance_ELB"
	targetProject := "ap-southeast-1"

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckScm(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckScmV3CertificatePushDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccScmCertificatePushV3_config(rName, targetService, targetProject),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScmV3CertificatePushExists(certResourceName, pushResourceName, targetService, targetProject),
					resource.TestCheckResourceAttrSet(pushResourceName, "id"),
					resource.TestCheckResourceAttrSet(pushResourceName, "certificate_id"),
				),
			},
		},
	})
}

func testAccCheckScmV3CertificatePushExists(certResourceName string, pushResourceName string, targetService string, targetProject string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("terraform.State: %#v", s)
		// Get the resource info by certificateResorceName
		certRs, ok := s.RootModule().Resources[certResourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", certResourceName)
		}

		if certRs.Primary.ID == "" {
			return fmt.Errorf("No id is set for the certificate resource: %s", certResourceName)
		}

		// Get the resource info by certificatePushResorceName
		pushRs, ok := s.RootModule().Resources[pushResourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", pushResourceName)
		}
		if pushRs.Primary.ID == "" {
			return fmt.Errorf("No id is set for the certificate push resource: %s", pushResourceName)
		}

		certId := certRs.Primary.Attributes["id"]
		pushCertId := pushRs.Primary.Attributes["id"]

		pushOpts := certificates.PushOpts{
			TargetProject: targetProject,
			TargetService: targetService,
		}
		calcId := generateCertPushId(certId, pushOpts)
		if strings.Compare(calcId, pushCertId) != 0 {
			return fmt.Errorf("Error the ID[%s] from the push service is inconsistent with "+
				"the calculated ID[%s]. ", pushCertId, calcId)
		}
		return nil
	}
}

// Push service does not support deletion.
// After the certificate is deleted, the deletion is complete.
func testAccCheckScmV3CertificatePushDestroy(s *terraform.State) error {
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

func testAccScmCertificatePushV3_config(resourceName string, targetService string, targetProject string) string {

	certResource := fmt.Sprintf(`
resource "huaweicloud_scm_certificate" "certificate_1" {
  name              = "%s"
  certificate       = file("%s")
  certificate_chain = file("%s")
  private_key       = file("%s")
}`, resourceName, HW_CERTIFICATE_KEY_PATH, HW_CERTIFICATE_CHAIN_PATH, HW_CERTIFICATE_PRIVATE_KEY_PATH)

	pushResource := fmt.Sprintf(`
resource "huaweicloud_scm_certificate_push" "push_1" {
  certificate_id = huaweicloud_scm_certificate.certificate_1.id
  target_service = "%s"
  target_project = "%s"
}
`, targetService, targetProject)

	return certResource + pushResource
}

package huaweicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/scm/v3/certificates"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
				Config: testAccScmCertificateV3_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScmV3CertificateExists(resourceName, &certInfo),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "UPLOAD"),
				),
			},
			{
				Config: testAccScmCertificateV3_basic(rUpdateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScmV3CertificateExists(resourceName, &certInfo),
					resource.TestCheckResourceAttr(resourceName, "name", rUpdateName),
					resource.TestCheckResourceAttr(resourceName, "status", "UPLOAD"),
				),
			},
		},
	})
}

func TestAccScmCertificationV3_push(t *testing.T) {
	var certInfo certificates.CertificateEscrowInfo
	resourceName := "huaweicloud_scm_certificate.certificate_2"

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	service := HW_CERTIFICATE_SERVICE
	defaultProject := HW_CERTIFICATE_PROJECT
	newProject := HW_CERTIFICATE_PROJECT_UPDATED

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckScm(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckScmV3CertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccScmCertificateV3_push(rName, defaultProject, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScmV3CertificateExists(resourceName, &certInfo),
					testAccCheckScmV3CertificatePushExists(resourceName, service, defaultProject),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "UPLOAD"),
				),
			},
			{
				Config: testAccScmCertificateV3_push(rName, newProject, service),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScmV3CertificateExists(resourceName, &certInfo),
					testAccCheckScmV3CertificatePushExists(resourceName, service, newProject),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "UPLOAD"),
				),
			},
		},
	})
}

func TestAccScmCertificationV3_batchPush(t *testing.T) {
	var certInfo certificates.CertificateEscrowInfo
	resourceName := "huaweicloud_scm_certificate.certificate_3"

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	service := HW_CERTIFICATE_SERVICE
	defaultProject := HW_CERTIFICATE_PROJECT

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckScm(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckScmV3CertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccScmCertificateV3_batchPush(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScmV3CertificateExists(resourceName, &certInfo),
					testAccCheckScmV3CertificatePushExists(resourceName, service, defaultProject),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
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

func testAccCheckScmV3CertificatePushExists(
	certResourceName string, service string, project string) resource.TestCheckFunc {
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

		stateService := certRs.Primary.Attributes["target.0.service"]
		stateProject := certRs.Primary.Attributes["target.0.project.0"]

		if strings.Compare(service, stateService) != 0 ||
			strings.Compare(project, stateProject) != 0 {
			return fmt.Errorf("Push certificate failed! service: %s, project: %s",
				service, project)
		}

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

func testAccScmCertificateV3_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_scm_certificate" "certificate_1" {
  name              = "%s"
  certificate       = file("%s")
  certificate_chain = file("%s")
  private_key       = file("%s")
}`, name, HW_CERTIFICATE_KEY_PATH, HW_CERTIFICATE_CHAIN_PATH, HW_CERTIFICATE_PRIVATE_KEY_PATH)
}

func testAccScmCertificateV3_push(name string, project string, service string) string {
	return fmt.Sprintf(`
resource "huaweicloud_scm_certificate" "certificate_2" {
  name              = "%s"
  certificate       = file("%s")
  certificate_chain = file("%s")
  private_key       = file("%s")

  target {
    project  = ["%s"]
    service  = "%s"
  }
}`, name, HW_CERTIFICATE_KEY_PATH, HW_CERTIFICATE_CHAIN_PATH, HW_CERTIFICATE_PRIVATE_KEY_PATH,
		project, service)
}

func testAccScmCertificateV3_batchPush(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_scm_certificate" "certificate_3" {
  name              = "%s"
  certificate       = file("%s")
  certificate_chain = file("%s")
  private_key       = file("%s")

  target {
    project  = ["%s", "%s"]
    service  = "%s"
  }
}`, name, HW_CERTIFICATE_KEY_PATH, HW_CERTIFICATE_CHAIN_PATH, HW_CERTIFICATE_PRIVATE_KEY_PATH,
		HW_CERTIFICATE_PROJECT, HW_CERTIFICATE_PROJECT_UPDATED, HW_CERTIFICATE_SERVICE)
}

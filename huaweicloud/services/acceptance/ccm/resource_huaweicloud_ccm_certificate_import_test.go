package ccm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/scm/v3/certificates"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getCertificateImportResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.ScmV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCM client: %s", err)
	}
	return certificates.Get(client, state.Primary.ID).Extract()
}

func TestAccCertificateImport_basic(t *testing.T) {
	var certInfo certificates.CertificateEscrowInfo
	resourceName := "huaweicloud_ccm_certificate_import.certificate_1"

	rName := acceptance.RandomAccResourceNameWithDash()
	rUpdateName := acceptance.RandomAccResourceNameWithDash()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&certInfo,
		getCertificateImportResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCCMCertificateImport(t)
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
				),
			},
			{
				Config: testAccCertificateImport_basic(rUpdateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rUpdateName),
					resource.TestCheckResourceAttr(resourceName, "status", "UPLOAD"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate", "certificate_chain", "private_key"},
			},
		},
	})
}

func TestAccCertificateImport_push(t *testing.T) {
	var certInfo certificates.CertificateEscrowInfo
	resourceName := "huaweicloud_ccm_certificate_import.certificate_2"

	rName := acceptance.RandomAccResourceNameWithDash()
	service := acceptance.HW_CERTIFICATE_SERVICE
	defaultProject := acceptance.HW_CERTIFICATE_PROJECT
	newProject := acceptance.HW_CERTIFICATE_PROJECT_UPDATED

	rc := acceptance.InitResourceCheck(
		resourceName,
		&certInfo,
		getCertificateImportResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCCMCertificateImport(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateImport_push(rName, defaultProject, service),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					testCertificateImportPushExists(resourceName, service, defaultProject),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "UPLOAD"),
				),
			},
			{
				Config: testAccCertificateImport_push(rName, newProject, service),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					testCertificateImportPushExists(resourceName, service, newProject),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "UPLOAD"),
				),
			},
		},
	})
}

func TestAccCertificateImport_batchPush(t *testing.T) {
	var certInfo certificates.CertificateEscrowInfo
	resourceName := "huaweicloud_ccm_certificate_import.certificate_3"

	rName := acceptance.RandomAccResourceNameWithDash()
	service := acceptance.HW_CERTIFICATE_SERVICE
	defaultProject := acceptance.HW_CERTIFICATE_PROJECT

	rc := acceptance.InitResourceCheck(
		resourceName,
		&certInfo,
		getCertificateImportResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCCMCertificateImport(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateImport_batchPush(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					testCertificateImportPushExists(resourceName, service, defaultProject),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "UPLOAD"),
				),
			},
		},
	})
}

func testCertificateImportPushExists(
	certResourceName string, service string, project string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Get the resource info by certificateResorceName
		certRs, ok := s.RootModule().Resources[certResourceName]
		if !ok {
			return fmt.Errorf("not found: %s", certResourceName)
		}

		if certRs.Primary.ID == "" {
			return fmt.Errorf("no ID is set for the certificate resource: %s", certResourceName)
		}

		stateService := certRs.Primary.Attributes["target.0.service"]
		stateProject := certRs.Primary.Attributes["target.0.project.0"]
		if strings.Compare(service, stateService) != 0 ||
			strings.Compare(project, stateProject) != 0 {
			return fmt.Errorf("push certificate failed! service: %s, project: %s",
				service, project)
		}
		return nil
	}
}

func testAccCertificateImport_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_certificate_import" "certificate_1" {
  name              = "%s"
  certificate       = file("%s")
  certificate_chain = file("%s")
  private_key       = file("%s")
}`, name, acceptance.HW_CERTIFICATE_KEY_PATH, acceptance.HW_CERTIFICATE_CHAIN_PATH,
		acceptance.HW_CERTIFICATE_PRIVATE_KEY_PATH)
}

func testAccCertificateImport_push(name string, project string, service string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_certificate_import" "certificate_2" {
  name              = "%s"
  certificate       = file("%s")
  certificate_chain = file("%s")
  private_key       = file("%s")

  target {
    project = ["%s"]
    service = "%s"
  }
}`, name, acceptance.HW_CERTIFICATE_KEY_PATH, acceptance.HW_CERTIFICATE_CHAIN_PATH,
		acceptance.HW_CERTIFICATE_PRIVATE_KEY_PATH, project, service)
}

func testAccCertificateImport_batchPush(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_certificate_import" "certificate_3" {
  name              = "%s"
  certificate       = file("%s")
  certificate_chain = file("%s")
  private_key       = file("%s")

  target {
    project = ["%s", "%s"]
    service = "%s"
  }
}`, name, acceptance.HW_CERTIFICATE_KEY_PATH, acceptance.HW_CERTIFICATE_CHAIN_PATH,
		acceptance.HW_CERTIFICATE_PRIVATE_KEY_PATH, acceptance.HW_CERTIFICATE_PROJECT,
		acceptance.HW_CERTIFICATE_PROJECT_UPDATED, acceptance.HW_CERTIFICATE_SERVICE)
}

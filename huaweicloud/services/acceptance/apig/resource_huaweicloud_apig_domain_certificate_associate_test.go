package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
)

func getDomainCertificateAssociateFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("apig", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG client: %s", err)
	}

	return apig.GetDomainAssociatedCertificateByCertificateId(
		client,
		state.Primary.Attributes["instance_id"],
		state.Primary.Attributes["group_id"],
		state.Primary.Attributes["domain_id"],
		state.Primary.Attributes["certificate_id"],
	)
}

func TestAccDomainCertificateAssociate_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		associatedCertificates interface{}

		resourceName = "huaweicloud_apig_domain_certificate_associate.test"
		rc           = acceptance.InitResourceCheck(resourceName, &associatedCertificates, getDomainCertificateAssociateFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
			acceptance.TestAccPreCheckCertificateBase(t)
			acceptance.TestAccPreCheckCertificateRootCA(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDomainCertificateAssociate_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttrPair(resourceName, "group_id", "huaweicloud_apig_group.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_id"),
					resource.TestCheckResourceAttrPair(resourceName, "certificate_id", "huaweicloud_apig_certificate.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "verified_client_certificate_enabled", "true"),
				),
			},
			{
				Config: testAccDomainCertificateAssociate_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "verified_client_certificate_enabled", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDomainCertificateAssociate_base(name string) string {
	return fmt.Sprintf(` 
resource "huaweicloud_apig_certificate" "test" {
  instance_id     = "%[1]s"
  type            = "instance"
  name            = "%[2]s"
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
  instance_id = "%[1]s"
  group_id    = huaweicloud_apig_group.test.id
  url_domain  = local.domain_name
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID,
		name,
		acceptance.HW_CERTIFICATE_CONTENT,
		acceptance.HW_CERTIFICATE_PRIVATE_KEY,
		acceptance.HW_CERTIFICATE_ROOT_CA,
	)
}

func testAccDomainCertificateAssociate_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_domain_certificate_associate" "test" {
  instance_id    = "%[2]s"
  group_id       = huaweicloud_apig_group.test.id
  domain_id      = huaweicloud_apig_group_domain_associate.test.domain_id
  certificate_id = huaweicloud_apig_certificate.test.id

  verified_client_certificate_enabled = true
}
`, testAccDomainCertificateAssociate_base(name), acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}

func testAccDomainCertificateAssociate_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_domain_certificate_associate" "test" {
  instance_id    = "%[2]s"
  group_id       = huaweicloud_apig_group.test.id
  domain_id      = huaweicloud_apig_group_domain_associate.test.domain_id
  certificate_id = huaweicloud_apig_certificate.test.id
}
`, testAccDomainCertificateAssociate_base(name), acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}

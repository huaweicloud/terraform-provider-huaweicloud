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

func getDomainBatchCertificatesAssociateFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("apig", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG client: %s", err)
	}

	return apig.GetDomainAssociatedCertificatesByDomainId(
		client,
		state.Primary.Attributes["instance_id"],
		state.Primary.Attributes["group_id"],
		state.Primary.Attributes["domain_id"],
	)
}

func TestAccDomainBatchCertificatesAssociate_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		associatedCertificates interface{}
		resourceName           = "huaweicloud_apig_domain_batch_certificates_associate.test"
		rc                     = acceptance.InitResourceCheck(resourceName, &associatedCertificates, getDomainBatchCertificatesAssociateFunc)
		withoutRootCA          = "huaweicloud_apig_domain_batch_certificates_associate.without_root_ca"
		rcWithoutRootCA        = acceptance.InitResourceCheck(withoutRootCA, &associatedCertificates, getDomainBatchCertificatesAssociateFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
			acceptance.TestAccPreCheckCertificateBase(t)
			acceptance.TestAccPreCheckCertificateRootCA(t)
			acceptance.TestAccPreCheckCertificateWithoutRootCA(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDomainBatchCertificatesAssociate_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttrPair(resourceName, "group_id", "huaweicloud_apig_group.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_id"),
					resource.TestCheckResourceAttr(resourceName, "certificate_ids.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "certificate_ids.0", "huaweicloud_apig_certificate.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "verified_client_certificate_enabled", "true"),
					rcWithoutRootCA.CheckResourceExists(),
					resource.TestCheckResourceAttr(withoutRootCA, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttrPair(withoutRootCA, "group_id", "huaweicloud_apig_group.test", "id"),
					resource.TestCheckResourceAttrSet(withoutRootCA, "domain_id"),
					resource.TestCheckResourceAttr(withoutRootCA, "certificate_ids.#", "1"),
					resource.TestCheckResourceAttrPair(withoutRootCA, "certificate_ids.0", "huaweicloud_apig_certificate.test.1", "id"),
					resource.TestCheckResourceAttr(withoutRootCA, "verified_client_certificate_enabled", "false"),
				),
			},
			{
				Config: testAccDomainBatchCertificatesAssociate_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "certificate_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "verified_client_certificate_enabled", "false"),
					rcWithoutRootCA.CheckResourceExists(),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      withoutRootCA,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDomainBatchCertificatesAssociate_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_certificate" "test" {
  count = 2

  instance_id     = "%[1]s"
  type            = "instance"
  name            = "%[2]s${count.index}"
  content         = count.index == 0 ? "%[3]s" : "%[6]s"
  private_key     = count.index == 0 ? "%[4]s" : "%[7]s"
  trusted_root_ca = count.index == 0 ? "%[5]s" : null
}

data "huaweicloud_apig_instance_ssl_certificates" "test" {
  count = 2

  instance_id = "%[1]s"
  name        = huaweicloud_apig_certificate.test[count.index].name

  depends_on = [huaweicloud_apig_certificate.test]
}

locals {
  domain_name = try([for v in data.huaweicloud_apig_instance_ssl_certificates.test[0].certificates : v.common_name
  if v.id == huaweicloud_apig_certificate.test[0].id][0], null)
  without_root_ca_domain_name = try([for v in data.huaweicloud_apig_instance_ssl_certificates.test[1].certificates : v.common_name
  if v.id == huaweicloud_apig_certificate.test[1].id][0], null)

  domain_names = [local.domain_name, local.without_root_ca_domain_name]
}

resource "huaweicloud_apig_group" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
}

resource "huaweicloud_apig_group_domain_associate" "test" {
  count = 2

  instance_id = "%[1]s"
  group_id    = huaweicloud_apig_group.test.id
  url_domain  = local.domain_names[count.index]
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID,
		name,
		acceptance.HW_CERTIFICATE_CONTENT,
		acceptance.HW_CERTIFICATE_PRIVATE_KEY,
		acceptance.HW_CERTIFICATE_ROOT_CA,
		acceptance.HW_NEW_CERTIFICATE_CONTENT,
		acceptance.HW_NEW_CERTIFICATE_PRIVATE_KEY,
	)
}

func testAccDomainBatchCertificatesAssociate_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_domain_batch_certificates_associate" "test" {
  instance_id     = "%[2]s"
  group_id        = huaweicloud_apig_group.test.id
  domain_id       = huaweicloud_apig_group_domain_associate.test[0].domain_id
  certificate_ids = [huaweicloud_apig_certificate.test[0].id]

  verified_client_certificate_enabled = true
}

resource "huaweicloud_apig_domain_batch_certificates_associate" "without_root_ca" {
  instance_id     = "%[2]s"
  group_id        = huaweicloud_apig_group.test.id
  domain_id       = huaweicloud_apig_group_domain_associate.test[1].domain_id
  certificate_ids = [huaweicloud_apig_certificate.test[1].id]
}
`, testAccDomainBatchCertificatesAssociate_base(name), acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}

func testAccDomainBatchCertificatesAssociate_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_domain_batch_certificates_associate" "test" {
  instance_id     = "%[2]s"
  group_id        = huaweicloud_apig_group.test.id
  domain_id       = huaweicloud_apig_group_domain_associate.test[0].domain_id
  certificate_ids = [huaweicloud_apig_certificate.test[0].id]
}

resource "huaweicloud_apig_domain_batch_certificates_associate" "without_root_ca" {
  instance_id     = "%[2]s"
  group_id        = huaweicloud_apig_group.test.id
  domain_id       = huaweicloud_apig_group_domain_associate.test[1].domain_id
  certificate_ids = [huaweicloud_apig_certificate.test[1].id]
}
`, testAccDomainBatchCertificatesAssociate_base(name), acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}

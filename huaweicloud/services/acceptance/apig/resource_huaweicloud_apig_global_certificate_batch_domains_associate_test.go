package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getGlobalCertificateBatchDomainsAssociateFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("apig", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG client: %s", err)
	}

	_, verifyDisabledDomains, err := apig.GetGlobalCertificateAssociatedDomains(
		client,
		state.Primary.ID,
		utils.ParseStateAttributeToListWithSeparator(state.Primary.Attributes["verify_disabled_domain_names_origin"], ","),
	)

	return verifyDisabledDomains, err
}

// Before running this test, please ensure that the certificate has at least four domains.
func TestAccGlobalCertificateBatchDomainsAssociate_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		associatedDomains interface{}

		rNamePart1 = "huaweicloud_apig_global_certificate_batch_domains_associate.part1"
		rcPart1    = acceptance.InitResourceCheck(rNamePart1, &associatedDomains, getGlobalCertificateBatchDomainsAssociateFunc)
		rNamePart2 = "huaweicloud_apig_global_certificate_batch_domains_associate.part2"
		rcPart2    = acceptance.InitResourceCheck(rNamePart2, &associatedDomains, getGlobalCertificateBatchDomainsAssociateFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
			acceptance.TestAccPreCheckCertificateBase(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcPart1.CheckResourceDestroy(),
			rcPart2.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccGlobalCertificateBatchDomainsAssociate_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rcPart1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNamePart1, "verify_disabled_domain_names.#", "2"),
					resource.TestCheckResourceAttr(rNamePart1, "verify_disabled_domain_names_origin.#", "2"),
					rcPart1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNamePart2, "verify_disabled_domain_names.#", "1"),
					resource.TestCheckResourceAttr(rNamePart2, "verify_disabled_domain_names_origin.#", "1"),
				),
			},
			{
				Config: testAccGlobalCertificateBatchDomainsAssociate_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rcPart1.CheckResourceExists(),
					// After resources refreshed, the verify_disabled_domain_names will be overridden as all domains under the same
					// certificate.
					resource.TestCheckResourceAttr(rNamePart1, "verify_disabled_domain_names.#", "3"),
					resource.TestCheckResourceAttr(rNamePart1, "verify_disabled_domain_names_origin.#", "2"),
					rcPart1.CheckResourceExists(),
					// After resources refreshed, the verify_disabled_domain_names will be overridden as all domains under the same
					// certificate.
					resource.TestCheckResourceAttr(rNamePart2, "verify_disabled_domain_names.#", "3"),
					resource.TestCheckResourceAttr(rNamePart2, "verify_disabled_domain_names_origin.#", "1"),
				),
			},
			{
				Config: testAccGlobalCertificateBatchDomainsAssociate_basic_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rcPart1.CheckResourceExists(),
					// When multiple resources are used to manage the same function, verify_disabled_domain_names will store the results
					// modified by other resources, resulting in verify_disabled_domain_names displaying all binding results except for the
					// first change.
					resource.TestMatchResourceAttr(rNamePart1, "verify_disabled_domain_names.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(rNamePart1, "verify_disabled_domain_names_origin.#", "1"),
					rcPart2.CheckResourceExists(),
					resource.TestMatchResourceAttr(rNamePart2, "verify_disabled_domain_names.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(rNamePart2, "verify_disabled_domain_names_origin.#", "2"),
				),
			},
			{
				ResourceName:      rNamePart1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"verify_disabled_domain_names_origin",
				},
			},
			{
				// After resource part1 is imported, then verify_disabled_domain_names will be overridden as all domains under the same
				// certificate.
				Config: testAccGlobalCertificateBatchDomainsAssociate_basic_step4(name),
				Check: resource.ComposeTestCheckFunc(
					rcPart1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNamePart1, "verify_disabled_domain_names.#", "3"),
					resource.TestCheckResourceAttr(rNamePart1, "verify_disabled_domain_names_origin.#", "1"),
					rcPart2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNamePart2, "verify_disabled_domain_names.#", "3"),
					resource.TestCheckResourceAttr(rNamePart2, "verify_disabled_domain_names_origin.#", "2"),
					resource.TestMatchResourceAttr(rNamePart1, "domains.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(rNamePart2, "domains.0.id"),
					resource.TestCheckResourceAttrSet(rNamePart2, "domains.0.url_domain"),
					resource.TestCheckResourceAttrSet(rNamePart2, "domains.0.instance_id"),
					resource.TestCheckResourceAttrSet(rNamePart2, "domains.0.status"),
					resource.TestCheckResourceAttrSet(rNamePart2, "domains.0.min_ssl_version"),
					resource.TestCheckResourceAttrSet(rNamePart2, "domains.0.api_group_id"),
					resource.TestCheckResourceAttrSet(rNamePart2, "domains.0.api_group_name"),
				),
			},
		},
	})
}

func testAccGlobalCertificateBatchDomainsAssociate_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_certificate" "test" {
  name        = "%[2]s"
  content     = "%[3]s"
  private_key = "%[4]s"
}

data "huaweicloud_apig_instance_ssl_certificates" "test" {
  instance_id = "%[1]s"
  name        = huaweicloud_apig_certificate.test.name

  depends_on = [huaweicloud_apig_certificate.test]
}

locals {
  domain_names = try([for v in data.huaweicloud_apig_instance_ssl_certificates.test.certificates : v.san
  if v.id == huaweicloud_apig_certificate.test.id][0], [])
}

resource "huaweicloud_apig_group" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
}

resource "huaweicloud_apig_group_domain_associate" "test" {
  count = 4

  instance_id = "%[1]s"
  group_id    = huaweicloud_apig_group.test.id
  url_domain  = local.domain_names[count.index]

  depends_on = [huaweicloud_apig_certificate.test]
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID,
		name, acceptance.HW_CERTIFICATE_CONTENT,
		acceptance.HW_CERTIFICATE_PRIVATE_KEY)
}

func testAccGlobalCertificateBatchDomainsAssociate_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_global_certificate_batch_domains_associate" "part1" {
  certificate_id               = huaweicloud_apig_certificate.test.id
  verify_disabled_domain_names = slice(local.domain_names, 0, 2)

  depends_on = [huaweicloud_apig_group_domain_associate.test]
}

resource "huaweicloud_apig_global_certificate_batch_domains_associate" "part2" {
  certificate_id               = huaweicloud_apig_certificate.test.id
  verify_disabled_domain_names = slice(local.domain_names, 3, 4)

  depends_on = [huaweicloud_apig_group_domain_associate.test]
}
`, testAccGlobalCertificateBatchDomainsAssociate_base(name))
}

func testAccGlobalCertificateBatchDomainsAssociate_basic_step2(name string) string {
	// Refresh the verify_disabled_domain_names for all bound domains.
	return testAccGlobalCertificateBatchDomainsAssociate_basic_step1(name)
}

func testAccGlobalCertificateBatchDomainsAssociate_basic_step3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_global_certificate_batch_domains_associate" "part1" {
  certificate_id               = huaweicloud_apig_certificate.test.id
  verify_disabled_domain_names = slice(local.domain_names, 0, 1)
}

resource "huaweicloud_apig_global_certificate_batch_domains_associate" "part2" {
  certificate_id               = huaweicloud_apig_certificate.test.id
  verify_disabled_domain_names = slice(local.domain_names, 2, 4)
}
`, testAccGlobalCertificateBatchDomainsAssociate_base(name))
}

// After importing the resources in part1, `verify_disabled_domain_names` will contain all
// the domain names bound to the certificate.
func testAccGlobalCertificateBatchDomainsAssociate_basic_step4(name string) string {
	return testAccGlobalCertificateBatchDomainsAssociate_basic_step3(name)
}

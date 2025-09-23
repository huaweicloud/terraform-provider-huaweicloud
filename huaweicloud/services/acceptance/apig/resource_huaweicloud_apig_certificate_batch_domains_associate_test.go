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

func getCertificateAssociatedDomainsFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("apig", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG client: %s", err)
	}

	enabledDomains, disabledDomains, err := apig.GetCertificateAssociatedDomainsByDomain(client,
		state.Primary.Attributes["certificate_id"],
		utils.ParseStateAttributeToListWithSeparator(state.Primary.Attributes["verify_enabled_domain_names_origin"], ","),
		utils.ParseStateAttributeToListWithSeparator(state.Primary.Attributes["verify_disabled_domain_names_origin"], ","),
	)
	if err != nil {
		return nil, err
	}

	return []interface{}{enabledDomains, disabledDomains}, err
}

// Before running this test, please ensure that the certificate has at least five domains.
func TestAccCertificateBatchDomainsAssociate_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		associatedDomains interface{}

		rNamePart1 = "huaweicloud_apig_certificate_batch_domains_associate.part1"
		rcPart1    = acceptance.InitResourceCheck(rNamePart1, &associatedDomains, getCertificateAssociatedDomainsFunc)
		rNamePart2 = "huaweicloud_apig_certificate_batch_domains_associate.part2"
		rcPart2    = acceptance.InitResourceCheck(rNamePart2, &associatedDomains, getCertificateAssociatedDomainsFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
			acceptance.TestAccPreCheckCertificateBase(t)
			acceptance.TestAccPreCheckCertificateRootCA(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rcPart1.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateBatchDomainsAssociate_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rcPart1.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rNamePart1, "certificate_id", "huaweicloud_apig_certificate.test", "id"),
					resource.TestCheckResourceAttr(rNamePart1, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttr(rNamePart1, "verify_enabled_domain_names.#", "2"),
					resource.TestCheckResourceAttr(rNamePart1, "verify_enabled_domain_names_origin.#", "2"),
					resource.TestCheckNoResourceAttr(rNamePart1, "verify_disabled_domain_names_origin.#"),
					rcPart2.CheckResourceExists(),
					resource.TestCheckNoResourceAttr(rNamePart2, "verify_enabled_domain_names_origin.#"),
					resource.TestCheckResourceAttr(rNamePart2, "verify_disabled_domain_names.#", "1"),
					resource.TestCheckResourceAttr(rNamePart2, "verify_disabled_domain_names_origin.#", "1"),
				),
			},
			// Refresh the `verify_disabled_domain_names` and `verify_enabled_domain_names` of the part1 and part2 for all bound domains.
			{
				Config: testAccCertificateBatchDomainsAssociate_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rcPart1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNamePart1, "verify_enabled_domain_names.#", "2"),
					resource.TestCheckResourceAttr(rNamePart1, "verify_enabled_domain_names_origin.#", "2"),
					resource.TestCheckResourceAttr(rNamePart1, "verify_disabled_domain_names.#", "1"),
					resource.TestCheckNoResourceAttr(rNamePart1, "verify_disabled_domain_names_origin.#"),
					rcPart2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNamePart2, "verify_enabled_domain_names.#", "2"),
					resource.TestCheckNoResourceAttr(rNamePart2, "verify_enabled_domain_names_origin.#"),
					resource.TestCheckResourceAttr(rNamePart2, "verify_disabled_domain_names.#", "1"),
					resource.TestCheckResourceAttr(rNamePart2, "verify_disabled_domain_names_origin.#", "1"),
				),
			},
			// When using multiple resources to manage domain bindings for the same certificate, the values ​​of `verify_enabled_domain_names`
			// and `verify_disabled_domain_names_origin` will be affected by the other resources.
			// Since the order in which the interfaces of part1 and part2 are sent is uncertain, the length of
			// verify_enabled_domain_names and verify_disabled_domain_names cannot be determined either.
			{
				Config: testAccCertificateBatchDomainsAssociate_basic_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rcPart1.CheckResourceExists(),
					resource.TestMatchResourceAttr(rNamePart1, "verify_enabled_domain_names.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(rNamePart1, "verify_enabled_domain_names_origin.#", "1"),
					resource.TestMatchResourceAttr(rNamePart1, "verify_disabled_domain_names.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(rNamePart1, "verify_disabled_domain_names_origin.#", "0"),
					rcPart2.CheckResourceExists(),
					resource.TestMatchResourceAttr(rNamePart1, "verify_enabled_domain_names.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(rNamePart2, "verify_enabled_domain_names_origin.#", "2"),
					resource.TestMatchResourceAttr(rNamePart1, "verify_disabled_domain_names.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(rNamePart2, "verify_disabled_domain_names_origin.#", "1"),
				),
			},
			{
				ResourceName:      rNamePart1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"verify_enabled_domain_names_origin",
					"verify_disabled_domain_names_origin",
				},
			},
			// After importing the resources in part1, `verify_enabled_domain_names` and `verify_disabled_domain_names` will contain all
			// the domain names bound to the certificate.
			{
				Config: testAccCertificateBatchDomainsAssociate_basic_step4(name),
				Check: resource.ComposeTestCheckFunc(
					rcPart1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNamePart1, "verify_enabled_domain_names.#", "2"),
					resource.TestCheckResourceAttr(rNamePart1, "verify_enabled_domain_names_origin.#", "1"),
					resource.TestCheckResourceAttr(rNamePart1, "verify_disabled_domain_names.#", "1"),
					resource.TestCheckResourceAttr(rNamePart1, "verify_disabled_domain_names_origin.#", "0"),
					rcPart2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNamePart2, "verify_enabled_domain_names.#", "2"),
					resource.TestCheckResourceAttr(rNamePart2, "verify_enabled_domain_names_origin.#", "2"),
					resource.TestCheckResourceAttr(rNamePart2, "verify_disabled_domain_names.#", "1"),
					resource.TestCheckResourceAttr(rNamePart2, "verify_disabled_domain_names_origin.#", "1"),
				),
			},
		},
	})
}

func testAccCertificateBatchDomainsAssociate_base(name string) string {
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
  domain_names = try([for v in data.huaweicloud_apig_instance_ssl_certificates.test.certificates : v.san
  if v.id == huaweicloud_apig_certificate.test.id][0], null)
}

resource "huaweicloud_apig_group" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
}

resource "huaweicloud_apig_group_domain_associate" "test" {
  count = 5

  instance_id = "%[1]s"
  group_id    = huaweicloud_apig_group.test.id
  url_domain  = local.domain_names[count.index]

  depends_on = [huaweicloud_apig_certificate.test]
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID,
		name, acceptance.HW_CERTIFICATE_CONTENT,
		acceptance.HW_CERTIFICATE_PRIVATE_KEY,
		acceptance.HW_CERTIFICATE_ROOT_CA)
}

func testAccCertificateBatchDomainsAssociate_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_certificate_batch_domains_associate" "part1" {
  certificate_id = huaweicloud_apig_certificate.test.id
  instance_id    = "%[2]s"

  verify_enabled_domain_names = slice(local.domain_names, 0, 2)

  depends_on = [huaweicloud_apig_group_domain_associate.test]
}

resource "huaweicloud_apig_certificate_batch_domains_associate" "part2" {
  certificate_id = huaweicloud_apig_certificate.test.id
  instance_id    = "%[2]s"

  verify_disabled_domain_names = slice(local.domain_names, 3, 4)

  depends_on = [huaweicloud_apig_group_domain_associate.test]
}
`, testAccCertificateBatchDomainsAssociate_base(name), acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}

func testAccCertificateBatchDomainsAssociate_basic_step2(name string) string {
	// Refresh the `verify_enabled_domain_names` and `verify_disabled_domain_names` for all bound domains.
	return testAccCertificateBatchDomainsAssociate_basic_step1(name)
}

func testAccCertificateBatchDomainsAssociate_basic_step3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_certificate_batch_domains_associate" "part1" {
  certificate_id = huaweicloud_apig_certificate.test.id
  instance_id    = "%[2]s"

  # Bind the first domain and delete the second domain.
  verify_enabled_domain_names  = slice(local.domain_names, 0, 1)
}

resource "huaweicloud_apig_certificate_batch_domains_associate" "part2" {
  certificate_id = huaweicloud_apig_certificate.test.id
  instance_id    = "%[2]s"

  # Bind the first two domains. Combined with the operations in part 1, after the update is completed,
  # the length of verify_enabled_domain_names is 2.
  verify_enabled_domain_names  = slice(local.domain_names, 1, 3)
  verify_disabled_domain_names = slice(local.domain_names, 4, 5)
}
`, testAccCertificateBatchDomainsAssociate_base(name), acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}

func testAccCertificateBatchDomainsAssociate_basic_step4(name string) string {
	// After importing the resources in part1, `verify_enabled_domain_names` and `verify_disabled_domain_names` will contain all
	// the domain names bound to the certificate.
	return testAccCertificateBatchDomainsAssociate_basic_step3(name)
}

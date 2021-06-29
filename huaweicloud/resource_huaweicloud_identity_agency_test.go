package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/identity/v3/agency"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccIdentityAgency_basic(t *testing.T) {
	var agency agency.Agency

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_identity_agency.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityAgencyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityAgency_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityAgencyExists(resourceName, &agency),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a test agency"),
					resource.TestCheckResourceAttr(resourceName, "delegated_service_name", "op_svc_evs"),
					resource.TestCheckResourceAttr(resourceName, "duration", "FOREVER"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIdentityAgency_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityAgencyExists(resourceName, &agency),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a updated test agency"),
					resource.TestCheckResourceAttr(resourceName, "delegated_service_name", "op_svc_evs"),
					resource.TestCheckResourceAttr(resourceName, "duration", "FOREVER"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "2"),
				),
			},
		},
	})
}

func TestAccIdentityAgency_domain(t *testing.T) {
	var agency agency.Agency

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_identity_agency.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityAgencyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityAgency_domain(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityAgencyExists(resourceName, &agency),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a test agency"),
					resource.TestCheckResourceAttr(resourceName, "delegated_domain_name", HW_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "duration", "FOREVER"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIdentityAgency_domainUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityAgencyExists(resourceName, &agency),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a updated test agency"),
					resource.TestCheckResourceAttr(resourceName, "delegated_domain_name", HW_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "duration", "FOREVER"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "2"),
				),
			},
		},
	})
}

func testAccCheckIdentityAgencyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.IAMV3Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud IAM client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_identity_agency" {
			continue
		}

		v, err := agency.Get(client, rs.Primary.ID).Extract()
		if err == nil && v.ID == rs.Primary.ID {
			return fmtp.Errorf("Identity Agency <%s> still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckIdentityAgencyExists(n string, ag *agency.Agency) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		client, err := config.IAMV3Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud Identity Agency: %s", err)
		}

		found, err := agency.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("Identity Agency <%s> not found", rs.Primary.ID)
		}
		ag = found

		return nil
	}
}

func testAccIdentityAgency_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_agency" "test" {
  name                   = "%s"
  description            = "This is a test agency"
  delegated_service_name = "op_svc_evs"

  domain_roles = [
    "OBS OperateAccess",
  ]
}
`, rName)
}

func testAccIdentityAgency_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_agency" "test" {
  name                   = "%s"
  description            = "This is a updated test agency"
  delegated_service_name = "op_svc_evs"

  domain_roles = [
    "OBS OperateAccess", "KMS Administrator",
  ]
}
`, rName)
}

func testAccIdentityAgency_domain(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_agency" "test" {
  name                  = "%s"
  description           = "This is a test agency"
  delegated_domain_name = "%s"

  domain_roles = [
    "Anti-DDoS Administrator",
  ]
}
`, rName, HW_DOMAIN_NAME)
}

func testAccIdentityAgency_domainUpdate(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_agency" "test" {
  name                  = "%s"
  description           = "This is a updated test agency"
  delegated_domain_name = "%s"

  domain_roles = [
    "Anti-DDoS Administrator",
    "Ticket Administrator",
  ]
}
`, rName, HW_DOMAIN_NAME)
}

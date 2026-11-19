package iam

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getV3TrustAgencyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam_no_version", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}
	getAgencyHttpUrl := "v5/agencies/{agency_id}"
	getAgencyPath := client.Endpoint + getAgencyHttpUrl
	getAgencyPath = strings.ReplaceAll(getAgencyPath, "{agency_id}", state.Primary.ID)
	getAgencyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getAgencyResp, err := client.Request("GET", getAgencyPath, &getAgencyOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IAM trust agency: %s", err)
	}
	return utils.FlattenResponse(getAgencyResp)
}

// TestAccV3TrustAgency_basic covers the basic CRUD flow for a trust agency
// whose trust policy principal is a HUAWEI CLOUD service. Step 3 updates the
// trust_policy in-place (now that it is no longer ForceNew) to verify the
// new `PUT /v5/agencies/{agency_id}/trust-policy` integration.
func TestAccV3TrustAgency_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_identity_trust_agency.test"
		rName        = acceptance.RandomAccResourceName()

		rc = acceptance.InitResourceCheck(resourceName, &obj, getV3TrustAgencyResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPreCheckIAMV5(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV3TrustAgency_basic_step1(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "policy_names.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "max_session_duration", "3600"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttrSet(resourceName, "trust_policy"),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				Config: testAccV3TrustAgency_basic_step2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "policy_names.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "max_session_duration", "7200"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo1", "bar1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "trust_policy", `{"Statement":[{"Action":["sts:agencies:assume"],"Effect":"Allow","Principal":{"Service":["service.APIG"]}}],"Version":"5.0"}`),
				),
			},
			{
				// Verify `trust_policy` update works in-place (no ForceNew).
				// Changing the service principal from service.APIG to service.CCE
				// must not destroy and recreate the trust agency.
				Config: testAccV3TrustAgency_basic_step3(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "trust_policy", `{"Statement":[{"Action":["sts:agencies:assume"],"Effect":"Allow","Principal":{"Service":["service.CCE"]}}],"Version":"5.0"}`),
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

// TestAccV3TrustAgency_federatedSamlIdp covers the federated scenario:
// the trust agency's trust policy grants the federated SAML identity
// provider permission to assume the agency.
func TestAccV3TrustAgency_federatedSamlIdp(t *testing.T) {
	var (
		obj interface{}

		agencyName = "huaweicloud_identity_trust_agency.test"
		samlName   = "huaweicloud_identityv5_saml_provider.test"
		rName      = acceptance.RandomAccResourceName()

		rc = acceptance.InitResourceCheck(agencyName, &obj, getV3TrustAgencyResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPreCheckIAMV5(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV3TrustAgency_federatedSamlIdp(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					// The trust policy must reference the SAML provider URN
					resource.TestMatchResourceAttr(agencyName, "trust_policy", regexp.MustCompile(`.*IdP.*`)),
					resource.TestCheckResourceAttrSet(samlName, "urn"),
				),
			},
		},
	})
}

func testAccV3TrustAgency_basic_step1(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_trust_agency" "test" {
  name         = "%s"
  policy_names = ["NATReadOnlyPolicy"]
  description  = "Created by terraform script"
  trust_policy = jsonencode(
    {
      Statement = [
        {
          Action = [
            "sts:agencies:assume",
          ]
          Effect = "Allow"
          Principal = {
            Service = [
              "service.APIG",
            ]
          }
        },
      ]
      Version = "5.0"
    }
  )

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName)
}

func testAccV3TrustAgency_basic_step2(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_trust_agency" "test" {
  name         = "%s"
  policy_names = ["NATReadOnlyPolicy", "RDSReadOnlyPolicy"]
  max_session_duration = 7200
  description  = "Updated by terraform script"
  trust_policy = jsonencode(
    {
      Statement = [
        {
          Action = [
            "sts:agencies:assume",
          ]
          Effect = "Allow"
          Principal = {
            Service = [
              "service.APIG",
            ]
          }
        },
      ]
      Version = "5.0"
    }
  )

  tags = {
    foo1 = "bar1"
    key1 = "value1"
  }
}
`, rName)
}

// testAccV3TrustAgency_basic_step3 changes the trust policy (Principal.Service)
// to a different service. This validates that trust_policy is updatable in-place
// rather than triggering a ForceNew recreation.
func testAccV3TrustAgency_basic_step3(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_trust_agency" "test" {
  name         = "%s"
  policy_names = ["NATReadOnlyPolicy", "RDSReadOnlyPolicy"]
  max_session_duration = 7200
  description  = "Updated by terraform script"
  trust_policy = jsonencode(
    {
      Statement = [
        {
          Action = [
            "sts:agencies:assume",
          ]
          Effect = "Allow"
          Principal = {
            Service = [
              "service.CCE",
            ]
          }
        },
      ]
      Version = "5.0"
    }
  )

  tags = {
    foo1 = "bar1"
    key1 = "value1"
  }
}
`, rName)
}

// testAccV3TrustAgency_federatedSamlIdp builds a SAML identity provider and a
// trust agency whose trust policy grants the SAML provider permission to
// assume the agency.
func testAccV3TrustAgency_federatedSamlIdp(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_saml_provider" "test" {
  name        = "%[1]s-idp"
  description = "Federated SAML IdP for agency trust test"

  saml_metadata_document = <<EOF
%s
EOF
}

resource "huaweicloud_identity_trust_agency" "test" {
  name         = "%[1]s-agency"
  policy_names = ["AdministratorAccess"]
  description  = "Trust agency for federated SAML IdP"
  trust_policy = jsonencode(
    {
      Statement = [
        {
          Action = [
            "sts:agencies:assume",
          ]
          Effect = "Allow"
          Principal = {
            IdP = [
              huaweicloud_identityv5_saml_provider.test.urn,
            ]
          }
        },
      ]
      Version = "5.0"
    }
  )
}
`, rName, testAccV5SamlMetadataDocument())
}

package iam

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// The agency URN is used for both the OIDC and SAML credential tests.
func testAccV5AgencyCredentialPreCheck(t *testing.T) {
	acceptance.TestAccPreCheck(t)
	acceptance.TestAccPreCheckIAMV5(t)
	if os.Getenv("HW_AGENCY_URN") == "" {
		t.Skip("This test requires HW_AGENCY_URN to be set (the URN of a trust agency whose trust principal is an identity provider).")
	}
}

// The ID token must be issued by an OIDC identity provider already registered in IAM,
// and the audience of the token should match the registered client ID.
func testAccV5OidcCredentialPreCheck(t *testing.T) {
	testAccV5AgencyCredentialPreCheck(t)
	if os.Getenv("HW_OIDC_ID_TOKEN") == "" {
		t.Skip("This test requires HW_OIDC_ID_TOKEN to be set (the ID token issued by the OIDC identity provider).")
	}
}

// The SAML response must be a base64 encoded assertion response issued by the SAML identity provider.
func testAccV5SamlCredentialPreCheck(t *testing.T) {
	testAccV5AgencyCredentialPreCheck(t)
	if os.Getenv("HW_SAML_RESPONSE") == "" {
		t.Skip("This test requires HW_SAML_RESPONSE to be set (the base64 encoded SAML authentication response).")
	}
}

func TestAccV5OidcCredential_basic(t *testing.T) {
	rName := "data.huaweicloud_identityv5_oidc_credential.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccV5OidcCredentialPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccV5OidcCredential_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "access_key_id"),
					resource.TestCheckResourceAttrSet(rName, "secret_access_key"),
					resource.TestCheckResourceAttrSet(rName, "security_token"),
					resource.TestCheckResourceAttrSet(rName, "expiration"),
					resource.TestCheckResourceAttrSet(rName, "assumed_agency_urn"),
				),
			},
		},
	})
}

func testAccV5OidcCredential_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_identityv5_oidc_credential" "test" {
  id_token            = "%s"
  agency_urn          = "%s"
  agency_session_name = "terraform-acc-test"
}
`, os.Getenv("HW_OIDC_ID_TOKEN"), os.Getenv("HW_AGENCY_URN"))
}

func TestAccV5SamlCredential_basic(t *testing.T) {
	rName := "data.huaweicloud_identityv5_saml_credential.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccV5SamlCredentialPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccV5SamlCredential_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "access_key_id"),
					resource.TestCheckResourceAttrSet(rName, "secret_access_key"),
					resource.TestCheckResourceAttrSet(rName, "security_token"),
					resource.TestCheckResourceAttrSet(rName, "expiration"),
					resource.TestCheckResourceAttrSet(rName, "assumed_agency_urn"),
					resource.TestCheckResourceAttrSet(rName, "subject"),
				),
			},
		},
	})
}

func testAccV5SamlCredential_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_identityv5_saml_credential" "test" {
  saml_response       = "%s"
  agency_urn          = "%s"
  agency_session_name = "terraform-acc-test"
}
`, os.Getenv("HW_SAML_RESPONSE"), os.Getenv("HW_AGENCY_URN"))
}

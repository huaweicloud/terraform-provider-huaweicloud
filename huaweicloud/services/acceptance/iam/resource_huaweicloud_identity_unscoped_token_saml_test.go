package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityUnscopedTokenSaml_basic(t *testing.T) {
	idpId := "YourIdpId"
	samlResponse := "YourSamlResponse"
	resourceName := "huaweicloud_identity_unscoped_token_saml.test"

	// Avoid CheckDestroy because the token can not be destroyed.
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityUnscopedTokenSaml(idpId, samlResponse),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "token"),
					resource.TestCheckResourceAttrSet(resourceName, "username"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.#"),
					resource.TestCheckResourceAttrSet(resourceName, "expires_at"),
				),
			},
		},
	})
}

func testAccIdentityUnscopedTokenSaml(idpId, samlResponse string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_unscoped_token_saml" "test" {
  idp_id             = "%[1]s"
  saml_response      = "%[2]s"
  with_global_domain = true
}
`, idpId, samlResponse)
}

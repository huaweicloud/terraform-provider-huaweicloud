package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityUnscopedTokenWithIdToken_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_identity_unscoped_token_with_id_token.test"

		name = acceptance.RandomAccResourceName()
	)

	// Avoid CheckDestroy because the token can not be destroyed.
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckIdentityIDPId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: "3.3.0",
			},
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityUnscopedTokenWithIdToken(name),
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

func testAccIdentityUnscopedTokenWithIdToken(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identity_unscoped_token_with_id_token" "test" {
  idp_id      = "%[2]s"
  protocol_id = "oidc"
  id_token    = huaweicloud_identity_temporary_access_key.test.securitytoken
}
`, testAccTemporaryAccessKey_basic_step1(name), acceptance.HW_IDENTITY_OIDC_IDP_ID)
}

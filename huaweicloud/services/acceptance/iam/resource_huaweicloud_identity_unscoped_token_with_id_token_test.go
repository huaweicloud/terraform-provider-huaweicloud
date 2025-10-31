package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityUnscopedTokenWithIdToken_basic(t *testing.T) {
	idpId := "YourIdpId"
	protocolId := "YourProtocolId"
	idToken := "YourIdToken"
	resourceName := "huaweicloud_identity_unscoped_token_with_id_token.test"

	// Avoid CheckDestroy because the token can not be destroyed.
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityUnscopedTokenWithIdToken(idpId, protocolId, idToken),
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

func testAccIdentityUnscopedTokenWithIdToken(idpId, protocolId, idToken string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_unscoped_token_with_id_token" "test" {
  idp_id      = "%[1]s"
  protocol_id = "%[2]s"
  id_token    = "%[3]s"
}
`, idpId, protocolId, idToken)
}

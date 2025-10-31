package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityTokenWithIdToken_basic(t *testing.T) {
	idpId := "YourIdpId"
	idToken := "YourIdToken"
	domainName := "YourDomainName"
	projectName := "cn-north-4"
	resourceName := "huaweicloud_identity_token_with_id_token.test"

	// Avoid CheckDestroy because the token can not be destroyed.
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityTokenWithIdToken(idpId, idToken),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "token"),
					resource.TestCheckResourceAttrSet(resourceName, "username"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.#"),
					resource.TestCheckResourceAttrSet(resourceName, "expires_at"),
				),
			},
			{
				Config: testAccIdentityTokenWithIdTokenWithDomainName(idpId, idToken, domainName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "token"),
					resource.TestCheckResourceAttrSet(resourceName, "username"),
					resource.TestCheckResourceAttrSet(resourceName, "groups.#"),
					resource.TestCheckResourceAttrSet(resourceName, "expires_at"),
				),
			},
			{
				Config: testAccIdentityTokenWithIdTokenWithProjectName(idpId, idToken, projectName),
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

func testAccIdentityTokenWithIdToken(idpId, idToken string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_token_with_id_token" "test" {
  idp_id   = "%[1]s"
  id_token = "%[2]s"
}
`, idpId, idToken)
}

func testAccIdentityTokenWithIdTokenWithDomainName(idpId, idToken, domainName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_token_with_id_token" "test" {
  idp_id      = "%[1]s"
  id_token    = "%[2]s"
  domain_name = "%[3]s"
}
`, idpId, idToken, domainName)
}

func testAccIdentityTokenWithIdTokenWithProjectName(idpId, idToken, projectName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_token_with_id_token" "test" {
  idp_id       = "%[1]s"
  id_token     = "%[2]s"
  project_name = "%[3]s"
}
`, idpId, idToken, projectName)
}

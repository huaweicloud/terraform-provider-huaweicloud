package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityCenterClinet_basic(t *testing.T) {
	rName := "huaweicloud_identitycenter_client.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testIdentityCenterClient_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "client_secret"),
					resource.TestCheckResourceAttrSet(rName, "client_secret_expires_at"),
				),
			},
		},
	})
}

func testIdentityCenterClient_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_identitycenter_client" "test"{
  client_name    = "client_test"
  client_type    = "public"
  token_endpoint_auth_method = "client_secret_post"
  grant_types 	= ["urn:ietf:params:oauth:grant-type:device_code"]
  response_types  = ["code"]
  scopes		= ["openid"]
}
`)
}

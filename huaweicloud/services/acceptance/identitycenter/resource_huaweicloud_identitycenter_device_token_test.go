package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityCenterDeviceToken_basic(t *testing.T) {
	rName := "huaweicloud_identitycenter_device_token.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testIdentityCenterDeviceToken_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "access_token"),
					resource.TestCheckResourceAttrSet(rName, "expires_in"),
				),
			},
		},
	})
}

func testIdentityCenterDeviceToken_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_identitycenter_device_token" "test"{
	client_id      = "%s"
	client_secret   = "%s"
	device_code	= "%s"
	grant_type		= "urn:ietf:params:oauth:grant-type:device_code"
}
`, acceptance.HW_IDENTITY_CENTER_CLIENT_ID, acceptance.HW_IDENTITY_CENTER_CLIENT_SECRET, acceptance.HW_IDENTITY_CENTER_VERIFIED_DEVICE_CODE)
}

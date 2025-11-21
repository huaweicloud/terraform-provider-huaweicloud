package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccPasswordReset_basic(t *testing.T) {
	userId := acceptance.HW_IDENTITY_CENTER_USER_ID

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckIdentityCenterUserId(t)
		},

		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testPasswordReset_basic(userId),
			},
		},
	})
}

func testPasswordReset_basic(userId string) string {
	return fmt.Sprintf(`
data "huaweicloud_identitycenter_instance" "test" {}
 
resource "huaweicloud_identitycenter_password_reset" "test"{
  mode              = "EMAIL"
  user_id           = "%s"
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
}
`, userId)
}

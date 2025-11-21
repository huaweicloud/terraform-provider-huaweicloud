package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccUserSessionDelete_basic(t *testing.T) {
	userId := acceptance.HW_IDENTITY_CENTER_USER_ID
	sessionId := acceptance.HW_IDENTITY_CENTER_SESSION_ID

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckIdentityCenterUserId(t)
			acceptance.TestAccPreCheckIdentityCenterSessionId(t)
		},

		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testUserSessionDelete_basic(userId, sessionId),
			},
		},
	})
}

func testUserSessionDelete_basic(userId string, sessionId string) string {
	return fmt.Sprintf(`
data "huaweicloud_identitycenter_instance" "test" {}
 
resource "huaweicloud_identitycenter_user_session_delete" "test" {
  user_id           = "%[1]s"
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
  session_ids       = ["%[2]s"]
}
`, userId, sessionId)
}

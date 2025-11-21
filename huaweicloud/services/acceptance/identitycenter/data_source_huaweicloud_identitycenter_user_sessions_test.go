package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceIdentityCenterUserSessions_basic(t *testing.T) {
	userId := acceptance.HW_IDENTITY_CENTER_USER_ID
	rName := "data.huaweicloud_identitycenter_user_sessions.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckIdentityCenterUserId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceIdentityCenterUserSessions_basic(userId),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "session_list.0.creation_time"),
					resource.TestCheckResourceAttrSet(rName, "session_list.0.ip_address"),
					resource.TestCheckResourceAttrSet(rName, "session_list.0.session_id"),
					resource.TestCheckResourceAttrSet(rName, "session_list.0.session_not_valid_after"),
					resource.TestCheckResourceAttrSet(rName, "session_list.0.user_agent"),
				),
			},
		},
	})
}

func testAccDatasourceIdentityCenterUserSessions_basic(userId string) string {
	return fmt.Sprintf(`
data "huaweicloud_identitycenter_instance" "test" {}

data "huaweicloud_identitycenter_user_sessions" "test"{
  user_id           = "%s"
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
}
`, userId)
}

package meeting

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/meeting/v1/assignments"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAdminAssignment_basic(t *testing.T) {
	var (
		administrator assignments.Administrator
		resourceName  = "huaweicloud_meeting_admin_assignment.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&administrator,
		getUserFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAppAuth(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAdminAssignment_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "app_id", acceptance.HW_MEETING_APP_ID),
					resource.TestCheckResourceAttr(resourceName, "app_key", acceptance.HW_MEETING_APP_KEY),
					resource.TestCheckResourceAttrPair(resourceName, "account", "huaweicloud_meeting_user.test", "account"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAdminAssignmentImportStateIdFunc(),
			},
		},
	})
}

func testAccAdminAssignmentImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var (
			accountName, password string
			appId, appKey         string
			account               string
		)
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "huaweicloud_meeting_admin_assignment" {
				accountName = rs.Primary.Attributes["account_name"]
				password = rs.Primary.Attributes["account_password"]
				appId = rs.Primary.Attributes["app_id"]
				appKey = rs.Primary.Attributes["app_key"]
				account = rs.Primary.ID
			}
		}
		if account != "" && accountName != "" && password != "" {
			return fmt.Sprintf("%s/%s/%s", account, accountName, password), nil
		}
		if account != "" && appId != "" && appKey != "" {
			return fmt.Sprintf("%s/%s/%s//", account, appId, appKey), nil
		}
		return "", fmt.Errorf("resource not found: %s", account)
	}
}

func testAccAdminAssignment_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_meeting_user" "test" {
  app_id  = "%[1]s"
  app_key = "%[2]s"

  name     = "Test Name"
  password = "HuaweiTest@123"
  country  = "chinaPR"
  email    = "123456789@example.com"
  phone    = "+8612345678987"
}

resource "huaweicloud_meeting_admin_assignment" "test" {
  app_id  = "%[1]s"
  app_key = "%[2]s"

  account = huaweicloud_meeting_user.test.account
}
`, acceptance.HW_MEETING_APP_ID, acceptance.HW_MEETING_APP_KEY)
}

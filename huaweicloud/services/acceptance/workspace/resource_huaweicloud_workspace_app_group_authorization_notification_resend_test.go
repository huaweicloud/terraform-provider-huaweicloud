package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAppGroupAuthorizationNotificationResend_basic(t *testing.T) {
	var (
		rName                  = "huaweicloud_workspace_app_group_authorization_notification_resend.test"
		withNotificationRecord = "huaweicloud_workspace_app_group_authorization_notification_resend.with_notification_record"

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAppGroupAuthorizationNotificationResend_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "records.#", "1"),
					resource.TestMatchResourceAttr(withNotificationRecord, "records.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
				),
			},
		},
	})
}

func testAccAppGroupAuthorizationNotificationResend_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_group" "test" {
  name = "%[1]s"
}

resource "huaweicloud_workspace_user" "test" {
  name  = "%[1]s"
  email = "%[1]s@example.com"

  account_expires            = "0"
  password_never_expires     = false
  enable_change_password     = true
  next_login_change_password = true
  disabled                   = false
}

resource "huaweicloud_workspace_app_group_authorization" "test" {
  app_group_id = huaweicloud_workspace_app_group.test.id

  accounts {
    id      = huaweicloud_workspace_user.test.id
    account = huaweicloud_workspace_user.test.name
    type    = "USER"
  }

  # Wait for the authorization notification record to be generated.
  provisioner "local-exec" {
    command = "sleep 20"
  }
}

data "huaweicloud_workspace_app_group_authorizations" "test" {
  app_group_id = huaweicloud_workspace_app_group.test.id
  account      = huaweicloud_workspace_user.test.name
  account_type = "USER"

  depends_on = [huaweicloud_workspace_app_group_authorization.test]
}

locals {
  authorization_id = try([for v in data.huaweicloud_workspace_app_group_authorizations.test.authorizations :
  v.id if v.account_id == huaweicloud_workspace_user.test.id][0], "")
}

# Resend by authorization record.
resource "huaweicloud_workspace_app_group_authorization_notification_resend" "test" {
  records {
    id = local.authorization_id
  }
}

data "huaweicloud_workspace_app_group_authorization_notification_records" "test" {
  app_group_id = huaweicloud_workspace_app_group.test.id
  account      = huaweicloud_workspace_user.test.name

  depends_on = [
    huaweicloud_workspace_app_group_authorization.test
  ]
}

# Resend by authorization notification record.
resource "huaweicloud_workspace_app_group_authorization_notification_resend" "with_notification_record" {
  is_notification_record = true

  dynamic "records" {
    for_each = data.huaweicloud_workspace_app_group_authorization_notification_records.test.records

    content {
      id                = records.value["id"]
      account           = records.value["account"]
      account_auth_type = records.value["account_auth_type"]
      account_auth_name = records.value["account_auth_name"]
      app_group_id      = records.value["app_group_id"]
      app_group_name    = records.value["app_group_name"]
      mail_send_type    = records.value["mail_send_type"]
    }
  }

  # A new record is generated each time a notification is successfully resent, so the changes need to be ignored.
  lifecycle {
    ignore_changes = [
      records
    ]
  }
}
`, name)
}

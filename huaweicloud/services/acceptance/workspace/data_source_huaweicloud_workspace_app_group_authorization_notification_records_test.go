package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAppGroupAuthorizationNotificationRecords_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_workspace_app_group_authorization_notification_records.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byAccount   = "data.huaweicloud_workspace_app_group_authorization_notification_records.filter_by_account"
		dcByAccount = acceptance.InitDataSourceCheck(byAccount)

		byMailSendType   = "data.huaweicloud_workspace_app_group_authorization_notification_records.filter_by_mail_send_type"
		dcByMailSendType = acceptance.InitDataSourceCheck(byMailSendType)

		byMailSendResult   = "data.huaweicloud_workspace_app_group_authorization_notification_records.filter_by_mail_send_result"
		dcByMailSendResult = acceptance.InitDataSourceCheck(byMailSendResult)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAppGroupAuthorizationNotificationRecords_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "records.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByAccount.CheckResourceExists(),
					resource.TestCheckOutput("is_account_useful", "true"),
					dcByMailSendType.CheckResourceExists(),
					resource.TestCheckOutput("is_mail_send_type_useful", "true"),
					dcByMailSendResult.CheckResourceExists(),
					resource.TestCheckOutput("is_mail_send_result_useful", "true"),
					// Check attributes.
					resource.TestCheckResourceAttrSet(byAccount, "records.0.id"),
					resource.TestCheckResourceAttrSet(byAccount, "records.0.account"),
					resource.TestCheckResourceAttrSet(byAccount, "records.0.account_auth_type"),
					resource.TestCheckResourceAttrSet(byAccount, "records.0.account_auth_name"),
					resource.TestCheckResourceAttrSet(byAccount, "records.0.app_group_id"),
					resource.TestCheckResourceAttrSet(byAccount, "records.0.app_group_name"),
					resource.TestCheckResourceAttrSet(byAccount, "records.0.mail_send_type"),
					resource.TestCheckResourceAttrSet(byAccount, "records.0.mail_send_result"),
					resource.TestCheckResourceAttrSet(byAccount, "records.0.send_at"),
				),
			},
		},
	})
}

func testAccDataAppGroupAuthorizationNotificationRecords_basic() string {
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

  # Wait for the authorization records to be generated.
  provisioner "local-exec" {
    command = "sleep 20"
  }
}

data "huaweicloud_workspace_app_group_authorization_notification_records" "test" {
  app_group_id = huaweicloud_workspace_app_group.test.id
  depends_on   = [huaweicloud_workspace_app_group_authorization.test]
}

# Filter by account name.
locals {
  account_name = huaweicloud_workspace_user.test.name
}

data "huaweicloud_workspace_app_group_authorization_notification_records" "filter_by_account" {
  app_group_id = huaweicloud_workspace_app_group.test.id
  account      = local.account_name

  depends_on = [huaweicloud_workspace_app_group_authorization.test]
}

locals {
  account_filter_result = [for v in data.huaweicloud_workspace_app_group_authorization_notification_records.filter_by_account.records :
  strcontains(v.account, local.account_name)]
}

output "is_account_useful" {
  value = length(local.account_filter_result) > 0 && alltrue(local.account_filter_result)
}

# Filter by mail send type.
locals {
  mail_send_type = "ADD_GROUP_AUTHORIZATION"
}

data "huaweicloud_workspace_app_group_authorization_notification_records" "filter_by_mail_send_type" {
  app_group_id   = huaweicloud_workspace_app_group.test.id
  mail_send_type = local.mail_send_type
  depends_on     = [huaweicloud_workspace_app_group_authorization.test]
}

locals {
  mail_send_type_filter_result = [for v in data.huaweicloud_workspace_app_group_authorization_notification_records.filter_by_mail_send_type.records :
  v.mail_send_type == local.mail_send_type]
}

output "is_mail_send_type_useful" {
  value = length(local.mail_send_type_filter_result) > 0 && alltrue(local.mail_send_type_filter_result)
}

# Filter by mail send result.
locals {
  mail_send_result = "SUCCESS"
}

data "huaweicloud_workspace_app_group_authorization_notification_records" "filter_by_mail_send_result" {
  app_group_id     = huaweicloud_workspace_app_group.test.id
  mail_send_result = local.mail_send_result
  depends_on       = [huaweicloud_workspace_app_group_authorization.test]
}

locals {
  mail_send_result_filter_records = data.huaweicloud_workspace_app_group_authorization_notification_records.filter_by_mail_send_result.records
  mail_send_result_filter_result  = [for v in local.mail_send_result_filter_records : v.mail_send_result == local.mail_send_result]
}

output "is_mail_send_result_useful" {
  value = length(local.mail_send_result_filter_result) > 0 && alltrue(local.mail_send_result_filter_result)
}
`, acceptance.RandomAccResourceNameWithDash())
}

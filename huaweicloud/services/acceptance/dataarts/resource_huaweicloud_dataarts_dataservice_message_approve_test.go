package dataarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataServiceMessageApprove_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dataarts_dataservice_messages.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsReviewerName(t)
			acceptance.TestAccPreCheckDataArtsRelatedDliQueueName(t)
		},

		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// Authorize the APP to access the API (do not need to approve).
				Config: testAccDataServiceMessageApprove_basic_step1(name),
			},
			{
				// Cancel the API access permission of the APP, then generate an approve message and inform the reviewer.
				Config: testAccDataServiceMessageApprove_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_any_message_pending_check", "true"),
				),
			},
		},
	})
}

func testAccDataServiceMessageApprove_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_dataservice_api_auth_action" "publish" {
  depends_on = [huaweicloud_dataarts_dataservice_api_publishment.test]

  workspace_id = "%[2]s"

  api_id      = huaweicloud_dataarts_dataservice_api.test.id
  instance_id = try(data.huaweicloud_dataarts_dataservice_instances.test.instances[0].id, "")
  app_id      = huaweicloud_dataarts_dataservice_app.test[0].id
  type        = "APPLY_TYPE_AUTHORIZE"
}
`, testAccDataServiceApiAuth_base(name),
		acceptance.HW_DATAARTS_WORKSPACE_ID)
}

func testAccDataServiceMessageApprove_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_dataservice_api_auth_action" "unpublish" {
  depends_on = [huaweicloud_dataarts_dataservice_api_publishment.test]

  workspace_id = "%[2]s"

  api_id      = huaweicloud_dataarts_dataservice_api.test.id
  instance_id = try(data.huaweicloud_dataarts_dataservice_instances.test.instances[0].id, "")
  app_id      = huaweicloud_dataarts_dataservice_app.test[0].id
  type        = "APPLY_TYPE_API_CANCEL_AUTHORIZE"
}

data "huaweicloud_dataarts_dataservice_messages" "test" {
  depends_on = [huaweicloud_dataarts_dataservice_api_auth_action.unpublish]

  workspace_id = "%[2]s"

  api_name = huaweicloud_dataarts_dataservice_api.test.name
}

locals {
  pending_check_message_ids = [
    for v in data.huaweicloud_dataarts_dataservice_messages.test.messages: v.id if v.api_apply_status == "STATUS_TYPE_PENDING_CHECK"
  ]
}

output "is_any_message_pending_check" {
  value = length(local.pending_check_message_ids) > 0
}

resource "huaweicloud_dataarts_dataservice_message_approve" "approve_immediately" {
  workspace_id = "%[2]s"

  message_id = try(local.pending_check_message_ids[0], "")
  action     = 0

  lifecycle {
    ignore_changes = [
      message_id,
    ]
  }
}
`, testAccDataServiceApiAuth_base(name),
		acceptance.HW_DATAARTS_WORKSPACE_ID)
}

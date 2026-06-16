package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccArchitectureApprovalsBatchAction_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccArchitectureApprovalsBatchAction_nonExistentWorkspace(),
				ExpectError: regexp.MustCompile(`error action withdrawing architecture approvals`),
			},
			{
				Config:      testAccArchitectureApprovalsBatchAction_nonExistentAction(),
				ExpectError: regexp.MustCompile("invalid action type"),
			},
			{
				Config: testAccArchitectureApprovalsBatchAction_basic(name),
			},
		},
	})
}

func testAccArchitectureApprovalsBatchAction_nonExistentWorkspace() string {
	randUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_approvals_batch_action" "non_existent_workspace" {
  workspace_id = "%[1]s"
  approval_ids = "%[1]s"
  action       = "recall"
}
`, randUUID.String())
}

func testAccArchitectureApprovalsBatchAction_nonExistentAction() string {
	randUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_approvals_batch_action" "non_existent_action" {
  workspace_id = "%[1]s"
  approval_ids = "%[1]s"
  action       = "%[1]s"
}
`, randUUID.String())
}

func testAccArchitectureApprovalsBatchAction_basic_base(name string) string {
	return fmt.Sprintf(`
variable "subjects" {
  description = "A list of subjects to be created."
  type        = list(map(string))

  default = [
    {
      name = "%[1]s-biz_theme"
      code = "%[1]s-biz_theme"
    },
    {
      name = "%[1]s-data_theme"
      code = "%[1]s-data_theme"
    },
    {
      name = "%[1]s-tech_theme"
      code = "%[1]s-tech_theme"
    }
  ]
}

resource "huaweicloud_dataarts_architecture_subject" "test" {
  count = length(var.subjects)

  workspace_id = "%[2]s"
  name         = lookup(var.subjects[count.index], "name", "")
  code         = lookup(var.subjects[count.index], "code", "")
  owner        = "%[1]s"
  level        = 1
  description  = "Created by Terraform with for_each"
}

resource "huaweicloud_dataarts_architecture_batch_publish" "test" {
  workspace_id       = "%[2]s"
  approver_user_id   = "%[3]s"
  approver_user_name = "%[4]s"
  fast_approval      = false

  dynamic "biz_infos" {
    for_each = huaweicloud_dataarts_architecture_subject.test[*].id

    content {
      biz_id   = biz_infos.value
      biz_type = "SUBJECT"
    }
  }
}

data "huaweicloud_dataarts_architecture_approvals" "test" {
  depends_on = [
    huaweicloud_dataarts_architecture_batch_publish.test,
  ]

  workspace_id    = "%[2]s"
  approval_status = "DEVELOPING"
}
`, name, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_ARCHITECTURE_USER_ID, acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME)
}

func testAccArchitectureApprovalsBatchAction_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_architecture_approvals_batch_action" "batch_recall" {
  workspace_id = "%[2]s"
  approval_ids = join(",", [for v in data.huaweicloud_dataarts_architecture_approvals.test.approvals : v.id])
  action       = "recall"
  
  lifecycle {
    ignore_changes = [approval_ids]
  }
}

resource "huaweicloud_dataarts_architecture_batch_publish" "recall_approvals_batch_publish" {
  depends_on = [
    huaweicloud_dataarts_architecture_approvals_batch_action.batch_recall,
  ]

  workspace_id       = "%[2]s"
  approver_user_id   = "%[3]s"
  approver_user_name = "%[4]s"
  fast_approval      = false

  dynamic "biz_infos" {
    for_each = huaweicloud_dataarts_architecture_subject.test[*].id

    content {
      biz_id   = biz_infos.value
      biz_type = "SUBJECT"
    }
  }
}

data "huaweicloud_dataarts_architecture_approvals" "test_batch_reject_approvals" {
  depends_on = [
    huaweicloud_dataarts_architecture_batch_publish.recall_approvals_batch_publish,
  ]

  workspace_id    = "%[2]s"
  approval_status = "DEVELOPING"
}

resource "huaweicloud_dataarts_architecture_approvals_batch_action" "batch_reject" {
  workspace_id = "%[2]s"
  approval_ids = join(",", [for v in data.huaweicloud_dataarts_architecture_approvals.test_batch_reject_approvals.approvals : v.id])
  action       = "reject"
  message      = "Batch Reject Approval"

  lifecycle {
    ignore_changes = [approval_ids]
  }
}

resource "huaweicloud_dataarts_architecture_batch_publish" "resolve_approvals_batch_publish" {
  depends_on = [
    huaweicloud_dataarts_architecture_approvals_batch_action.batch_reject,
  ]

  workspace_id       = "%[2]s"
  approver_user_id   = "%[3]s"
  approver_user_name = "%[4]s"
  fast_approval      = false

  dynamic "biz_infos" {
    for_each = huaweicloud_dataarts_architecture_subject.test[*].id

    content {
      biz_id   = biz_infos.value
      biz_type = "SUBJECT"
    }
  }
}

data "huaweicloud_dataarts_architecture_approvals" "test_batch_resolve_approvals" {
  depends_on = [
    huaweicloud_dataarts_architecture_batch_publish.resolve_approvals_batch_publish,
  ]

  workspace_id    = "%[2]s"
  approval_status = "DEVELOPING"
}

resource "huaweicloud_dataarts_architecture_approvals_batch_action" "batch_resolve" {
  workspace_id = "%[2]s"
  approval_ids = join(",", [for v in data.huaweicloud_dataarts_architecture_approvals.test_batch_resolve_approvals.approvals : v.id])
  action       = "resolve"
  message      = "Batch Resolve Approval"

  lifecycle {
    ignore_changes = [approval_ids]
  }
}

resource "huaweicloud_dataarts_architecture_batch_unpublish" "batch_unpublish" {
  depends_on = [
    huaweicloud_dataarts_architecture_approvals_batch_action.batch_resolve,
  ]

  workspace_id       = "%[2]s"
  approver_user_id   = "%[3]s"
  approver_user_name = "%[4]s"
  fast_approval      = true

  dynamic "biz_infos" {
    for_each = huaweicloud_dataarts_architecture_subject.test[*].id

    content {
      biz_id   = biz_infos.value
      biz_type = "SUBJECT"
    }
  }
}
`, testAccArchitectureApprovalsBatchAction_basic_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID,
		acceptance.HW_DATAARTS_ARCHITECTURE_USER_ID, acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME)
}

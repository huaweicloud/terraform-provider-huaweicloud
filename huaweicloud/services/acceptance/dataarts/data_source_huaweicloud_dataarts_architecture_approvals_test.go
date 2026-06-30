package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataArchitectureApprovals_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dataarts_architecture_approvals.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByBizId   = "data.huaweicloud_dataarts_architecture_approvals.filter_by_biz_id"
		dcFilterByBizId = acceptance.InitDataSourceCheck(filterByBizId)

		filterByName   = "data.huaweicloud_dataarts_architecture_approvals.filter_by_name"
		dcFilterByName = acceptance.InitDataSourceCheck(filterByName)

		filterByCreateBy   = "data.huaweicloud_dataarts_architecture_approvals.filter_by_create_by"
		dcFilterByCreateBy = acceptance.InitDataSourceCheck(filterByCreateBy)

		filterByApprover   = "data.huaweicloud_dataarts_architecture_approvals.filter_by_approver"
		dcFilterByApprover = acceptance.InitDataSourceCheck(filterByApprover)

		filterByApprovalStatus   = "data.huaweicloud_dataarts_architecture_approvals.filter_by_approval_status"
		dcFilterByApprovalStatus = acceptance.InitDataSourceCheck(filterByApprovalStatus)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataArchitectureApprovals_nonExistentWorkspace(),
				ExpectError: regexp.MustCompile("error querying DataArts Architecture approvals"),
			},
			{
				Config: testAccDataSourceArchitectureApprovals_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "approvals.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),

					// filter by biz ID
					dcFilterByBizId.CheckResourceExists(),
					resource.TestCheckOutput("is_biz_id_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(filterByBizId, "approvals.0.id"),
					resource.TestCheckResourceAttrPair(filterByBizId, "approvals.0.name_ch",
						"huaweicloud_dataarts_architecture_subject.test", "name"),
					resource.TestCheckResourceAttrPair(filterByBizId, "approvals.0.name_en",
						"huaweicloud_dataarts_architecture_subject.test", "name"),
					resource.TestCheckResourceAttrPair(filterByBizId, "approvals.0.biz_id",
						"huaweicloud_dataarts_architecture_subject.test", "id"),
					resource.TestCheckResourceAttrSet(filterByBizId, "approvals.0.biz_type"),
					resource.TestCheckResourceAttrSet(filterByBizId, "approvals.0.biz_info"),
					resource.TestCheckResourceAttrSet(filterByBizId, "approvals.0.biz_status"),
					resource.TestCheckResourceAttrSet(filterByBizId, "approvals.0.approval_status"),
					resource.TestCheckResourceAttrSet(filterByBizId, "approvals.0.approval_type"),
					resource.TestCheckResourceAttrSet(filterByBizId, "approvals.0.submit_time"),
					resource.TestCheckResourceAttrPair(filterByBizId, "approvals.0.create_by",
						"huaweicloud_dataarts_architecture_subject.test", "created_by"),
					resource.TestCheckResourceAttrSet(filterByBizId, "approvals.0.approval_time"),
					resource.TestCheckResourceAttrSet(filterByBizId, "approvals.0.approver"),
					resource.TestCheckResourceAttrSet(filterByBizId, "approvals.0.msg"),

					// filter by name
					dcFilterByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),

					// filter by create by
					dcFilterByCreateBy.CheckResourceExists(),
					resource.TestCheckOutput("is_create_by_filter_useful", "true"),

					// filter by approver
					dcFilterByApprover.CheckResourceExists(),
					resource.TestCheckOutput("is_approver_filter_useful", "true"),

					// filter by approval status
					dcFilterByApprovalStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_approval_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataArchitectureApprovals_nonExistentWorkspace() string {
	randUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
data "huaweicloud_dataarts_architecture_approvals" "test" {
  workspace_id = "%[1]s"
}
`, randUUID.String())
}

func testAccDataSourceArchitectureApprovals_basic_base(name string) string {
	return fmt.Sprintf(` 
resource "huaweicloud_dataarts_architecture_subject" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  code         = "%[2]s"
  owner        = "tf"
  level        = 1
  description  = "Created by Terraform."
}

resource "huaweicloud_dataarts_architecture_batch_publish" "test" {
  workspace_id       = "%[1]s"
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

resource "huaweicloud_dataarts_architecture_batch_unpublish" "test" {
  depends_on = [
    huaweicloud_dataarts_architecture_batch_publish.test,
  ]

  workspace_id       = "%[1]s"
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
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name, acceptance.HW_DATAARTS_ARCHITECTURE_USER_ID, acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME)
}

func testAccDataSourceArchitectureApprovals_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dataarts_architecture_approvals" "all" {
  depends_on = [
    huaweicloud_dataarts_architecture_batch_unpublish.test,
  ]

  workspace_id = "%[2]s"
}

# Filter by biz ID
locals {
  biz_id = huaweicloud_dataarts_architecture_subject.test.id
}

data "huaweicloud_dataarts_architecture_approvals" "filter_by_biz_id" {
  depends_on = [
    huaweicloud_dataarts_architecture_batch_unpublish.test,
  ]

  workspace_id = "%[2]s"
  biz_id       = local.biz_id
}

locals {
  biz_id_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_approvals.filter_by_biz_id.approvals : v.biz_id == local.biz_id
  ]
}

output "is_biz_id_filter_useful" {
  value = length(local.biz_id_filter_result) > 0 && alltrue(local.biz_id_filter_result)
}

# Filter by name
locals {
  name = huaweicloud_dataarts_architecture_subject.test.name
}

data "huaweicloud_dataarts_architecture_approvals" "filter_by_name" {
  depends_on = [
    huaweicloud_dataarts_architecture_batch_unpublish.test,
  ]

  workspace_id = "%[2]s"
  name         = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_approvals.filter_by_name.approvals : v.name_ch == local.name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by create by
locals {
  create_by = huaweicloud_dataarts_architecture_subject.test.created_by
}

data "huaweicloud_dataarts_architecture_approvals" "filter_by_create_by" {
  depends_on = [
    huaweicloud_dataarts_architecture_batch_unpublish.test,
  ]

  workspace_id = "%[2]s"
  create_by    = local.create_by
}

locals {
  create_by_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_approvals.filter_by_create_by.approvals : v.create_by == local.create_by
  ]
}

output "is_create_by_filter_useful" {
  value = length(local.create_by_filter_result) > 0 && alltrue(local.create_by_filter_result)
}

# Filter by approver
locals {
  approver = huaweicloud_dataarts_architecture_batch_publish.test.approver_user_name
}

data "huaweicloud_dataarts_architecture_approvals" "filter_by_approver" {
  depends_on = [
    huaweicloud_dataarts_architecture_batch_unpublish.test,
  ]

  workspace_id = "%[2]s"
  approver     = local.approver
}

locals {
  approver_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_approvals.filter_by_approver.approvals : v.approver == local.approver
  ]
}

output "is_approver_filter_useful" {
  value = length(local.approver_filter_result) > 0 && alltrue(local.approver_filter_result)
}

# Filter by approval status
locals {
  approval_status = "FINISHED"
}

data "huaweicloud_dataarts_architecture_approvals" "filter_by_approval_status" {
  depends_on = [
    huaweicloud_dataarts_architecture_batch_unpublish.test,
  ]

  workspace_id 	  = "%[2]s"
  approval_status = local.approval_status
}

locals {
  approval_status_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_approvals.filter_by_approval_status.approvals : v.approval_status == "APPROVED"
  ]
}

output "is_approval_status_filter_useful" {
  value = length(local.approval_status_filter_result) > 0 && alltrue(local.approval_status_filter_result)
}
`, testAccDataSourceArchitectureApprovals_basic_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}

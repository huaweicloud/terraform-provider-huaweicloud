package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceArchitectureBatchPublishment_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsArchitectureReviewer(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccArchitectureBatchPublishment_basic(name),
			},
			{
				Config:      testAccArchitectureBatchPublishment_expectErr(name),
				ExpectError: regexp.MustCompile(`The definition has been released and does not support resubmission`),
			},
		},
	})
}

func testAccArchitectureBatchPublishment_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_subject" "test" {
  count        = 2
  workspace_id = "%[1]s"
  name         = "%[2]s_${count.index}"
  code         = "%[2]s_${count.index}"
  owner        = "%[2]s"
  level        = 1
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testAccArchitectureBatchPublishment_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_architecture_batch_publishment" "test" {
  workspace_id       = "%[2]s"
  approver_user_id   = "%[3]s"
  approver_user_name = "%[4]s"

  dynamic "biz_infos" {
    for_each = huaweicloud_dataarts_architecture_subject.test[*].id
    content {
      biz_id   = biz_infos.value
      biz_type = "SUBJECT"
    }
  }
}
`, testAccArchitectureBatchPublishment_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID,
		acceptance.HW_DATAARTS_ARCHITECTURE_USER_ID, acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME)
}

func testAccArchitectureBatchPublishment_expectErr(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_architecture_batch_publishment" "test_error" {
  workspace_id       = "%[2]s"
  approver_user_id   = "%[3]s"
  approver_user_name = "%[4]s"

  biz_infos {
    biz_id   = huaweicloud_dataarts_architecture_subject.test[0].id
    biz_type = "SUBJECT"
  }
}
`, testAccArchitectureBatchPublishment_basic(name), acceptance.HW_DATAARTS_WORKSPACE_ID,
		acceptance.HW_DATAARTS_ARCHITECTURE_USER_ID, acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME)
}

func TestAccResourceArchitectureBatchPublishment_action(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsArchitectureReviewer(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			// Publish the subjects.
			{
				Config: testAccArchitectureBatchPublishment_action_step1(name),
			},
			// Publish a subject that has already been published, assert an error message.
			{
				Config:      testAccArchitectureBatchPublishment_action_step2(name),
				ExpectError: regexp.MustCompile(`The definition has been released and does not support resubmission`),
			},
			// Take subjects offline.
			{
				Config: testAccArchitectureBatchPublishment_action_step3(name),
			},
			// Take a subject offline that was offline, assert an error message.
			{
				Config:      testAccArchitectureBatchPublishment_action_step4(name),
				ExpectError: regexp.MustCompile(fmt.Sprintf("Definition %s_0  is not online", name)),
			},
		},
	})
}

func testAccArchitectureBatchPublishment_action_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_architecture_batch_publish" "test" {
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
`, testAccArchitectureBatchPublishment_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID,
		acceptance.HW_DATAARTS_ARCHITECTURE_USER_ID, acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME)
}

func testAccArchitectureBatchPublishment_action_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_architecture_batch_publish" "test_error" {
  workspace_id       = "%[2]s"
  approver_user_id   = "%[3]s"
  approver_user_name = "%[4]s"
  fast_approval      = true

  biz_infos {
    biz_id   = huaweicloud_dataarts_architecture_subject.test[0].id
    biz_type = "SUBJECT"
  }
}
`, testAccArchitectureBatchPublishment_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID,
		acceptance.HW_DATAARTS_ARCHITECTURE_USER_ID, acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME)
}

func testAccArchitectureBatchPublishment_action_step3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_architecture_batch_unpublish" "test" {
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
`, testAccArchitectureBatchPublishment_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID,
		acceptance.HW_DATAARTS_ARCHITECTURE_USER_ID, acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME)
}

func testAccArchitectureBatchPublishment_action_step4(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_architecture_batch_unpublish" "test_error" {
  workspace_id       = "%[2]s"
  approver_user_id   = "%[3]s"
  approver_user_name = "%[4]s"
  fast_approval      = true

  biz_infos {
    biz_id   = huaweicloud_dataarts_architecture_subject.test[0].id
    biz_type = "SUBJECT"
  }
}
`, testAccArchitectureBatchPublishment_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID,
		acceptance.HW_DATAARTS_ARCHITECTURE_USER_ID, acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME)
}

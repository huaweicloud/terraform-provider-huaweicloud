package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceArchitectureBatchPublishment_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsArchitectureReviewer(t)
			acceptance.TestAccPreCheckDataArtsArchitectureBatchPublishment(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccArchitectureBatchPublishment_basic(),
			},
			{
				Config:      testAccArchitectureBatchPublishment_expectErr(),
				ExpectError: regexp.MustCompile(`The definition has been released and does not support resubmission`),
			},
		},
	})
}

func testAccArchitectureBatchPublishment_base() string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_reviewer" "test" {
  workspace_id = "%[1]s"
  user_id      = "%[2]s"
  user_name    = "%[3]s"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_ARCHITECTURE_USER_ID,
		acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME)
}

func testAccArchitectureBatchPublishment_basic() string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dataarts_architecture_batch_publishment" "test" {
  workspace_id       = "%s"
  approver_user_id   = huaweicloud_dataarts_architecture_reviewer.test.user_id
  approver_user_name = huaweicloud_dataarts_architecture_reviewer.test.user_name
  fast_approval      = true

  biz_infos {
    biz_id   = "%s"
    biz_type = "SUBJECT"
  }
}
`, testAccArchitectureBatchPublishment_base(), acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_ARCHITECTURE_SUBJECT_BIZ_ID)
}

func testAccArchitectureBatchPublishment_expectErr() string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dataarts_architecture_batch_publishment" "test_error" {
  workspace_id       = "%s"
  approver_user_id   = huaweicloud_dataarts_architecture_reviewer.test.user_id
  approver_user_name = huaweicloud_dataarts_architecture_reviewer.test.user_name
  fast_approval      = true

  biz_infos {
    biz_id   = "%s"
    biz_type = "SUBJECT"
  }
}
`, testAccArchitectureBatchPublishment_base(), acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_ARCHITECTURE_SUBJECT_BIZ_ID)
}

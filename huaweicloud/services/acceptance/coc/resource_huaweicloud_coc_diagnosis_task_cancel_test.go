package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceCocDiagnosisTaskCancel_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocDiagnosisTaskID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCocDiagnosisTaskCancel_basic(),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testCocDiagnosisTaskCancel_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_coc_diagnosis_task_cancel" "test" {
  task_id = "%s"
}
`, acceptance.HW_COC_DIAGNOSIS_TASK_ID)
}

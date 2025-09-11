package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceCocDiagnosisTaskRetry_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocInstanceID(t)
			acceptance.TestAccPreCheckCocDiagnosisTaskID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCocDiagnosisTaskRetry_basic(),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testCocDiagnosisTaskRetry_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_coc_diagnosis_task_retry" "test" {
  task_id     = "%[1]s"
  instance_id = "%[2]s"
}
`, acceptance.HW_COC_DIAGNOSIS_TASK_ID, acceptance.HW_COC_INSTANCE_ID)
}

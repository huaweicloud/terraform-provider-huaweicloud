package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccObjectCompare_basic(t *testing.T) {
	resourceName := "huaweicloud_drs_object_compare.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testObjectCompare_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "job_id"),
					resource.TestCheckResourceAttr(resourceName, "compare_task_num", "2"),
				),
			},
		},
	})
}

func testObjectCompare_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_object_compare" "test" {
  job_id           = "%s"
  compare_task_num = 2
}
`, acceptance.HW_DRS_JOB_ID)
}

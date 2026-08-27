package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBatchSelectObjects_basic(t *testing.T) {
	resourceName := "huaweicloud_drs_batch_select_objects.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testBatchSelectObjects_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "all_counts"),
					resource.TestCheckResourceAttrSet(resourceName, "results.#"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.job_id"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.status"),
				),
			},
		},
	})
}

func testBatchSelectObjects_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_batch_select_objects" "test" {
  jobs {
    job_id   = "%s"
    selected = false
  }
}
`, acceptance.HW_DRS_JOB_ID)
}

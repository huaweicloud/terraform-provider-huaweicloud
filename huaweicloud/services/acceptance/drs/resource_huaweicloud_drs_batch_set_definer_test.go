package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBatchSetDefiner_basic(t *testing.T) {
	resourceName := "huaweicloud_drs_batch_set_definer.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testBatchSetDefiner_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "results.#"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.status"),
				),
			},
		},
	})
}

func testBatchSetDefiner_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_batch_set_definer" "test" {
  jobs {
    job_id          = "%s"
    replace_definer = true
  }
}
`, acceptance.HW_DRS_JOB_ID)
}

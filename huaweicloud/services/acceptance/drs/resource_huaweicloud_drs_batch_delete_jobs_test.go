package drs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceBatchDeleteJobs_basic(t *testing.T) {
	resourceName := "huaweicloud_drs_batch_delete_jobs.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobIds(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceBatchDeleteJobs_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "results.#"),
					resource.TestMatchResourceAttr(resourceName, "results.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.status"),
				),
			},
		},
	})
}

func testAccResourceBatchDeleteJobs_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_batch_delete_jobs" "test" {
  jobs = split(",", "%s")
}
`, acceptance.HW_DRS_JOB_IDS)
}

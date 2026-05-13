package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDrsJobClone_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_drs_job_clone.test"
		rName        = acceptance.RandomAccResourceName()
		expectedName = rName + "-copy"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDrsJobClone_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", expectedName),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttr(resourceName, "is_clone_job", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
				),
			},
		},
	})
}

func testAccDrsJobClone_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_job_clone" "test" {
  job_id = "%[1]s"
  name   = "%[2]s-copy"
}
`, acceptance.HW_DRS_JOB_ID, rName)
}

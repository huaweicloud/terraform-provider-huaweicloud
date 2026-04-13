package dli

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDliFlinkJobExport_basic(t *testing.T) {
	resourceName := "huaweicloud_dli_flink_job_export.test"

	// This resource is a one-time action resource, so only need to test the creation.
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliFlinkSQLOBSPath(t)
			acceptance.TestAccPreCheckDliFlinkSQLJobIds(t, 2)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccFlinkJobExport_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "job_ids.#", "1"),
				),
			},
			{
				Config: testAccFlinkJobExport_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "job_ids.#", "2"),
				),
			},
		},
	})
}

func testAccFlinkJobExport_basic() string {
	return fmt.Sprintf(`
# In this case, HW_DLI_FLINK_SQL_OBS_PATH no need to startwith "obs://"
resource "huaweicloud_dli_flink_job_export" "test" {
  obs_path = "%[1]s"
  job_ids  = slice([for id in split(",", "%[2]s") : tonumber(trimspace(id))], 0, 1)
}
`, acceptance.HW_DLI_FLINK_SQL_OBS_PATH, acceptance.HW_DLI_FLINK_SQL_JOB_IDS)
}

func testAccFlinkJobExport_basic_update() string {
	return fmt.Sprintf(`
# In this case, HW_DLI_FLINK_SQL_OBS_PATH no need to startwith "obs://"
resource "huaweicloud_dli_flink_job_export" "test" {
  obs_path = "%[1]s"
  job_ids  = slice([for id in split(",", "%[2]s") : tonumber(trimspace(id))], 0, 2)

  enable_force_new = "true"
}
`, acceptance.HW_DLI_FLINK_SQL_OBS_PATH, acceptance.HW_DLI_FLINK_SQL_JOB_IDS)
}

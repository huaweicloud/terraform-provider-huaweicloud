package dli

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDliFlinkJobImport_basic(t *testing.T) {
	resourceName := "huaweicloud_dli_flink_job_import.test"

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
				Config: testAccFlinkJobImport_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "obs_path", acceptance.HW_DLI_FLINK_SQL_OBS_PATH),
					resource.TestCheckResourceAttr(resourceName, "is_cover", "true"),
				),
			},
			{
				Config: testAccFlinkJobImport_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "obs_path", acceptance.HW_DLI_FLINK_SQL_OBS_PATH),
					resource.TestCheckResourceAttr(resourceName, "is_cover", "false"),
				),
			},
		},
	})
}

func testAccFlinkJobImport_base() string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_flink_job_export" "test" {
  obs_path = "%[1]s"
  job_ids  = [for id in split(",", "%[2]s") : tonumber(trimspace(id))]
}

`, acceptance.HW_DLI_FLINK_SQL_OBS_PATH, acceptance.HW_DLI_FLINK_SQL_JOB_IDS)
}

func testAccFlinkJobImport_basic() string {
	return fmt.Sprintf(`
%[1]s

# In this case, HW_DLI_FLINK_SQL_OBS_PATH no need to startwith "obs://"
resource "huaweicloud_dli_flink_job_import" "test" {
  obs_path = "%[2]s"
  is_cover = true

  depends_on = [
    huaweicloud_dli_flink_job_export.test
  ]
}
`, testAccFlinkJobImport_base(), acceptance.HW_DLI_FLINK_SQL_OBS_PATH)
}

func testAccFlinkJobImport_basic_update() string {
	return fmt.Sprintf(`
%[1]s

# In this case, HW_DLI_FLINK_SQL_OBS_PATH no need to startwith "obs://"
resource "huaweicloud_dli_flink_job_import" "test" {
  obs_path = "%[2]s"
  is_cover = false

  enable_force_new = "true"

  depends_on = [
    huaweicloud_dli_flink_job_export.test
  ]
}
`, testAccFlinkJobImport_base(), acceptance.HW_DLI_FLINK_SQL_OBS_PATH)
}

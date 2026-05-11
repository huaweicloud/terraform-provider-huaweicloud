package dli

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceFlinkSqlJobSavepointImport_basic(t *testing.T) {
	var (
		rName = "huaweicloud_dli_flinksql_job_savepoint_import.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliFlinkSQLJobIds(t, 2)
			acceptance.TestAccPreCheckDliFlinkSQLOBSPath(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccFlinkSqlJobSavepointImport_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "job_id"),
					resource.TestCheckResourceAttr(rName, "savepoint_path", acceptance.HW_DLI_FLINK_SQL_OBS_PATH),
				),
			},
			{
				Config: testAccFlinkSqlJobSavepointImport_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "job_id"),
					resource.TestCheckResourceAttr(rName, "savepoint_path", acceptance.HW_DLI_FLINK_SQL_OBS_PATH),
				),
			},
		},
	})
}

func testAccFlinkSqlJobSavepointImport_base() string {
	return fmt.Sprintf(`
locals {
  flink_sql_job_ids = split(",", "%[1]s")
}
`, acceptance.HW_DLI_FLINK_SQL_JOB_IDS)
}

func testAccFlinkSqlJobSavepointImport_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dli_flinksql_job_savepoint_import" "test" {
  job_id         = local.flink_sql_job_ids[0]
  savepoint_path = "%[2]s"
}
`, testAccFlinkSqlJobSavepointImport_base(), acceptance.HW_DLI_FLINK_SQL_OBS_PATH)
}

func testAccFlinkSqlJobSavepointImport_basic_update() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dli_flinksql_job_savepoint_import" "test" {
  job_id         = local.flink_sql_job_ids[1]
  savepoint_path = "%[2]s"

  enable_force_new = "true"
}
`, testAccFlinkSqlJobSavepointImport_base(), acceptance.HW_DLI_FLINK_SQL_OBS_PATH)
}

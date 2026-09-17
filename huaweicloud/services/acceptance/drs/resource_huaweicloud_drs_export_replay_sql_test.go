package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDrsExportReplaySql_basic(t *testing.T) {
	resourceName := "huaweicloud_drs_export_replay_sql.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDrsExportReplaySql_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "job_id", acceptance.HW_DRS_JOB_ID),
					resource.TestCheckResourceAttr(resourceName, "file_type", "abnormal_sql"),
				),
			},
		},
	})
}

func testAccDrsExportReplaySql_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_export_replay_sql" "test" {
  job_id    = "%s"
  file_type = "abnormal_sql"
}
`, acceptance.HW_DRS_JOB_ID)
}

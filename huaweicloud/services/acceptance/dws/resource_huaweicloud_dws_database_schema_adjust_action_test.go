package dws

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDatabaseSchemaAdjustAction_basic(t *testing.T) {
	var (
		permSpace = 20480
		rName     = "huaweicloud_dws_database_schema_adjust_action.test"
	)

	// Avoid CheckDestroy because this resource is a one-time action resource.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
			acceptance.TestAccPreCheckDwsSchemaAdjust(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataDatabaseSchemaAdjustAction_invalidCluster(),
				ExpectError: regexp.MustCompile("Cluster does not exist or has been deleted"),
			},
			{
				// HTTP error details often mention database/schema or a DWS code; some failures surface as non-zero ret_code.
				Config: testAccDataDatabaseSchemaAdjustAction_invalidDatabase(),
				ExpectError: regexp.MustCompile(
					`(?i)error adjusting database schema space:.+`),
			},
			{
				Config: testAccDataDatabaseSchemaAdjustAction_basic(permSpace),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "cluster_id", acceptance.HW_DWS_CLUSTER_ID),
					resource.TestCheckResourceAttr(rName, "database", acceptance.HW_DWS_GRANT_DATABASE_NAME),
					resource.TestCheckResourceAttr(rName, "schema", acceptance.HW_DWS_GRANT_SCHEMA_NAME),
					resource.TestCheckResourceAttr(rName, "perm_space", strconv.Itoa(permSpace)),
				),
			},
		},
	})
}

func testAccDataDatabaseSchemaAdjustAction_basic(permSpace int) string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_database_schema_adjust_action" "test" {
  cluster_id = "%[1]s"
  database   = "%[2]s"
  schema     = "%[3]s"
  perm_space = %[4]d
}
`, acceptance.HW_DWS_CLUSTER_ID,
		acceptance.HW_DWS_GRANT_DATABASE_NAME,
		acceptance.HW_DWS_GRANT_SCHEMA_NAME,
		permSpace)
}

func testAccDataDatabaseSchemaAdjustAction_invalidCluster() string {
	randomUUID, _ := uuid.NewRandom()
	return fmt.Sprintf(`
resource "huaweicloud_dws_database_schema_adjust_action" "invalid_cluster" {
  cluster_id = "%[1]s"
  database   = "gaussdb"
  schema     = "nonexistent_schema_for_acc_test"
  perm_space = 10240
}
`, randomUUID.String())
}

func testAccDataDatabaseSchemaAdjustAction_invalidDatabase() string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_database_schema_adjust_action" "invalid_database" {
  cluster_id = "%[1]s"
  database   = "database_name_not_exist_for_acc_test"
  schema     = "%[2]s"
  perm_space = 10240
}
`, acceptance.HW_DWS_CLUSTER_ID,
		acceptance.HW_DWS_GRANT_SCHEMA_NAME)
}

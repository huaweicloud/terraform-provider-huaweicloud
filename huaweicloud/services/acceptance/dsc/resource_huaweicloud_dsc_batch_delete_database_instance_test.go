package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceBatchDeleteDatabaseInstance_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDscDbId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceBatchDeleteDatabaseInstance_basic(),
			},
		},
	})
}

func testAccResourceBatchDeleteDatabaseInstance_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dsc_batch_delete_database_instance" "test" {
  db_ids = split(",", "%s")
}
`, acceptance.HW_DSC_DB_ID)
}

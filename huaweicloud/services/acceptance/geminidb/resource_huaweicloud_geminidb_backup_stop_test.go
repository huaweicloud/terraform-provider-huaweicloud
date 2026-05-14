package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGeminiDBBackupStop_basic(t *testing.T) {
	resourceName := "huaweicloud_geminidb_backup_stop.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBBackupStop_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "backup_id"),
				),
			},
		},
	})
}

func testAccGeminiDBBackupStop_basic() string {
	return fmt.Sprintf(`

resource "huaweicloud_geminidb_backup_stop" "test" {
  backup_id = "%s"
}
`, acceptance.HW_GEMINIDB_BACKUP_ID)
}

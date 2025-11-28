package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAppServerBatchMigrate_basic(t *testing.T) {
	resourceName := "huaweicloud_workspace_app_server_batch_migrate.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerBatchMigrate(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServerBatchMigrate_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "server_ids.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(resourceName, "host_id"),
				),
			},
		},
	})
}

func testAccAppServerBatchMigrate_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_server_batch_migrate" "test" {
  server_ids = split(",", "%[1]s")
  host_id    = "%[2]s"
}
`, acceptance.HW_WORKSPACE_APP_SERVER_BATCH_MIGRATE_SERVER_IDS, acceptance.HW_WORKSPACE_APP_SERVER_BATCH_MIGRATE_HOST_ID)
}

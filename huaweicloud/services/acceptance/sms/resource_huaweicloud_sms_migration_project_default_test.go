package sms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceSmsMigrateProjectDefault_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testSmsMigrateProjectDefault_basic(name),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testSmsMigrateProjectDefault_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_sms_migration_project" "test" {
  name          = "%[1]s"
  region        = "cn-north-4"
  use_public_ip = false
  exist_server  = false
  type          = "MIGRATE_FILE"
  syncing       = false
  is_default    = true

  lifecycle {
    ignore_changes = [
      is_default
    ]
  }
}

data "huaweicloud_sms_migration_projects" "test" {}

locals {
  system_mig_project_id = [for v in data.huaweicloud_sms_migration_projects.test.migprojects: v.id if v.name == "SystemProject"][0]
}

resource "huaweicloud_sms_migration_project_default" "test" {
  mig_project_id = local.system_mig_project_id

  depends_on = [huaweicloud_sms_migration_project.test]

  lifecycle {
    ignore_changes = [
      mig_project_id
    ]
  }
}
`, name)
}

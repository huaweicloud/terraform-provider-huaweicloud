package cbr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceVaultMigrateResources_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceVaultMigrateResources_basic(),
			},
		},
	})
}

func testResourceVaultMigrateResources_base() string {
	name := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  name              = "%[1]s"
  size              = 50
  description       = "test description"
  volume_type       = "GPSSD"
  multiattach       = false
  charging_mode     = "postPaid"
}

resource "huaweicloud_cbr_vault" "test1" {
  name            = "%[1]s_1"
  type            = "disk"
  protection_type = "backup"
  size            = 100

  resources {
    includes = [huaweicloud_evs_volume.test.id]
  }

  lifecycle {
    ignore_changes = [
      resources,
    ]
  }
}

resource "huaweicloud_cbr_vault" "test2" {
  name            = "%[1]s_2"
  type            = "disk"
  protection_type = "backup"
  size            = 100
}
`, name)
}

func testResourceVaultMigrateResources_basic() string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cbr_vault_migrate_resources" "test" {
  vault_id             = huaweicloud_cbr_vault.test1.id
  destination_vault_id = huaweicloud_cbr_vault.test2.id
  resource_ids         = [huaweicloud_evs_volume.test.id]
}`, testResourceVaultMigrateResources_base())
}

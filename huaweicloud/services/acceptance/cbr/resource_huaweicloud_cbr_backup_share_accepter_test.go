package cbr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cbr/v3/backups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getBackupShareAccepterResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.CbrV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CBR v3 client: %s", err)
	}

	return backups.Get(client, state.Primary.ID)
}

func TestAccBackupShareAccepter_basic(t *testing.T) {
	var (
		obj          *backups.BackupResp
		name         = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_cbr_backup_share_accepter.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getBackupShareAccepterResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAcceptBackup(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccBackupShareAccepter_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "backup_id", acceptance.HW_SHARED_BACKUP_ID),
					resource.TestCheckResourceAttrPair(resourceName, "vault_id", "huaweicloud_cbr_vault.test", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBackupShareAccepter_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_vault" "test" {
  name             = "%[1]s"
  type             = "server"
  consistent_level = "crash_consistent"
  protection_type  = "backup"
  size             = 100
}

resource "huaweicloud_cbr_backup_share_accepter" "test" {
  backup_id = "%[2]s"
  vault_id  = huaweicloud_cbr_vault.test.id
}
`, name, acceptance.HW_SHARED_BACKUP_ID)
}

package deprecated

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/vbs/v2/backups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccVBSBackupV2_basic(t *testing.T) {
	var config backups.Backup
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_vbs_backup.backup_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVBSBackupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVBSBackupV2_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVBSBackupV2Exists(resourceName, &config),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
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

func testAccCheckVBSBackupV2Destroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	vbsClient, err := config.VbsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud vbs client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_vbs_backup" {
			continue
		}

		_, err := backups.Get(vbsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("VBS backup still exists")
		}
	}

	return nil
}

func testAccCheckVBSBackupV2Exists(n string, configs *backups.Backup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		vbsClient, err := config.VbsV2Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating Huaweicloud vbs client: %s", err)
		}

		found, err := backups.Get(vbsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmtp.Errorf("VBS backup not found")
		}

		*configs = *found

		return nil
	}
}

func testAccVBSBackupV2_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "zones" {}

resource "huaweicloud_evs_volume" "volume" {
  name              = "%[1]s"
  description       = "my volume"
  volume_type       = "SAS"
  size              = 20
  availability_zone = data.huaweicloud_availability_zones.zones.names[0]
}
  
resource "huaweicloud_evs_snapshot" "snapshot_1" {
  name        = "%[1]s"
  description = "for vbs backup"
  volume_id   = huaweicloud_evs_volume.volume.id
}
 
resource "huaweicloud_vbs_backup" "backup_1" {
  volume_id   = huaweicloud_evs_volume.volume.id
  snapshot_id = huaweicloud_evs_snapshot.snapshot_1.id
  name        = "%[1]s"
}
`, rName)
}

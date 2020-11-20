package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/vbs/v2/backups"
)

func TestAccVBSBackupV2_basic(t *testing.T) {
	var config backups.Backup
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVBSBackupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVBSBackupV2_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVBSBackupV2Exists("huaweicloud_vbs_backup.backup_1", &config),
					resource.TestCheckResourceAttr(
						"huaweicloud_vbs_backup.backup_1", "name", rName),
					resource.TestCheckResourceAttr(
						"huaweicloud_vbs_backup.backup_1", "status", "available"),
				),
			},
			{
				ResourceName:      "huaweicloud_vbs_backup.backup_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVBSBackupV2_timeout(t *testing.T) {
	var config backups.Backup
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVBSBackupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVBSBackupV2_timeout(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVBSBackupV2Exists("huaweicloud_vbs_backup.backup_1", &config),
				),
			},
		},
	})
}

func testAccCheckVBSBackupV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	vbsClient, err := config.vbsV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud vbs client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_vbs_backup" {
			continue
		}

		_, err := backups.Get(vbsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("VBS backup still exists")
		}
	}

	return nil
}

func testAccCheckVBSBackupV2Exists(n string, configs *backups.Backup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		vbsClient, err := config.vbsV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating Huaweicloud vbs client: %s", err)
		}

		found, err := backups.Get(vbsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("VBS backup not found")
		}

		*configs = *found

		return nil
	}
}

func testAccVBSBackupV2_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_evs_volume" "volume" {
  name              = "%s"
  description       = "my volume"
  volume_type       = "SAS"
  size              = 20
  availability_zone = "%s"
}
  
resource "huaweicloud_evs_snapshot" "snapshot_1" {
  name        = "%s"
  description = "for vbs backup"
  volume_id   = huaweicloud_evs_volume.volume.id
}
  
resource "huaweicloud_vbs_backup" "backup_1" {
  volume_id   = huaweicloud_evs_volume.volume.id
  snapshot_id = huaweicloud_evs_snapshot.snapshot_1.id
  name        = "%s"
}
`, rName, OS_AVAILABILITY_ZONE, rName, rName)
}

func testAccVBSBackupV2_timeout(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_evs_volume" "volume" {
  name              = "%s"
  description       = "my volume"
  volume_type       = "SAS"
  size              = 20
  availability_zone = "%s"
}

resource "huaweicloud_evs_snapshot" "snapshot_1" {
  name        = "%s"
  description = "for vbs backup"
  volume_id   = huaweicloud_evs_volume.volume.id
}

resource "huaweicloud_vbs_backup" "backup_1" {
  volume_id   = huaweicloud_evs_volume.volume.id
  snapshot_id = huaweicloud_evs_snapshot.snapshot_1.id
  name        = "%s"
  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`, rName, OS_AVAILABILITY_ZONE, rName, rName)
}

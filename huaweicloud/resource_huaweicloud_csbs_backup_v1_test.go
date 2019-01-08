package huaweicloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"testing"

	"github.com/huaweicloud/golangsdk/openstack/csbs/v1/backup"
)

func TestAccCSBSBackupV1_basic(t *testing.T) {
	var backups backup.Backup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSBSBackupV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCSBSBackupV1Exists("huaweicloud_csbs_backup_v1.csbs", &backups),
					resource.TestCheckResourceAttr(
						"huaweicloud_csbs_backup_v1.csbs", "backup_name", "csbs-test1"),
					resource.TestCheckResourceAttr(
						"huaweicloud_csbs_backup_v1.csbs", "resource_type", "OS::Nova::Server"),
				),
			},
		},
	})
}

func TestAccCSBSBackupV1_timeout(t *testing.T) {
	var backups backup.Backup

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSBSBackupV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupV1_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCSBSBackupV1Exists("huaweicloud_csbs_backup_v1.csbs", &backups),
				),
			},
		},
	})
}

func testAccCSBSBackupV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	backupClient, err := config.csbsV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating csbs client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_csbs_backup_v1" {
			continue
		}

		_, err := backup.Get(backupClient, rs.Primary.ID).ExtractBackup()
		if err == nil {
			return fmt.Errorf("Backup still exists")
		}
	}

	return nil
}

func testAccCSBSBackupV1Exists(n string, backups *backup.Backup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Backup not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		backupClient, err := config.csbsV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating csbs client: %s", err)
		}

		found, err := backup.Get(backupClient, rs.Primary.ID).ExtractBackup()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("backup not found")
		}

		*backups = *found

		return nil
	}
}

var testAccCSBSBackupV1_basic = fmt.Sprintf(`
resource "huaweicloud_compute_instance_v2" "instance_1" {
  name = "instance_1"
  image_id = "%s"
  security_groups = ["default"]
  availability_zone = "%s"
  flavor_id = "%s"
  metadata {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
resource "huaweicloud_csbs_backup_v1" "csbs" {
  backup_name      = "csbs-test1"
  description      = "test-code"
  resource_id = "${huaweicloud_compute_instance_v2.instance_1.id}"
  resource_type = "OS::Nova::Server"
}
`, OS_IMAGE_ID, OS_AVAILABILITY_ZONE, OS_FLAVOR_ID, OS_NETWORK_ID)

var testAccCSBSBackupV1_timeout = fmt.Sprintf(`
resource "huaweicloud_compute_instance_v2" "instance_1" {
  name = "instance_1"
  image_id = "%s"
  security_groups = ["default"]
  availability_zone = "%s"
  flavor_id = "%s"
  metadata {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
resource "huaweicloud_csbs_backup_v1" "csbs" {
  backup_name      = "csbs-test1"
  description      = "test-code"
  resource_id = "${huaweicloud_compute_instance_v2.instance_1.id}"
  resource_type = "OS::Nova::Server"
}
`, OS_IMAGE_ID, OS_AVAILABILITY_ZONE, OS_FLAVOR_ID, OS_NETWORK_ID)

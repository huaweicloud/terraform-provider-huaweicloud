package huaweicloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"testing"

	"github.com/huaweicloud/golangsdk/openstack/csbs/v1/backup"
)

func TestAccCSBSBackupV1_basic(t *testing.T) {
	var backups backup.Backup
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSBSBackupV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupV1_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCSBSBackupV1Exists("huaweicloud_csbs_backup.csbs", &backups),
					resource.TestCheckResourceAttr(
						"huaweicloud_csbs_backup.csbs", "backup_name", rName),
					resource.TestCheckResourceAttr(
						"huaweicloud_csbs_backup.csbs", "resource_type", "OS::Nova::Server"),
				),
			},
			{
				ResourceName:      "huaweicloud_csbs_backup.csbs",
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccCSBSBackupV1_timeout(t *testing.T) {
	var backups backup.Backup
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSBSBackupV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupV1_timeout(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCSBSBackupV1Exists("huaweicloud_csbs_backup.csbs", &backups),
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
		if rs.Type != "huaweicloud_csbs_backup" {
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

func testAccCSBSBackupV1_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_compute_instance_v2" "instance_1" {
  name              = "%s"
  image_id          = "%s"
  security_groups   = ["default"]
  availability_zone = "%s"
  flavor_id         = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
resource "huaweicloud_csbs_backup" "csbs" {
  backup_name      = "%s"
  description      = "test-code"
  resource_id      = huaweicloud_compute_instance_v2.instance_1.id
  resource_type    = "OS::Nova::Server"
}
`, rName, OS_IMAGE_ID, OS_AVAILABILITY_ZONE, OS_FLAVOR_ID, OS_NETWORK_ID, rName)
}

func testAccCSBSBackupV1_timeout(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_compute_instance_v2" "instance_1" {
  name              = "%s"
  image_id          = "%s"
  security_groups   = ["default"]
  availability_zone = "%s"
  flavor_id         = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
resource "huaweicloud_csbs_backup" "csbs" {
  backup_name      = "%s"
  description      = "test-code"
  resource_id      = huaweicloud_compute_instance_v2.instance_1.id
  resource_type    = "OS::Nova::Server"
}
`, rName, OS_IMAGE_ID, OS_AVAILABILITY_ZONE, OS_FLAVOR_ID, OS_NETWORK_ID, rName)
}

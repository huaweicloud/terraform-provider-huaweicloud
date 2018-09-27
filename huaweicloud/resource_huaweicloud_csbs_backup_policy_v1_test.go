package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/huaweicloud/golangsdk/openstack/csbs/v1/policies"
)

func TestAccCSBSBackupPolicyV1_basic(t *testing.T) {
	var policy policies.BackupPolicy

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCSBSBackupPolicyV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCSBSBackupPolicyV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupPolicyV1Exists("huaweicloud_csbs_backup_policy_v1.backup_policy_v1", &policy),
					resource.TestCheckResourceAttr(
						"huaweicloud_csbs_backup_policy_v1.backup_policy_v1", "name", "backup-policy"),
					resource.TestCheckResourceAttr(
						"huaweicloud_csbs_backup_policy_v1.backup_policy_v1", "status", "suspended"),
				),
			},
			resource.TestStep{
				Config: testAccCSBSBackupPolicyV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupPolicyV1Exists("huaweicloud_csbs_backup_policy_v1.backup_policy_v1", &policy),
					resource.TestCheckResourceAttr(
						"huaweicloud_csbs_backup_policy_v1.backup_policy_v1", "name", "backup-policy-update"),
				),
			},
		},
	})
}

func TestAccCSBSBackupPolicyV1_timeout(t *testing.T) {
	var policy policies.BackupPolicy

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCSBSBackupPolicyV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCSBSBackupPolicyV1_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupPolicyV1Exists("huaweicloud_csbs_backup_policy_v1.backup_policy_v1", &policy),
				),
			},
		},
	})
}

func testAccCheckCSBSBackupPolicyV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	policyClient, err := config.csbsV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating csbs client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_csbs_backup_policy_v1" {
			continue
		}

		_, err := policies.Get(policyClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("backup policy still exists")
		}
	}

	return nil
}

func testAccCheckCSBSBackupPolicyV1Exists(n string, policy *policies.BackupPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		policyClient, err := config.csbsV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating csbs client: %s", err)
		}

		found, err := policies.Get(policyClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("backup policy not found")
		}

		*policy = *found

		return nil
	}
}

var testAccCSBSBackupPolicyV1_basic = fmt.Sprintf(`
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
resource "huaweicloud_csbs_backup_policy_v1" "backup_policy_v1" {
	name            = "backup-policy"
  	resource {
      id = "${huaweicloud_compute_instance_v2.instance_1.id}"
      type = "OS::Nova::Server"
      name = "resource4"
  	}
  	scheduled_operation {
      name ="mybackup"
      enabled = true
      operation_type ="backup"
      max_backups = "2"
      trigger_pattern = "BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\nRRULE:FREQ=WEEKLY;BYDAY=TH;BYHOUR=12;BYMINUTE=27\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n"
  	}
}
`, OS_IMAGE_ID, OS_AVAILABILITY_ZONE, OS_FLAVOR_ID, OS_NETWORK_ID)

var testAccCSBSBackupPolicyV1_update = fmt.Sprintf(`
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
resource "huaweicloud_csbs_backup_policy_v1" "backup_policy_v1" {
	name            = "backup-policy-update"
  	resource {
      id = "${huaweicloud_compute_instance_v2.instance_1.id}"
      type = "OS::Nova::Server"
      name = "resource4"
  	}
  	scheduled_operation {
      name ="mybackup"
      enabled = true
      operation_type ="backup"
      max_backups = "2"
      trigger_pattern = "BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\nRRULE:FREQ=WEEKLY;BYDAY=TH;BYHOUR=12;BYMINUTE=27\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n"
  	}
}
`, OS_IMAGE_ID, OS_AVAILABILITY_ZONE, OS_FLAVOR_ID, OS_NETWORK_ID)

var testAccCSBSBackupPolicyV1_timeout = fmt.Sprintf(`
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
resource "huaweicloud_csbs_backup_policy_v1" "backup_policy_v1" {
	name            = "backup-policy"
  	resource {
      id = "${huaweicloud_compute_instance_v2.instance_1.id}"
      type = "OS::Nova::Server"
      name = "resource4"
  	}
  	scheduled_operation {
      name ="mybackup"
      enabled = true
      operation_type ="backup"
      max_backups = "2"
      trigger_pattern = "BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\nRRULE:FREQ=WEEKLY;BYDAY=TH;BYHOUR=12;BYMINUTE=27\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n"
  	}

	timeouts {
    create = "5m"
    delete = "5m"
  }
}
`, OS_IMAGE_ID, OS_AVAILABILITY_ZONE, OS_FLAVOR_ID, OS_NETWORK_ID)

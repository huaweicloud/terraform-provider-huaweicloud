package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/csbs/v1/policies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccCSBSBackupPolicyV1_basic(t *testing.T) {
	var policy policies.BackupPolicy
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	updateName := fmt.Sprintf("tf-acc-test-update-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCSBSBackupPolicyV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupPolicyV1_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupPolicyV1Exists("huaweicloud_csbs_backup_policy.backup_policy", &policy),
					resource.TestCheckResourceAttr(
						"huaweicloud_csbs_backup_policy.backup_policy", "name", rName),
					resource.TestCheckResourceAttr(
						"huaweicloud_csbs_backup_policy.backup_policy", "status", "suspended"),
				),
			},
			{
				ResourceName:      "huaweicloud_csbs_backup_policy.backup_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCSBSBackupPolicyV1_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupPolicyV1Exists("huaweicloud_csbs_backup_policy.backup_policy", &policy),
					resource.TestCheckResourceAttr(
						"huaweicloud_csbs_backup_policy.backup_policy", "name", updateName),
				),
			},
		},
	})
}

func TestAccCSBSBackupPolicyV1_timeout(t *testing.T) {
	var policy policies.BackupPolicy
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCSBSBackupPolicyV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSBSBackupPolicyV1_timeout(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCSBSBackupPolicyV1Exists("huaweicloud_csbs_backup_policy.backup_policy", &policy),
				),
			},
		},
	})
}

func testAccCheckCSBSBackupPolicyV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	policyClient, err := config.CsbsV1Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating csbs client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_csbs_backup_policy" {
			continue
		}

		_, err := policies.Get(policyClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("backup policy still exists")
		}
	}

	return nil
}

func testAccCheckCSBSBackupPolicyV1Exists(n string, policy *policies.BackupPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		policyClient, err := config.CsbsV1Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating csbs client: %s", err)
		}

		found, err := policies.Get(policyClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("backup policy not found")
		}

		*policy = *found

		return nil
	}
}

func testAccCSBSBackupPolicyV1_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

resource "huaweicloud_compute_instance_v2" "instance_1" {
  name               = "%s"
  image_id           = "%s"
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone  = "%s"
  flavor_id          = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
resource "huaweicloud_csbs_backup_policy" "backup_policy" {
  name = "%s"
  resource {
    id   = huaweicloud_compute_instance_v2.instance_1.id
    type = "OS::Nova::Server"
    name = "resource4"
  }
  scheduled_operation {
    name            ="mybackup"
    enabled         = true
    operation_type  ="backup"
    max_backups     = "2"
    trigger_pattern = "BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\nRRULE:FREQ=WEEKLY;BYDAY=TH;BYHOUR=12;BYMINUTE=27\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n"
  }
}
`, rName, HW_IMAGE_ID, HW_AVAILABILITY_ZONE, HW_FLAVOR_ID, HW_NETWORK_ID, rName)
}

func testAccCSBSBackupPolicyV1_update(rName, updateName string) string {
	return fmt.Sprintf(`
data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

resource "huaweicloud_compute_instance_v2" "instance_1" {
  name               = "%s"
  image_id           = "%s"
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone  = "%s"
  flavor_id          = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
resource "huaweicloud_csbs_backup_policy" "backup_policy" {
  name = "%s"
  resource {
    id   = huaweicloud_compute_instance_v2.instance_1.id
    type = "OS::Nova::Server"
    name = "resource4"
  }
  scheduled_operation {
    name            ="mybackup"
    enabled         = true
    operation_type  ="backup"
    max_backups     = "2"
    trigger_pattern = "BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\nRRULE:FREQ=WEEKLY;BYDAY=TH;BYHOUR=12;BYMINUTE=27\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n"
  }
}
`, rName, HW_IMAGE_ID, HW_AVAILABILITY_ZONE, HW_FLAVOR_ID, HW_NETWORK_ID, updateName)
}

func testAccCSBSBackupPolicyV1_timeout(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

resource "huaweicloud_compute_instance_v2" "instance_1" {
  name               = "%s"
  image_id           = "%s"
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone  = "%s"
  flavor_id          = "%s"
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
resource "huaweicloud_csbs_backup_policy" "backup_policy" {
  name = "%s"
  resource {
    id   = huaweicloud_compute_instance_v2.instance_1.id
    type = "OS::Nova::Server"
    name = "resource4"
  }
  scheduled_operation {
    name            ="mybackup"
    enabled         = true
    operation_type  ="backup"
    max_backups     = "2"
    trigger_pattern = "BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\nRRULE:FREQ=WEEKLY;BYDAY=TH;BYHOUR=12;BYMINUTE=27\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n"
  }

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`, rName, HW_IMAGE_ID, HW_AVAILABILITY_ZONE, HW_FLAVOR_ID, HW_NETWORK_ID, rName)
}

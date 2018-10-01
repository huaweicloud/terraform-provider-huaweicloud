package huaweicloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"testing"

	"github.com/huaweicloud/golangsdk/openstack/vbs/v2/policies"
)

func TestAccVBSBackupPolicyV2_basic(t *testing.T) {
	var policy policies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckRequiredEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccVBSBackupPolicyV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVBSBackupPolicyV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccVBSBackupPolicyV2Exists("huaweicloud_vbs_backup_policy_v2.vbs", &policy),
					resource.TestCheckResourceAttr(
						"huaweicloud_vbs_backup_policy_v2.vbs", "name", "policy_001"),
					resource.TestCheckResourceAttr(
						"huaweicloud_vbs_backup_policy_v2.vbs", "status", "ON"),
				),
			},
			resource.TestStep{
				Config: testAccVBSBackupPolicyV2_update,
				Check: resource.ComposeTestCheckFunc(
					testAccVBSBackupPolicyV2Exists("huaweicloud_vbs_backup_policy_v2.vbs", &policy),
					resource.TestCheckResourceAttr(
						"huaweicloud_vbs_backup_policy_v2.vbs", "name", "policy_002"),
					resource.TestCheckResourceAttr(
						"huaweicloud_vbs_backup_policy_v2.vbs", "status", "ON"),
				),
			},
		},
	})
}

func TestAccVBSBackupPolicyV2_timeout(t *testing.T) {
	var policy policies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccVBSBackupPolicyV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVBSBackupPolicyV2_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccVBSBackupPolicyV2Exists("huaweicloud_vbs_backup_policy_v2.vbs", &policy),
				),
			},
		},
	})
}

func testAccVBSBackupPolicyV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	vbsClient, err := config.vbsV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating huaweicloud sfs client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_vbs_backup_policy_v2" {
			continue
		}

		_, err := policies.List(vbsClient, policies.ListOpts{ID: rs.Primary.ID})
		if err != nil {
			return fmt.Errorf("Backup Policy still exists")
		}
	}

	return nil
}

func testAccVBSBackupPolicyV2Exists(n string, policy *policies.Policy) resource.TestCheckFunc {
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
			return fmt.Errorf("Error creating huaweicloud vbs client: %s", err)
		}

		policyList, err := policies.List(vbsClient, policies.ListOpts{ID: rs.Primary.ID})
		if err != nil {
			return err
		}
		found := policyList[0]
		if found.ID != rs.Primary.ID {
			return fmt.Errorf("backup policy not found")
		}

		*policy = found

		return nil
	}
}

var testAccVBSBackupPolicyV2_basic = fmt.Sprintf(`
resource "huaweicloud_vbs_backup_policy_v2" "vbs" {
  name = "policy_001"
  start_time  = "12:00"
  status  = "ON"
  retain_first_backup = "N"
  rentention_num = 2
  frequency = 1
      tags =[
        {
          key = "k2"
          value = "v2"
          }]
}
`)

var testAccVBSBackupPolicyV2_update = fmt.Sprintf(`
resource "huaweicloud_vbs_backup_policy_v2" "vbs" {
  name = "policy_002"
  start_time  = "12:00"
  status  = "ON"
  retain_first_backup = "N"
  rentention_num = 2
  frequency = 1
      tags =[
        {
          key = "k2"
          value = "v2"
          }] 
}
`)

var testAccVBSBackupPolicyV2_timeout = fmt.Sprintf(`
resource "huaweicloud_vbs_backup_policy_v2" "vbs" {
  name = "policy_002"
  start_time  = "12:00"
  status  = "ON"
  retain_first_backup = "N"
  rentention_num = 2
  frequency = 1
      tags =[
        {
          key = "k2"
          value = "v2"
          }]

  timeouts {
    create = "5m"
    delete = "5m"
  }
}`)

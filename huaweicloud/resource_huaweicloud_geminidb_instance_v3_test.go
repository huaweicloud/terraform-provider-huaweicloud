package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/geminidb/v3/instances"
)

func TestGeminiDBInstance_basic(t *testing.T) {
	var instance instances.GeminiDBInstance
	name := fmt.Sprintf("acc-instance-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGeminiDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBInstanceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGeminiDBInstanceExists("huaweicloud_geminidb_instance.instance_acc", &instance),
					resource.TestCheckResourceAttr(
						"huaweicloud_geminidb_instance.instance_acc", "name", name),
					resource.TestCheckResourceAttr(
						"huaweicloud_geminidb_instance.instance_acc", "status", "normal"),
				),
			},
		},
	})
}

func testAccCheckGeminiDBInstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.GeminiDBV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud GeminiDB client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_geminidb_instance" {
			continue
		}

		found, err := instances.GetInstanceByID(client, rs.Primary.ID)
		if err != nil {
			return err
		}
		if found.Id != "" {
			return fmt.Errorf("Instance <%s> still exists.", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckGeminiDBInstanceExists(n string, instance *instances.GeminiDBInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s.", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set.")
		}

		config := testAccProvider.Meta().(*Config)
		client, err := config.GeminiDBV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud GeminiDB client: %s", err)
		}

		found, err := instances.GetInstanceByID(client, rs.Primary.ID)
		if err != nil {
			return err
		}
		if found.Id == "" {
			return fmt.Errorf("Instance <%s> not found.", rs.Primary.ID)
		}
		instance = &found

		return nil
	}
}

func testAccGeminiDBInstanceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup_v2" "secgroup_acc" {
  name = "secgroup_acc"
}

resource "huaweicloud_geminidb_instance" "instance_acc" {
  name        = "%s"
  password    = "Test@123"
  flavor      = "geminidb.cassandra.xlarge.4"
  volume_size = 100
  vpc_id      = "%s"
  subnet_id   = "%s"
  security_group_id = huaweicloud_networking_secgroup_v2.secgroup_acc.id
  availability_zone = "%s"

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 14
  }
}
	`, name, OS_VPC_ID, OS_NETWORK_ID, OS_AVAILABILITY_ZONE)
}

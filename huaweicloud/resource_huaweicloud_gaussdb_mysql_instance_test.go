package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/taurusdb/v3/instances"
)

func TestGaussDBInstance_basic(t *testing.T) {
	var instance instances.TaurusDBInstance
	name := fmt.Sprintf("gauss-instance-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGaussDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBInstanceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGaussDBInstanceExists("huaweicloud_gaussdb_instance.instance_acc", &instance),
					resource.TestCheckResourceAttr(
						"huaweicloud_gaussdb_instance.instance_acc", "name", name),
				),
			},
		},
	})
}

func testAccCheckGaussDBInstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.initServiceClient("gaussdb", OS_REGION_NAME, "mysql/v3")
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud GaussDB client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_gaussdb_instance" {
			continue
		}

		v, err := instances.Get(client, rs.Primary.ID).Extract()
		if err == nil && v.Id == rs.Primary.ID {
			return fmt.Errorf("Instance <%s> still exists.", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckGaussDBInstanceExists(n string, instance *instances.TaurusDBInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s.", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set.")
		}

		config := testAccProvider.Meta().(*Config)
		client, err := config.initServiceClient("gaussdb", OS_REGION_NAME, "mysql/v3")
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud GaussDB client: %s", err)
		}

		found, err := instances.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Instance <%s> not found.", rs.Primary.ID)
		}
		instance = found

		return nil
	}
}

func testAccGaussDBInstanceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup_v2" "secgroup_acc" {
  name = "secgroup_acc"
}

resource "huaweicloud_gaussdb_instance" "instance_acc" {
  name        = "%s"
  password    = "Test@123"
  flavor      = "gaussdb.mysql.4xlarge.x86.4"
  vpc_id      = "%s"
  subnet_id   = "%s"
  security_group_id = huaweicloud_networking_secgroup_v2.secgroup_acc.id
  enterprise_project_id = "0"
}
`, name, OS_VPC_ID, OS_NETWORK_ID)
}

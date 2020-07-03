package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/opengauss/v3/instances"
)

func TestOpenGaussInstance_basic(t *testing.T) {
	var instance instances.GaussDBInstance
	name := fmt.Sprintf("opengauss-instance-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOpenGaussInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenGaussInstanceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOpenGaussInstanceExists("huaweicloud_gaussdb_opengauss_instance.instance_acc", &instance),
					resource.TestCheckResourceAttr(
						"huaweicloud_gaussdb_opengauss_instance.instance_acc", "name", name),
				),
			},
		},
	})
}

func testAccCheckOpenGaussInstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.initServiceClient("gaussdb", OS_REGION_NAME, "opengauss/v3")
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud GaussDB client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_gaussdb_opengauss_instance" {
			continue
		}

		v, err := instances.GetInstanceByID(client, rs.Primary.ID)
		if err == nil && v.Id == rs.Primary.ID {
			return fmt.Errorf("Instance <%s> still exists.", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckOpenGaussInstanceExists(n string, instance *instances.GaussDBInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s.", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set.")
		}

		config := testAccProvider.Meta().(*Config)
		client, err := config.initServiceClient("gaussdb", OS_REGION_NAME, "opengauss/v3")
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud GaussDB client: %s", err)
		}

		found, err := instances.GetInstanceByID(client, rs.Primary.ID)
		if err != nil {
			return err
		}
		if found.Id != rs.Primary.ID {
			return fmt.Errorf("Instance <%s> not found.", rs.Primary.ID)
		}
		instance = &found

		return nil
	}
}

func testAccOpenGaussInstanceConfig_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_networking_secgroup_v2" "secgroup_1" {
  name = "default"
}

resource "huaweicloud_gaussdb_opengauss_instance" "instance_acc" {
  name        = "%s"
  password    = "Test@123"
  flavor      = "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"
  vpc_id      = "%s"
  subnet_id   = "%s"
  availability_zone = "cn-north-4a,cn-north-4a,cn-north-4a"

  ha {
    mode = "enterprise"
    replication_mode = "sync"
    consistency = "strong"
  }

  volume {
    type = "ULTRAHIGH"
    size = 120
  }
  security_group_id = data.huaweicloud_networking_secgroup_v2.secgroup_1.id
  sharding_num = 1
  coordinator_num = 1
}
`, name, OS_VPC_ID, OS_NETWORK_ID)
}

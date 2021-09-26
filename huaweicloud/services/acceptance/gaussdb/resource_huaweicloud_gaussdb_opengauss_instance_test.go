package gaussdb

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/opengauss/v3/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccOpenGaussInstance_basic(t *testing.T) {
	var instance instances.GaussDBInstance

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	password := fmt.Sprintf("%s@123", acctest.RandString(5))
	newPassword := fmt.Sprintf("%sUpdate@123", acctest.RandString(5))
	resourceName := "huaweicloud_gaussdb_opengauss_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: testAccCheckOpenGaussInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenGaussInstanceConfig_basic(rName, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOpenGaussInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "password", password),
				),
			},
			{
				Config: testAccOpenGaussInstanceConfig_update(rName, newPassword),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOpenGaussInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "password", newPassword),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.start_time", "08:30-09:30"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "8"),
				),
			},
		},
	})
}

func testAccCheckOpenGaussInstanceDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := config.OpenGaussV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud GaussDB client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_gaussdb_opengauss_instance" {
			continue
		}

		v, err := instances.GetInstanceByID(client, rs.Primary.ID)
		if err == nil && v.Id == rs.Primary.ID {
			return fmtp.Errorf("Instance <%s> still exists.", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckOpenGaussInstanceExists(n string, instance *instances.GaussDBInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s.", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set.")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.OpenGaussV3Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud GaussDB client: %s", err)
		}

		found, err := instances.GetInstanceByID(client, rs.Primary.ID)
		if err != nil {
			return err
		}
		if found.Id != rs.Primary.ID {
			return fmtp.Errorf("Instance <%s> not found.", rs.Primary.ID)
		}
		instance = &found

		return nil
	}
}

func testAccOpenGaussInstanceConfig_base(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "myaz" {}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}
`, testAccVpcConfig_Base(rName))
}

func testAccOpenGaussInstanceConfig_basic(rName, password string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_gaussdb_opengauss_instance" "test" {
  name        = "%s"
  password    = "%s"
  flavor      = "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id

  availability_zone = "${data.huaweicloud_availability_zones.myaz.names[0]},${data.huaweicloud_availability_zones.myaz.names[0]},${data.huaweicloud_availability_zones.myaz.names[0]}"
  security_group_id = data.huaweicloud_networking_secgroup.test.id

  ha {
    mode             = "enterprise"
    replication_mode = "sync"
    consistency      = "strong"
  }
  volume {
    type = "ULTRAHIGH"
    size = 40
  }

  sharding_num = 1
  coordinator_num = 2
}
`, testAccOpenGaussInstanceConfig_base(rName), rName, password)
}

func testAccOpenGaussInstanceConfig_update(rName, password string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_gaussdb_opengauss_instance" "test" {
  name        = "%s-update"
  password    = "%s"
  flavor      = "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id

  availability_zone = "${data.huaweicloud_availability_zones.myaz.names[0]},${data.huaweicloud_availability_zones.myaz.names[0]},${data.huaweicloud_availability_zones.myaz.names[0]}"
  security_group_id = data.huaweicloud_networking_secgroup.test.id

  ha {
    mode             = "enterprise"
    replication_mode = "sync"
    consistency      = "strong"
  }
  volume {
    type = "ULTRAHIGH"
    size = 40
  }
  backup_strategy {
    start_time = "08:30-09:30"
    keep_days  = 8
  }

  sharding_num = 1
  coordinator_num = 2
}
`, testAccOpenGaussInstanceConfig_base(rName), rName, password)
}

package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/rds/v3/instances"
)

func TestAccRdsReadReplicaInstance_basic(t *testing.T) {
	var resourceName string = "huaweicloud_rds_read_replica_instance.replica_instance"
	var replica instances.RdsInstanceResponse
	nameSuffix := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRdsReplicaInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccReadRdsReplicaInstanceBasic(nameSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsReplicaInstanceV3Exists(resourceName, &replica),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "type", "Replica"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "ULTRAHIGH"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"db",
				},
			},
			{
				Config: testAccReadRdsReplicaInstanceUpdate(nameSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsReplicaInstanceV3Exists(resourceName, &replica),
					// check modification of flavor
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.c2.large.rr"),
					// check modification of volume
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "ULTRAHIGH"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "120"),
					// check modification of tags
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar2"),
				),
			},
		},
	})
}

func testAccCheckRdsReplicaInstanceV3Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.RdsV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating huaweicloud rds client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_rds_read_replica_instance" {
			continue
		}

		id := rs.Primary.ID
		instance, err := getRdsInstanceByID(client, id)
		if err != nil {
			return err
		}
		if instance.Id != "" {
			return fmt.Errorf("huaweicloud_rds_read_replica_instance (%s) still exists", id)
		}
	}

	return nil
}

func testAccCheckRdsReplicaInstanceV3Exists(n string, instance *instances.RdsInstanceResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s. ", n)
		}

		id := rs.Primary.ID
		if id == "" {
			return fmt.Errorf("No ID is set. ")
		}

		config := testAccProvider.Meta().(*Config)
		client, err := config.RdsV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating huaweicloud rds client: %s", err)
		}

		found, err := getRdsInstanceByID(client, id)
		if err != nil {
			return fmt.Errorf("Error checking %s exist, err=%s", n, err)
		}
		if found.Id == "" {
			return fmt.Errorf("resource %s does not exist", n)
		}

		instance = found
		return nil
	}
}

func testAccReadRdsReplicaInstanceBasic(val string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "secgroup" {
  name         = "acctest_sg_rds"
  description  = "security group for rds read replica instance test"
}

resource "huaweicloud_rds_instance" "instance" {
  name                  = "rds_instance_%s"
  flavor                = "rds.pg.c2.medium"
  availability_zone     = ["%s"]
  security_group_id     = huaweicloud_networking_secgroup.secgroup.id
  vpc_id                = "%s"
  subnet_id             = "%s"
  enterprise_project_id = "%s"

  db {
    password    = "Huangwei!120521"
    type        = "PostgreSQL"
    version     = "10"
    port        = "8635"
  }
  volume {
    type = "ULTRAHIGH"
    size = 50
  }
  backup_strategy {
    start_time  = "08:00-09:00"
    keep_days   = 1
  }
  tags = {
    key = "value"
    foo = "bar"
  }
}

resource "huaweicloud_rds_read_replica_instance" "replica_instance" {
  name                  = "replica_instance_%s"
  flavor                = "rds.pg.c2.medium.rr"
  primary_instance_id   = huaweicloud_rds_instance.instance.id
  availability_zone     = "%s"
  enterprise_project_id = "%s"
  volume {
    type = "ULTRAHIGH"
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, val, OS_AVAILABILITY_ZONE, OS_VPC_ID, OS_NETWORK_ID, OS_ENTERPRISE_PROJECT_ID_TEST, val, OS_AVAILABILITY_ZONE, OS_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccReadRdsReplicaInstanceUpdate(val string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "secgroup" {
  name          = "acctest_sg_rds"
  description   = "security group for rds read replica instance test"
}

resource "huaweicloud_rds_instance" "instance" {
  name                  = "rds_instance_%s"
  flavor                = "rds.pg.c2.medium"
  availability_zone     = ["%s"]
  security_group_id     = huaweicloud_networking_secgroup.secgroup.id
  vpc_id                = "%s"
  subnet_id             = "%s"
  enterprise_project_id = "%s"

  db {
    password    = "Huangwei!120521"
    type        = "PostgreSQL"
    version     = "10"
    port        = "8635"
  }
  volume {
    type = "ULTRAHIGH"
    size = 50
  }
  backup_strategy {
    start_time  = "08:00-09:00"
    keep_days   = 1
  }
  tags = {
    key = "value"
    foo = "bar"
  }
}

resource "huaweicloud_rds_read_replica_instance" "replica_instance" {
  name                  = "replica_instance_%s"
  flavor                = "rds.pg.c2.large.rr"
  primary_instance_id   = huaweicloud_rds_instance.instance.id
  availability_zone     = "%s"
  enterprise_project_id = "%s"
  volume {
	type = "ULTRAHIGH"
	size = 120
  }

  tags = {
    key = "value1"
    foo = "bar2"
  }
}
`, val, OS_AVAILABILITY_ZONE, OS_VPC_ID, OS_NETWORK_ID, OS_ENTERPRISE_PROJECT_ID_TEST, val, OS_AVAILABILITY_ZONE, OS_ENTERPRISE_PROJECT_ID_TEST)
}

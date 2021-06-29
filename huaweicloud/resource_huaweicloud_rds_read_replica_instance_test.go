package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/huaweicloud/golangsdk/openstack/rds/v3/instances"
)

func TestAccRdsReadReplicaInstance_basic(t *testing.T) {
	var replica instances.RdsInstanceResponse
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceType := "huaweicloud_rds_read_replica_instance"
	resourceName := "huaweicloud_rds_read_replica_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRdsInstanceV3Destroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccReadRdsReplicaInstance_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(resourceName, &replica),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.n1.large.2.rr"),
					resource.TestCheckResourceAttr(resourceName, "type", "Replica"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "CLOUDSSD"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
				),
			},
			{
				Config: testAccReadRdsReplicaInstance_update(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(resourceName, &replica),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.n1.xlarge.2.rr"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "CLOUDSSD"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar2"),
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
		},
	})
}

func TestAccRdsReadReplicaInstance_withEpsId(t *testing.T) {
	var replica instances.RdsInstanceResponse
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceType := "huaweicloud_rds_read_replica_instance"
	resourceName := "huaweicloud_rds_read_replica_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEpsID(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRdsInstanceV3Destroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccReadRdsReplicaInstance_withEpsId(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(resourceName, &replica),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccReadRdsReplicaInstanceV3_base(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_instance" "test" {
  name              = "%s"
  flavor            = "rds.pg.n1.large.2"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = huaweicloud_networking_secgroup.test.id
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id

  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "12"
    port     = 8635
  }
  volume {
    type = "CLOUDSSD"
    size = 50
  }
}
`, testAccRdsInstanceV3_base(name), name)
}

func testAccReadRdsReplicaInstance_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_read_replica_instance" "test" {
  name                = "%s"
  flavor              = "rds.pg.n1.large.2.rr"
  primary_instance_id = huaweicloud_rds_instance.test.id
  availability_zone   = data.huaweicloud_availability_zones.test.names[0]

  volume {
    type = "CLOUDSSD"
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, testAccReadRdsReplicaInstanceV3_base(name), name)
}

func testAccReadRdsReplicaInstance_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_read_replica_instance" "test" {
  name                = "%s"
  flavor              = "rds.pg.n1.xlarge.2.rr"
  primary_instance_id = huaweicloud_rds_instance.test.id
  availability_zone   = data.huaweicloud_availability_zones.test.names[0]

  volume {
	type = "CLOUDSSD"
  }

  tags = {
    key1 = "value"
    foo = "bar2"
  }
}
`, testAccReadRdsReplicaInstanceV3_base(name), name)
}

func testAccReadRdsReplicaInstance_withEpsId(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_read_replica_instance" "test" {
  name                  = "%s"
  flavor                = "rds.pg.n1.large.2.rr"
  primary_instance_id   = huaweicloud_rds_instance.test.id
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  enterprise_project_id = "%s"

  volume {
    type = "CLOUDSSD"
  }
}
`, testAccReadRdsReplicaInstanceV3_base(name), name, HW_ENTERPRISE_PROJECT_ID_TEST)
}

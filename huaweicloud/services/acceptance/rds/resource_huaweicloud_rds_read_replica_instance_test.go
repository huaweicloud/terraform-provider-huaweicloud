package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk/openstack/rds/v3/instances"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccReadReplicaInstance_basic(t *testing.T) {
	var replica instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "huaweicloud_rds_read_replica_instance"
	resourceName := "huaweicloud_rds_read_replica_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccReadReplicaInstance_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &replica),
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
				Config: testAccReadReplicaInstance_update(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &replica),
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

func TestAccReadReplicaInstance_withEpsId(t *testing.T) {
	var replica instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "huaweicloud_rds_read_replica_instance"
	resourceName := "huaweicloud_rds_read_replica_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccReadReplicaInstance_withEpsId(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &replica),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccReadReplicaInstance_base(name string) string {
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
`, testAccRdsInstance_base(name), name)
}

func testAccReadReplicaInstance_basic(name string) string {
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
`, testAccReadReplicaInstance_base(name), name)
}

func testAccReadReplicaInstance_update(name string) string {
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
`, testAccReadReplicaInstance_base(name), name)
}

func testAccReadReplicaInstance_withEpsId(name string) string {
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
`, testAccReadReplicaInstance_base(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

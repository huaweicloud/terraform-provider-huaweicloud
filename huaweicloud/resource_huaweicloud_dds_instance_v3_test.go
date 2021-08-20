package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/golangsdk/openstack/dds/v3/instances"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccDDSV3Instance_basic(t *testing.T) {
	var instance instances.InstanceResponse
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_dds_instance.instance"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDDSV3InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDDSInstanceV3Config_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDDSV3InstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "ssl", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.start_time", "08:00-09:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "8"),
				),
			},
			{
				Config: testAccDDSInstanceV3Config_updateBackupStrategy(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDDSV3InstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.start_time", "00:00-01:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "7"),
				),
			},
			{
				Config: testAccDDSInstanceV3Config_updateFlavorNum(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDDSV3InstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckDDSV3InstanceFlavor(&instance, "shard", "num", 3),
				),
			},
			{
				Config: testAccDDSInstanceV3Config_updateFlavorSize(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDDSV3InstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckDDSV3InstanceFlavor(&instance, "shard", "size", "30"),
				),
			},
			{
				Config: testAccDDSInstanceV3Config_updateFlavorSpecCode(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDDSV3InstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckDDSV3InstanceFlavor(&instance, "mongos", "spec_code", "dds.mongodb.c6.large.4.mongos"),
				),
			},
		},
	})
}

func TestAccDDSV3Instance_withEpsId(t *testing.T) {
	var instance instances.InstanceResponse
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_dds_instance.instance"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEpsID(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDDSV3InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDDSInstanceV3Config_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDDSV3InstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccCheckDDSV3InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.DdsV3Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud DDS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_dds_instance" {
			continue
		}

		opts := instances.ListInstanceOpts{
			Id: rs.Primary.ID,
		}
		allPages, err := instances.List(client, &opts).AllPages()
		if err != nil {
			return err
		}
		instances, err := instances.ExtractInstances(allPages)
		if err != nil {
			return err
		}

		if instances.TotalCount > 0 {
			return fmtp.Errorf("Instance still exists. ")
		}
	}

	return nil
}

func testAccCheckDDSV3InstanceExists(n string, instance *instances.InstanceResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s. ", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set. ")
		}

		config := testAccProvider.Meta().(*config.Config)
		client, err := config.DdsV3Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud DDS client: %s ", err)
		}

		opts := instances.ListInstanceOpts{
			Id: rs.Primary.ID,
		}
		allPages, err := instances.List(client, &opts).AllPages()
		if err != nil {
			return err
		}
		instances, err := instances.ExtractInstances(allPages)
		if err != nil {
			return err
		}
		if instances.TotalCount == 0 {
			return fmtp.Errorf("dds instance not found.")
		}

		insts := instances.Instances
		found := insts[0]
		*instance = found

		return nil
	}
}

func testAccCheckDDSV3InstanceFlavor(instance *instances.InstanceResponse, groupType, key string, v interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if key == "num" {
			if groupType == "mongos" {
				for _, group := range instance.Groups {
					if group.Type == "mongos" {
						if len(group.Nodes) != v.(int) {
							return fmtp.Errorf(
								"Error updating HuaweiCloud DDS instance: num of mongos nodes expect %d, but got %d",
								v.(int), len(group.Nodes))
						}
						return nil
					}
				}
			} else {
				groupIDs := make([]string, 0)
				for _, group := range instance.Groups {
					if group.Type == "shard" {
						groupIDs = append(groupIDs, group.Id)
					}
				}
				if len(groupIDs) != v.(int) {
					return fmtp.Errorf(
						"Error updating HuaweiCloud DDS instance: num of shard groups expect %d, but got %d",
						v.(int), len(groupIDs))
				}
				return nil
			}
		}

		if key == "size" {
			for _, group := range instance.Groups {
				if group.Type == groupType {
					if group.Volume.Size != v.(string) {
						return fmtp.Errorf(
							"Error updating HuaweiCloud DDS instance: size expect %s, but got %s",
							v.(string), group.Volume.Size)
					}
					return nil
				}
			}
		}

		if key == "spec_code" {
			for _, group := range instance.Groups {
				if group.Type == groupType {
					for _, node := range group.Nodes {
						if node.SpecCode != v.(string) {
							return fmtp.Errorf(
								"Error updating HuaweiCloud DDS instance: spec_code expect %s, but got %s",
								v.(string), node.SpecCode)
						}
					}
					return nil
				}
			}
		}
		return nil
	}
}

func testAccDDSInstanceV3Config_Base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_networking_secgroup" "secgroup_acc" {
  name = "%s"
}`, rName)
}

func testAccDDSInstanceV3Config_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dds_instance" "instance" {
  name              = "%s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = data.huaweicloud_vpc.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.secgroup_acc.id
  password          = "Test@123"
  mode              = "Sharding"

  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.c6.large.2.mongos"
  }
  flavor {
    type      = "shard"
    num       = 2
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.c6.large.2.shard"
  }
  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.c6.large.2.config"
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = "8"
  }

  tags = {
	foo   = "bar"
    owner = "terraform"
  }
}`, testAccDDSInstanceV3Config_Base(rName), rName)
}

func testAccDDSInstanceV3Config_updateBackupStrategy(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dds_instance" "instance" {
  name              = "%s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = data.huaweicloud_vpc.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.secgroup_acc.id
  password          = "Test@123"
  mode              = "Sharding"

  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.c6.large.2.mongos"
  }
  flavor {
    type      = "shard"
    num       = 2
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.c6.large.2.shard"
  }
  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.c6.large.2.config"
  }

  backup_strategy {
    start_time = "00:00-01:00"
    keep_days  = "7"
  }

  tags = {
	foo   = "bar"
    owner = "terraform"
  }
}`, testAccDDSInstanceV3Config_Base(rName), rName)
}

func testAccDDSInstanceV3Config_updateFlavorNum(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dds_instance" "instance" {
  name              = "%s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = data.huaweicloud_vpc.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.secgroup_acc.id
  password          = "Test@123"
  mode              = "Sharding"

  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.c6.large.2.mongos"
  }
  flavor {
    type      = "shard"
    num       = 3
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.c6.large.2.shard"
  }
  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.c6.large.2.config"
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = "8"
  }

  tags = {
	foo   = "bar"
    owner = "terraform"
  }
}`, testAccDDSInstanceV3Config_Base(rName), rName)
}

func testAccDDSInstanceV3Config_updateFlavorSize(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dds_instance" "instance" {
  name              = "%s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = data.huaweicloud_vpc.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.secgroup_acc.id
  password          = "Test@123"
  mode              = "Sharding"

  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.c6.large.2.mongos"
  }
  flavor {
    type      = "shard"
    num       = 3
    storage   = "ULTRAHIGH"
    size      = 30
    spec_code = "dds.mongodb.c6.large.2.shard"
  }
  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.c6.large.2.config"
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = "8"
  }

  tags = {
	foo   = "bar"
    owner = "terraform"
  }
}`, testAccDDSInstanceV3Config_Base(rName), rName)
}

func testAccDDSInstanceV3Config_updateFlavorSpecCode(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dds_instance" "instance" {
  name              = "%s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = data.huaweicloud_vpc.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.secgroup_acc.id
  password          = "Test@123"
  mode              = "Sharding"

  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.c6.large.4.mongos"
  }
  flavor {
    type      = "shard"
    num       = 3
    storage   = "ULTRAHIGH"
    size      = 30
    spec_code = "dds.mongodb.c6.large.2.shard"
  }
  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.c6.large.2.config"
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = "8"
  }

  tags = {
	foo   = "bar"
    owner = "terraform"
  }
}`, testAccDDSInstanceV3Config_Base(rName), rName)
}

func testAccDDSInstanceV3Config_withEpsId(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dds_instance" "instance" {
  name                  = "%s"
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  vpc_id                = data.huaweicloud_vpc.test.id
  subnet_id             = data.huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.secgroup_acc.id
  password              = "Test@123"
  mode                  = "Sharding"
  enterprise_project_id = "%s"

  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.c6.large.2.mongos"
  }
  flavor {
    type      = "shard"
    num       = 2
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.c6.large.2.shard"
  }
  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.c6.large.2.config"
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = "8"
  }

  tags = {
	foo = "bar"
    owner = "terraform"
  }
}`, testAccDDSInstanceV3Config_Base(rName), rName, HW_ENTERPRISE_PROJECT_ID_TEST)
}

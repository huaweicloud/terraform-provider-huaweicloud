package dds

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/dds/v3/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getDdsResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DdsV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("Error creating HuaweiCloud DDS client: %s ", err)
	}

	opts := instances.ListInstanceOpts{
		Id: state.Primary.ID,
	}
	allPages, err := instances.List(client, &opts).AllPages()
	if err != nil {
		return nil, err
	}
	instances, err := instances.ExtractInstances(allPages)
	if err != nil {
		return nil, err
	}
	if instances.TotalCount == 0 {
		return nil, fmtp.Errorf("dds instance not found.")
	}

	insts := instances.Instances
	found := insts[0]
	return &found, nil

}

func TestAccDDSV3Instance_basic(t *testing.T) {
	var instance instances.InstanceResponse
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dds_instance.instance"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDdsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDDSInstanceV3Config_basic(rName, 8800),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "ssl", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.start_time", "08:00-09:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "8"),
					resource.TestCheckResourceAttr(resourceName, "port", "8800"),
				),
			},
			{
				Config: testAccDDSInstanceV3Config_basic(rName, 8635),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "port", "8635"),
				),
			},
			{
				Config: testAccDDSInstanceV3Config_updateBackupStrategy(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.start_time", "00:00-01:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "7"),
				),
			},
			{
				Config: testAccDDSInstanceV3Config_updateFlavorNum(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckDDSV3InstanceFlavor(&instance, "shard", "num", 3),
				),
			},
			{
				Config: testAccDDSInstanceV3Config_updateFlavorSize(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckDDSV3InstanceFlavor(&instance, "shard", "size", "30"),
				),
			},
			{
				Config: testAccDDSInstanceV3Config_updateFlavorSpecCode(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckDDSV3InstanceFlavor(&instance, "mongos", "spec_code", "dds.mongodb.c6.large.4.mongos"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"availability_zone", "flavor", "password",
				},
			},
		},
	})
}

func TestAccDDSV3Instance_withEpsId(t *testing.T) {
	var instance instances.InstanceResponse
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dds_instance.instance"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDdsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDDSInstanceV3Config_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccDDSV3Instance_prePaid(t *testing.T) {
	var instance instances.InstanceResponse
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dds_instance.instance"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDdsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDDSInstanceV3Config_prePaid(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.start_time", "08:00-09:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "8"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccDDSInstanceV3Config_prePaidUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.start_time", "00:00-01:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "7"),
					testAccCheckDDSV3InstanceFlavor(&instance, "shard", "num", 3),
					testAccCheckDDSV3InstanceFlavor(&instance, "shard", "size", "30"),
					testAccCheckDDSV3InstanceFlavor(&instance, "mongos", "spec_code", "dds.mongodb.c6.large.4.mongos"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
				),
			},
		},
	})
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

func testAccDDSInstanceV3Config_basic(rName string, port int) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dds_instance" "instance" {
  name              = "%s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = data.huaweicloud_vpc.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.secgroup_acc.id
  password          = "Terraform@123"
  mode              = "Sharding"
  port              = %d

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
}`, testAccDDSInstanceV3Config_Base(rName), rName, port)
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
  password          = "Terraform@123"
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
  password          = "Terraform@123"
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
  password          = "Terraform@123"
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
  password          = "Terraform@123"
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
  password              = "Terraform@123"
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
    foo   = "bar"
    owner = "terraform"
  }
}`, testAccDDSInstanceV3Config_Base(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccDDSInstanceV3Config_prePaid(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dds_instance" "instance" {
  name              = "%s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = data.huaweicloud_vpc.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.secgroup_acc.id
  password          = "Terraform@123"
  mode              = "Sharding"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = false

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
    keep_days  = 8
  }
}`, testAccDDSInstanceV3Config_Base(rName), rName)
}

func testAccDDSInstanceV3Config_prePaidUpdate(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dds_instance" "instance" {
  name              = "%s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = data.huaweicloud_vpc.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.secgroup_acc.id
  password          = "Terraform@123"
  mode              = "Sharding"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"

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
    start_time = "00:00-01:00"
    keep_days  = 7
  }
}`, testAccDDSInstanceV3Config_Base(rName), rName)
}

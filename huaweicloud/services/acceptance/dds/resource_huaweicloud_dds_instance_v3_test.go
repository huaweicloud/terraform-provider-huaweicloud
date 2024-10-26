package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dds/v3/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getDdsResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DdsV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error creating DDS client: %s", err)
	}

	opts := instances.ListInstanceOpts{
		Id: state.Primary.ID,
	}
	allPages, err := instances.List(client, &opts).AllPages()
	if err != nil {
		return nil, err
	}
	instanceList, err := instances.ExtractInstances(allPages)
	if err != nil {
		return nil, err
	}
	if instanceList.TotalCount == 0 {
		return nil, fmt.Errorf("dds instance not found")
	}

	insts := instanceList.Instances
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
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.period", "1,3,5"),
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
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.period", "2,4,6"),
				),
			},
			{
				Config: testAccDDSInstanceV3Config_updateFlavorNum(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckDDSV3InstanceFlavor(&instance, "shard", "num", 3),
					resource.TestCheckResourceAttr(resourceName, "slow_log_desensitization", "off"),
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
					testAccCheckDDSV3InstanceFlavor(&instance, "config", "spec_code", "dds.mongodb.s6.xlarge.2.config"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"availability_zone", "flavor", "password", "groups", "nodes",
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
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDDSInstanceV3Config_baseEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
				),
			},
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
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
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
					testAccCheckDDSV3InstanceFlavor(&instance, "mongos", "spec_code", "dds.mongodb.s6.large.4.mongos"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
		},
	})
}

func TestAccDDSV3Instance_withConfigurationSharding(t *testing.T) {
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
				Config: testAccDDSInstanceV3Config_withShardingConfiguration(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.type", "mongos"),
					resource.TestCheckResourceAttrPair(resourceName, "configuration.0.id", "huaweicloud_dds_parameter_template.mongos1", "id"),
					resource.TestCheckResourceAttr(resourceName, "configuration.1.type", "shard"),
					resource.TestCheckResourceAttrPair(resourceName, "configuration.1.id", "huaweicloud_dds_parameter_template.shard1", "id"),
					resource.TestCheckResourceAttr(resourceName, "configuration.2.type", "config"),
					resource.TestCheckResourceAttrPair(resourceName, "configuration.2.id", "huaweicloud_dds_parameter_template.config1", "id"),
					testAccCheckDDSV3InstanceFlavor(&instance, "replica", "num", 3),
					resource.TestCheckResourceAttr(resourceName, "balancer_active_begin", "00:00"),
					resource.TestCheckResourceAttr(resourceName, "balancer_active_end", "02:00"),
				),
			},
			{
				Config: testAccDDSInstanceV3Config_withShardingConfigurationUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.type", "mongos"),
					resource.TestCheckResourceAttrPair(resourceName, "configuration.0.id", "huaweicloud_dds_parameter_template.mongos2", "id"),
					resource.TestCheckResourceAttr(resourceName, "configuration.1.type", "shard"),
					resource.TestCheckResourceAttrPair(resourceName, "configuration.1.id", "huaweicloud_dds_parameter_template.shard2", "id"),
					resource.TestCheckResourceAttr(resourceName, "configuration.2.type", "config"),
					resource.TestCheckResourceAttrPair(resourceName, "configuration.2.id", "huaweicloud_dds_parameter_template.config2", "id"),
					testAccCheckDDSV3InstanceFlavor(&instance, "replica", "num", 5),
					resource.TestCheckResourceAttr(resourceName, "balancer_active_begin", ""),
					resource.TestCheckResourceAttr(resourceName, "balancer_status", "stop"),
				),
			},
		},
	})
}

func TestAccDDSV3Instance_withConfigurationReplicaSet(t *testing.T) {
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
				Config: testAccDDSInstanceV3Config_withReplicaSetConfiguration(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "replica_set_name", "replicaName"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.type", "replica"),
					resource.TestCheckResourceAttr(resourceName, "client_network_ranges.0", "192.168.0.0/24"),
					resource.TestCheckResourceAttrPair(resourceName, "configuration.0.id", "huaweicloud_dds_parameter_template.replica1", "id"),
				),
			},
			{
				Config: testAccDDSInstanceV3Config_withReplicaSetConfigurationUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.type", "replica"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "02:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "03:00"),
					resource.TestCheckResourceAttrPair(resourceName, "configuration.0.id", "huaweicloud_dds_parameter_template.replica2", "id"),
				),
			},
		},
	})
}

func TestAccDDSV3Instance_withSecondLevelMonitoring(t *testing.T) {
	var instance instances.InstanceResponse
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dds_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDdsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSSecondLevelMonitoringEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDDSInstanceV3Config_secondLevelMonitoring(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "second_level_monitoring_enabled", "true"),
				),
			},
			{
				Config: testAccDDSInstanceV3Config_secondLevelMonitoringUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "second_level_monitoring_enabled", "false"),
				),
			},
		},
	})
}

func TestAccDDSV3Instance_updateAZ(t *testing.T) {
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
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDDSInstanceV3Config_withReplicaSetConfiguration(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone", "data.huaweicloud_availability_zones.test", "names.0"),
				),
			},
			{
				Config: testAccDDSInstanceV3Config_updateAZ(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone", "data.huaweicloud_availability_zones.test", "names.1"),
				),
			},
		},
	})
}

func testAccCheckDDSV3InstanceFlavor(instance *instances.InstanceResponse, groupType, key string, v interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if key == "num" {
			if groupType == "shard" {
				groupIDs := make([]string, 0)
				for _, group := range instance.Groups {
					if group.Type == "shard" {
						groupIDs = append(groupIDs, group.Id)
					}
				}
				if len(groupIDs) != v.(int) {
					return fmt.Errorf(
						"Error updating DDS instance: num of shard groups expect %d, but got %d",
						v.(int), len(groupIDs))
				}
				return nil
			}

			for _, group := range instance.Groups {
				if group.Type == groupType {
					if len(group.Nodes) != v.(int) {
						return fmt.Errorf(
							"Error updating DDS instance: num of %s nodes expect %d, but got %d", groupType, v.(int), len(group.Nodes))
					}
					return nil
				}
			}
		}

		if key == "size" {
			for _, group := range instance.Groups {
				if group.Type == groupType {
					if group.Volume.Size != v.(string) {
						return fmt.Errorf(
							"Error updating DDS instance: size expect %s, but got %s",
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
							return fmt.Errorf(
								"Error updating DDS instance: spec_code expect %s, but got %s",
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

func testAccDDSInstanceV3Config_basic(rName string, port int) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dds_instance" "instance" {
  name              = "%s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "Terraform@123"
  mode              = "Sharding"
  port              = %d

  datastore {
    type           = "DDS-Community"
    version        = "4.0"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.s6.large.2.mongos"
  }
  flavor {
    type      = "shard"
    num       = 2
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s6.large.2.shard"
  }
  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s6.large.2.config"
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = "8"
    period     = "1,3,5"
  }

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}`, common.TestBaseNetwork(rName), rName, port)
}

func testAccDDSInstanceV3Config_updateBackupStrategy(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dds_instance" "instance" {
  name              = "%s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "Terraform@123"
  mode              = "Sharding"

  datastore {
    type           = "DDS-Community"
    version        = "4.0"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.s6.large.2.mongos"
  }
  flavor {
    type      = "shard"
    num       = 2
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s6.large.2.shard"
  }
  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s6.large.2.config"
  }

  backup_strategy {
    start_time = "00:00-01:00"
    keep_days  = "7"
    period     = "2,4,6"
  }

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}`, common.TestBaseNetwork(rName), rName)
}

func testAccDDSInstanceV3Config_updateFlavorNum(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dds_instance" "instance" {
  name                     = "%s"
  availability_zone        = data.huaweicloud_availability_zones.test.names[0]
  vpc_id                   = huaweicloud_vpc.test.id
  subnet_id                = huaweicloud_vpc_subnet.test.id
  security_group_id        = huaweicloud_networking_secgroup.test.id
  password                 = "Terraform@123"
  mode                     = "Sharding"
  slow_log_desensitization = "off"

  datastore {
    type           = "DDS-Community"
    version        = "4.0"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.s6.large.2.mongos"
  }
  flavor {
    type      = "shard"
    num       = 3
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s6.large.2.shard"
  }
  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s6.large.2.config"
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = "8"
    period     = "1,5"
  }

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}`, common.TestBaseNetwork(rName), rName)
}

func testAccDDSInstanceV3Config_updateFlavorSize(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dds_instance" "instance" {
  name              = "%s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "Terraform@123"
  mode              = "Sharding"

  datastore {
    type           = "DDS-Community"
    version        = "4.0"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.s6.large.2.mongos"
  }
  flavor {
    type      = "shard"
    num       = 3
    storage   = "ULTRAHIGH"
    size      = 30
    spec_code = "dds.mongodb.s6.large.2.shard"
  }
  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s6.large.2.config"
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = "8"
    period     = "1,5"
  }

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}`, common.TestBaseNetwork(rName), rName)
}

func testAccDDSInstanceV3Config_updateFlavorSpecCode(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dds_instance" "instance" {
  name              = "%s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "Terraform@123"
  mode              = "Sharding"

  datastore {
    type           = "DDS-Community"
    version        = "4.0"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.s6.large.2.mongos"
  }
  flavor {
    type      = "shard"
    num       = 3
    storage   = "ULTRAHIGH"
    size      = 30
    spec_code = "dds.mongodb.s6.large.2.shard"
  }
  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s6.xlarge.2.config"
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = "8"
    period     = "1,5"
  }

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}`, common.TestBaseNetwork(rName), rName)
}

func testAccDDSInstanceV3Config_withEpsId(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dds_instance" "instance" {
  name                  = "%s"
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  password              = "Terraform@123"
  mode                  = "Sharding"
  enterprise_project_id = "%s"

  datastore {
    type           = "DDS-Community"
    version        = "4.0"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.s6.large.2.mongos"
  }
  flavor {
    type      = "shard"
    num       = 2
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s6.large.2.shard"
  }
  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s6.large.2.config"
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = "8"
    period     = "1,5"
  }

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}`, common.TestBaseNetwork(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccDDSInstanceV3Config_baseEpsId(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dds_instance" "instance" {
  name                  = "%s"
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  password              = "Terraform@123"
  mode                  = "Sharding"
  enterprise_project_id = "0"

  datastore {
    type           = "DDS-Community"
    version        = "4.0"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.s6.large.2.mongos"
  }
  flavor {
    type      = "shard"
    num       = 2
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s6.large.2.shard"
  }
  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s6.large.2.config"
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = "8"
    period     = "1,5"
  }

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}`, common.TestBaseNetwork(rName), rName)
}

func testAccDDSInstanceV3Config_prePaid(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dds_instance" "instance" {
  name              = "%s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "Terraform@123"
  mode              = "Sharding"
  description       = "test description"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = false

  datastore {
    type           = "DDS-Community"
    version        = "4.0"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.s6.large.2.mongos"
  }
  flavor {
    type      = "shard"
    num       = 2
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s6.large.2.shard"
  }
  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s6.large.2.config"
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 8
    period     = "1,5"
  }
}`, common.TestBaseNetwork(rName), rName)
}

func testAccDDSInstanceV3Config_prePaidUpdate(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dds_instance" "instance" {
  name              = "%s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "Terraform@123"
  mode              = "Sharding"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"

  datastore {
    type           = "DDS-Community"
    version        = "4.0"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.s6.large.4.mongos"
  }
  flavor {
    type      = "shard"
    num       = 3
    storage   = "ULTRAHIGH"
    size      = 30
    spec_code = "dds.mongodb.s6.large.2.shard"
  }
  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s6.large.2.config"
  }

  backup_strategy {
    start_time = "00:00-01:00"
    keep_days  = 7
    period     = "1,3"
  }
}`, common.TestBaseNetwork(rName), rName)
}

func testAAccDDSInstance_templete(rName, nodeType string, i, value int) string {
	return fmt.Sprintf(`
resource "huaweicloud_dds_parameter_template" "%[3]s%[1]v" {
  name         = "%[2]s_%[3]s%[1]v"
  description  = "test description"
  node_type    = "%[3]s"
  node_version = "4.0"

  parameter_values = {
    connPoolMaxConnsPerHost        = %[4]v
    connPoolMaxShardedConnsPerHost = %[4]v
  }
}
`, i, rName, nodeType, value)
}

func testAccDDSInstanceV3Config_withShardingConfiguration(rName string) string {
	templateMongos1 := testAAccDDSInstance_templete(rName, "mongos", 1, 800)
	templateShard1 := testAAccDDSInstance_templete(rName, "shard", 1, 1000)
	templateConfig1 := testAAccDDSInstance_templete(rName, "config", 1, 400)
	return fmt.Sprintf(`
%[1]s

%[2]s

%[3]s

%[4]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dds_instance" "instance" {
  name                  = "%[5]s"
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  password              = "Terraform@123"
  mode                  = "Sharding"
  balancer_active_begin = "00:00"
  balancer_active_end   = "02:00"

  datastore {
    type           = "DDS-Community"
    version        = "4.0"
    storage_engine = "wiredTiger"
  }

  configuration {
    type = "mongos"
    id   = huaweicloud_dds_parameter_template.mongos1.id
  }

  configuration {
    type = "shard"
    id   = huaweicloud_dds_parameter_template.shard1.id
  }

  configuration {
    type = "config"
    id   = huaweicloud_dds_parameter_template.config1.id
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.s6.large.2.mongos"
  }

  flavor {
    type      = "shard"
    num       = 2
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s6.large.2.shard"
  }

  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s6.large.2.config"
  }
}`, common.TestBaseNetwork(rName), templateMongos1, templateShard1, templateConfig1, rName)
}

func testAccDDSInstanceV3Config_withShardingConfigurationUpdate(rName string) string {
	templateMongos2 := testAAccDDSInstance_templete(rName, "mongos", 2, 500)
	templateShard2 := testAAccDDSInstance_templete(rName, "shard", 2, 800)
	templateConfig2 := testAAccDDSInstance_templete(rName, "config", 2, 600)
	return fmt.Sprintf(`
%[1]s

%[2]s

%[3]s

%[4]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dds_instance" "instance" {
  name              = "%[5]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "Terraform@123"
  mode              = "Sharding"
  balancer_status   = "stop"

  datastore {
    type           = "DDS-Community"
    version        = "4.0"
    storage_engine = "wiredTiger"
  }

  configuration {
    type = "mongos"
    id   = huaweicloud_dds_parameter_template.mongos2.id
  }

  configuration {
    type = "shard"
    id   = huaweicloud_dds_parameter_template.shard2.id
  }

  configuration {
    type = "config"
    id   = huaweicloud_dds_parameter_template.config2.id
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.s6.large.4.mongos"
  }

  flavor {
    type      = "shard"
    num       = 2
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s6.large.2.shard"
  }

  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s6.large.2.config"
  }
}`, common.TestBaseNetwork(rName), templateMongos2, templateShard2, templateConfig2, rName)
}

func testAccDDSInstanceV3Config_withReplicaSetConfiguration(rName string) string {
	templateRreplica1 := testAAccDDSInstance_templete(rName, "replica", 1, 400)
	return fmt.Sprintf(`
%[1]s

%[2]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dds_instance" "instance" {
  name                  = "%[3]s"
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  password              = "Terraform@123"
  mode                  = "ReplicaSet"
  replica_set_name      = "replicaName"
  client_network_ranges = ["192.168.0.0/24"]
  maintain_begin        = "01:00"
  maintain_end          = "03:00"

  datastore {
    type           = "DDS-Community"
    version        = "4.0"
    storage_engine = "wiredTiger"
  }
  configuration {
    type = "replica"
    id   = huaweicloud_dds_parameter_template.replica1.id
  }
  flavor {
    type      = "replica"
    storage   = "ULTRAHIGH"
    num       = 3
    size      = 20
    spec_code = "dds.mongodb.s6.large.2.repset"
  }
}`, common.TestBaseNetwork(rName), templateRreplica1, rName)
}

func testAccDDSInstanceV3Config_withReplicaSetConfigurationUpdate(rName string) string {
	templateRreplica2 := testAAccDDSInstance_templete(rName, "replica", 2, 700)
	return fmt.Sprintf(`
%[1]s

%[2]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dds_instance" "instance" {
  name              = "%[3]s-update"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "Terraform@1234"
  mode              = "ReplicaSet"
  replica_set_name  = "replicaNameUpdate"
  maintain_begin    = "02:00"
  maintain_end      = "03:00"
  description       = "test"

  datastore {
    type           = "DDS-Community"
    version        = "4.0"
    storage_engine = "wiredTiger"
  }
  configuration {
    type = "replica"
    id   = huaweicloud_dds_parameter_template.replica2.id
  }
  flavor {
    type      = "replica"
    storage   = "ULTRAHIGH"
    num       = 5
    size      = 30
    spec_code = "dds.mongodb.s6.large.4.repset"
  }
}`, common.TestBaseNetwork(rName), templateRreplica2, rName)
}

func testAccDDSInstanceV3Config_secondLevelMonitoring(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dds_instance" "test" {
  name                            = "%s"
  availability_zone               = data.huaweicloud_availability_zones.test.names[0]
  vpc_id                          = huaweicloud_vpc.test.id
  subnet_id                       = huaweicloud_vpc_subnet.test.id
  security_group_id               = huaweicloud_networking_secgroup.test.id
  password                        = "Terraform@123"
  mode                            = "Sharding"
  second_level_monitoring_enabled = true

  datastore {
    type           = "DDS-Community"
    version        = "4.0"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.s6.xlarge.2.mongos"
  }

  flavor {
    type      = "shard"
    num       = 2
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s6.xlarge.2.shard"
  }

  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s6.xlarge.2.config"
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = "8"
    period     = "1,3"
  }
}`, common.TestBaseNetwork(rName), rName)
}

func testAccDDSInstanceV3Config_secondLevelMonitoringUpdate(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dds_instance" "test" {
  name                            = "%s"
  availability_zone               = data.huaweicloud_availability_zones.test.names[0]
  vpc_id                          = huaweicloud_vpc.test.id
  subnet_id                       = huaweicloud_vpc_subnet.test.id
  security_group_id               = huaweicloud_networking_secgroup.test.id
  password                        = "Terraform@123"
  mode                            = "Sharding"
  second_level_monitoring_enabled = false

  datastore {
    type           = "DDS-Community"
    version        = "4.0"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.s6.xlarge.2.mongos"
  }
  
  flavor {
    type      = "shard"
    num       = 2
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s6.xlarge.2.shard"
  }

  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.s6.xlarge.2.config"
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = "8"
    period     = "1,3"
  }
}`, common.TestBaseNetwork(rName), rName)
}

func testAccDDSInstanceV3Config_updateAZ(rName string) string {
	templateRreplica1 := testAAccDDSInstance_templete(rName, "replica", 1, 400)
	return fmt.Sprintf(`
%[1]s

%[2]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dds_instance" "instance" {
  name                  = "%[3]s"
  availability_zone     = data.huaweicloud_availability_zones.test.names[1]
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  password              = "Terraform@123"
  mode                  = "ReplicaSet"
  replica_set_name      = "replicaName"
  client_network_ranges = ["192.168.0.0/24"]
  maintain_begin        = "01:00"
  maintain_end          = "03:00"

  datastore {
    type           = "DDS-Community"
    version        = "4.0"
    storage_engine = "wiredTiger"
  }
  configuration {
    type = "replica"
    id   = huaweicloud_dds_parameter_template.replica1.id
  }
  flavor {
    type      = "replica"
    storage   = "ULTRAHIGH"
    num       = 3
    size      = 20
    spec_code = "dds.mongodb.s6.large.2.repset"
  }
}`, common.TestBaseNetwork(rName), templateRreplica1, rName)
}

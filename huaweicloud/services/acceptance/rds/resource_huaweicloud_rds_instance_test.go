package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/rds/v3/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rds"
)

func TestAccRdsInstance_basic(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "huaweicloud_rds_instance"
	resourceName := "huaweicloud_rds_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "test_description"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "1"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.n1.large.2"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC+08:00"),
					resource.TestCheckResourceAttr(resourceName, "fixed_ip", "192.168.0.210"),
					resource.TestCheckResourceAttr(resourceName, "private_ips.0", "192.168.0.210"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8634"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "06:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "09:00"),
				),
			},
			{
				Config: testAccRdsInstance_update(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", name)),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "2"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.n1.large.2"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "100"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_updated"),
					resource.TestCheckResourceAttr(resourceName, "fixed_ip", "192.168.0.230"),
					resource.TestCheckResourceAttr(resourceName, "private_ips.0", "192.168.0.230"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8636"),
					resource.TestCheckResourceAttr(resourceName, "maintain_begin", "15:00"),
					resource.TestCheckResourceAttr(resourceName, "maintain_end", "17:00"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttrSet(resourceName, "db.0.password"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"db",
					"status",
					"availability_zone",
				},
			},
		},
	})
}

func TestAccRdsInstance_ha(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "huaweicloud_rds_instance"
	resourceName := "huaweicloud_rds_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_ha(name, "async", "availability"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "1"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.n1.large.2.ha"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC+08:00"),
					resource.TestCheckResourceAttr(resourceName, "ha_replication_mode", "async"),
					resource.TestCheckResourceAttr(resourceName, "switch_strategy", "availability"),
				),
			},
			{
				Config: testAccRdsInstance_ha(name, "sync", "reliability"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "1"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.n1.large.2.ha"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC+08:00"),
					resource.TestCheckResourceAttr(resourceName, "ha_replication_mode", "sync"),
					resource.TestCheckResourceAttr(resourceName, "switch_strategy", "reliability"),
				),
			},
		},
	})
}

func TestAccRdsInstance_mysql(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	resourceType := "huaweicloud_rds_instance"
	resourceName := "huaweicloud_rds_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_mysql_step1(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_rds_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.limit_size", "400"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.trigger_threshold", "15"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.name", "div_precision_increment"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.value", "12"),
					resource.TestCheckResourceAttr(resourceName, "binlog_retention_hours", "12"),
				),
			},
			{
				Config: testAccRdsInstance_mysql_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_rds_flavors.test", "flavors.1.name"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.limit_size", "500"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.trigger_threshold", "20"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "3308"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.name", "connect_timeout"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.value", "14"),
					resource.TestCheckResourceAttr(resourceName, "binlog_retention_hours", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "db.0.password"),
				),
			},
			{
				Config: testAccRdsInstance_mysql_step3(updateName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "volume.0.limit_size", "0"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.trigger_threshold", "0"),
					resource.TestCheckResourceAttr(resourceName, "binlog_retention_hours", "6"),
				),
			},
		},
	})
}

func TestAccRdsInstance_sqlserver(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "huaweicloud_rds_instance"
	resourceName := "huaweicloud_rds_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_sqlserver(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "collation", "Chinese_PRC_CI_AS"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8634"),
				),
			},
			{
				Config: testAccRdsInstance_sqlserver_update(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "collation", "Chinese_PRC_CI_AI"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8634"),
				),
			},
		},
	})
}

func TestAccRdsInstance_mariadb(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "huaweicloud_rds_instance"
	resourceName := "huaweicloud_rds_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_mariadb(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_rds_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "3306"),
					resource.TestCheckResourceAttrSet(resourceName, "db.0.password"),
				),
			},
		},
	})
}

func TestAccRdsInstance_prePaid(t *testing.T) {
	var (
		instance instances.RdsInstanceResponse

		resourceType = "huaweicloud_rds_instance"
		resourceName = "huaweicloud_rds_instance.test"
		name         = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_prePaid(name, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
				),
			},
			{
				Config: testAccRdsInstance_prePaid_update(name, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "60"),
				),
			},
		},
	})
}

func TestAccRdsInstance_restore_mysql(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "huaweicloud_rds_instance"
	resourceName := "huaweicloud_rds_instance.test_backup"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_restore_mysql(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_rds_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.limit_size", "400"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.trigger_threshold", "15"),
					resource.TestCheckResourceAttr(resourceName, "ssl_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "3306"),
				),
			},
		},
	})
}

func TestAccRdsInstance_restore_sqlserver(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "huaweicloud_rds_instance"
	resourceName := "huaweicloud_rds_instance.test_backup"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_restore_sqlserver(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_rds_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "CLOUDSSD"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8634"),
				),
			},
		},
	})
}

func TestAccRdsInstance_restore_pg(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "huaweicloud_rds_instance"
	resourceName := "huaweicloud_rds_instance.test_backup"
	pwd := fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(2, "!#%^*"),
		acctest.RandIntRange(10, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_restore_pg(name, pwd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_rds_flavors.test", "flavors.0.name"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "CLOUDSSD"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8732"),
				),
			},
		},
	})
}

func testAccCheckRdsInstanceDestroy(rsType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.RdsV3Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating rds client: %s", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != rsType {
				continue
			}

			id := rs.Primary.ID
			instance, err := rds.GetRdsInstanceByID(client, id)
			if err != nil {
				return err
			}
			if instance.Id != "" {
				return fmt.Errorf("%s (%s) still exists", rsType, id)
			}
		}
		return nil
	}
}

func testAccCheckRdsInstanceExists(name string, instance *instances.RdsInstanceResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id := rs.Primary.ID
		if id == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.RdsV3Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating rds client: %s", err)
		}

		found, err := rds.GetRdsInstanceByID(client, id)
		if err != nil {
			return fmt.Errorf("error checking %s exist, err=%s", name, err)
		}
		if found.Id == "" {
			return fmt.Errorf("resource %s does not exist", name)
		}

		instance = found
		return nil
	}
}

func testAccRdsInstance_base() string {
	return `
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}`
}

func testAccRdsInstance_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  description       = "test_description"
  flavor            = "rds.pg.n1.large.2"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  time_zone         = "UTC+08:00"
  fixed_ip          = "192.168.0.210"
  maintain_begin    = "06:00"
  maintain_end      = "09:00"

  db {
    type     = "PostgreSQL"
    version  = "12"
    port     = 8634
  }
  volume {
    type = "CLOUDSSD"
    size = 50
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, testAccRdsInstance_base(), name)
}

// name, volume.size, backup_strategy, flavor, tags and password will be updated
func testAccRdsInstance_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test" {
  name                  = "%[2]s-update"
  flavor                = "rds.pg.n1.large.2"
  availability_zone     = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id     = data.huaweicloud_networking_secgroup.test.id
  subnet_id             = data.huaweicloud_vpc_subnet.test.id
  vpc_id                = data.huaweicloud_vpc.test.id
  enterprise_project_id = "%[3]s"
  time_zone             = "UTC+08:00"
  fixed_ip              = "192.168.0.230"
  maintain_begin        = "15:00"
  maintain_end          = "17:00"

  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "12"
    port     = 8636
  }
  volume {
    type = "CLOUDSSD"
    size = 100
  }
  backup_strategy {
    start_time = "09:00-10:00"
    keep_days  = 2
  }

  tags = {
    key1 = "value"
    foo  = "bar_updated"
  }
}
`, testAccRdsInstance_base(), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccRdsInstance_ha(name, replicationMode, switchStrategy string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test" {
  name                = "%[2]s"
  flavor              = "rds.pg.n1.large.2.ha"
  security_group_id   = data.huaweicloud_networking_secgroup.test.id
  subnet_id           = data.huaweicloud_vpc_subnet.test.id
  vpc_id              = data.huaweicloud_vpc.test.id
  time_zone           = "UTC+08:00"
  ha_replication_mode = "%[3]s"
  switch_strategy     = "%[4]s"
  availability_zone   = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[1],
  ]

  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "12"
    port     = 8634
  }
  volume {
    type = "CLOUDSSD"
    size = 50
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, testAccRdsInstance_base(), name, replicationMode, switchStrategy)
}

// if the instance flavor has been changed, then a temp instance will be kept for 12 hours,
// the binding relationship between instance and security group or subnet cannot be unbound
// when deleting the instance in this period time, so we cannot create a new vpc, subnet and
// security group in the test case, otherwise, they cannot be deleted when destroy the resource
func testAccRdsInstance_mysql_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
}

resource "huaweicloud_rds_instance" "test" {
  name                   = "%[2]s"
  flavor                 = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id      = data.huaweicloud_networking_secgroup.test.id
  subnet_id              = data.huaweicloud_vpc_subnet.test.id
  vpc_id                 = data.huaweicloud_vpc.test.id
  availability_zone      = slice(sort(data.huaweicloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)
  ssl_enable             = true  
  binlog_retention_hours = "12"

  db {
    type     = "MySQL"
    version  = "8.0"
    port     = 3306
  }

  backup_strategy {
    start_time = "08:15-09:15"
    keep_days  = 3
    period     = 1
  }

  volume {
    type              = "CLOUDSSD"
    size              = 40
    limit_size        = 400
    trigger_threshold = 15
  }

  parameters {
    name  = "div_precision_increment"
    value = "12"
  }
}
`, testAccRdsInstance_base(), name)
}

func testAccRdsInstance_mysql_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

data "huaweicloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
}

resource "huaweicloud_rds_instance" "test" {
  name                   = "%[3]s"
  flavor                 = data.huaweicloud_rds_flavors.test.flavors[1].name
  security_group_id      = data.huaweicloud_networking_secgroup.test.id
  subnet_id              = data.huaweicloud_vpc_subnet.test.id
  vpc_id                 = data.huaweicloud_vpc.test.id
  availability_zone      = slice(sort(data.huaweicloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)
  ssl_enable             = false
  param_group_id         = huaweicloud_rds_parametergroup.pg_1.id
  binlog_retention_hours = "0"

  db {
    password = "Huangwei!120521"
    type     = "MySQL"
    version  = "8.0"
    port     = 3308
  }

  backup_strategy {
    start_time = "18:15-19:15"
    keep_days  = 5
    period     = 3
  }

  volume {
    type              = "CLOUDSSD"
    size              = 40
    limit_size        = 500
    trigger_threshold = 20
  }

  parameters {
    name  = "connect_timeout"
    value = "14"
  }
}
`, testAccRdsInstance_base(), testAccRdsConfig_basic(name), name)
}

func testAccRdsInstance_mysql_step3(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

data "huaweicloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
}

resource "huaweicloud_rds_instance" "test" {
  name                   = "%[3]s"
  flavor                 = data.huaweicloud_rds_flavors.test.flavors[1].name
  security_group_id      = data.huaweicloud_networking_secgroup.test.id
  subnet_id              = data.huaweicloud_vpc_subnet.test.id
  vpc_id                 = data.huaweicloud_vpc.test.id
  availability_zone      = slice(sort(data.huaweicloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)
  ssl_enable             = false
  param_group_id         = huaweicloud_rds_parametergroup.pg_1.id
  binlog_retention_hours = "6"

  db {
    password = "Huangwei!120521"
    type     = "MySQL"
    version  = "8.0"
    port     = 3308
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }

  parameters {
    name  = "connect_timeout"
    value = "14"
  }
}
`, testAccRdsInstance_base(), testAccRdsConfig_basic(name), name)
}

func testAccRdsInstance_sqlserver(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_networking_secgroup_rule" "ingress" {
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = 8634
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = huaweicloud_networking_secgroup.test.id
}

data "huaweicloud_rds_flavors" "test" {
  db_type       = "SQLServer"
  db_version    = "2017_EE"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
  collation         = "Chinese_PRC_CI_AS"

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0],
  ]

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2017_EE"
    port     = 8634
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, common.TestBaseNetwork(name), name)
}

func testAccRdsInstance_sqlserver_update(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_networking_secgroup_rule" "ingress" {
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = 8634
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = huaweicloud_networking_secgroup.test.id
}

data "huaweicloud_rds_flavors" "test" {
  db_type       = "SQLServer"
  db_version    = "2017_EE"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
  collation         = "Chinese_PRC_CI_AI"

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0],
  ]

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2017_EE"
    port     = 8634
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, common.TestBaseNetwork(name), name)
}

func testAccRdsInstance_mariadb(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_flavors" "test" {
  db_type       = "MariaDB"
  db_version    = "10.5"
  instance_mode = "single"
  group_type    = "dedicated"
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  availability_zone = slice(sort(data.huaweicloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)

  db {
    password = "Huangwei!120521"
    type     = "MariaDB"
    version  = "10.5"
    port     = 3306
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, testAccRdsInstance_base(), name)
}

func testAccRdsInstance_prePaid(name string, isAutoRenew bool) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_flavors" "test" {
  db_type       = "SQLServer"
  db_version    = "2019_SE"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "huaweicloud_rds_instance" "test" {
  vpc_id            = data.huaweicloud_vpc.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  
  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0],
  ]

  name      = "%[2]s"
  flavor    = data.huaweicloud_rds_flavors.test.flavors[0].name
  collation = "Chinese_PRC_CI_AS"

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2019_SE"
    port     = 8638
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "%[3]v"
}
`, testAccRdsInstance_base(), name, isAutoRenew)
}

func testAccRdsInstance_prePaid_update(name string, isAutoRenew bool) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_flavors" "test" {
  db_type       = "SQLServer"
  db_version    = "2019_SE"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "huaweicloud_rds_instance" "test" {
  vpc_id            = data.huaweicloud_vpc.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  
  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0],
  ]

  name      = "%[2]s"
  flavor    = data.huaweicloud_rds_flavors.test.flavors[0].name
  collation = "Chinese_PRC_CI_AS"

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2019_SE"
    port     = 8638
  }

  volume {
    type = "CLOUDSSD"
    size = 60
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "%[3]v"
}
`, testAccRdsInstance_base(), name, isAutoRenew)
}

func testAccRdsInstance_configuration(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "5.7"
  instance_mode = "single"
  group_type    = "dedicated"
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0],
  ]

  db {
    password = "Huangwei!120521"
    type     = "MySQL"
    version  = "5.7"
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }

  parameters {
    name  = "div_precision_increment"
    value = "12"
  }
}
`, testAccRdsInstance_base(), name)
}

func testAccRdsInstance_configuration_update(name string) string {
	return fmt.Sprintf(`
%s

%s

data "huaweicloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "5.7"
  instance_mode = "single"
  group_type    = "dedicated"
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  param_group_id    = huaweicloud_rds_parametergroup.pg_1.id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0],
  ]

  db {
    password = "Huangwei!120521"
    type     = "MySQL"
    version  = "5.7"
    port     = 3306
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }

  parameters {
    name  = "div_precision_increment"
    value = "12"
  }
}
`, testAccRdsInstance_base(), testAccRdsConfig_basic(name), name)
}

func testAccRdsInstance_parameters(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_instance" "test" {
  name                = "%s"
  flavor              = "rds.mysql.sld4.large.ha"
  security_group_id   = data.huaweicloud_networking_secgroup.test.id
  subnet_id           = data.huaweicloud_vpc_subnet.test.id
  vpc_id              = data.huaweicloud_vpc.test.id
  ha_replication_mode = "semisync"

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[3],
  ]

  db {
    password = "Huangwei!120521"
    type     = "MySQL"
    version  = "5.7"
    port     = 3306
  }

  volume {
    type = "LOCALSSD"
    size = 40
  }

  parameters {
    name  = "div_precision_increment"
    value = "12"
  }
}
`, testAccRdsInstance_base(), name)
}

func testAccRdsInstance_newParameters(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_instance" "test" {
  name                = "%s"
  flavor              = "rds.mysql.sld4.large.ha"
  security_group_id   = data.huaweicloud_networking_secgroup.test.id
  subnet_id           = data.huaweicloud_vpc_subnet.test.id
  vpc_id              = data.huaweicloud_vpc.test.id
  ha_replication_mode = "semisync"

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[3],
  ]

  db {
    password = "Huangwei!120521"
    type     = "MySQL"
    version  = "5.7"
    port     = 3306
  }

  volume {
    type = "LOCALSSD"
    size = 40
  }

  parameters {
    name  = "connect_timeout"
    value = "14"
  }
}
`, testAccRdsInstance_base(), name)
}

func testAccRdsInstance_restore_mysql(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test_backup" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  availability_zone = slice(sort(data.huaweicloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)
  ssl_enable        = true  

  restore {
    instance_id = huaweicloud_rds_backup.test.instance_id
    backup_id   = huaweicloud_rds_backup.test.id
  }

  db {
    password = "Huangwei!120521"
    type     = "MySQL"
    version  = "8.0"
    port     = 3306
  }

  backup_strategy {
    start_time = "08:15-09:15"
    keep_days  = 3
    period     = 1
  }

  volume {
    type              = "CLOUDSSD"
    size              = 50
    limit_size        = 400
    trigger_threshold = 15
  }
}
`, testBackup_mysql_basic(name), name)
}

func testAccRdsInstance_restore_sqlserver(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test_backup" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  availability_zone = slice(sort(data.huaweicloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)

  restore {
    instance_id = huaweicloud_rds_backup.test.instance_id
    backup_id   = huaweicloud_rds_backup.test.id
  }

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2019_SE"
    port     = 8634
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }
}
`, testBackup_sqlserver_basic(name), name)
}

func testAccRdsInstance_restore_pg(name, pwd string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test_backup" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  availability_zone = slice(sort(data.huaweicloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)

  restore {
    instance_id = huaweicloud_rds_backup.test.instance_id
    backup_id   = huaweicloud_rds_backup.test.id
  }

  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "14"
    port     = 8732
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }
}
`, testBackup_pg_basic(name), name, pwd)
}

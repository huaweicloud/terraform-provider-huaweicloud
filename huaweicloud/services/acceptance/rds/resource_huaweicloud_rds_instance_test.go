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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rds"
)

func TestAccRdsInstance_basic(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "huaweicloud_rds_instance"
	resourceName := "huaweicloud_rds_instance.test"
	pwd := fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(2, "!#%^*"),
		acctest.RandIntRange(10, 99))
	newPwd := fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(2, "!#%^*"),
		acctest.RandIntRange(10, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_basic(name, pwd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "1"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.n1.large.2"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC+08:00"),
					resource.TestCheckResourceAttr(resourceName, "fixed_ip", "192.168.0.58"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8635"),
					resource.TestCheckResourceAttr(resourceName, "db.0.password", pwd),
				),
			},
			{
				Config: testAccRdsInstance_update(name, newPwd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", name)),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "2"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.n1.xlarge.2"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "100"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_updated"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8636"),
					resource.TestCheckResourceAttr(resourceName, "db.0.password", newPwd),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"db",
					"status",
				},
			},
		},
	})
}

func TestAccRdsInstance_withEpsId(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "huaweicloud_rds_instance"
	resourceName := "huaweicloud_rds_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_epsId(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
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
				Config: testAccRdsInstance_ha(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "1"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.n1.large.2.ha"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "time_zone", "UTC+08:00"),
					resource.TestCheckResourceAttr(resourceName, "fixed_ip", "192.168.0.58"),
					resource.TestCheckResourceAttr(resourceName, "ha_replication_mode", "async"),
				),
			},
		},
	})
}

func TestAccRdsInstance_mysql(t *testing.T) {
	var instance instances.RdsInstanceResponse
	name := acceptance.RandomAccResourceName()
	resourceType := "huaweicloud_rds_instance"
	resourceName := "huaweicloud_rds_instance.test"
	pwd := fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(2, "!#%^*"),
		acctest.RandIntRange(10, 99))
	newPwd := fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(2, "!#%^*"),
		acctest.RandIntRange(10, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_mysql(name, pwd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.mysql.sld4.large.ha"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "fixed_ip", "192.168.0.58"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "db.0.password", pwd),
				),
			},
			{
				Config: testAccRdsInstance_mysqlUpdate(name, newPwd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.mysql.sld4.large.ha"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "fixed_ip", "192.168.0.58"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "3308"),
					resource.TestCheckResourceAttr(resourceName, "db.0.password", newPwd),
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
	pwd := fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(2, "!#%^*"),
		acctest.RandIntRange(10, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceDestroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstance_sqlserver(name, pwd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "collation", "Chinese_PRC_CI_AS"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "db.0.port", "8635"),
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
		password     = fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(2, "!#%^*"),
			acctest.RandIntRange(10, 99))
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
				Config: testAccRdsInstance_prePaid(name, password, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccRdsInstance_prePaid(name, password, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
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

func testAccRdsInstance_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name          = "%s"
  cidr          = "192.168.0.0/24"
  gateway_ip    = "192.168.0.1"
  primary_dns   = "100.125.1.250"
  secondary_dns = "100.125.21.250"
  vpc_id        = huaweicloud_vpc.test.id

  timeouts {
    delete = "30m"
  }
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%s"

  timeouts {
    delete = "30m"
  }
}
`, name, name, name)
}

func testAccRdsInstance_basic(name, password string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_instance" "test" {
  name              = "%s"
  flavor            = "rds.pg.n1.large.2"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
  time_zone         = "UTC+08:00"
  fixed_ip          = "192.168.0.58"

  db {
    password = "%s"
    type     = "PostgreSQL"
    version  = "12"
    port     = 8635
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
`, testAccRdsInstance_base(name), name, password)
}

// name, volume.size, backup_strategy, flavor, tags and password will be updated
func testAccRdsInstance_update(name, password string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_instance" "test" {
  name              = "%s-update"
  flavor            = "rds.pg.n1.xlarge.2"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
  time_zone         = "UTC+08:00"

  db {
    password = "%s"
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
`, testAccRdsInstance_base(name), name, password)
}

func testAccRdsInstance_epsId(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_instance" "test" {
  name                  = "%s"
  flavor                = "rds.pg.n1.large.2"
  availability_zone     = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id     = huaweicloud_networking_secgroup.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  vpc_id                = huaweicloud_vpc.test.id
  enterprise_project_id = "%s"

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
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
}
`, testAccRdsInstance_base(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccRdsInstance_ha(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_instance" "test" {
  name                = "%s"
  flavor              = "rds.pg.n1.large.2.ha"
  security_group_id   = huaweicloud_networking_secgroup.test.id
  subnet_id           = huaweicloud_vpc_subnet.test.id
  vpc_id              = huaweicloud_vpc.test.id
  time_zone           = "UTC+08:00"
  fixed_ip            = "192.168.0.58"
  ha_replication_mode = "async"
  availability_zone   = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[1],
  ]

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
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, testAccRdsInstance_base(name), name)
}

func testAccRdsInstance_mysql(name, pwd string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_instance" "test" {
  name                = "%s"
  flavor              = "rds.mysql.sld4.large.ha"
  security_group_id   = huaweicloud_networking_secgroup.test.id
  subnet_id           = huaweicloud_vpc_subnet.test.id
  vpc_id              = huaweicloud_vpc.test.id
  fixed_ip            = "192.168.0.58"
  ha_replication_mode = "semisync"

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[3],
  ]

  db {
    password = "%s"
    type     = "MySQL"
    version  = "5.7"
    port     = 3306
  }

  volume {
    type = "LOCALSSD"
    size = 40
  }
}
`, testAccRdsInstance_base(name), name, pwd)
}

func testAccRdsInstance_mysqlUpdate(name, pwd string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_instance" "test" {
  name                = "%s"
  flavor              = "rds.mysql.sld4.large.ha"
  security_group_id   = huaweicloud_networking_secgroup.test.id
  subnet_id           = huaweicloud_vpc_subnet.test.id
  vpc_id              = huaweicloud_vpc.test.id
  fixed_ip            = "192.168.0.58"
  ha_replication_mode = "semisync"

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[3],
  ]

  db {
    password = "%s"
    type     = "MySQL"
    version  = "5.7"
    port     = 3308
  }

  volume {
    type = "LOCALSSD"
    size = 40
  }
}
`, testAccRdsInstance_base(name), name, pwd)
}

func testAccRdsInstance_sqlserver(name, pwd string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_rds_instance" "test" {
  name                = "%s"
  flavor              = "rds.mssql.spec.se.c6.large.4"
  security_group_id   = huaweicloud_networking_secgroup.test.id
  subnet_id           = data.huaweicloud_vpc_subnet.test.id
  vpc_id              = data.huaweicloud_vpc.test.id
  collation           = "Chinese_PRC_CI_AS"

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0],
  ]

  db {
    password = "%s"
    type     = "SQLServer"
    version  = "2014_SE"
    port     = 8635
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
`, testAccRdsInstance_base(name), name, pwd)
}

func testAccRdsInstance_prePaid(name, pwd string, isAutoRenew bool) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = try(slice(data.huaweicloud_availability_zones.test.names, 0, 1), [])

  name      = "%[2]s"
  flavor    = "rds.mssql.spec.se.c6.large.4"
  collation = "Chinese_PRC_CI_AS"

  db {
    password = "%[3]s"
    type     = "SQLServer"
    version  = "2014_SE"
    port     = 8635
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "%[4]v"
}
`, testAccRdsInstance_base(name), name, pwd, isAutoRenew)
}

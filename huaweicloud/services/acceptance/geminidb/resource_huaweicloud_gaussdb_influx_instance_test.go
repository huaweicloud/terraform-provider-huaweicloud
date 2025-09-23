package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/geminidb/v3/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getNosqlInstance(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.GeminiDBV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error creating HuaweiCloud GaussDB client: %s", err)
	}

	found, err := instances.GetInstanceByID(client, state.Primary.ID)
	if err != nil {
		return nil, err
	}
	if found.Id == "" {
		return nil, fmt.Errorf("instance <%s> not found", state.Primary.ID)
	}

	return &found, nil
}

func TestAccGaussInfluxInstance_basic(t *testing.T) {
	var instance instances.GeminiDBInstance
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_gaussdb_influx_instance.test"
	password := acceptance.RandomPassword()
	newPassword := acceptance.RandomPassword()
	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getNosqlInstance,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGaussInfluxInstanceConfig_basic(rName, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "node_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "volume_size", "100"),
					resource.TestCheckResourceAttr(resourceName, "status", "normal"),
				),
			},
			{
				Config: testAccGaussInfluxInstanceConfig_update(rName, newPassword),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "password", newPassword),
					resource.TestCheckResourceAttr(resourceName, "node_num", "4"),
					resource.TestCheckResourceAttr(resourceName, "volume_size", "200"),
					resource.TestCheckResourceAttr(resourceName, "status", "normal"),
				),
			},
		},
	})
}

func TestAccGaussInfluxInstance_prePaid(t *testing.T) {
	var (
		instance instances.GeminiDBInstance

		resourceName = "huaweicloud_gaussdb_influx_instance.test"
		rName        = acceptance.RandomAccResourceName()
		password     = acceptance.RandomPassword()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getNosqlInstance,
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
				Config: testAccGaussInfluxInstanceConfig_prePaid(rName, password, false),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccGaussInfluxInstanceConfig_prePaid(rName, password, true),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
				),
			},
		},
	})
}

func TestAccGaussInfluxInstance_updateWithEpsId(t *testing.T) {
	var instance instances.GeminiDBInstance
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_gaussdb_influx_instance.test"
	password := acceptance.RandomPassword()
	newPassword := acceptance.RandomPassword()
	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getNosqlInstance,
	)
	srcEPS := acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
	destEPS := acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccGaussInfluxInstanceConfig_withEpsId(rName, password, srcEPS),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", srcEPS),
				),
			},
			{
				Config: testAccGaussInfluxInstanceConfig_withEpsId(rName, newPassword, destEPS),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "password", newPassword),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", destEPS),
				),
			},
		},
	})
}

func testAccGaussInfluxInstanceConfig_basic(rName, password string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus             = 2
  engine            = "influxdb"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_gaussdb_influx_instance" "test" {
  name        = "%s"
  password    = "%s"
  flavor      = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[0].name
  volume_size = 100
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id
  node_num    = 3

  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 14
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestBaseNetwork(rName), rName, password)
}

func testAccGaussInfluxInstanceConfig_update(rName, password string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus             = 2
  engine            = "influxdb"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_gaussdb_influx_instance" "test" {
  name        = "%s-update"
  password    = "%s"
  flavor      = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[0].name
  volume_size = 200
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id
  node_num    = 4

  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 14
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestBaseNetwork(rName), rName, password)
}

func testAccGaussInfluxInstanceConfig_prePaid(rName, password string, isAutoRenew bool) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus             = 2
  engine            = "influxdb"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_gaussdb_influx_instance" "test" {
  name        = "%s"
  password    = "%s"
  flavor      = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[0].name
  volume_size = 200
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id
  node_num    = 4

  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "%v"
}
`, common.TestBaseNetwork(rName), rName, password, isAutoRenew)
}

func testAccGaussInfluxInstanceConfig_withEpsId(rName, password, epsId string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus             = 2
  engine            = "influxdb"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_gaussdb_influx_instance" "test" {
  name                  = "%s"
  password              = "%s"
  flavor                = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[0].name
  volume_size           = 100
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  node_num              = 3
  enterprise_project_id = "%s"

  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 14
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestBaseNetwork(rName), rName, password, epsId)
}

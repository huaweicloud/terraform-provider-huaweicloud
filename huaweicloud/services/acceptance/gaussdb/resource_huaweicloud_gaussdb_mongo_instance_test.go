package gaussdb

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/chnsz/golangsdk/openstack/geminidb/v3/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccGaussMongoInstance_basic(t *testing.T) {
	var instance instances.GeminiDBInstance
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_gaussdb_mongo_instance.test"
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
				Config: testAccGaussMongoInstanceConfig_basic(rName, password),
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
				Config: testAccGaussMongoInstanceConfig_update(rName, newPassword),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "password", newPassword),
					resource.TestCheckResourceAttr(resourceName, "node_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "volume_size", "200"),
					resource.TestCheckResourceAttr(resourceName, "status", "normal"),
				),
			},
		},
	})
}

func TestAccGaussMongoInstance_prePaid(t *testing.T) {
	var (
		instance instances.GeminiDBInstance

		resourceName = "huaweicloud_gaussdb_mongo_instance.test"
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
				Config: testAccGaussMongoInstanceConfig_prePaid(rName, password, false),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccGaussMongoInstanceConfig_prePaid(rName, password, true),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
				),
			},
		},
	})
}

func testAccGaussMongoInstanceConfig_basic(rName, password string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus             = 4
  engine            = "mongodb"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_gaussdb_mongo_instance" "test" {
  name        = "%s"
  password    = "%s"
  flavor      = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[1].name
  volume_size = 100
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id
  node_num    = 3

  security_group_id = data.huaweicloud_networking_secgroup.test.id
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
`, testAccVpcConfig_Base(rName), rName, password)
}

func testAccGaussMongoInstanceConfig_update(rName, password string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus             = 4
  engine            = "mongodb"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_gaussdb_mongo_instance" "test" {
  name        = "%s-update"
  password    = "%s"
  flavor      = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[1].name
  volume_size = 200
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id
  node_num    = 3

  security_group_id = data.huaweicloud_networking_secgroup.test.id
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
`, testAccVpcConfig_Base(rName), rName, password)
}

func testAccGaussMongoInstanceConfig_prePaid(rName, password string, isAutoRenew bool) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus             = 4
  engine            = "mongodb"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_gaussdb_mongo_instance" "test" {
  name        = "%[2]s"
  password    = "%[3]s"
  flavor      = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[1].name
  volume_size = 200
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id
  node_num    = 3

  security_group_id = data.huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "%[4]v"
}
`, testAccVpcConfig_Base(rName), rName, password, isAutoRenew)
}

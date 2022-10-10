package gaussdb

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/chnsz/golangsdk/openstack/opengauss/v3/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getOpenGaussInstanceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.OpenGaussV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error creating HuaweiCloud GaussDB client: %s", err)
	}
	return instances.GetInstanceByID(client, state.Primary.ID)
}

func TestAccOpenGaussInstance_basic(t *testing.T) {
	var (
		instance     instances.GaussDBInstance
		resourceName = "huaweicloud_gaussdb_opengauss_instance.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
		password     = fmt.Sprintf("%s@123", acctest.RandString(5))
		newPassword  = fmt.Sprintf("%sUpdate@123", acctest.RandString(5))
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getOpenGaussInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccOpenGaussInstance_basic(rName, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "ha.0.mode", "enterprise"),
					resource.TestCheckResourceAttr(resourceName, "ha.0.replication_mode", "sync"),
					resource.TestCheckResourceAttr(resourceName, "ha.0.consistency", "strong"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "ULTRAHIGH"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "sharding_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "coordinator_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
				),
			},
			{
				Config: testAccOpenGaussInstance_update(rName, newPassword),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "password", newPassword),
					resource.TestCheckResourceAttr(resourceName, "sharding_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "coordinator_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "80"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.start_time", "08:00-09:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "8"),
				),
			},
		},
	})
}

func TestAccOpenGaussInstance_prepaid(t *testing.T) {
	var (
		instance     instances.GaussDBInstance
		resourceName = "huaweicloud_gaussdb_opengauss_instance.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
		password     = acceptance.RandomPassword()
		newPassword  = acceptance.RandomPassword()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getOpenGaussInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccOpenGaussInstance_prepaid(rName, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "ha.0.mode", "enterprise"),
					resource.TestCheckResourceAttr(resourceName, "ha.0.replication_mode", "sync"),
					resource.TestCheckResourceAttr(resourceName, "ha.0.consistency", "strong"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "ULTRAHIGH"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "sharding_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "coordinator_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
				),
			},
			{
				Config: testAccOpenGaussInstance_prepaidUpdate(rName, newPassword),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "password", newPassword),
					resource.TestCheckResourceAttr(resourceName, "sharding_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "coordinator_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "80"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.start_time", "08:00-09:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "8"),
				),
			},
		},
	})
}

func TestAccOpenGaussInstance_haModeCentralized(t *testing.T) {
	var (
		instance     instances.GaussDBInstance
		resourceName = "huaweicloud_gaussdb_opengauss_instance.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
		password     = fmt.Sprintf("%s@123", acctest.RandString(5))
		newPassword  = fmt.Sprintf("%sUpdate@123", acctest.RandString(5))
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getOpenGaussInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccOpenGaussInstance_haModeCentralized(rName, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id", "huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "flavor", "gaussdb.opengauss.ee.m6.2xlarge.x868.ha"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "ha.0.mode", "centralization_standard"),
					resource.TestCheckResourceAttr(resourceName, "ha.0.replication_mode", "sync"),
					resource.TestCheckResourceAttr(resourceName, "ha.0.consistency", "strong"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "ULTRAHIGH"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "replica_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "40"),
				),
			},
			{
				Config: testAccOpenGaussInstance_haModeCentralizedUpdate(rName, newPassword),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "password", newPassword),
					resource.TestCheckResourceAttr(resourceName, "replica_num", "3"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "80"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.start_time", "08:00-09:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "8"),
				),
			},
		},
	})
}

func testAccOpenGaussInstance_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id = huaweicloud_vpc.test.id

  name       = "%[1]s"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1), 1)

  timeouts {
    delete = "40m"
  }
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[1]s"

  timeouts {
    delete = "40m"
  }
}

resource "huaweicloud_networking_secgroup_rule" "in_v4_tcp_opengauss" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  direction         = "ingress"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
}
`, rName)
}

func testAccOpenGaussInstance_basic(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_gaussdb_opengauss_instance" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"
  name              = "%[2]s"
  password          = "%[3]s"
  sharding_num      = 1
  coordinator_num   = 2
  availability_zone = "${data.huaweicloud_availability_zones.test.names[0]},${data.huaweicloud_availability_zones.test.names[0]},${data.huaweicloud_availability_zones.test.names[0]}"

  ha {
    mode             = "enterprise"
    replication_mode = "sync"
    consistency      = "strong"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
`, testAccOpenGaussInstance_base(rName), rName, password)
}

func testAccOpenGaussInstance_update(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_gaussdb_opengauss_instance" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"
  name              = "%[2]s-update"
  password          = "%[3]s"
  sharding_num      = 2
  coordinator_num   = 3
  availability_zone = "${data.huaweicloud_availability_zones.test.names[0]},${data.huaweicloud_availability_zones.test.names[0]},${data.huaweicloud_availability_zones.test.names[0]}"

  ha {
    mode             = "enterprise"
    replication_mode = "sync"
    consistency      = "strong"
  }

  volume {
    type = "ULTRAHIGH"
    size = 80
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 8
  }
}
`, testAccOpenGaussInstance_base(rName), rName, password)
}

func testAccOpenGaussInstance_prepaid(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_gaussdb_opengauss_instance" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"
  name              = "%[2]s"
  password          = "%[3]s"
  sharding_num      = 1
  coordinator_num   = 2
  availability_zone = "${data.huaweicloud_availability_zones.test.names[0]},${data.huaweicloud_availability_zones.test.names[0]},${data.huaweicloud_availability_zones.test.names[0]}"

  ha {
    mode             = "enterprise"
    replication_mode = "sync"
    consistency      = "strong"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"

  timeouts {
    update = "2h"
  }
}
`, testAccOpenGaussInstance_base(rName), rName, password)
}

func testAccOpenGaussInstance_prepaidUpdate(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_gaussdb_opengauss_instance" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"
  name              = "%[2]s-update"
  password          = "%[3]s"
  sharding_num      = 2
  coordinator_num   = 3
  availability_zone = "${data.huaweicloud_availability_zones.test.names[0]},${data.huaweicloud_availability_zones.test.names[0]},${data.huaweicloud_availability_zones.test.names[0]}"

  ha {
    mode             = "enterprise"
    replication_mode = "sync"
    consistency      = "strong"
  }

  volume {
    type = "ULTRAHIGH"
    size = 80
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 8
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "false"

  timeouts {
    update = "2h"
  }
}
`, testAccOpenGaussInstance_base(rName), rName, password)
}

func testAccOpenGaussInstance_haModeCentralized(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_gaussdb_opengauss_instance" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = "gaussdb.opengauss.ee.m6.2xlarge.x868.ha"
  name              = "%[2]s"
  password          = "%[3]s"
  replica_num       = 3
  availability_zone = "${data.huaweicloud_availability_zones.test.names[0]},${data.huaweicloud_availability_zones.test.names[0]},${data.huaweicloud_availability_zones.test.names[0]}"

  ha {
    mode             = "centralization_standard"
    replication_mode = "sync"
    consistency      = "strong"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
`, testAccOpenGaussInstance_base(rName), rName, password)
}

func testAccOpenGaussInstance_haModeCentralizedUpdate(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_gaussdb_opengauss_instance" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = "gaussdb.opengauss.ee.m6.2xlarge.x868.ha"
  name              = "%[2]s-update"
  password          = "%[3]s"
  replica_num       = 3
  availability_zone = "${data.huaweicloud_availability_zones.test.names[0]},${data.huaweicloud_availability_zones.test.names[0]},${data.huaweicloud_availability_zones.test.names[0]}"

  ha {
    mode             = "centralization_standard"
    replication_mode = "sync"
    consistency      = "strong"
  }

  volume {
    type = "ULTRAHIGH"
    size = 80
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 8
  }
}
`, testAccOpenGaussInstance_base(rName), rName, password)
}

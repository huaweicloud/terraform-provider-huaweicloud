package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/drs/v3/jobs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getDrsJobResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.DrsV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DRS client, err: %s", err)
	}
	detailResp, err := jobs.Get(client, jobs.QueryJobReq{Jobs: []string{state.Primary.ID}})
	if err != nil {
		return nil, err
	}
	status := detailResp.Results[0].Status
	if status == "DELETED" {
		return nil, golangsdk.ErrDefault404{}
	}
	return detailResp, nil
}

func TestAccResourceDrsJob_basic(t *testing.T) {
	var obj jobs.BatchCreateJobReq
	resourceName := "huaweicloud_drs_job.test"
	name := acceptance.RandomAccResourceName()
	dbName := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	pwd := "TestDrs@123"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDrsJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDrsJob_migrate_mysql(name, dbName, pwd, ""),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", "migration"),
					resource.TestCheckResourceAttr(resourceName, "direction", "up"),
					resource.TestCheckResourceAttr(resourceName, "net_type", "eip"),
					resource.TestCheckResourceAttr(resourceName, "destination_db_readnoly", "true"),
					resource.TestCheckResourceAttr(resourceName, "migration_type", "FULL_INCR_TRANS"),
					resource.TestCheckResourceAttr(resourceName, "description", name),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.ip", "192.168.0.58"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.user", "root"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.ip", "192.168.0.59"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.user", "root"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.instance_id",
						"huaweicloud_rds_instance.test2", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.region",
						"huaweicloud_rds_instance.test2", "region"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", name),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "postPaid"),
				),
			},
			{
				Config: testAccDrsJob_migrate_mysql(updateName, dbName, pwd, "stop"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "type", "migration"),
					resource.TestCheckResourceAttr(resourceName, "direction", "up"),
					resource.TestCheckResourceAttr(resourceName, "net_type", "eip"),
					resource.TestCheckResourceAttr(resourceName, "destination_db_readnoly", "true"),
					resource.TestCheckResourceAttr(resourceName, "migration_type", "FULL_INCR_TRANS"),
					resource.TestCheckResourceAttr(resourceName, "description", updateName),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.ip", "192.168.0.58"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.user", "root"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.ip", "192.168.0.59"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.user", "root"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.instance_id",
						"huaweicloud_rds_instance.test2", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.region",
						"huaweicloud_rds_instance.test2", "region"),
					resource.TestCheckResourceAttr(resourceName, "status", "PAUSING"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", updateName),
				),
			},
			{
				Config: testAccDrsJob_migrate_mysql(updateName, dbName, pwd, "restart"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"source_db.0.password", "destination_db.0.password",
					"expired_days", "migrate_definer", "force_destroy", "action", "updated_at",
					"source_db.0.ip", "destination_db.0.ip", "engine_type"},
			},
		},
	})
}

func testAccDrsJob_mysql(index int, name, pwd, ip string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_instance" "test%d" {
  depends_on = [
    huaweicloud_networking_secgroup_rule.ingress,
    huaweicloud_networking_secgroup_rule.egress,
  ]
  name                = "%s%d"
  flavor              = "rds.mysql.x1.large.2.ha"
  security_group_id   = huaweicloud_networking_secgroup.test.id
  subnet_id           = huaweicloud_vpc_subnet.test.id
  vpc_id              = huaweicloud_vpc.test.id
  fixed_ip            = "%s"
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
    type = "CLOUDSSD"
    size = 40
  }
}
`, index, name, index, ip, pwd)
}

func testAccDrsJob_kafka(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_dms_kafka_flavors" "test" {
  type      = "cluster"
  flavor_id = "c6.2u4g.cluster"
}

locals {
  flavor = data.huaweicloud_dms_kafka_flavors.test.flavors[0]
  ipList = split(",", huaweicloud_dms_kafka_instance.test.connect_address)
  port   = "9092"
  ips    = format("%%s:%%s,%%s:%%s", local.ipList[0], local.port, local.ipList[1], local.port)
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name              = "%[1]s"
  vpc_id            = huaweicloud_vpc.test.id
  network_id        = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor_id          = local.flavor.id
  storage_spec_code  = local.flavor.ios[0].storage_spec_code
  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[1],
    data.huaweicloud_availability_zones.test.names[2]
  ]
  engine_version = "2.7"
  storage_space  = local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node
  broker_num     = 3
  arch_type      = "X86"
  ssl_enable     = false
}

resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id = huaweicloud_dms_kafka_instance.test.id
  name        = "%[1]s"
  partitions  = 20
  aging_time  = 72
}
`, name)
}

const testAccSecgroupRule string = `
resource "huaweicloud_networking_secgroup_rule" "ingress" {
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = "3306,9092"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = huaweicloud_networking_secgroup.test.id
}

resource "huaweicloud_networking_secgroup_rule" "egress" {
  direction         = "egress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = huaweicloud_networking_secgroup.test.id
}
`

func testAccDrsJob_migrate_mysql(name, dbName, pwd, action string) string {
	netConfig := common.TestBaseNetwork(name)
	sourceDb := testAccDrsJob_mysql(1, dbName, pwd, "192.168.0.58")
	destDb := testAccDrsJob_mysql(2, dbName, pwd, "192.168.0.59")

	return fmt.Sprintf(`
%s

%s

data "huaweicloud_availability_zones" "test" {}

%s
%s

resource "huaweicloud_drs_job" "test" {
  name           = "%s"
  type           = "migration"
  engine_type    = "mysql"
  direction      = "up"
  net_type       = "eip"
  migration_type = "FULL_INCR_TRANS"
  description    = "%s"
  force_destroy  = true

  source_db {
    engine_type = "mysql"
    ip          = huaweicloud_rds_instance.test1.fixed_ip
    port        = 3306
    user        = "root"
    password    = "%s"
  }


  destination_db {
    region      = huaweicloud_rds_instance.test2.region
    ip          = huaweicloud_rds_instance.test2.fixed_ip
    port        = 3306
    engine_type = "mysql"
    user        = "root"
    password    = "%s"
    instance_id = huaweicloud_rds_instance.test2.id
    subnet_id   = huaweicloud_rds_instance.test2.subnet_id
  }

  tags = {
    key = "%s"
  }

  action = "%s"

  lifecycle {
    ignore_changes = [
      source_db.0.password, destination_db.0.password, force_destroy,
    ]
  }
}
`, netConfig, testAccSecgroupRule, sourceDb, destDb, name, name, pwd, pwd, name, action)
}

func TestAccResourceDrsJob_eip(t *testing.T) {
	var obj jobs.BatchCreateJobReq
	resourceName := "huaweicloud_drs_job.test"
	name := acceptance.RandomAccResourceName()
	dbName := acceptance.RandomAccResourceName()
	pwd := "TestDrs@123"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDrsJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDrsJob_eip(name, dbName, pwd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", "migration"),
					resource.TestCheckResourceAttr(resourceName, "direction", "up"),
					resource.TestCheckResourceAttr(resourceName, "net_type", "eip"),
					resource.TestCheckResourceAttr(resourceName, "destination_db_readnoly", "true"),
					resource.TestCheckResourceAttr(resourceName, "migration_type", "FULL_INCR_TRANS"),
					resource.TestCheckResourceAttr(resourceName, "description", name),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.ip", "192.168.0.58"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.user", "root"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.ip", "192.168.0.59"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.user", "root"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.instance_id",
						"huaweicloud_rds_instance.test2", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.region",
						"huaweicloud_rds_instance.test2", "region"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "postPaid"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"source_db.0.password", "destination_db.0.password",
					"expired_days", "migrate_definer", "force_destroy", "action", "updated_at",
					"source_db.0.ip", "destination_db.0.ip", "engine_type", "public_ip_list", "status"},
			},
		},
	})
}

func testAccDrsJob_eip(name, dbName, pwd string) string {
	netConfig := common.TestBaseNetwork(name)
	sourceDb := testAccDrsJob_mysql(1, dbName, pwd, "192.168.0.58")
	destDb := testAccDrsJob_mysql(2, dbName, pwd, "192.168.0.59")

	return fmt.Sprintf(`
%[1]s

%[2]s

data "huaweicloud_availability_zones" "test" {}

%[3]s

%[4]s

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    name        = "%[5]s"
    size        = 5
    charge_mode = "traffic"
  }
}

resource "huaweicloud_drs_job" "test" {
  name           = "%[5]s"
  type           = "migration"
  engine_type    = "mysql"
  direction      = "up"
  net_type       = "eip"
  migration_type = "FULL_INCR_TRANS"
  description    = "%[5]s"
  force_destroy  = true

  source_db {
    engine_type = "mysql"
    ip          = huaweicloud_rds_instance.test1.fixed_ip
    port        = 3306
    user        = "root"
    password    = "%[6]s"
  }

  destination_db {
    region      = huaweicloud_rds_instance.test2.region
    ip          = huaweicloud_rds_instance.test2.fixed_ip
    port        = 3306
    engine_type = "mysql"
    user        = "root"
    password    = "%[6]s"
    instance_id = huaweicloud_rds_instance.test2.id
    subnet_id   = huaweicloud_rds_instance.test2.subnet_id
  }

  public_ip_list {
    id        = huaweicloud_vpc_eip.test.id
    public_ip = huaweicloud_vpc_eip.test.address
    type      = "master"
  }

  lifecycle {
    ignore_changes = [
      source_db.0.password, destination_db.0.password, force_destroy,
    ]
  }
}
`, netConfig, testAccSecgroupRule, sourceDb, destDb, name, pwd)
}

func TestAccResourceDrsJob_sync(t *testing.T) {
	var obj jobs.BatchCreateJobReq
	resourceName := "huaweicloud_drs_job.test"
	name := acceptance.RandomAccResourceName()
	dbName := acceptance.RandomAccResourceName()
	pwd := "TestDrs@123"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDrsJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDrsJob_synchronize_mysql(name, dbName, pwd, false),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", "sync"),
					resource.TestCheckResourceAttr(resourceName, "direction", "up"),
					resource.TestCheckResourceAttr(resourceName, "net_type", "vpc"),
					resource.TestCheckResourceAttr(resourceName, "destination_db_readnoly", "true"),
					resource.TestCheckResourceAttr(resourceName, "migration_type", "FULL_INCR_TRANS"),
					resource.TestCheckResourceAttr(resourceName, "description", name),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.ip", "192.168.0.58"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.user", "root"),
					resource.TestCheckResourceAttrPair(resourceName, "source_db.0.vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "source_db.0.subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.ip", "192.168.0.59"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.user", "root"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.instance_id",
						"huaweicloud_rds_instance.test2", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.region",
						"huaweicloud_rds_instance.test2", "region"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "progress"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "month"),
					resource.TestCheckResourceAttr(resourceName, "period", "1"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "order_id"),
				),
			},
			{
				Config: testAccDrsJob_synchronize_mysql(name, dbName, pwd, true),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"source_db.0.password", "destination_db.0.password",
					"expired_days", "migrate_definer", "force_destroy", "status", "auto_renew", "updated_at", "policy_config",
					"source_db.0.ip", "destination_db.0.ip", "engine_type"},
			},
		},
	})
}

func testAccRdsMysqlDatabse(dbname string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_mysql_database" "test" {
  instance_id   = huaweicloud_rds_instance.test1.id
  name          = "%[1]s"
  character_set = "utf8"
}
`, dbname)
}

func testAccDrsJob_synchronize_mysql(name, dbName, pwd string, autoRenew bool) string {
	netConfig := common.TestBaseNetwork(name)
	sourceDb := testAccDrsJob_mysql(1, dbName, pwd, "192.168.0.58")
	destDb := testAccDrsJob_mysql(2, dbName, pwd, "192.168.0.59")
	database := testAccRdsMysqlDatabse(dbName)

	return fmt.Sprintf(`
%[1]s

%[2]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_drs_node_types" "test" {
  engine_type = "mysql"
  type        = "sync"
  direction   = "up"
}

%[3]s

%[4]s

%[5]s

resource "huaweicloud_drs_job" "test" {
  name           = "%[6]s"
  type           = "sync"
  engine_type    = "mysql"
  direction      = "up"
  node_type      = data.huaweicloud_drs_node_types.test.node_types[0]
  net_type       = "vpc"
  migration_type = "FULL_INCR_TRANS"
  description    = "%[6]s"
  force_destroy  = true

  source_db {
    engine_type = "mysql"
    ip          = huaweicloud_rds_instance.test1.fixed_ip
    port        = 3306
    user        = "root"
    password    = "%[7]s"
    vpc_id      = huaweicloud_rds_instance.test1.vpc_id
    subnet_id   = huaweicloud_rds_instance.test1.subnet_id
  }

  destination_db {
    region      = huaweicloud_rds_instance.test2.region
    ip          = huaweicloud_rds_instance.test2.fixed_ip
    port        = 3306
    engine_type = "mysql"
    user        = "root"
    password    = "%[7]s"
    instance_id = huaweicloud_rds_instance.test2.id
    subnet_id   = huaweicloud_rds_instance.test2.subnet_id
  }

  databases = [huaweicloud_rds_mysql_database.test.name]

  policy_config {
    filter_ddl_policy = "drop_database"
    conflict_policy   = "overwrite"
    index_trans       = true
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "%[8]v"

  limit_speed {
    speed      = "15"
    start_time = "16:00"
    end_time   = "17:59"
  }

  lifecycle {
    ignore_changes = [
      source_db.0.password, destination_db.0.password, force_destroy,
    ]
  }
}
`, netConfig, testAccSecgroupRule, sourceDb, destDb, database, name, pwd, autoRenew)
}

func TestAccResourceDrsJob_dualAZ(t *testing.T) {
	var obj jobs.BatchCreateJobReq
	resourceName := "huaweicloud_drs_job.test"
	name := acceptance.RandomAccResourceName()
	dbName := acceptance.RandomAccResourceName()
	pwd := "TestDrs@123"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDrsJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDrsJob_dualAZ(name, dbName, pwd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", "sync"),
					resource.TestCheckResourceAttr(resourceName, "direction", "up"),
					resource.TestCheckResourceAttr(resourceName, "net_type", "vpc"),
					resource.TestCheckResourceAttr(resourceName, "destination_db_readnoly", "true"),
					resource.TestCheckResourceAttr(resourceName, "migration_type", "FULL_INCR_TRANS"),
					resource.TestCheckResourceAttr(resourceName, "description", name),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.ip", "192.168.0.58"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.user", "root"),
					resource.TestCheckResourceAttrPair(resourceName, "source_db.0.vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "source_db.0.subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.ip", "192.168.0.59"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.user", "root"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.instance_id",
						"huaweicloud_rds_instance.test2", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.region",
						"huaweicloud_rds_instance.test2", "region"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "progress"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttrSet(resourceName, "master_job_id"),
					resource.TestCheckResourceAttrSet(resourceName, "slave_job_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"source_db.0.password", "destination_db.0.password",
					"expired_days", "migrate_definer", "force_destroy", "status", "updated_at",
					"source_db.0.ip", "destination_db.0.ip", "engine_type"},
			},
		},
	})
}

func testAccDrsJob_dualAZ(name, dbName, pwd string) string {
	netConfig := common.TestBaseNetwork(name)
	sourceDb := testAccDrsJob_mysql(1, dbName, pwd, "192.168.0.58")
	destDb := testAccDrsJob_mysql(2, dbName, pwd, "192.168.0.59")

	return fmt.Sprintf(`
%[1]s

%[2]s

%[3]s

%[4]s

%[5]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_drs_availability_zones" "test" {
  engine_type = "mysql"
  type        = "sync"
  direction   = "up"
  node_type   = "high"
}

resource "huaweicloud_drs_job" "test" {
  name           = "%[6]s"
  type           = "sync"
  engine_type    = "mysql"
  direction      = "up"
  net_type       = "vpc"
  migration_type = "FULL_INCR_TRANS"
  description    = "%[6]s"
  force_destroy  = true
  master_az      = data.huaweicloud_drs_availability_zones.test.names[0]
  slave_az       = data.huaweicloud_drs_availability_zones.test.names[1]

  source_db {
    engine_type = "mysql"
    ip          = huaweicloud_rds_instance.test1.fixed_ip
    port        = 3306
    user        = "root"
    password    = "%[7]s"
    vpc_id      = huaweicloud_rds_instance.test1.vpc_id
    subnet_id   = huaweicloud_rds_instance.test1.subnet_id
  }

  destination_db {
    region      = huaweicloud_rds_instance.test2.region
    ip          = huaweicloud_rds_instance.test2.fixed_ip
    port        = 3306
    engine_type = "mysql"
    user        = "root"
    password    = "%[7]s"
    instance_id = huaweicloud_rds_instance.test2.id
    subnet_id   = huaweicloud_rds_instance.test2.subnet_id
  }

  databases = [huaweicloud_rds_mysql_database.test.name]

  lifecycle {
    ignore_changes = [
      source_db.0.password, destination_db.0.password, force_destroy,
    ]
  }
}
`, netConfig, testAccSecgroupRule, sourceDb, destDb, testAccRdsMysqlDatabse(dbName), name, pwd)
}

func TestAccResourceDrsJob_mysql_to_kafka(t *testing.T) {
	var obj jobs.BatchCreateJobReq
	resourceName := "huaweicloud_drs_job.test"
	name := acceptance.RandomAccResourceName()
	dbName := acceptance.RandomAccResourceName()
	pwd := "TestDrs@123"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDrsJobResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDrsJob_mysql_to_kafka(name, dbName, pwd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", "sync"),
					resource.TestCheckResourceAttr(resourceName, "engine_type", "mysql-to-kafka"),
					resource.TestCheckResourceAttr(resourceName, "direction", "down"),
					resource.TestCheckResourceAttr(resourceName, "net_type", "vpc"),
					resource.TestCheckResourceAttr(resourceName, "destination_db_readnoly", "false"),
					resource.TestCheckResourceAttr(resourceName, "migration_type", "FULL_INCR_TRANS"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.ip", "192.168.0.58"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.port", "3306"),
					resource.TestCheckResourceAttr(resourceName, "source_db.0.user", "root"),
					resource.TestCheckResourceAttrPair(resourceName, "source_db.0.instance_id",
						"huaweicloud_rds_instance.test1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "source_db.0.vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "source_db.0.subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "destination_db.0.engine_type", "kafka"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "destination_db.0.subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "progress"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "postPaid"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"source_db.0.password", "destination_db.0.password",
					"expired_days", "migrate_definer", "force_destroy", "status", "updated_at", "policy_config",
					"source_db.0.ip", "destination_db.0.ip", "engine_type"},
			},
		},
	})
}

func testAccDrsJob_mysql_to_kafka(name, dbName, pwd string) string {
	netConfig := common.TestBaseNetwork(name)
	sourceDb := testAccDrsJob_mysql(1, dbName, pwd, "192.168.0.58")
	destDb := testAccDrsJob_kafka(dbName)

	return fmt.Sprintf(`
%[1]s

%[2]s

%[3]s

%[4]s

%[5]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_drs_job" "test" {
  name                    = "%[6]s"
  type                    = "sync"
  engine_type             = "mysql-to-kafka"
  direction               = "down"
  net_type                = "vpc"
  migration_type          = "FULL_INCR_TRANS"
  force_destroy           = true
  destination_db_readnoly = false

  source_db {
    engine_type = "mysql"
    ip          = huaweicloud_rds_instance.test1.fixed_ip
    port        = 3306
    user        = "root"
    password    = "%[7]s"
    instance_id = huaweicloud_rds_instance.test1.id
    vpc_id      = huaweicloud_rds_instance.test1.vpc_id
    subnet_id   = huaweicloud_rds_instance.test1.subnet_id
  }

  destination_db {
    ip          = local.ips
    engine_type = "kafka"
    vpc_id      = huaweicloud_dms_kafka_instance.test.vpc_id
    subnet_id   = huaweicloud_dms_kafka_instance.test.network_id

    kafka_security_config {
      type = "PLAINTEXT"
    }
  }

  policy_config {
    topic_policy     = "0"
    topic            = huaweicloud_dms_kafka_topic.test.name
    partition_policy = "1"
  }

  databases = [huaweicloud_rds_mysql_database.test.name]

  lifecycle {
    ignore_changes = [
      source_db.0.password, destination_db.0.password, force_destroy,
    ]
  }
}
`, netConfig, testAccSecgroupRule, sourceDb, destDb, testAccRdsMysqlDatabse(dbName), name, pwd)
}

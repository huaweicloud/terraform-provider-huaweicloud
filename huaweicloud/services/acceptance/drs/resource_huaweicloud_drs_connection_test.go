package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/drs"
)

func getConnectionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("drs", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DRS client: %s", err)
	}

	return drs.GetConnectionById(client, state.Primary.ID)
}

func TestAccConnection_mysql(t *testing.T) {
	var (
		connection interface{}
		rName      = "huaweicloud_drs_connection.test"
		name       = acceptance.RandomAccResourceName()
		password   = "Test@123456"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&connection,
		getConnectionResourceFunc,
	)

	resource.ParallelTest(
		t, resource.TestCase{
			PreCheck: func() {
				acceptance.TestAccPreCheck(t)
			},
			ProviderFactories: acceptance.TestAccProviderFactories,
			CheckDestroy:      rc.CheckResourceDestroy(),
			Steps: []resource.TestStep{
				{
					Config: testAccConnection_mysql(name, password),
					Check: resource.ComposeTestCheckFunc(
						rc.CheckResourceExists(),
						resource.TestCheckResourceAttr(rName, "name", name),
						resource.TestCheckResourceAttr(rName, "db_type", "mysql"),
						resource.TestCheckResourceAttr(rName, "description", "Test DRS connection for MySQL"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.endpoint_name", "mysql"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.ip", "192.168.0.100"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.db_port", "3306"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.db_user", "root"),
						resource.TestCheckResourceAttr(rName, "ssl.0.ssl_link", "false"),
						resource.TestCheckResourceAttrSet(rName, "create_time"),
						resource.TestCheckResourceAttrSet(rName, "enterprise_project_id"),
					),
				},
				{
					Config: testAccConnection_mysql_update(name, password),
					Check: resource.ComposeTestCheckFunc(
						rc.CheckResourceExists(),
						resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s-update", name)),
						resource.TestCheckResourceAttr(rName, "description", "Updated DRS connection for MySQL"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.endpoint_name", "mysql"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.ip", "192.168.0.200"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.db_port", "3307"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.db_user", "admin"),
						resource.TestCheckResourceAttr(rName, "ssl.0.ssl_link", "false"),
						resource.TestCheckResourceAttr(rName, "config.0.driver_name", "mysql"),
						resource.TestCheckResourceAttrSet(rName, "enterprise_project_id"),
					),
				},
				{
					ResourceName:      rName,
					ImportState:       true,
					ImportStateVerify: true,
					ImportStateVerifyIgnore: []string{
						"endpoint.0.db_password",
					},
				},
			},
		},
	)
}

func testAccConnection_mysql(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_connection" "test" {
  name        = "%s"
  db_type     = "mysql"
  description = "Test DRS connection for MySQL"

  endpoint {
    endpoint_name = "mysql"
    ip            = "192.168.0.100"
    db_port       = "3306"
    db_user       = "root"
    db_password   = "%s"
  }

  ssl {
    ssl_link = false
  }

  lifecycle {
    ignore_changes = [
      endpoint.0.db_password,
   ]
  }
}
`, name, password)
}

func testAccConnection_mysql_update(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_connection" "test" {
  name        = "%s-update"
  db_type     = "mysql"
  description = "Updated DRS connection for MySQL"

  endpoint {
    endpoint_name = "mysql"
    ip            = "192.168.0.200"
    db_port       = "3307"
    db_user       = "admin"
    db_password   = "%s"
  }

  ssl {
    ssl_link = false
  }

  config {
    driver_name = "mysql"
  }
  
  lifecycle {
    ignore_changes = [
      endpoint.0.db_password,
   ]
  }
}
`, name, password)
}

func TestAccConnection_withRdsInstance(t *testing.T) {
	var (
		connection interface{}
		rName      = "huaweicloud_drs_connection.test"
		name       = acceptance.RandomAccResourceName()
		password   = "TestDrs@123"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&connection,
		getConnectionResourceFunc,
	)

	resource.ParallelTest(
		t, resource.TestCase{
			PreCheck: func() {
				acceptance.TestAccPreCheck(t)
			},
			ProviderFactories: acceptance.TestAccProviderFactories,
			CheckDestroy:      rc.CheckResourceDestroy(),
			Steps: []resource.TestStep{
				{
					Config: testAccConnection_withRdsInstance(name, password),
					Check: resource.ComposeTestCheckFunc(
						rc.CheckResourceExists(),
						resource.TestCheckResourceAttr(rName, "name", name),
						resource.TestCheckResourceAttr(rName, "db_type", "mysql"),
						resource.TestCheckResourceAttr(rName, "description", "Test DRS connection with RDS instance"),
						resource.TestCheckResourceAttrPair(
							rName, "endpoint.0.instance_id",
							"huaweicloud_rds_instance.test", "id",
						),
						resource.TestCheckResourceAttr(rName, "endpoint.0.endpoint_name", "cloud_mysql"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.db_port", "3306"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.db_user", "root"),
						resource.TestCheckResourceAttrPair(
							rName, "vpc.0.vpc_id",
							"huaweicloud_rds_instance.test", "vpc_id",
						),
						resource.TestCheckResourceAttrPair(
							rName, "vpc.0.subnet_id",
							"huaweicloud_rds_instance.test", "subnet_id",
						),
						resource.TestCheckResourceAttr(rName, "ssl.0.ssl_link", "false"),
						resource.TestCheckResourceAttrSet(rName, "create_time"),
						resource.TestCheckResourceAttrSet(rName, "enterprise_project_id"),
					),
				},
				{
					Config: testAccConnection_withRdsInstance_update(name, password),
					Check: resource.ComposeTestCheckFunc(
						rc.CheckResourceExists(),
						resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s-update", name)),
						resource.TestCheckResourceAttr(
							rName, "description", "Updated DRS connection with RDS instance",
						),
						resource.TestCheckResourceAttrPair(
							rName, "endpoint.0.instance_id",
							"huaweicloud_rds_instance.test", "id",
						),
						resource.TestCheckResourceAttr(rName, "endpoint.0.endpoint_name", "cloud_mysql"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.db_port", "3307"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.db_user", "admin"),
						resource.TestCheckResourceAttrPair(
							rName, "vpc.0.vpc_id",
							"huaweicloud_rds_instance.test", "vpc_id",
						),
						resource.TestCheckResourceAttrPair(
							rName, "vpc.0.subnet_id",
							"huaweicloud_rds_instance.test", "subnet_id",
						),
						resource.TestCheckResourceAttr(rName, "ssl.0.ssl_link", "false"),
						resource.TestCheckResourceAttr(rName, "config.0.driver_name", "mysql"),
						resource.TestCheckResourceAttrSet(rName, "enterprise_project_id"),
					),
				},
				{
					ResourceName:      rName,
					ImportState:       true,
					ImportStateVerify: true,
					ImportStateVerifyIgnore: []string{
						"endpoint.0.db_password",
					},
				},
			},
		},
	)
}

func testAccConnection_withRdsInstance(name, password string) string {
	return fmt.Sprintf(
		`
%s

%s

resource "huaweicloud_drs_connection" "test" {
  name        = "%s"
  db_type     = "mysql"
  description = "Test DRS connection with RDS instance"

  endpoint {
    endpoint_name = "cloud_mysql"
    instance_id   = huaweicloud_rds_instance.test.id
    db_port       = "3306"
    db_user       = "root"
    db_password   = "%s"
  }

  vpc {
    vpc_id    = huaweicloud_rds_instance.test.vpc_id
    subnet_id = huaweicloud_rds_instance.test.subnet_id
  }

  ssl {
    ssl_link = false
  }

  lifecycle {
    ignore_changes = [
      endpoint.0.db_password,
    ]
  }
}
`, common.TestBaseNetwork(name), testAccConnection_rdsInstance(name, password), name, password)
}

func testAccConnection_withRdsInstance_update(name, password string) string {
	return fmt.Sprintf(
		`
%s

%s

resource "huaweicloud_drs_connection" "test" {
  name        = "%s-update"
  db_type     = "mysql"
  description = "Updated DRS connection with RDS instance"

  endpoint {
    endpoint_name = "cloud_mysql"
    instance_id   = huaweicloud_rds_instance.test.id
    db_port       = "3307"
    db_user       = "admin"
    db_password   = "%s"
  }

  vpc {
    vpc_id    = huaweicloud_rds_instance.test.vpc_id
    subnet_id = huaweicloud_rds_instance.test.subnet_id
  }

  ssl {
    ssl_link = false
  }

  config {
    driver_name = "mysql"
  }

  lifecycle {
    ignore_changes = [
      endpoint.0.db_password,
    ]
  }
}
`, common.TestBaseNetwork(name), testAccConnection_rdsInstance(name, password), name, password)
}

func testAccConnection_rdsInstance(name, password string) string {
	return fmt.Sprintf(
		`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_rds_instance" "test" {
  depends_on = [
    huaweicloud_networking_secgroup_rule.ingress,
  ]

  name                = "%s"
  flavor              = "rds.mysql.x1.large.2"
  security_group_id   = huaweicloud_networking_secgroup.test.id
  subnet_id           = huaweicloud_vpc_subnet.test.id
  vpc_id              = huaweicloud_vpc.test.id
  fixed_ip            = "192.168.0.100"

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0],
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

resource "huaweicloud_networking_secgroup_rule" "ingress" {
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = "3306,9092"
  protocol          = "tcp"
  remote_ip_prefix  = "192.168.0.0/16"
  security_group_id = huaweicloud_networking_secgroup.test.id
}

resource "huaweicloud_networking_secgroup_rule" "egress" {
  direction         = "egress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  remote_ip_prefix  = "192.168.0.0/16"
  security_group_id = huaweicloud_networking_secgroup.test.id
}
`, name, password)
}

func TestAccConnection_mongoDB(t *testing.T) {
	var (
		connection interface{}
		rName      = "huaweicloud_drs_connection.test"
		name       = acceptance.RandomAccResourceName()
		password   = "Test@123456"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&connection,
		getConnectionResourceFunc,
	)

	resource.ParallelTest(
		t, resource.TestCase{
			PreCheck: func() {
				acceptance.TestAccPreCheck(t)
			},
			ProviderFactories: acceptance.TestAccProviderFactories,
			CheckDestroy:      rc.CheckResourceDestroy(),
			Steps: []resource.TestStep{
				{
					Config: testAccConnection_mongoDB(name, password),
					Check: resource.ComposeTestCheckFunc(
						rc.CheckResourceExists(),
						resource.TestCheckResourceAttr(rName, "name", name),
						resource.TestCheckResourceAttr(rName, "db_type", "mongodb"),
						resource.TestCheckResourceAttr(rName, "description", "Test DRS connection for MongoDb"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.endpoint_name", "mongodb"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.ip", "192.168.0.1:8080"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.db_user", "mog"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.db_name", "root"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.source_sharding.#", "2"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.source_sharding.0.ip", "192.168.0.1:8000"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.source_sharding.0.db_user", "mog"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.source_sharding.0.db_name", "root"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.source_sharding.1.ip", "192.168.0.2:8000"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.source_sharding.1.db_user", "mog"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.source_sharding.1.db_name", "root"),
						resource.TestCheckResourceAttr(rName, "ssl.0.ssl_link", "false"),
						resource.TestCheckResourceAttrSet(rName, "create_time"),
						resource.TestCheckResourceAttrSet(rName, "enterprise_project_id"),
					),
				},
				{
					Config: testAccConnection_mongoDB_update(name, password),
					Check: resource.ComposeTestCheckFunc(
						rc.CheckResourceExists(),
						resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s-update", name)),
						resource.TestCheckResourceAttr(rName, "db_type", "mongodb"),
						resource.TestCheckResourceAttr(rName, "description", "Updated DRS connection for MongoDB"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.endpoint_name", "mongodb"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.ip", "192.168.0.0.8080"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.db_user", "admin"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.db_name", "admin"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.source_sharding.#", "2"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.source_sharding.0.ip", "192.168.0.1.8010"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.source_sharding.0.db_user", "admin"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.source_sharding.0.db_name", "admin"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.source_sharding.1.ip", "192.168.0.2.8010"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.source_sharding.1.db_user", "admin"),
						resource.TestCheckResourceAttr(rName, "endpoint.0.source_sharding.1.db_name", "admin"),
						resource.TestCheckResourceAttr(rName, "ssl.0.ssl_link", "false"),
						resource.TestCheckResourceAttr(rName, "config.0.driver_name", "mongodb"),
						resource.TestCheckResourceAttrSet(rName, "enterprise_project_id"),
					),
				},
				{
					ResourceName:      rName,
					ImportState:       true,
					ImportStateVerify: true,
					ImportStateVerifyIgnore: []string{
						"endpoint.0.db_password",
						"endpoint.0.source_sharding.0.db_password",
						"endpoint.0.source_sharding.0.endpoint_name",
						"endpoint.0.source_sharding.1.db_password",
						"endpoint.0.source_sharding.1.endpoint_name",
					},
				},
			},
		},
	)
}

func testAccConnection_mongoDB(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_connection" "test" {
  name        = "%[1]s"
  db_type     = "mongodb"
  description = "Test DRS connection for MongoDb"

  endpoint {
    endpoint_name = "mongodb"
    ip            = "192.168.0.1:8080"
    db_user       = "mog"
    db_password   = "%[2]s"
    db_name       = "root"

    source_sharding {
      endpoint_name = "mongodb"
      ip            = "192.168.0.1:8000"
      db_user       = "mog"
      db_password   = "%[2]s"
      db_name       = "root"
    }

    source_sharding {
      endpoint_name = "mongodb"
      ip            = "192.168.0.2:8000"
      db_user       = "mog"
      db_password   = "%[2]s"
      db_name       = "root"
    }
  }

  ssl {
    ssl_link = false
  }

  lifecycle {
    ignore_changes = [
      endpoint.0.db_password,
      endpoint.0.source_sharding.0.db_password,
      endpoint.0.source_sharding.0.endpoint_name,
      endpoint.0.source_sharding.1.db_password,
      endpoint.0.source_sharding.1.endpoint_name,
   ]
  }
}
`, name, password)
}

func testAccConnection_mongoDB_update(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_connection" "test" {
  name        = "%[1]s-update"
  db_type     = "mongodb"
  description = "Updated DRS connection for MongoDB"

  endpoint {
    endpoint_name = "mongodb"
    ip            = "192.168.0.0.8080"
    db_user       = "admin"
    db_password   = "%[2]s"
    db_name       = "admin"

    source_sharding {
      endpoint_name = "mongodb"
      ip            = "192.168.0.1.8010"
      db_user       = "admin"
      db_password   = "%[2]s"
      db_name       = "admin"
    }

    source_sharding {
      endpoint_name = "mongodb"
      ip            = "192.168.0.2.8010"
      db_user       = "admin"
      db_password   = "%[2]s"
      db_name       = "admin"
    }
  }

  ssl {
    ssl_link = false
  }

  config {
    driver_name = "mongodb"
  }

  lifecycle {
    ignore_changes = [
      endpoint.0.db_password,
      endpoint.0.source_sharding.0.db_password,
      endpoint.0.source_sharding.0.endpoint_name,
      endpoint.0.source_sharding.1.db_password,
      endpoint.0.source_sharding.1.endpoint_name,
   ]
  }
}
`, name, password)
}

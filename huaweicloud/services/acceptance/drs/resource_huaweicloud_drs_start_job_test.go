package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccStartJob_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testStartJob_with_update_data_progress_rules_basic(),
			},
		},
	})
}

func testStartJob_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%[1]s"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1)
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = "%[1]s"
  delete_default_rules = true
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

resource "huaweicloud_rds_instance" "test1" {
  depends_on = [
    huaweicloud_networking_secgroup_rule.ingress,
    huaweicloud_networking_secgroup_rule.egress,
  ]

  name                = "%[1]s_1"
  flavor              = "rds.mysql.x1.large.2.ha"
  security_group_id   = huaweicloud_networking_secgroup.test.id
  subnet_id           = huaweicloud_vpc_subnet.test.id
  vpc_id              = huaweicloud_vpc.test.id
  fixed_ip            = "192.168.0.50"
  ha_replication_mode = "semisync"

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[3],
  ]

  db {
    password = "TestDrs@123"
    type     = "MySQL"
    version  = "5.7"
    port     = 3306
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}

resource "huaweicloud_rds_instance" "test2" {
  depends_on = [
    huaweicloud_networking_secgroup_rule.ingress,
    huaweicloud_networking_secgroup_rule.egress,
  ]

  name                = "%[1]s_2"
  flavor              = "rds.mysql.x1.large.2.ha"
  security_group_id   = huaweicloud_networking_secgroup.test.id
  subnet_id           = huaweicloud_vpc_subnet.test.id
  vpc_id              = huaweicloud_vpc.test.id
  fixed_ip            = "192.168.0.51"
  ha_replication_mode = "semisync"

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[3],
  ]

  db {
    password = "TestDrs@123"
    type     = "MySQL"
    version  = "5.7"
    port     = 3306
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}

resource "huaweicloud_rds_mysql_database" "test" {
  instance_id   = huaweicloud_rds_instance.test1.id
  name          = "%[1]s"
  character_set = "utf8"
}

data "huaweicloud_drs_node_types" "test" {
  engine_type = "mysql"
  type        = "sync"
  direction   = "up"
}

resource "huaweicloud_drs_job" "test" {
  name           = "%[1]s"
  type           = "sync"
  engine_type    = "mysql"
  direction      = "up"
  node_type      = data.huaweicloud_drs_node_types.test.node_types[0]
  net_type       = "vpc"
  migration_type = "FULL_INCR_TRANS"
  description    = "description test"
  force_destroy  = true

  source_db {
    engine_type = "mysql"
    ip          = huaweicloud_rds_instance.test1.fixed_ip
    port        = 3306
    user        = "root"
    password    = "TestDrs@123"
    vpc_id      = huaweicloud_rds_instance.test1.vpc_id
    subnet_id   = huaweicloud_rds_instance.test1.subnet_id
  }

  destination_db {
    region      = huaweicloud_rds_instance.test2.region
    ip          = huaweicloud_rds_instance.test2.fixed_ip
    port        = 3306
    engine_type = "mysql"
    user        = "root"
    password    = "TestDrs@123"
    instance_id = huaweicloud_rds_instance.test2.id
    subnet_id   = huaweicloud_rds_instance.test2.subnet_id
  }

  databases = [huaweicloud_rds_mysql_database.test.name]

  policy_config {
    filter_ddl_policy               = "drop_database"
    conflict_policy                 = "overwrite"
    index_trans                     = true
    transformation_name_case_policy = "lowercase"
  }

  is_pre_check  = false
  is_start_job  = false
  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"

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
`, name)
}

func testStartJob_with_update_data_progress_rules_basic() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_drs_update_data_progress_rules" "test" {
  job_id = huaweicloud_drs_job.test.id

  data_process_info {
    add_columns {
      column_type  = "ADDITIONALCOLUMN,create_time"
      column_name  = "name"
      column_value = "__create_timestamp"
      data_type    = "long"
    }

    db_object {
      object_scope = "database"

      object_info = jsonencode({
        (huaweicloud_rds_mysql_database.test.name) = {
          "all" = true
        }
      })
    }
  }
}

resource "huaweicloud_drs_start_job" "test" {
  depends_on = [huaweicloud_drs_update_data_progress_rules.test]

  job_id = huaweicloud_drs_job.test.id
}
`, testStartJob_base(name))
}

# Prerequisite resources: VPC, subnet and security group
resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr == "" ? cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0) : var.subnet_cidr
  gateway_ip = var.gateway_ip == "" ? cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1) : var.gateway_ip
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.security_group_name
  delete_default_rules = true
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  count = 2

  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  remote_ip_prefix  = "192.168.0.0/16"
  protocol          = "tcp"
  direction         = count.index == 0 ? "ingress" : "egress"
  ports             = count.index == 0 ? "3306" : null
}

data "huaweicloud_availability_zones" "test" {}

# Prerequisite resources: two RDS MySQL instances (0: first, 1: second)
resource "huaweicloud_rds_instance" "test" {
  count = 2

  depends_on = [
    huaweicloud_networking_secgroup_rule.test,
  ]

  name                = count.index == 0 ? var.rds1_name : var.rds2_name
  flavor              = var.rds_flavor
  security_group_id   = huaweicloud_networking_secgroup.test.id
  subnet_id           = huaweicloud_vpc_subnet.test.id
  vpc_id              = huaweicloud_vpc.test.id
  fixed_ip            = count.index == 0 ? var.rds1_fixed_ip : var.rds2_fixed_ip
  ha_replication_mode = "semisync"

  availability_zone = [
    try(data.huaweicloud_availability_zones.test.names[0], ""),
    try(data.huaweicloud_availability_zones.test.names[3], ""),
  ]

  db {
    password = var.db_password
    type     = "MySQL"
    version  = "5.7"
    port     = 3306
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}

# The database to be synchronized is created on both instances, so that each
# instance has data to be synchronized as the source database.
resource "huaweicloud_rds_mysql_database" "test" {
  count = 2

  instance_id   = huaweicloud_rds_instance.test[count.index].id
  name          = var.db_name
  character_set = "utf8"
}

# Node types of the DRS jobs, they are different for different data directions
# (0: up/forward, 1: down/reverse)
data "huaweicloud_drs_node_types" "test" {
  count = 2

  engine_type = "mysql"
  type        = "sync"
  direction   = count.index == 0 ? "up" : "down"
}

# Two DRS jobs (0: forward/up from the first RDS to the second, 1: reverse/down
# from the second RDS to the first). The jobs are created without starting so
# that the objects can be selected while they are in the CONFIGURATION status.
resource "huaweicloud_drs_job" "test" {
  count = 2

  name                    = count.index == 0 ? var.forward_job_name : var.reverse_job_name
  type                    = "sync"
  engine_type             = "mysql"
  direction               = count.index == 0 ? "up" : "down"
  node_type               = try(data.huaweicloud_drs_node_types.test[count.index].node_types[0], null)
  net_type                = "vpc"
  migration_type          = "FULL_INCR_TRANS"
  description             = count.index == 0 ? var.forward_job_description : var.reverse_job_description
  force_destroy           = true
  is_start_job            = false
  is_pre_check            = false
  destination_db_readnoly = false
  charging_mode           = "postPaid"

  source_db {
    engine_type = "mysql"
    ip          = huaweicloud_rds_instance.test[count.index].fixed_ip
    port        = 3306
    user        = var.db_user
    password    = var.db_password
    vpc_id      = count.index == 0 ? huaweicloud_rds_instance.test[count.index].vpc_id : null
    instance_id = count.index == 1 ? huaweicloud_rds_instance.test[count.index].id : null
    subnet_id   = huaweicloud_rds_instance.test[count.index].subnet_id
  }

  destination_db {
    region      = count.index == 0 ? huaweicloud_rds_instance.test[1 - count.index].region : null
    ip          = huaweicloud_rds_instance.test[1 - count.index].fixed_ip
    port        = 3306
    engine_type = "mysql"
    user        = var.db_user
    password    = var.db_password
    instance_id = count.index == 0 ? huaweicloud_rds_instance.test[1 - count.index].id : null
    vpc_id      = count.index == 1 ? huaweicloud_rds_instance.test[1 - count.index].vpc_id : null
    subnet_id   = huaweicloud_rds_instance.test[1 - count.index].subnet_id
  }

  policy_config {
    filter_ddl_policy = "drop_database"
    conflict_policy   = "overwrite"
    index_trans       = true
  }

  dynamic "limit_speed" {
    for_each = var.limit_speed == null ? [] : [var.limit_speed]

    content {
      speed      = limit_speed.value.speed
      start_time = limit_speed.value.start_time
      end_time   = limit_speed.value.end_time
    }
  }

  lifecycle {
    ignore_changes = [
      source_db.0.password, destination_db.0.password, force_destroy,
    ]
  }
}

# Select the objects to be synchronized for each job.
resource "huaweicloud_drs_batch_select_objects" "test" {
  count = 2

  jobs {
    job_id        = huaweicloud_drs_job.test[count.index].id
    selected      = true
    sync_database = true

    job {
      id          = huaweicloud_rds_mysql_database.test[count.index].name
      object_type = "database"
      object_name = huaweicloud_rds_mysql_database.test[count.index].name
      select      = "true"
    }
  }
}

# Start each job after the objects are selected.
resource "huaweicloud_drs_start_job" "test" {
  count = 2

  depends_on = [huaweicloud_drs_batch_select_objects.test]

  job_id = huaweicloud_drs_job.test[count.index].id
}

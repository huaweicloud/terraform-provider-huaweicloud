# VPC
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

data "huaweicloud_rds_flavors" "test" {
  db_type       = var.rds_db_type
  db_version    = var.rds_db_version
  instance_mode = var.rds_instance_mode
}

# RDS MySQL instances (0: source, 1: destination)
resource "huaweicloud_rds_instance" "test" {
  count = 2

  name                = count.index == 0 ? var.source_rds_name : var.dest_rds_name
  flavor              = var.rds_flavor != "" ? var.rds_flavor : try(data.huaweicloud_rds_flavors.test.flavors[0].name, null)
  security_group_id   = huaweicloud_networking_secgroup.test.id
  subnet_id           = huaweicloud_vpc_subnet.test.id
  vpc_id              = huaweicloud_vpc.test.id
  fixed_ip            = count.index == 0 ? var.source_rds_fixed_ip : var.dest_rds_fixed_ip
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

# DRS migration job
resource "huaweicloud_drs_job" "test" {
  name           = var.job_name
  type           = "migration"
  engine_type    = "mysql"
  direction      = "up"
  net_type       = "eip"
  migration_type = "FULL_INCR_TRANS"
  description    = var.description
  force_destroy  = true

  source_db {
    engine_type = "mysql"
    ip          = huaweicloud_rds_instance.test[0].fixed_ip
    port        = 3306
    user        = "root"
    password    = var.db_password
    ssl_enabled = false
  }

  destination_db {
    region      = huaweicloud_rds_instance.test[1].region
    ip          = huaweicloud_rds_instance.test[1].fixed_ip
    port        = 3306
    engine_type = "mysql"
    user        = "root"
    password    = var.db_password
    instance_id = huaweicloud_rds_instance.test[1].id
    subnet_id   = huaweicloud_rds_instance.test[1].subnet_id
  }

  lifecycle {
    ignore_changes = [
      source_db.0.password, destination_db.0.password, force_destroy, action,
    ]
  }
}

# LTS log group
resource "huaweicloud_lts_group" "test" {
  group_name  = var.lts_group_name
  ttl_in_days = var.lts_ttl_in_days
}

# LTS log stream
resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = var.lts_stream_name
  is_favorite = true
}

# DRS LTS config
resource "huaweicloud_drs_lts_config" "test" {
  job_id        = huaweicloud_drs_job.test.id
  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id
}

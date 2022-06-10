variable "param_group_id" {
  type    = string
  default = "qwert"
}

variable "password" {
  type    = string
  default = "12345"
}

data "huaweicloud_vpc" "vpc" {
  id = "VPCID"
}

data "huaweicloud_networking_secgroup" "sg" {
  secgroup_id = "SGID"
}

data "huaweicloud_vpc_subnet" "subnet" {
  id = "SubnetID"
}

locals {
  rg = csvdecode(file("mysql.csv"))
}

resource "huaweicloud_rds_instance" "mysql" {
  for_each            = { for rg in local.rg : rg.name => rg }
  name                = each.value.name
  flavor              = each.value.flavor
  ha_replication_mode = each.value.ha
  vpc_id              = data.huaweicloud_vpc.vpc.id
  subnet_id           = data.huaweicloud_vpc_subnet.subnet.id
  security_group_id   = data.huaweicloud_networking_secgroup.sg.secgroup_id
  param_group_id      = var.param_group_id
  availability_zone   = [each.value.az]

  db {
    type     = "MySQL"
    version  = "5.7"
    password = var.password
  }

  volume {
    type = "CLOUDSSD"
    size = each.value.size
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 7
  }
}

resource "huaweicloud_drs_job" "test" {
  for_each       = { for rg in local.rg : rg.name => rg }
  name           = each.value.name1
  type           = "migration"
  engine_type    = "mysql"
  direction      = "down"
  net_type       = "vpc"
  migration_type = "FULL_TRANS"
  description    = "TEST"
  force_destroy  = "true"

  source_db {
    engine_type = "mysql"
    ip          = each.value.sip
    port        = "3306"
    user        = each.value.suser
    password    = each.value.spass
  }

  destination_db {
    region      = "cn-east-3"
    ip          = huaweicloud_rds_instance.mysql[each.value.name].fixed_ip
    port        = 3306
    engine_type = "mysql"
    user        = "root"
    password    = var.password
    instance_id = huaweicloud_rds_instance.mysql[each.value.name].id
    subnet_id   = data.huaweicloud_vpc_subnet.subnet.id
  }
}

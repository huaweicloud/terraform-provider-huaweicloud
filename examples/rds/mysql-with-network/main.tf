data "huaweicloud_availability_zones" "myaz" {}

data "huaweicloud_vpc" "myvpc" {
  name = var.vpc_name
}

data "huaweicloud_vpc_subnet" "mysubnet" {
  vpc_id = data.huaweicloud_vpc.myvpc.id
  name   = var.subnet_name
}

data "huaweicloud_networking_secgroup" "mysecgroup" {
  name = var.secgroup_name
}

resource "huaweicloud_rds_instance" "myinstance" {
  name                = "mysql_instance"
  flavor              = "rds.mysql.c2.large.ha"
  ha_replication_mode = "async"
  vpc_id              = data.huaweicloud_vpc.myvpc.id
  subnet_id           = data.huaweicloud_vpc_subnet.mysubnet.id
  security_group_id   = data.huaweicloud_networking_secgroup.mysecgroup.id
  availability_zone = [
    data.huaweicloud_availability_zones.myaz.names[0],
    data.huaweicloud_availability_zones.myaz.names[1]
  ]

  db {
    type     = "MySQL"
    version  = "8.0"
    password = var.rds_password
  }
  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}

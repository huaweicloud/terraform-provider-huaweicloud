data "huaweicloud_availability_zones" "myaz" {}

resource "huaweicloud_vpc" "myvpc" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "mysubnet" {
  vpc_id      = huaweicloud_vpc.myvpc.id
  name        = var.subnet_name
  cidr        = var.subnet_cidr
  gateway_ip  = var.subnet_gateway
  primary_dns = var.primary_dns
}

resource "huaweicloud_networking_secgroup" "mysecgroup" {
  name        = "mysecgroup"
  description = "a basic security group"
}

resource "huaweicloud_networking_secgroup_rule" "allow_rds" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 3306
  port_range_max    = 3306
  remote_ip_prefix  = var.allow_cidr
  security_group_id = huaweicloud_networking_secgroup.mysecgroup.id
}

resource "random_password" "mypassword" {
  length           = 12
  special          = true
  override_special = "!@#%^*-_=+"
}

resource "huaweicloud_rds_instance" "myinstance" {
  name                = "mysql_instance"
  flavor              = "rds.mysql.c2.large.ha"
  ha_replication_mode = "async"
  vpc_id              = huaweicloud_vpc.myvpc.id
  subnet_id           = huaweicloud_vpc_subnet.mysubnet.id
  security_group_id   = huaweicloud_networking_secgroup.mysecgroup.id
  availability_zone = [
    data.huaweicloud_availability_zones.myaz.names[0],
    data.huaweicloud_availability_zones.myaz.names[1]
  ]

  db {
    type     = "MySQL"
    version  = "8.0"
    password = random_password.mypassword.result
  }
  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}

resource "huaweicloud_vpc_eip" "myeip" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

# get the port of rds instance by private_ip
data "huaweicloud_networking_port" "rds_port" {
  network_id = huaweicloud_vpc_subnet.mysubnet.id
  fixed_ip   = huaweicloud_rds_instance.myinstance.private_ips[0]
}

resource "huaweicloud_vpc_eip_associate" "associated" {
  public_ip = huaweicloud_vpc_eip.myeip.address
  port_id   = data.huaweicloud_networking_port.rds_port.id
}

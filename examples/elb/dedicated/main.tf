data "huaweicloud_availability_zones" "default" {}

data "huaweicloud_compute_flavors" "default" {
  availability_zone = data.huaweicloud_availability_zones.default.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_images_image" "default" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

resource "huaweicloud_vpc" "default" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "default" {
  vpc_id      = huaweicloud_vpc.default.id
  name        = var.subnet_name
  cidr        = var.subnet_cidr
  gateway_ip  = var.subnet_gateway
  ipv6_enable = true
}

resource "huaweicloud_network_acl" "default" {
  name = var.network_acl_name
  subnets = [
    huaweicloud_vpc_subnet.default.id
  ]

  inbound_rules = [
    huaweicloud_network_acl_rule.default.id
  ]
}

resource "huaweicloud_network_acl_rule" "default" {
  name                   = var.network_acl_rule_name
  protocol               = "tcp"
  action                 = "allow"
  source_ip_address      = huaweicloud_vpc.default.cidr
  source_port            = "8080"
  destination_ip_address = "0.0.0.0/0"
  destination_port       = "8081"
}

resource "huaweicloud_networking_secgroup" "default" {
  name                 = var.security_group_name
  delete_default_rules = true
}

resource "huaweicloud_networking_secgroup_rule" "in_v4_tcp_3389" {
  depends_on = [
    huaweicloud_compute_eip_associate.default
  ]

  security_group_id = huaweicloud_networking_secgroup.default.id
  ethertype         = "IPv4"
  direction         = "ingress"
  protocol          = "tcp"
  ports             = "3389"
  remote_ip_prefix  = format("%s/32", huaweicloud_compute_instance.default.access_ip_v4)
}

resource "huaweicloud_networking_secgroup_rule" "in_v4_icmp_all" {
  security_group_id = huaweicloud_networking_secgroup.default.id
  ethertype         = "IPv4"
  direction         = "ingress"
  protocol          = "icmp"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_networking_secgroup_rule" "in_v4_elb_member" {
  security_group_id = huaweicloud_networking_secgroup.default.id
  ethertype         = "IPv4"
  direction         = "ingress"
  protocol          = "tcp"
  ports             = "80,8081"
  remote_ip_prefix  = huaweicloud_vpc.default.cidr
}

resource "huaweicloud_networking_secgroup_rule" "in_v4_all_group" {
  security_group_id = huaweicloud_networking_secgroup.default.id
  ethertype         = "IPv4"
  direction         = "ingress"
  remote_group_id   = huaweicloud_networking_secgroup.default.id
}

resource "huaweicloud_networking_secgroup_rule" "in_v6_all_group" {
  security_group_id = huaweicloud_networking_secgroup.default.id
  ethertype         = "IPv6"
  direction         = "ingress"
  remote_group_id   = huaweicloud_networking_secgroup.default.id
}

resource "huaweicloud_networking_secgroup_rule" "out_v4_all" {
  security_group_id = huaweicloud_networking_secgroup.default.id
  ethertype         = "IPv4"
  direction         = "egress"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_networking_secgroup_rule" "out_v6_all" {
  security_group_id = huaweicloud_networking_secgroup.default.id
  ethertype         = "IPv6"
  direction         = "egress"
  remote_ip_prefix  = "::/0"
}

resource "huaweicloud_compute_instance" "default" {
  name              = var.ecs_instance_name
  image_id          = data.huaweicloud_images_image.default.id
  flavor_id         = data.huaweicloud_compute_flavors.default.ids[0]
  availability_zone = data.huaweicloud_availability_zones.default.names[0]
  security_groups   = [var.security_group_name]

  network {
    uuid = huaweicloud_vpc_subnet.default.id
  }
}

resource "huaweicloud_vpc_eip" "default" {
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

resource "huaweicloud_compute_eip_associate" "default" {
  public_ip   = huaweicloud_vpc_eip.default.address
  instance_id = huaweicloud_compute_instance.default.id
}

resource "huaweicloud_elb_loadbalancer" "default" {
  name            = var.elb_loadbalancer_name
  description     = "Created by terraform"
  vpc_id          = huaweicloud_vpc.default.id
  ipv4_subnet_id  = huaweicloud_vpc_subnet.default.ipv4_subnet_id
  ipv6_network_id = huaweicloud_vpc_subnet.default.id

  availability_zone = [
    data.huaweicloud_availability_zones.default.names[0]
  ]

  tags = {
    owner = "terraform"
  }
}

resource "huaweicloud_elb_listener" "default" {
  name            = var.elb_listener_name
  description     = "Created by terraform"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = huaweicloud_elb_loadbalancer.default.id

  idle_timeout     = 60
  request_timeout  = 60
  response_timeout = 60

  tags = {
    owner = "terraform"
  }
}

resource "huaweicloud_elb_pool" "default" {
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = huaweicloud_elb_listener.default.id

  persistence {
    type = "HTTP_COOKIE"
  }
}

resource "huaweicloud_elb_monitor" "default" {
  protocol    = "HTTP"
  interval    = 20
  timeout     = 15
  max_retries = 10
  url_path    = "/"
  port        = 8080
  pool_id     = huaweicloud_elb_pool.default.id
}

resource "huaweicloud_elb_member" "default" {
  address       = huaweicloud_compute_instance.default.access_ip_v4
  protocol_port = 8080
  pool_id       = huaweicloud_elb_pool.default.id
  subnet_id     = huaweicloud_vpc_subnet.default.ipv4_subnet_id
}

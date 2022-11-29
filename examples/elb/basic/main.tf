data "huaweicloud_availability_zones" "myaz" {}

data "huaweicloud_compute_flavors" "myflavor" {
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_images_image" "myimage" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

resource "huaweicloud_vpc" "vpc_1" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "subnet_1" {
  vpc_id      = huaweicloud_vpc.vpc_1.id
  name        = var.subnet_name
  cidr        = var.subnet_cidr
  gateway_ip  = var.subnet_gateway
  primary_dns = var.primary_dns
}

resource "huaweicloud_networking_secgroup" "secgroup_1" {
  name        = var.secgroup_name
  description = "basic security group"
}

# allow http
resource "huaweicloud_networking_secgroup_rule" "allow_http" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 80
  port_range_max    = 80
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = huaweicloud_networking_secgroup.secgroup_1.id
}

resource "huaweicloud_compute_instance" "instance" {
  name              = "instance_${count.index}"
  image_id          = data.huaweicloud_images_image.myimage.id
  flavor_id         = data.huaweicloud_compute_flavors.myflavor.ids[0]
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
  security_groups   = [var.secgroup_name]

  network {
    uuid = huaweicloud_vpc_subnet.subnet_1.id
  }
  count = 2
}

resource "huaweicloud_vpc_eip" "eip_1" {
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

resource "huaweicloud_lb_loadbalancer" "elb_1" {
  name          = var.loadbalancer_name
  vip_subnet_id = huaweicloud_vpc_subnet.subnet_1.ipv4_subnet_id
}

# associate eip with loadbalancer
resource "huaweicloud_vpc_eip_associate" "associate_1" {
  public_ip = huaweicloud_vpc_eip.eip_1.address
  port_id   = huaweicloud_lb_loadbalancer.elb_1.vip_port_id
}

resource "huaweicloud_lb_listener" "listener_1" {
  name            = "listener_http"
  protocol        = "HTTP"
  protocol_port   = 80
  loadbalancer_id = huaweicloud_lb_loadbalancer.elb_1.id
}

resource "huaweicloud_lb_pool" "group_1" {
  name        = "group_1"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = huaweicloud_lb_listener.listener_1.id
}

resource "huaweicloud_lb_monitor" "health_check" {
  name           = "health_check"
  type           = "HTTP"
  url_path       = "/"
  expected_codes = "200-202"
  delay          = 10
  timeout        = 5
  max_retries    = 3
  pool_id        = huaweicloud_lb_pool.group_1.id
}

resource "huaweicloud_lb_member" "member_1" {
  address       = huaweicloud_compute_instance.instance[0].access_ip_v4
  protocol_port = 80
  weight        = 1
  pool_id       = huaweicloud_lb_pool.group_1.id
  subnet_id     = huaweicloud_vpc_subnet.subnet_1.ipv4_subnet_id
}

resource "huaweicloud_lb_member" "member_2" {
  address       = huaweicloud_compute_instance.instance[1].access_ip_v4
  protocol_port = 80
  weight        = 1
  pool_id       = huaweicloud_lb_pool.group_1.id
  subnet_id     = huaweicloud_vpc_subnet.subnet_1.ipv4_subnet_id
}

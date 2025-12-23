# Create VPC
resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr
  gateway_ip = var.subnet_gateway_ip == "" ? cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1) : var.subnet_gateway_ip
}

# Create NAT Gateway And SNAT Rule
resource "huaweicloud_nat_gateway" "test" {
  name        = var.nat_gateway_name
  description = var.gateway_description
  spec        = var.gateway_spec
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id
}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = var.eip_bandwidth_name
    size        = var.eip_bandwidth_size
    share_type  = var.eip_bandwidth_share_type
    charge_mode = var.eip_bandwidth_charge_mode
  }
}

resource "huaweicloud_nat_snat_rule" "test" {
  nat_gateway_id = huaweicloud_nat_gateway.test.id
  floating_ip_id = huaweicloud_vpc_eip.test.id
  source_type    = var.snat_source_type
  subnet_id      = var.snat_source_type == 0 ? huaweicloud_vpc_subnet.test.id : null
  cidr           = var.snat_source_type == 1 ? var.snat_cidr : null
  description    = var.snat_description
}

# Create ECS instance
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], null)
  performance_type  = var.ecs_flavor_performance_type
  cpu_core_count    = var.ecs_flavor_cpu_core_count
  memory_size       = var.ecs_flavor_memory_size
}

data "huaweicloud_images_images" "test" {
  flavor_id  = var.ecs_flavor_id != "" ? var.ecs_flavor_id : try(data.huaweicloud_compute_flavors.test.flavors[0].id, null)
  visibility = var.ecs_image_visibility
  os         = var.ecs_image_os
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.ecs_security_group_name
  delete_default_rules = true
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = "egress"
  ethertype         = "IPv4"
  remote_ip_prefix  = "0.0.0.0/0"
  priority          = 1
}

resource "huaweicloud_compute_instance" "test" {
  name               = var.ecs_instance_name
  availability_zone  = try(data.huaweicloud_availability_zones.test.names[0], null)
  flavor_id          = var.ecs_flavor_id != "" ? var.ecs_flavor_id : try(data.huaweicloud_compute_flavors.test.flavors[0].id, null)
  image_id           = var.ecs_image_id != "" ? var.ecs_image_id : try(data.huaweicloud_images_images.test.images[0].id, null)
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  system_disk_type   = var.ecs_system_disk_type
  system_disk_size   = var.ecs_system_disk_size

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }

  admin_pass = var.ecs_admin_password
  tags       = var.ecs_instance_tags

  depends_on = [
    huaweicloud_nat_gateway.test,
    huaweicloud_vpc_eip.test
  ]
}

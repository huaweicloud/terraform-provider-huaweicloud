# Create a VPC endpoint service
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], null)
  performance_type  = var.instance_flavor_performance_type
  cpu_core_count    = var.instance_flavor_cpu_core_count
  memory_size       = var.instance_flavor_memory_size
}

data "huaweicloud_images_image" "test" {
  name        = var.instance_image_name
  most_recent = var.instance_image_most_recent
}

resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr == "" ? cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0) : var.subnet_cidr
  gateway_ip = var.subnet_gateway_ip == "" ? cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1) : var.subnet_gateway_ip
}

resource "huaweicloud_networking_secgroup" "test" {
  name = var.security_group_name
}

resource "huaweicloud_compute_instance" "test" {
  name               = var.instance_name
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = var.instance_flavor_id == "" ? try(data.huaweicloud_compute_flavors.test.ids[0], "") : var.instance_flavor_id
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = try(data.huaweicloud_availability_zones.test.names[0], null)

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_vpcep_service" "test" {
  name        = var.endpoint_service_name
  server_type = var.endpoint_service_type
  vpc_id      = huaweicloud_vpc.test.id
  port_id     = huaweicloud_compute_instance.test.network[0].port

  dynamic "port_mapping" {
    for_each = var.endpoint_service_port_mapping

    content {
      service_port  = port_mapping.value.service_port
      terminal_port = port_mapping.value.terminal_port
    }
  }
}

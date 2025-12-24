data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_images_image" "test" {
  name        = var.image_name
  visibility  = var.image_visibility
  most_recent = var.image_most_recent
}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], null)
  performance_type  = var.flavor_performance_type
  cpu_core_count    = var.flavor_cpu_core_count
  memory_size       = var.flavor_memory_size
}

resource "huaweicloud_networking_secgroup" "test" {
  name        = var.secgroup_name
  description = "Security group for ECS instances"
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  count = length(var.security_group_rules)

  direction         = var.security_group_rules[count.index].direction
  description       = var.security_group_rules[count.index].description
  ethertype         = var.security_group_rules[count.index].ethertype
  protocol          = var.security_group_rules[count.index].protocol
  ports             = var.security_group_rules[count.index].ports
  remote_ip_prefix  = var.security_group_rules[count.index].remote_ip_prefix
  security_group_id = huaweicloud_networking_secgroup.test.id
}

# Create VPCs
resource "huaweicloud_vpc" "test" {
  count = length(var.vpcs)

  name = var.vpcs[count.index].vpc_name
  cidr = var.vpcs[count.index].vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  count = length(var.vpcs)

  vpc_id     = huaweicloud_vpc.test[count.index].id
  name       = var.vpcs[count.index].subnet_name
  cidr       = var.vpcs[count.index].subnet_cidr
  gateway_ip = var.vpcs[count.index].subnet_gateway_ip != null && var.vpcs[count.index].subnet_gateway_ip != "" ? var.vpcs[count.index].subnet_gateway_ip : cidrhost(cidrsubnet(huaweicloud_vpc.test[count.index].cidr, 8, 0), 1)
}

# Create ECS instances
resource "huaweicloud_compute_instance" "test" {
  count = length(var.vpcs)

  name               = var.vpcs[count.index].instance_name
  availability_zone  = try(data.huaweicloud_availability_zones.test.names[0], null)
  flavor_id          = var.vpcs[count.index].ecs_flavor_id != null && var.vpcs[count.index].ecs_flavor_id != "" ? var.vpcs[count.index].ecs_flavor_id : try(data.huaweicloud_compute_flavors.test.ids[0], null)
  image_id           = var.vpcs[count.index].ecs_image_id != null && var.vpcs[count.index].ecs_image_id != "" ? var.vpcs[count.index].ecs_image_id : try(data.huaweicloud_images_image.test.id, null)
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  system_disk_type   = var.vpcs[count.index].ecs_system_disk_type != null ? var.vpcs[count.index].ecs_system_disk_type : "SSD"
  system_disk_size   = var.vpcs[count.index].ecs_system_disk_size != null ? var.vpcs[count.index].ecs_system_disk_size : 40

  network {
    uuid = huaweicloud_vpc_subnet.test[count.index].id
  }

  admin_pass = var.vpcs[count.index].ecs_admin_password != null && var.vpcs[count.index].ecs_admin_password != "" ? var.vpcs[count.index].ecs_admin_password : var.ecs_admin_password
  tags       = var.vpcs[count.index].ecs_instance_tags != null ? var.vpcs[count.index].ecs_instance_tags : {}

  depends_on = [
    huaweicloud_nat_gateway.test,
    huaweicloud_vpc_eip.test
  ]
}

# Cross-VPC Peering Configuration
# ST.001 Disable
resource "huaweicloud_vpc_peering_connection" "cross_vpc_peering" {
  name        = var.peering_connection_name
  vpc_id      = huaweicloud_vpc.test[0].id
  peer_vpc_id = huaweicloud_vpc.test[1].id
}

resource "huaweicloud_vpc_route" "route_to_peer_vpc" {
  vpc_id      = huaweicloud_vpc.test[0].id
  destination = try(var.vpcs[1].vpc_cidr, null)
  type        = var.route_type
  nexthop     = huaweicloud_vpc_peering_connection.cross_vpc_peering.id
}

resource "huaweicloud_vpc_route" "route_other_to_internet_via_peering" {
  vpc_id      = huaweicloud_vpc.test[1].id
  destination = "0.0.0.0/0"
  type        = var.route_type
  nexthop     = huaweicloud_vpc_peering_connection.cross_vpc_peering.id
}
# ST.001 Enable

# Create EIP for NAT Gateway
resource "huaweicloud_vpc_eip" "test" {
  count = length(var.eips)

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    name        = "${var.nat_gateway_name}-${var.eips[count.index].name_suffix}-shared-bandwidth"
    size        = var.bandwidth_size
    charge_mode = var.bandwidth_charge_mode
  }
}

# Create NAT Gateway
resource "huaweicloud_nat_gateway" "test" {
  name        = var.nat_gateway_name
  description = "NAT gateway for VPC One to provide internet access"
  spec        = var.nat_gateway_spec
  vpc_id      = huaweicloud_vpc.test[0].id
  subnet_id   = huaweicloud_vpc_subnet.test[0].id
}

resource "huaweicloud_nat_snat_rule" "test" {
  count = length(var.snat_rules)

  nat_gateway_id = huaweicloud_nat_gateway.test.id
  floating_ip_id = huaweicloud_vpc_eip.test[var.snat_rules[count.index].eip_index].id
  subnet_id      = var.snat_rules[count.index].subnet_index != null ? huaweicloud_vpc_subnet.test[var.snat_rules[count.index].subnet_index].id : null
  source_type    = var.snat_rules[count.index].source_type
  cidr           = var.snat_rules[count.index].vpc_index != null ? var.vpcs[var.snat_rules[count.index].vpc_index].subnet_cidr : null
}

resource "huaweicloud_nat_dnat_rule" "test" {
  count = length(var.dnat_rules)

  nat_gateway_id        = huaweicloud_nat_gateway.test.id
  floating_ip_id        = huaweicloud_vpc_eip.test[var.dnat_rules[count.index].eip_index].id
  private_ip            = try(huaweicloud_compute_instance.test[var.dnat_rules[count.index].instance_index].network[0].fixed_ip_v4, null)
  protocol              = var.dnat_protocol
  internal_service_port = var.dnat_internal_port
  external_service_port = var.dnat_rules[count.index].external_service_port
}

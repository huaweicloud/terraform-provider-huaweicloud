resource "huaweicloud_compute_instance" "newCompute_Example" {
  name              = var.ecs_name
  image_id          = data.huaweicloud_images_image.newIMS_Example.id
  flavor_id         = "s6.small.1"
  security_groups   = [huaweicloud_networking_secgroup.newSecgroup_Example.name]
  admin_pass        = var.user_password
  availability_zone = data.huaweicloud_availability_zones.newAZ_Example.names[0]

  system_disk_type  = "SSD"
  system_disk_size  = 40

  data_disks {
    type = "SSD"
    size = "10"
  }
  # multi_data_disks

  network {
    fixed_ip_v4 = var.ecs_ipaddress
    uuid        = huaweicloud_vpc_subnet.newSubnet_Example.id
  }
}

data "huaweicloud_availability_zones" "newAZ_Example" {}

data "huaweicloud_images_image" "newIMS_Example" {
  name        = var.ims_name
  visibility  = "public" # "private"、"share"、"community"
  most_recent = true
}

resource "huaweicloud_vpc_eip" "newEIP_Example" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = var.bandwidth_name
    size        = 5
    share_type  = "PER" # "WHOLE"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_nat_gateway" "newNet_gateway_Example" {
  name                = var.net_gateway_name
  description         = "example for net test"
  spec                = "1"
  router_id           = huaweicloud_vpc.newVPC_Example.id
  internal_network_id = huaweicloud_vpc_subnet.newSubnet_Example.id
}

resource "huaweicloud_nat_dnat_rule" "newDNATRule_Example" {
  count = length(var.example_dnat_rule)

  floating_ip_id        = huaweicloud_vpc_eip.newEIP_Example.id
  nat_gateway_id        = huaweicloud_nat_gateway.newNet_gateway_Example.id
  port_id               = huaweicloud_compute_instance.newCompute_Example.network[0].port     

  internal_service_port = lookup(var.example_dnat_rule[count.index], "internal_service_port", null)
  protocol              = lookup(var.example_dnat_rule[count.index], "protocol", null)
  external_service_port = lookup(var.example_dnat_rule[count.index], "external_service_port", null)
}

resource "huaweicloud_vpc" "newVPC_Example" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "newSubnet_Example" {
  name          = var.subnet_name
  cidr          = var.subnet_cidr
  gateway_ip    = var.subnet_gateway_ip
  vpc_id        = huaweicloud_vpc.newVPC_Example.id
  primary_dns   = "100.125.129.250"
  secondary_dns = "100.125.1.250"
}

resource "huaweicloud_networking_secgroup" "newSecgroup_Example" {
  name        = var.secgroup_name
  description = "This is a security group"
}

resource "huaweicloud_networking_secgroup_rule" "newSecgroup_GressRule_Example" {
  count             = length(var.example_security_group)

  direction         = lookup(var.example_security_group[count.index], "direction", null)
  ethertype         = lookup(var.example_security_group[count.index], "ethertype", null)
  protocol          = lookup(var.example_security_group[count.index], "protocol", null)
  port_range_min    = lookup(var.example_security_group[count.index], "port_range_min", null)
  port_range_max    = lookup(var.example_security_group[count.index], "port_range_max", null)
  remote_ip_prefix  = lookup(var.example_security_group[count.index], "remote_ip_prefix", null)
  security_group_id = huaweicloud_networking_secgroup.newSecgroup_Example.id
}

resource "huaweicloud_evs_volume" "newVolume_Example" {
  name              = var.evs_name
  availability_zone = data.huaweicloud_availability_zones.newAZ_Example.names[0]
  volume_type       = "SSD"
  size              = 20
}

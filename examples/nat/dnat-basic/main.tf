data "huaweicloud_availability_zones" "newAZ_Example" {}

data "huaweicloud_images_image" "newIMS_Example" {
  name        = var.ims_name
  visibility  = "public"
  most_recent = true
}

resource "random_password" "password" {
  length           = 16
  special          = true
  override_special = "!@$%^-_=+[{}]:,./?~#*"
  min_numeric      = 1
  min_lower        = 1
  min_special      = 1
}

resource "huaweicloud_compute_instance" "newCompute_Example" {
  name              = var.ecs_name
  image_id          = data.huaweicloud_images_image.newIMS_Example.id
  flavor_id         = "s6.small.1"
  security_groups   = [huaweicloud_networking_secgroup.newSecgroup_Example.name]
  admin_pass        = random_password.password.result
  availability_zone = data.huaweicloud_availability_zones.newAZ_Example.names[0]

  system_disk_type = "SSD"
  system_disk_size = 40

  network {
    uuid = huaweicloud_vpc_subnet.newSubnet_Example.id
  }
}

resource "huaweicloud_vpc_eip" "newEIP_Example" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = var.bandwidth_name
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
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

resource "huaweicloud_networking_secgroup_rule" "newSecgroupRule_Example" {
  count = length(var.security_group_rule)

  direction         = lookup(var.security_group_rule[count.index], "direction", null)
  ethertype         = lookup(var.security_group_rule[count.index], "ethertype", null)
  protocol          = lookup(var.security_group_rule[count.index], "protocol", null)
  port_range_min    = lookup(var.security_group_rule[count.index], "port_range_min", null)
  port_range_max    = lookup(var.security_group_rule[count.index], "port_range_max", null)
  remote_ip_prefix  = lookup(var.security_group_rule[count.index], "remote_ip_prefix", null)
  security_group_id = huaweicloud_networking_secgroup.newSecgroup_Example.id
}

resource "huaweicloud_nat_gateway" "newNet_gateway_Example" {
  name                = var.net_gateway_name
  description         = "example for net test"
  spec                = "1"
  router_id           = huaweicloud_vpc.newVPC_Example.id
  internal_network_id = huaweicloud_vpc_subnet.newSubnet_Example.id
}

resource "huaweicloud_nat_snat_rule" "newSNATRule_Example" {
  nat_gateway_id = huaweicloud_nat_gateway.newNet_gateway_Example.id
  network_id     = huaweicloud_vpc_subnet.newSubnet_Example.id
  floating_ip_id = huaweicloud_vpc_eip.newEIP_Example.id
}

resource "huaweicloud_nat_dnat_rule" "newDNATRule_Example" {
  count = length(var.example_dnat_rule)

  floating_ip_id = huaweicloud_vpc_eip.newEIP_Example.id
  nat_gateway_id = huaweicloud_nat_gateway.newNet_gateway_Example.id
  port_id        = huaweicloud_compute_instance.newCompute_Example.network[0].port

  internal_service_port = lookup(var.example_dnat_rule[count.index], "internal_service_port", null)
  protocol              = lookup(var.example_dnat_rule[count.index], "protocol", null)
  external_service_port = lookup(var.example_dnat_rule[count.index], "external_service_port", null)
}

resource "null_resource" "provision" {
  depends_on = [huaweicloud_nat_snat_rule.newSNATRule_Example, huaweicloud_nat_dnat_rule.newDNATRule_Example]

  provisioner "remote-exec" {
    connection {
      user     = "root"
      password = random_password.password.result
      host     = huaweicloud_vpc_eip.newEIP_Example.address
      port     = var.ecs_ssh_port
    }
    inline = [
      "yum -y install nginx",
      "systemctl enable nginx",
      "systemctl start nginx",
      "systemctl status nginx",
    ]
  }
}

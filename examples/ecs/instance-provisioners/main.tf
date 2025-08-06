data "huaweicloud_availability_zones" "test" {
  count = var.availability_zone == "" ? 1 : 0
}

data "huaweicloud_compute_flavors" "test" {
  count = var.instance_flavor_id == "" ? 1 : 0

  availability_zone = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test[0].names[0], null) : var.availability_zone
  performance_type  = var.instance_performance_type
  cpu_core_count    = var.instance_cpu_core_count
  memory_size       = var.instance_memory_size
}

data "huaweicloud_images_images" "test" {
  count = var.instance_image_id == "" ? 1 : 0

  flavor_id  = var.instance_flavor_id == "" ? try(data.huaweicloud_compute_flavors.test[0].flavors[0].id, null) : var.instance_flavor_id
  visibility = var.instance_image_visibility
  os         = var.instance_image_os
}

resource "huaweicloud_kps_keypair" "test" {
  name     = var.keypair_name
  key_file = var.private_key_path
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
  count = var.security_group_name != "" ? 1 : 0

  name = var.security_group_name
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  count = var.security_group_name != "" ? 1 : 0

  security_group_id = huaweicloud_networking_secgroup.test[0].id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = "22"
  port_range_max    = "22"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_vpc_eip" "test" {
  count = var.associate_eip_address == "" ? 1 : 0

  publicip {
    type = var.eip_type
  }

  bandwidth {
    name        = var.bandwidth_name
    size        = var.bandwidth_size
    share_type  = var.bandwidth_share_type
    charge_mode = var.bandwidth_charge_mode
  }

  lifecycle {
    precondition {
      condition     = var.associate_eip_address != "" || var.bandwidth_name != ""
      error_message = "The bandwidth name must be a non-empty string if the EIP address is not provided."
    }
  }
}

resource "huaweicloud_compute_instance" "test" {
  name               = var.instance_name
  image_id           = var.instance_image_id == "" ? try(data.huaweicloud_images_images.test[0].images[0].id, null) : var.instance_image_id
  flavor_id          = var.instance_flavor_id == "" ? try(data.huaweicloud_compute_flavors.test[0].flavors[0].id, null) : var.instance_flavor_id
  availability_zone  = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test[0].names[0], null) : var.availability_zone
  key_pair           = huaweicloud_kps_keypair.test.name
  user_data          = var.instance_user_data
  system_disk_type   = var.instance_system_disk_type
  system_disk_size   = var.instance_system_disk_size
  security_group_ids = length(var.security_group_ids) == 0 ? huaweicloud_networking_secgroup.test[*].id : var.security_group_ids

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }

  lifecycle {
    precondition {
      condition     = length(var.security_group_ids) != 0 || var.security_group_name != ""
      error_message = "The security_group_ids must be a non-empty list if the security_group_name is not provided."
    }
  }
}

resource "huaweicloud_compute_eip_associate" "test" {
  public_ip   = var.associate_eip_address == "" ? huaweicloud_vpc_eip.test[0].address : var.associate_eip_address
  instance_id = huaweicloud_compute_instance.test.id
}

resource "null_resource" "test" {
  depends_on = [huaweicloud_compute_eip_associate.test]

  provisioner "remote-exec" {
    connection {
      user        = "root"
      private_key = file(var.private_key_path)
      host        = huaweicloud_compute_eip_associate.test.public_ip
    }

    inline = var.instance_remote_exec_inline
  }
}

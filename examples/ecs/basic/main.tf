data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_images_images" "test" {
  flavor_id  = try(data.huaweicloud_compute_flavors.test.flavors[0].id, "")
  visibility = "public"
  os         = "Ubuntu"
}

resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1)
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.security_group_name
  delete_default_rules = true
}

resource "huaweicloud_compute_instance" "basic" {
  name               = var.instance_name
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  flavor_id          = try(data.huaweicloud_compute_flavors.test.flavors[0].id, "")
  image_id           = try(data.huaweicloud_images_images.test.images[0].id, "")
  security_group_ids = [
    huaweicloud_networking_secgroup.test.id
  ]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
  
  admin_pass = var.administrator_password

  lifecycle {
    ignore_changes = [
      admin_pass,
    ]
  }
}

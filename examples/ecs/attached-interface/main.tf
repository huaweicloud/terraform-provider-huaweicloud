data "huaweicloud_availability_zones" "myaz" {}

data "huaweicloud_compute_flavors" "myflavor" {
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_vpc_subnet" "mynet" {
  name = "subnet-default"
}

data "huaweicloud_images_image" "myimage" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

resource "huaweicloud_compute_instance" "myinstance" {
  name              = "basic"
  image_id          = data.huaweicloud_images_image.myimage.id
  flavor_id         = data.huaweicloud_compute_flavors.myflavor.ids[0]
  security_groups   = ["default"]
  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet.mynet.id
  }
}

data "huaweicloud_vpc" "myvpc" {
  name = "vpc-default"
}

resource "huaweicloud_vpc_subnet" "attach" {
  name       = "subnet-attach"
  cidr       = "192.168.1.0/24"
  gateway_ip = "192.168.1.1"
  vpc_id     = data.huaweicloud_vpc.myvpc.id

  availability_zone = data.huaweicloud_availability_zones.myaz.names[0]
}

resource "huaweicloud_compute_interface_attach" "attached" {
  instance_id = huaweicloud_compute_instance.myinstance.id
  network_id  = huaweicloud_vpc_subnet.attach.id

  # This is optional
  fixed_ip = "192.168.1.100"
}

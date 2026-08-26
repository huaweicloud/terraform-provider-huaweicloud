resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = var.eip_bandwidth_name
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_servicestagev3_environment" "test" {
  name        = var.environment_name
  description = var.environment_description
  vpc_id      = huaweicloud_vpc.test.id
}

resource "huaweicloud_servicestagev3_environment_associate" "test" {
  environment_id = huaweicloud_servicestagev3_environment.test.id

  resources {
    id   = huaweicloud_vpc_eip.test.id
    type = "eip"
  }
}

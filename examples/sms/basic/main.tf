#sms template
resource "huaweicloud_sms_template" "test" {
  name                = "template_test_old"
  is_template         = var.is_template
  region              = "ap-southeast-1"
  projectid           = "0e1dae200a00f3772f62c00030321606"
  availability_zone   = "ap-southeast-1a"
  target_server_name  = "serverName"
  flavor              = "flavor"
  volumetype          = var.volumetype
  data_volume_type    = var.data_volume_type
  target_password     = "123456"


  vpc_id = "ee3ef9f0-0c49-46d6-8da3-5ba364adc791"
  vpc_name = "vpc123"
  vpc_cidr = "192.168.0.0/16"

  nics {
    id = "autoCreate"
    name = "nics1"
    cidr = var.nics_cidr
  }

  nics {
    id = "autoCreate"
    name = "nics2"
    cidr = var.nics_cidr
  }

  security_groups {
    id = "autoCreate"
    name = "sg1"
  }

  security_groups {
    id = "autoCreate"
    name = "sg2"
  }

  publicip_type = var.publicip_type
  publicip_bandwidth_size = 1

  disk {
    index = 0
    name = "disk0"
    disktype = var.volumetype
    size = 40
  }

  disk {
    index = 1
    name = "disk1"
    disktype = var.volumetype
    size = 40
  }
}

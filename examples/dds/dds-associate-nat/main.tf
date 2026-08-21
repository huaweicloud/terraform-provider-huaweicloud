data "huaweicloud_availability_zones" "test" {
  count = var.availability_zone == "" ? 1 : 0
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
  name                 = var.security_group_name
  delete_default_rules = true
}

resource "huaweicloud_nat_gateway" "test" {
  name      = var.gateway_name
  spec      = var.gateway_spec
  vpc_id    = huaweicloud_vpc.test.id
  subnet_id = huaweicloud_vpc_subnet.test.id
}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = var.eip_type
  }

  bandwidth {
    name        = var.eip_bandwidth_name
    share_type  = "PER"
    size        = var.eip_bandwidth_size
    charge_mode = var.eip_bandwidth_charge_mode
  }
}

resource "huaweicloud_dds_instance" "test" {
  name              = var.instance_name
  availability_zone = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test[0].names[0], null) : var.availability_zone
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  mode              = var.instance_mode

  datastore {
    type           = var.database_type
    version        = var.database_version
    storage_engine = var.storage_engine
  }

  flavor {
    type      = var.node_type
    num       = var.node_number
    spec_code = var.node_spec_code
    storage   = var.node_storage_type
    size      = var.node_size
  }
}

data "huaweicloud_dds_instances" "test" {
  name = huaweicloud_dds_instance.test.name

  depends_on = [huaweicloud_dds_instance.test]
}

locals {
  nodeId = try([for v in flatten(data.huaweicloud_dds_instances.test.instances[*].groups[*].nodes) : v if v.role == "Primary"][0].id, "")
}

resource "huaweicloud_dds_bind_gateway" "test"{
  instance_id           = huaweicloud_dds_instance.test.id
  node_id               = local.nodeId
  nat_gateway_id        = huaweicloud_nat_gateway.test.id
  public_ip_id          = huaweicloud_vpc_eip.test.id
  external_service_port = var.external_service_port
}

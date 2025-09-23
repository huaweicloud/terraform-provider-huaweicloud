# Create VPC for virtual interface
resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

# Create virtual gateway
resource "huaweicloud_dc_virtual_gateway" "test" {
  vpc_id = huaweicloud_vpc.test.id
  name   = var.virtual_gateway_name

  local_ep_group = [
    huaweicloud_vpc.test.cidr,
  ]
}

# Create a DC virtual interface with VGW service type
resource "huaweicloud_dc_virtual_interface" "test" {
  direct_connect_id = var.direct_connect_id
  vgw_id            = huaweicloud_dc_virtual_gateway.test.id
  name              = var.virtual_interface_name
  description       = var.virtual_interface_description
  type              = var.virtual_interface_type
  route_mode        = var.route_mode
  vlan              = var.vlan
  bandwidth         = var.bandwidth

  remote_ep_group = var.remote_ep_group

  address_family       = var.address_family
  local_gateway_v4_ip  = var.local_gateway_v4_ip
  remote_gateway_v4_ip = var.remote_gateway_v4_ip

  enable_bfd = var.enable_bfd
  enable_nqa = var.enable_nqa

  tags = var.virtual_interface_tags
}

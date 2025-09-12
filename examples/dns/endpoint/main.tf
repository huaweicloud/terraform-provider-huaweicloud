# Create a DNS endpoint
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

resource "huaweicloud_dns_endpoint" "test" {
  name      = var.dns_endpoint_name
  direction = var.dns_endpoint_direction

  ip_addresses {
    subnet_id = huaweicloud_vpc_subnet.test.id
  }

  ip_addresses {
    subnet_id = huaweicloud_vpc_subnet.test.id
  }
}

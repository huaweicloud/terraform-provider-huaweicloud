resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr != "" ? var.subnet_cidr : cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0)
  gateway_ip = var.subnet_gateway_ip != "" ? var.subnet_gateway_ip : cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1)
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.security_group_name
  delete_default_rules = true
}

# ST.001 Disable
resource "huaweicloud_networking_secgroup_rule" "ingress" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  remote_ip_prefix  = var.vpc_cidr
}

resource "huaweicloud_networking_secgroup_rule" "egress" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = "egress"
  ethertype         = "IPv4"
  remote_ip_prefix  = "0.0.0.0/0"
  priority          = 1
}
# ST.001 Enable

resource "huaweicloud_ges_graph" "test" {
  name                  = var.graph_name
  graph_size_type_index = var.graph_size_type_index
  product_type          = "InMemory"
  cpu_arch              = "x86_64"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  crypt_algorithm       = "generalCipher"
  enable_https          = var.enable_https
  enable_multi_az       = var.enable_multi_az

  dynamic "lts_operation_trace" {
    for_each = var.enable_audit ? [1] : []

    content {
      enable_audit         = true
      audit_log_group_name = var.audit_log_group_name
    }
  }

  tags = var.graph_tags
}

# Create a VPC for the GES graph
resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

# Create a subnet for the GES graph
resource "huaweicloud_vpc_subnet" "test" {
  name       = var.subnet_name
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = var.subnet_cidr
  gateway_ip = var.gateway_ip
}

# Create a security group for the GES graph
resource "huaweicloud_networking_secgroup" "test" {
  name = var.security_group_name
}

# Create a GES graph instance
resource "huaweicloud_ges_graph" "test" {
  name                  = var.graph_name
  graph_size_type_index = var.graph_size_type_index
  cpu_arch              = var.graph_cpu_arch
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  crypt_algorithm       = var.graph_crypt_algorithm
  enable_https          = var.graph_enable_https

  tags = var.graph_tags
}

# Create a GES graph backup
resource "huaweicloud_ges_backup" "test" {
  graph_id = huaweicloud_ges_graph.test.id
}

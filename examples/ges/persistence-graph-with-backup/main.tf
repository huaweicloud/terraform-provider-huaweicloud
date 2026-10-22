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

resource "huaweicloud_ges_graph" "test" {
  name                  = var.graph_name
  graph_size_type_index = "6"
  product_type          = "Persistence"
  cpu_arch              = "x86_64"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  crypt_algorithm       = "generalCipher"
  enable_https          = false
  keep_backup           = var.keep_backup

  vertex_id_type {
    id_type = var.vertex_id_type
  }

  dynamic "encryption" {
    for_each = var.enable_encryption ? [1] : []

    content {
      enable        = true
      master_key_id = var.kms_key_id
    }
  }

  tags = var.graph_tags
}

resource "huaweicloud_ges_backup" "test" {
  depends_on = [
    huaweicloud_ges_graph.test,
  ]

  graph_id = huaweicloud_ges_graph.test.id
}

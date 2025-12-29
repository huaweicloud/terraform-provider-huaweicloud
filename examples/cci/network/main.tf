resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.security_group_name
  delete_default_rules = true
}

resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr
  gateway_ip = var.subnet_gateway_ip
}

resource "huaweicloud_cciv2_namespace" "test" {
  name = var.namespace_name
}

resource "huaweicloud_cciv2_network" "test" {
  depends_on = [huaweicloud_cciv2_namespace.test]

  namespace = huaweicloud_cciv2_namespace.test.name
  name      = var.network_name

  annotations = {
    "yangtse.io/project-id"                 = huaweicloud_cciv2_namespace.test.annotations["tenant.kubernetes.io/project-id"]
    "yangtse.io/domain-id"                  = huaweicloud_cciv2_namespace.test.annotations["tenant.kubernetes.io/domain-id"]
    "yangtse.io/warm-pool-size"             = var.warm_pool_size
    "yangtse.io/warm-pool-recycle-interval" = var.warm_pool_recycle_interval
  }

  subnets {
    subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id
  }

  security_group_ids = [huaweicloud_networking_secgroup.test.id]
}

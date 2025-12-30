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

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  depends_on = [huaweicloud_cciv2_namespace.test]

  name              = var.elb_name
  cross_vpc_backend = true
  vpc_id            = huaweicloud_vpc.test.id
  ipv4_subnet_id    = huaweicloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    try(data.huaweicloud_availability_zones.test.names[1], "")
  ]
}

resource "huaweicloud_cciv2_service" "test" {
  depends_on = [huaweicloud_elb_loadbalancer.test]

  namespace = huaweicloud_cciv2_namespace.test.name
  name      = var.service_name

  annotations = {
    "kubernetes.io/elb.class" = "elb"
    "kubernetes.io/elb.id"    = huaweicloud_elb_loadbalancer.test.id
  }

  ports {
    name         = "test"
    app_protocol = "TCP"
    protocol     = "TCP"
    port         = 87
    target_port  = 65529
  }

  selector = {
    app = var.selector_app
  }

  type = var.service_type

  lifecycle {
    ignore_changes = [
      annotations,
    ]
  }
}

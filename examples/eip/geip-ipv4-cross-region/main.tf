data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  count = var.instance_flavor_id == "" ? 1 : 0

  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], null)
  performance_type  = var.instance_flavor_performance_type
  cpu_core_count    = var.instance_flavor_cpu_core_count
  memory_size       = var.instance_flavor_memory_size
}

data "huaweicloud_images_images" "test" {
  count = var.instance_image_id == "" ? 1 : 0

  flavor_id  = var.instance_flavor_id == "" ? try(data.huaweicloud_compute_flavors.test[0].flavors[0].id, null) : var.instance_flavor_id
  visibility = var.instance_image_visibility
  os         = var.instance_image_os
}

resource "huaweicloud_vpc" "test" {
  name                  = var.vpc_name
  cidr                  = var.vpc_cidr
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr == "" ? cidrsubnet(huaweicloud_vpc.test.cidr, 4, 0) : var.subnet_cidr
  gateway_ip = var.subnet_gateway_ip == "" ? cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 0), 1) : var.subnet_gateway_ip
}

resource "huaweicloud_networking_secgroup" "test" {
  name = var.security_group_name
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  count = length(var.security_group_rule_configurations)

  direction         = lookup(var.security_group_rule_configurations[count.index], "direction", "ingress")
  ethertype         = lookup(var.security_group_rule_configurations[count.index], "ethertype", "IPv4")
  protocol          = lookup(var.security_group_rule_configurations[count.index], "protocol", null)
  ports             = lookup(var.security_group_rule_configurations[count.index], "ports", null)
  remote_ip_prefix  = lookup(var.security_group_rule_configurations[count.index], "remote_ip_prefix", "0.0.0.0/0")
  security_group_id = huaweicloud_networking_secgroup.test.id
}

resource "huaweicloud_compute_instance" "test" {
  name               = var.instance_name
  availability_zone  = try(data.huaweicloud_availability_zones.test.names[0], null)
  flavor_id          = var.instance_flavor_id == "" ? try(data.huaweicloud_compute_flavors.test[0].flavors[0].id, "") : var.instance_flavor_id
  image_id           = var.instance_image_id == "" ? try(data.huaweicloud_images_images.test[0].images[0].id, "") : var.instance_image_id
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  admin_pass         = var.instance_administrator_password

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }

  depends_on = [huaweicloud_networking_secgroup_rule.test]

  lifecycle {
    ignore_changes = [
      availability_zone,
      flavor_id,
      image_id,
      admin_pass,
    ]
  }
}

resource "huaweicloud_vpc_internet_gateway" "test" {
  vpc_id    = huaweicloud_vpc.test.id
  subnet_id = huaweicloud_vpc_subnet.test.id
  name      = var.internet_gateway_name
  add_route = var.internet_gateway_add_route
}

data "huaweicloud_global_eip_pools" "test" {
  access_site = var.global_eip_access_site != "" ? var.global_eip_access_site : null
  ip_version  = var.global_eip_ip_version
}

resource "huaweicloud_global_internet_bandwidth" "test" {
  access_site           = try(data.huaweicloud_global_eip_pools.test.geip_pools[0].access_site, null)
  charge_mode           = var.internet_bandwidth_charge_mode
  size                  = var.internet_bandwidth_size
  isp                   = try(data.huaweicloud_global_eip_pools.test.geip_pools[0].isp, null)
  name                  = var.internet_bandwidth_name
  ingress_size          = var.internet_bandwidth_ingress_size
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
  tags                  = var.internet_bandwidth_tags
}

resource "huaweicloud_global_eip" "test" {
  access_site           = huaweicloud_global_internet_bandwidth.test.access_site
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
  geip_pool_name        = try(data.huaweicloud_global_eip_pools.test.geip_pools[0].name, null)
  internet_bandwidth_id = huaweicloud_global_internet_bandwidth.test.id
  name                  = var.global_eip_name
  description           = var.global_eip_description
  tags                  = var.global_eip_tags
}

data "huaweicloud_identity_projects" "test" {
  name = huaweicloud_compute_instance.test.region
}

resource "huaweicloud_global_eip_associate" "test" {
  global_eip_id  = huaweicloud_global_eip.test.id
  is_reserve_gcb = false

  associate_instance {
    region        = huaweicloud_compute_instance.test.region
    project_id    = try(data.huaweicloud_identity_projects.test.projects[0].id, null)
    instance_type = "ECS"
    instance_id   = huaweicloud_compute_instance.test.id
  }

  gc_bandwidth {
    name                  = var.gc_bandwidth_name
    charge_mode           = var.gc_bandwidth_charge_mode
    size                  = var.gc_bandwidth_size
    enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
  }

  depends_on = [huaweicloud_vpc_internet_gateway.test]
}

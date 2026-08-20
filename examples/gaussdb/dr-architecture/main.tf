# ST.001 Disable
# Create VPC network
resource "huaweicloud_vpc" "primary" {
  name = var.primary_vpc_name
  cidr = var.primary_vpc_cidr
}

resource "huaweicloud_vpc_subnet" "primary" {
  count = length(var.primary_subnet_names)

  vpc_id            = huaweicloud_vpc.primary.id
  name              = try(var.primary_subnet_names[count.index], null)
  cidr              = cidrsubnet(huaweicloud_vpc.primary.cidr, 8, count.index)
  gateway_ip        = cidrhost(cidrsubnet(huaweicloud_vpc.primary.cidr, 8, count.index), 1)
  availability_zone = try(var.primary_availability_zones[count.index], null)
}

resource "huaweicloud_networking_secgroup" "primary" {
  name                 = var.primary_security_group_name
  delete_default_rules = true
}

resource "huaweicloud_networking_secgroup_rule" "primary" {
  count = length(var.secgroup_rules)

  security_group_id = huaweicloud_networking_secgroup.primary.id
  direction         = "ingress"
  ethertype         = "IPv4"
  remote_ip_prefix  = var.secgroup_rules[count.index].source == "local" ? huaweicloud_vpc.primary.cidr : huaweicloud_vpc.dr.cidr
  ports             = try(var.secgroup_rules[count.index].ports, null)
  protocol          = "tcp"
}

resource "huaweicloud_vpc" "dr" {
  provider = huaweicloud.dr

  name = var.dr_vpc_name
  cidr = var.dr_vpc_cidr
}

resource "huaweicloud_vpc_subnet" "dr" {
  provider = huaweicloud.dr

  vpc_id            = huaweicloud_vpc.dr.id
  name              = var.dr_subnet_name
  cidr              = cidrsubnet(huaweicloud_vpc.dr.cidr, 8, 0)
  gateway_ip        = cidrhost(cidrsubnet(huaweicloud_vpc.dr.cidr, 8, 0), 1)
  availability_zone = var.dr_availability_zone
}

resource "huaweicloud_networking_secgroup" "dr" {
  provider             = huaweicloud.dr

  name                 = var.dr_security_group_name
  delete_default_rules = true
}

resource "huaweicloud_networking_secgroup_rule" "dr" {
  provider = huaweicloud.dr

  count = length(var.secgroup_rules)

  security_group_id = huaweicloud_networking_secgroup.dr.id
  direction         = "ingress"
  ethertype         = "IPv4"
  remote_ip_prefix  = var.secgroup_rules[count.index].source == "local" ? huaweicloud_vpc.dr.cidr : huaweicloud_vpc.primary.cidr
  ports             = try(var.secgroup_rules[count.index].ports, null)
  protocol          = "tcp"
}

# Establish Cross region connection
data "huaweicloud_identity_projects" "primary" {
  name = var.region_name
}

data "huaweicloud_identity_projects" "dr" {
  provider = huaweicloud.dr

  name = var.dr_region_name
}

resource "huaweicloud_cc_connection" "dr" {
  name        = "${var.primary_vpc_name}_cc"
  description = "Cloud connection for GaussDB cross-region DR"
}

resource "huaweicloud_cc_network_instance" "primary" {
  cloud_connection_id = huaweicloud_cc_connection.dr.id
  instance_id         = huaweicloud_vpc.primary.id
  project_id          = try(data.huaweicloud_identity_projects.primary.projects[0].id, null)
  region_id           = var.region_name
  type                = "vpc"
  cidrs               = [huaweicloud_vpc.primary.cidr]
}

resource "huaweicloud_cc_network_instance" "dr" {
  provider            = huaweicloud.dr

  cloud_connection_id = huaweicloud_cc_connection.dr.id
  instance_id         = huaweicloud_vpc.dr.id
  project_id          = try(data.huaweicloud_identity_projects.dr.projects[0].id, null)
  region_id           = var.dr_region_name
  type                = "vpc"
  cidrs               = [huaweicloud_vpc.dr.cidr]
}

resource "huaweicloud_cc_bandwidth_package" "dr" {
  name           = "${var.primary_vpc_name}_cc_bp"
  local_area_id  = "Chinese-Mainland"
  remote_area_id = "Chinese-Mainland"
  bandwidth      = var.cc_bandwidth
  billing_mode   = 3
  charge_mode    = "bandwidth"
  resource_id    = huaweicloud_cc_connection.dr.id
  resource_type  = "cloud_connection"
}

resource "huaweicloud_cc_inter_region_bandwidth" "dr" {
  cloud_connection_id  = huaweicloud_cc_connection.dr.id
  bandwidth_package_id = huaweicloud_cc_bandwidth_package.dr.id
  bandwidth            = var.cc_bandwidth
  inter_region_ids     = [var.region_name, var.dr_region_name]
}

# Create Instance password
resource "random_password" "instance" {
  count = length(var.instance_passwords)

  length           = 12
  special          = true
  override_special = "!@%^*-_=+"

  keepers = {
    password_needed = try(var.instance_passwords[count.index], "") == "" ? "true" : "false"
  }

  lifecycle {
    create_before_destroy = true
  }
}

# Create GaussDB instance
resource "huaweicloud_gaussdb_instance" "primary" {
  name                  = var.primary_instance_name
  flavor                = var.instance_flavor
  password              = try(random_password.instance[0].result, null)
  vpc_id                = huaweicloud_vpc.primary.id
  subnet_id             = try(huaweicloud_vpc_subnet.primary[0].id, null)
  security_group_id     = huaweicloud_networking_secgroup.primary.id
  availability_zone     = var.primary_instance_availability_zones
  port                  = var.instance_db_port
  enterprise_project_id = var.enterprise_project_id

  ha {
    mode             = "centralization_standard"
    replication_mode = "sync"
    consistency      = "strong"
  }

  replica_num = 3

  volume {
    type = var.primary_instance_volume_type
    size = var.primary_instance_volume_size
  }

  lifecycle {
    ignore_changes = [
      flavor,
    ]
  }
}

resource "huaweicloud_gaussdb_instance" "dr" {
  provider              = huaweicloud.dr

  name                  = var.dr_instance_name
  flavor                = var.instance_flavor
  password              = try(random_password.instance[1].result, null)
  vpc_id                = huaweicloud_vpc.dr.id
  subnet_id             = huaweicloud_vpc_subnet.dr.id
  security_group_id     = huaweicloud_networking_secgroup.dr.id
  availability_zone     = var.dr_instance_availability_zones
  port                  = var.instance_db_port
  enterprise_project_id = var.enterprise_project_id

  ha {
    mode             = "centralization_standard"
    replication_mode = "sync"
    consistency      = "strong"
  }

  replica_num = 3

  volume {
    type = var.dr_instance_volume_type
    size = var.dr_instance_volume_size
  }

  lifecycle {
    ignore_changes = [
      flavor,
    ]
  }
}

# DR configuration
resource "huaweicloud_gaussdb_dr_configuration_reset" "primary" {
  instance_id        = huaweicloud_gaussdb_instance.primary.id
  opposite_data_cidr = huaweicloud_vpc_subnet.dr.cidr
}

resource "huaweicloud_gaussdb_dr_configuration_reset" "dr" {
  provider           = huaweicloud.dr

  instance_id        = huaweicloud_gaussdb_instance.dr.id
  opposite_data_cidr = try(huaweicloud_vpc_subnet.primary[0].cidr, null)
}
# ST.001 Enable

resource "huaweicloud_gaussdb_dr_relationship" "test" {
  instance_id      = huaweicloud_gaussdb_instance.primary.id
  disaster_type    = var.dr_disaster_type
  dr_ip            = try(huaweicloud_gaussdb_instance.dr.private_ips[0], null)
  dr_user_name     = var.dr_user_name
  dr_user_password = var.dr_user_password

  depends_on = [
    huaweicloud_gaussdb_dr_configuration_reset.primary,
    huaweicloud_gaussdb_dr_configuration_reset.dr,
    huaweicloud_cc_network_instance.primary,
    huaweicloud_cc_network_instance.dr,
    huaweicloud_cc_inter_region_bandwidth.dr,
  ]
}

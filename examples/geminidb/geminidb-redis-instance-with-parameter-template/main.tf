# Create VPC and subnet for GeminiDB Redis instance
resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr == "" ? cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0) : var.subnet_cidr
  gateway_ip = var.gateway_ip == "" ? cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1) : var.gateway_ip
}

# Create security group for GeminiDB Redis instance
resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.security_group_name
  delete_default_rules = true
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = "ingress"
  ethertype         = "IPv4"
  remote_ip_prefix  = huaweicloud_vpc.test.cidr
  ports             = var.instance_db_port
  protocol          = "tcp"
}

data "huaweicloud_availability_zones" "test" {
  count = var.availability_zone == "" ? 1 : 0
}

resource "random_password" "test" {
  length           = 16
  min_upper        = 1
  min_lower        = 1
  min_numeric      = 1
  min_special      = 1
  special          = true
  override_special = "~!@#%^*-_=+?"
}

# Create GeminiDB Redis parameter template
resource "huaweicloud_geminidb_parameter_template" "test" {
  name        = var.parameter_template_name
  description = var.parameter_template_description

  datastore {
    type    = var.datastore_type
    version = var.datastore_version
  }

  values = var.parameter_template_values
}

# Create GeminiDB Redis instance
resource "huaweicloud_geminidb_instance" "test" {
  name              = var.instance_name
  availability_zone = var.availability_zone != "" ? var.availability_zone : join(",", slice(data.huaweicloud_availability_zones.test[0].names, 0, 1))
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = var.instance_password != "" ? var.instance_password : random_password.test.result
  mode              = var.instance_mode
  port              = var.instance_db_port
  ssl_option        = var.instance_ssl_option
  configuration_id  = huaweicloud_geminidb_parameter_template.test.id

  datastore {
    type           = var.datastore_type
    version        = var.datastore_version
    storage_engine = var.datastore_storage_engine
  }

  flavor {
    num       = var.instance_flavor_num
    size      = var.instance_flavor_size
    storage   = var.instance_flavor_storage
    spec_code = var.instance_flavor_spec_code
  }

  backup_strategy {
    start_time = var.instance_backup_time_window
    keep_days  = var.instance_backup_keep_days
  }

  tags = var.tags

  lifecycle {
    ignore_changes = [
      flavor.0.spec_code,
    ]
  }
}

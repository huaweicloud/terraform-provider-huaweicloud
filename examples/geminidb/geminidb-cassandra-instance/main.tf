# Create VPC and subnet for GeminiDB Cassandra instance
resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

data "huaweicloud_availability_zones" "test" {
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr == "" ? cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0) : var.subnet_cidr
  gateway_ip = var.gateway_ip == "" ? cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 8, 0), 1) : var.gateway_ip
}

# Query GeminiDB NoSQL flavors for Cassandra
data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus             = var.vcpus
  engine            = "cassandra"
  availability_zone = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test.names[0], null) : var.availability_zone
}

# Create security group for GeminiDB Cassandra instance
resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.security_group_name
  delete_default_rules = true
}

# Generate random password if not provided
resource "random_password" "test" {
  count            = var.instance_password == "" ? 1 : 0

  length           = 12
  special          = true
  override_special = "!@%^*-_=+"
}

# Create GeminiDB Cassandra instance
resource "huaweicloud_geminidb_instance" "test" {
  name              = var.instance_name
  availability_zone = var.availability_zone == "" ? try(data.huaweicloud_availability_zones.test.names[0], null) : var.availability_zone
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "Qi@123456"
  mode              = var.instance_mode
  port              = var.instance_db_port
  ssl_option        = var.instance_ssl_option

  datastore {
    type           = "cassandra"
    version        = ""
    storage_engine = "rocksDB"
  }

  flavor {
    num       = var.instance_flavor_num
    size      = var.instance_flavor_size
    storage   = var.instance_flavor_storage
    spec_code = var.instance_flavor_spec_code != "" ? var.instance_flavor_spec_code : try(data.huaweicloud_gaussdb_nosql_flavors.test.flavors[0].name, null)
  }

  backup_strategy {
    start_time = var.instance_backup_time_window
    keep_days  = var.instance_backup_keep_days
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  auto_renew    = "true"
  period        = 1

  tags = var.tags

  lifecycle {
    ignore_changes = [
      flavor.0.spec_code,
    ]
  }
}

# Create GeminiDB backup
resource "huaweicloud_geminidb_backup" "test" {
  instance_id = huaweicloud_geminidb_instance.test.id
  name        = var.backup_name
  description = var.backup_description

  depends_on = [huaweicloud_geminidb_instance.test]
}

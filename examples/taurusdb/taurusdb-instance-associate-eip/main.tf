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

resource "huaweicloud_networking_secgroup" "test" {
  name                 = var.security_group_name
  delete_default_rules = true
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = "ingress"
  ethertype         = "IPv4"
  remote_ip_prefix  = var.vpc_cidr
  ports             = var.instance_db_port
  protocol          = "tcp"
}

data "huaweicloud_taurusdb_flavors" "test" {
  engine                 = "gaussdb-mysql"
  version                = "8.0"
  availability_zone_mode = var.availability_zone_mode
}

# Query existing EIP if associate_eip_address is specified
data "huaweicloud_vpc_eip" "test" {
  count = var.associate_eip_address != "" ? 1 : 0

  public_ip = var.associate_eip_address
}

# Create a new EIP if associate_eip_address is not specified
resource "huaweicloud_vpc_eip" "test" {
  count = var.associate_eip_address == "" ? 1 : 0

  publicip {
    type = var.eip_type
  }

  bandwidth {
    name        = var.bandwidth_name
    size        = var.bandwidth_size
    share_type  = var.bandwidth_share_type
    charge_mode = var.bandwidth_charge_mode
  }
}

locals {
  # Get the first available AZ from the flavor's az_status
  available_azs = try([for k, v in data.huaweicloud_taurusdb_flavors.test.flavors[0].az_status : k if v == "normal"], [])
  master_az     = var.master_availability_zone != "" ? var.master_availability_zone : try(local.available_azs[0], "")

  public_ip    = var.associate_eip_address != "" ? var.associate_eip_address : huaweicloud_vpc_eip.test[0].address
  public_ip_id = var.associate_eip_address != "" ? data.huaweicloud_vpc_eip.test[0].id : huaweicloud_vpc_eip.test[0].id
}

resource "random_password" "test" {
  count = var.instance_password == "" ? 1 : 0

  length           = 12
  special          = true
  override_special = "!@%^*-_=+"
}

resource "huaweicloud_taurusdb_instance" "test" {
  name                     = var.instance_name
  flavor                   = var.instance_flavor_ref != "" ? var.instance_flavor_ref : try(data.huaweicloud_taurusdb_flavors.test.flavors[0].name, "")
  vpc_id                   = huaweicloud_vpc.test.id
  subnet_id                = huaweicloud_vpc_subnet.test.id
  security_group_id        = huaweicloud_networking_secgroup.test.id
  password                 = var.instance_password != "" ? var.instance_password : try(random_password.test[0].result, null)
  mode                     = var.instance_mode
  availability_zone_mode   = var.availability_zone_mode
  master_availability_zone = local.master_az
  read_replicas            = var.read_replicas
  enterprise_project_id    = var.enterprise_project_id

  datastore {
    engine  = "gaussdb-mysql"
    version = "8.0"
  }

  backup_strategy {
    start_time = var.instance_backup_time_window
    keep_days  = tostring(var.instance_backup_keep_days)
  }

  lifecycle {
    ignore_changes = [
      auto_scaling.0.scaling_strategy, password,
    ]
  }
}

# Associate EIP to the TaurusDB instance
# If associate_eip_address is specified, use the test EIP; otherwise, use the newly created EIP
resource "huaweicloud_taurusdb_eip_associate" "test" {
  instance_id  = huaweicloud_taurusdb_instance.test.id
  public_ip    = local.public_ip
  public_ip_id = local.public_ip_id
}

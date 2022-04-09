# create vpc and subnet
resource "huaweicloud_vpc" "vpc_1" {
  name = var.vpc_name
  cidr = "192.168.0.0/24"
}

resource "huaweicloud_vpc_subnet" "vpc_subnet_1" {
  name       = var.subnet_name
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.vpc_1.id
}

data "huaweicloud_availability_zones" "zones" {}

# 1.1 Query the available flavors of the specified capacity.
data "huaweicloud_dcs_flavors" "single_flavors" {
  cache_mode = "single"
  capacity   = 0.125
}

# 1.2 Create Single-node Redis instance
resource "huaweicloud_dcs_instance" "instance_1" {
  name           = "single_instance"
  engine         = "redis"
  engine_version = "5.0"
  capacity       = data.huaweicloud_dcs_flavors.single_flavors.capacity
  flavor         = data.huaweicloud_dcs_flavors.single_flavors.flavors[0].name

  availability_zones = [
    data.huaweicloud_availability_zones.zones.names[0]
  ]
  vpc_id        = huaweicloud_vpc.vpc_1.id
  subnet_id     = huaweicloud_vpc_subnet.vpc_subnet_1.id
  password      = "Abc@123"
  charging_mode = "postPaid"
}

# 2.1 Query the available Master/Standby flavors of the specified capacity.
data "huaweicloud_dcs_flavors" "master_standby_flavors" {
  cache_mode = "ha"
  capacity   = 0.125
}

# 2.2 Create Master/Standby Redis instances
resource "huaweicloud_dcs_instance" "master_standby_instance" {
  engine         = "redis"
  name           = "master_standby_instance"
  engine_version = "5.0"
  capacity       = data.huaweicloud_dcs_flavors.master_standby_flavors.capacity
  flavor         = data.huaweicloud_dcs_flavors.master_standby_flavors.flavors[0].name

  availability_zones = [
    data.huaweicloud_availability_zones.zones.names[0],
    data.huaweicloud_availability_zones.zones.names[1]
  ]
  vpc_id        = huaweicloud_vpc.vpc_1.id
  subnet_id     = huaweicloud_vpc_subnet.vpc_subnet_1.id
  password      = "Abc@123"
  charging_mode = "postPaid"
}

# 3.1 Query the available Proxy Cluster flavors of the specified capacity.
data "huaweicloud_dcs_flavors" "proxy_cluster_flavors" {
  cache_mode = "proxy"
  capacity   = 4
}

# 3.2 Create Proxy Cluster Redis instances
resource "huaweicloud_dcs_instance" "proxy_cluster_instance" {
  engine         = "redis"
  name           = "proxy_cluster_instance"
  engine_version = "5.0"
  capacity       = data.huaweicloud_dcs_flavors.proxy_cluster_flavors.capacity
  flavor         = data.huaweicloud_dcs_flavors.proxy_cluster_flavors.flavors[0].name

  availability_zones = [
    data.huaweicloud_availability_zones.zones.names[0],
    data.huaweicloud_availability_zones.zones.names[1]
  ]
  vpc_id        = huaweicloud_vpc.vpc_1.id
  subnet_id     = huaweicloud_vpc_subnet.vpc_subnet_1.id
  password      = "Abc@123"
  charging_mode = "postPaid"
}

# 4.1 Query the available Cluster flavors of the specified capacity.
data "huaweicloud_dcs_flavors" "redis_cluster_flavors" {
  cache_mode = "cluster"
  capacity   = 4
}

# 4.2 Create Cluster Redis instances, and configure Backup Policy, Whitelists and Tags.
resource "huaweicloud_dcs_instance" "redis_cluster_instance" {
  engine         = "redis"
  name           = "redis_cluster_instance"
  engine_version = "5.0"
  capacity       = data.huaweicloud_dcs_flavors.redis_cluster_flavors.capacity
  flavor         = data.huaweicloud_dcs_flavors.redis_cluster_flavors.flavors[0].name

  availability_zones = [
    data.huaweicloud_availability_zones.zones.names[0],
    data.huaweicloud_availability_zones.zones.names[1]
  ]
  vpc_id        = huaweicloud_vpc.vpc_1.id
  subnet_id     = huaweicloud_vpc_subnet.vpc_subnet_1.id
  password      = "Abc@123"
  charging_mode = "postPaid"

  backup_policy {
    backup_type = "auto"
    save_days   = 3
    period_type = "weekly"
    backup_at   = [1, 2, 3, 4, 5, 6, 7]
    begin_at    = "02:00-04:00"
  }

  whitelists {
    group_name = "group_1"
    ip_address = ["192.168.1.0/24"]
  }
  whitelists {
    group_name = "group_2"
    ip_address = ["10.11.3.0/24"]
  }

  rename_commands = {
    "command" : "cmd",
    "keys" : "key",
    "flushdb" : "flshdb",
    "flushall" : "flusall",
    "hgetall" : "getall"
  }

  tags = {
    "level" : "A"
  }
}

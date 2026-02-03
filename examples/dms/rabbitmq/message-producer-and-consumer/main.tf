data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_rabbitmq_flavors" "test" {
  type               = var.instance_flavor_type
  storage_spec_code  = var.instance_storage_spec_code
  availability_zones = try(slice(data.huaweicloud_availability_zones.test.names, 0, 1), null)
}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], null)
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_images_images" "test" {
  name       = var.ecs_image_name
  visibility = "public"
}

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
  name = var.security_group_name
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  count = length(var.security_group_rule_configurations)

  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = lookup(var.security_group_rule_configurations[count.index], "direction", null)
  ethertype         = lookup(var.security_group_rule_configurations[count.index], "ethertype", null)
  protocol          = lookup(var.security_group_rule_configurations[count.index], "protocol", null)
  port_range_min    = lookup(var.security_group_rule_configurations[count.index], "port_range_min", null)
  port_range_max    = lookup(var.security_group_rule_configurations[count.index], "port_range_max", null)
  remote_ip_prefix  = lookup(var.security_group_rule_configurations[count.index], "remote_ip_prefix", null)
  description       = lookup(var.security_group_rule_configurations[count.index], "description", null)
}

resource "huaweicloud_dms_rabbitmq_instance" "test" {
  name                  = var.instance_name
  engine_version        = var.instance_engine_version
  flavor_id             = try(data.huaweicloud_dms_rabbitmq_flavors.test.flavors[0].id, null)
  vpc_id                = huaweicloud_vpc.test.id
  network_id            = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  availability_zones    = try([data.huaweicloud_availability_zones.test.names[0]], null)
  broker_num            = var.instance_broker_num
  storage_space         = var.instance_storage_space
  storage_spec_code     = var.instance_storage_spec_code
  ssl_enable            = var.instance_ssl_enable
  access_user           = var.instance_access_user_name
  password              = var.instance_password
  description           = var.instance_description
  enterprise_project_id = var.enterprise_project_id
  tags                  = var.instance_tags
  charging_mode         = var.charging_mode
  period_unit           = var.period_unit
  period                = var.period
  auto_renew            = var.auto_renew

  lifecycle {
    ignore_changes = [
      flavor_id,
      availability_zones,
    ]
  }
}

resource "huaweicloud_dms_rabbitmq_vhost" "test" {
  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  name        = var.vhost_name
}

resource "huaweicloud_dms_rabbitmq_exchange" "test" {
  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  vhost       = var.vhost_name
  name        = var.exchange_name
  type        = var.exchange_type
  auto_delete = false
  durable     = true
  internal    = false

  depends_on = [
    huaweicloud_dms_rabbitmq_instance.test,
    huaweicloud_dms_rabbitmq_vhost.test
  ]
}

resource "huaweicloud_dms_rabbitmq_queue" "test" {
  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  vhost       = var.vhost_name
  name        = var.queue_name
  auto_delete = false
  durable     = true

  depends_on = [
    huaweicloud_dms_rabbitmq_instance.test,
    huaweicloud_dms_rabbitmq_vhost.test
  ]
}

resource "huaweicloud_dms_rabbitmq_exchange_associate" "test" {
  instance_id      = huaweicloud_dms_rabbitmq_instance.test.id
  vhost            = var.vhost_name
  exchange         = var.exchange_name
  destination_type = "Queue"
  destination      = var.queue_name
  routing_key      = var.queue_name

  depends_on = [
    huaweicloud_dms_rabbitmq_instance.test,
    huaweicloud_dms_rabbitmq_vhost.test,
    huaweicloud_dms_rabbitmq_exchange.test,
    huaweicloud_dms_rabbitmq_queue.test
  ]
}

# ST.001 Disable
resource "huaweicloud_compute_instance" "producer" {
  name              = var.producer_instance_name
  image_id          = try(data.huaweicloud_images_images.test.images[0].id, null)
  flavor_id         = try(data.huaweicloud_compute_flavors.test.flavors[0].id, null)
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], null)
  admin_pass        = var.instance_password
  eip_type          = var.eip_type

  security_group_ids = [huaweicloud_networking_secgroup.test.id]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }

  bandwidth {
    share_type  = var.eip_share_type
    size        = var.eip_size
    charge_mode = var.eip_charge_mode
  }

  user_data = templatefile("${path.module}/templates/user_data_producer_tpl", {
    producer_script   = file("${path.module}/apps/producer.py")
    rabbitmq_host     = huaweicloud_dms_rabbitmq_instance.test.connect_address
    rabbitmq_user     = var.instance_access_user_name
    rabbitmq_password = var.instance_password
    rabbitmq_vhost    = var.vhost_name
    queue_name        = var.queue_name
    exchange_name     = var.exchange_name
    exchange_type     = var.exchange_type
    routing_key       = var.queue_name
    message_interval  = var.message_interval
  })

  depends_on = [
    huaweicloud_dms_rabbitmq_instance.test,
    huaweicloud_dms_rabbitmq_vhost.test,
    huaweicloud_dms_rabbitmq_queue.test,
    huaweicloud_dms_rabbitmq_exchange_associate.test
  ]

  lifecycle {
    ignore_changes = [
      image_id,
      flavor_id,
      availability_zone
    ]
  }
}

resource "huaweicloud_compute_instance" "consumer" {
  name              = var.consumer_instance_name
  image_id          = try(data.huaweicloud_images_images.test.images[0].id, null)
  flavor_id         = try(data.huaweicloud_compute_flavors.test.flavors[0].id, null)
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], null)
  admin_pass        = var.instance_password
  eip_type          = var.eip_type

  security_group_ids = [huaweicloud_networking_secgroup.test.id]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }

  bandwidth {
    share_type  = var.eip_share_type
    size        = var.eip_size
    charge_mode = var.eip_charge_mode
  }

  user_data = templatefile("${path.module}/templates/user_data_consumer_tpl", {
    consumer_script   = file("${path.module}/apps/consumer.py")
    rabbitmq_host     = huaweicloud_dms_rabbitmq_instance.test.connect_address
    rabbitmq_user     = var.instance_access_user_name
    rabbitmq_password = var.instance_password
    rabbitmq_vhost    = var.vhost_name
    queue_name        = var.queue_name
  })

  depends_on = [
    huaweicloud_dms_rabbitmq_instance.test,
    huaweicloud_dms_rabbitmq_vhost.test,
    huaweicloud_dms_rabbitmq_queue.test,
    huaweicloud_dms_rabbitmq_exchange_associate.test
  ]

  lifecycle {
    ignore_changes = [
      image_id,
      flavor_id,
      availability_zone
    ]
  }
}
# ST.001 Enable

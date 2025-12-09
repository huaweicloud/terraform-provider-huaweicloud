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

data "huaweicloud_availability_zones" "test" {
  count = length(var.availability_zones) > 0 ? 0 : 1
}

resource "huaweicloud_apig_instance" "test" {
  name                  = var.instance_name
  edition               = var.instance_edition
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = var.enterprise_project_id
  availability_zones    = length(var.availability_zones) > 0 ? var.availability_zones : try(slice(data.huaweicloud_availability_zones.test[0].names, 0, var.availability_zones_count), null)
}

data "huaweicloud_dms_kafka_flavors" "test" {
  count = var.kafka_instance_flavor_id != "" ? 0 : 1

  type               = var.kafka_instance_flavor_type
  storage_spec_code  = var.kafka_instance_storage_spec_code
  availability_zones = length(var.availability_zones) > 0 ? var.availability_zones : try(slice(data.huaweicloud_availability_zones.test[0].names, 0, 3))
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name               = var.kafka_instance_name
  description        = var.kafka_instance_description
  availability_zones = length(var.availability_zones) > 0 ? var.availability_zones : try(slice(data.huaweicloud_availability_zones.test[0].names, 0, 3))
  vpc_id             = huaweicloud_vpc.test.id
  network_id         = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  ssl_enable         = var.kafka_instance_ssl_enable
  flavor_id          = var.kafka_instance_flavor_id != "" ? var.kafka_instance_flavor_id : try(data.huaweicloud_dms_kafka_flavors.test[0].flavors[0].id, null)
  engine_version     = var.kafka_instance_engine_version
  storage_spec_code  = var.kafka_instance_storage_spec_code
  storage_space      = var.kafka_instance_storage_space
  broker_num         = var.kafka_instance_broker_num
  charging_mode      = var.kafka_charging_mode
  period_unit        = var.kafka_period_unit
  period             = var.kafka_period
  auto_renew         = var.kafka_auto_new
  access_user        = var.kafka_instance_user_name
  password           = var.kafka_instance_user_password

  # If you want to change some of the following parameters, you need to remove the corresponding fields from "lifecycle.ignore_changes".
  lifecycle {
    ignore_changes = [
      access_user,
      availability_zones,
      flavor_id,
    ]
  }
}

resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id = huaweicloud_dms_kafka_instance.test.id
  name        = var.kafka_topic_name
  partitions  = var.kafka_topic_partitions
}

resource "huaweicloud_apig_plugin" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  name        = var.plugin_name
  description = var.plugin_description
  type        = "kafka_log"

  content = jsonencode({
    broker_list     = var.kafka_security_protocol == "PLAINTEXT" ? (split(",", huaweicloud_dms_kafka_instance.test.port_protocol[0].private_plain_address)) : var.kafka_security_protocol == "SASL_PLAINTEXT" ? (split(",",huaweicloud_dms_kafka_instance.test.port_protocol[0].private_sasl_plaintext_address)) : (split(",", huaweicloud_dms_kafka_instance.test.port_protocol[0].private_sasl_ssl_address))
    topic           = var.kafka_topic_name
    key             = var.kafka_message_key
    max_retry_count = var.kafka_max_retry_count
    retry_backoff   = var.kafka_retry_backoff

    sasl_config = {
      security_protocol = var.kafka_security_protocol
      sasl_mechanisms   = var.kafka_sasl_mechanisms
      sasl_username     = var.kafka_sasl_username != "" ? nonsensitive(var.kafka_sasl_username) : (var.kafka_security_protocol == "PLAINTEXT" ? "" : nonsensitive(var.kafka_access_user))
      sasl_password     = var.kafka_sasl_password != "" ? nonsensitive(var.kafka_sasl_password) : (var.kafka_security_protocol == "PLAINTEXT" ? "" : nonsensitive(var.kafka_password))
      ssl_ca_content    = var.kafka_ssl_ca_content != "" ? nonsensitive(var.kafka_ssl_ca_content) : ""
    }
  })

  lifecycle {
    ignore_changes = [
      content,
    ]
  }
}

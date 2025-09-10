resource "huaweicloud_vpc" "test" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = var.subnet_name
  cidr       = var.subnet_cidr
  gateway_ip = var.subnet_gateway
}

resource "huaweicloud_networking_secgroup" "test" {
  name = var.security_group_name
}

resource "huaweicloud_obs_bucket" "test" {
  bucket        = var.bucket_name
  acl           = var.bucket_acl
  force_destroy = true
}

data "huaweicloud_availability_zones" "test" {
  count = length(var.availability_zones) == 0 ? 1 : 0
}

data "huaweicloud_dms_kafka_flavors" "test" {
  count = var.instance_flavor_id == "" ? 1 : 0

  type               = var.instance_flavor_type
  availability_zones = length(var.availability_zones) == 0 ? try(slice(data.huaweicloud_availability_zones.test[0].names, 0, 3)) : var.availability_zones
  storage_spec_code  = var.instance_storage_spec_code
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name               = var.instance_name
  availability_zones = length(var.availability_zones) == 0 ? try(slice(data.huaweicloud_availability_zones.test[0].names, 0, 3)) : var.availability_zones
  engine_version     = var.instance_engine_version
  flavor_id          = var.instance_flavor_id == "" ? try(data.huaweicloud_dms_kafka_flavors.test[0].flavors[0].id, null) : var.instance_flavor_id
  storage_spec_code  = var.instance_storage_spec_code
  storage_space      = var.instance_storage_space
  broker_num         = var.instance_broker_num
  vpc_id             = huaweicloud_vpc.test.id
  network_id         = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  ssl_enable         = var.instance_ssl_enable
  description        = var.instance_description
  security_protocol  = var.instance_security_protocol
  charging_mode      = var.charging_mode

  lifecycle {
    ignore_changes = [
      availability_zones,
      flavor_id,
    ]
  }
}

resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id = huaweicloud_dms_kafka_instance.test.id
  name        = var.topic_name
  partitions  = var.topic_partitions
}

locals {
  kafka_connect_with_port = join(
    ",",
    formatlist(
      "%s:${huaweicloud_dms_kafka_instance.test.port}",
      split(",", huaweicloud_dms_kafka_instance.test.connect_address)
    )
  )
}

resource "huaweicloud_eg_connection" "test" {
  name      = var.connection_name
  type      = "KAFKA"
  vpc_id    = huaweicloud_vpc.test.id
  subnet_id = huaweicloud_vpc_subnet.test.id

  kafka_detail {
    instance_id     = huaweicloud_dms_kafka_instance.test.id
    connect_address = local.kafka_connect_with_port
    acks            = var.connection_acks
  }

  lifecycle {
    ignore_changes = [
      kafka_detail[0].user_name,
      kafka_detail[0].password,
    ]
  }

  depends_on = [
    huaweicloud_dms_kafka_topic.test
  ]
}

# The time sleep used to wait for Kafka topic and EG connection to be ready
resource "time_sleep" "test" {
  create_duration = "5s"

  depends_on = [
    huaweicloud_eg_connection.test
  ]
}

data "huaweicloud_eg_event_channels" "test" {
  provider_type = "OFFICIAL"
  name          = "default"
}

resource "huaweicloud_eg_event_subscription" "test" {
  channel_id = try(data.huaweicloud_eg_event_channels.test.channels[0].id, "")
  name       = try(data.huaweicloud_eg_event_channels.test.channels[0].name, "")

  sources {
    name          = "HC.OBS"
    provider_type = "OFFICIAL"

    filter_rule = jsonencode({
      "source" : [
        {
          "op" : "StringIn",
          "values" : ["HC.OBS"]
        }
      ],
      "type" : [
        {
          "op" : "StringIn",
          "values" : var.subscription_source_values
        }
      ],
    })
  }

  targets {
    name          = "HC.Kafka"
    provider_type = "OFFICIAL"
    connection_id = huaweicloud_eg_connection.test.id

    transform = jsonencode({
      "type" : "ORIGINAL",
    })

    detail_name = "kafka_detail"
    detail      = jsonencode({
      "topic" : huaweicloud_dms_kafka_topic.test.name
      "key_transform" : {
        "type" : "ORIGINAL",
      }
    })
  }

  lifecycle {
    ignore_changes = [
      sources
    ]
  }

  depends_on = [
    time_sleep.test
  ]
}

resource "huaweicloud_obs_bucket_object" "test" {
  bucket       = huaweicloud_obs_bucket.test.id
  key          = var.object_extension_name != "" ? format("%s%s", var.object_name, var.object_extension_name) : var.object_name
  content_type = "application/xml"
  content      = var.object_upload_content
}

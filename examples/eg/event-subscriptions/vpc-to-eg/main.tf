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

resource "huaweicloud_eg_custom_event_channel" "test" {
  name = var.channel_name
}

data "huaweicloud_eg_event_channels" "test" {
  provider_type = "OFFICIAL"
  name          = "default"
}

# The time sleep used to wait for EG channel to be ready
resource "time_sleep" "test" {
  create_duration = "3s"

  depends_on = [
    huaweicloud_eg_custom_event_channel.test
  ]
}

resource "huaweicloud_eg_event_subscription" "test" {
  channel_id = try(data.huaweicloud_eg_event_channels.test.channels[0].id, "")
  name       = try(data.huaweicloud_eg_event_channels.test.channels[0].name, "")

  sources {
    name          = var.sources_name
    provider_type = var.sources_provider_type

    filter_rule = jsonencode({
      "source" : [
        {
          "op" : var.source_op,
          "values" : ["HC.VPC"]
        }
      ],
      "type" : [
        {
          "op" : var.type_op,
          "values" : var.subscription_source_values
        }
      ],
    })
  }

  targets {
    name          = var.targets_name
    provider_type = var.targets_provider_type
    detail_name   = var.detail_name
    transform     = jsonencode(var.transform)
    detail        = jsonencode({
      "agency_name" : var.agency_name
      "target_channel_id" : huaweicloud_eg_custom_event_channel.test.id
      "target_project_id" : var.target_project_id
      "target_region" : var.target_region_name != "" ? var.target_region_name : var.region_name
    })
  }

  lifecycle {
    ignore_changes = [
      sources, targets
    ]
  }

  depends_on = [
    time_sleep.test
  ]
}

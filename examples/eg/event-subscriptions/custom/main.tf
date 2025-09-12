resource "huaweicloud_eg_custom_event_channel" "test" {
  name = var.channel_name
}

resource "huaweicloud_eg_custom_event_source" "test" {
  channel_id = huaweicloud_eg_custom_event_channel.test.id
  name       = var.source_name
  type       = var.source_type
}

data "huaweicloud_eg_connections" "test" {
  name = var.connection_name
}

# The time sleep used to wait for custom channel and source to be ready
resource "time_sleep" "test" {
  create_duration = "3s"

  depends_on = [
    huaweicloud_eg_custom_event_channel.test
  ]
}

resource "huaweicloud_eg_event_subscription" "test" {
  channel_id = huaweicloud_eg_custom_event_channel.test.id
  name       = var.subscription_name

  sources {
    name          = huaweicloud_eg_custom_event_channel.test.name
    provider_type = var.sources_provider_type

    filter_rule = jsonencode({
      "source" : [
        {
          "op" : var.source_op,
          "values" : [huaweicloud_eg_custom_event_channel.test.name]
        }
      ]
    })
  }

  targets {
    name          = var.targets_name
    provider_type = var.targets_provider_type
    connection_id = try(data.huaweicloud_eg_connections.test.connections[0].id, "")
    transform     = jsonencode(var.transform)
    detail_name   = var.detail_name
    detail        = jsonencode({
      "url" : var.target_url
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

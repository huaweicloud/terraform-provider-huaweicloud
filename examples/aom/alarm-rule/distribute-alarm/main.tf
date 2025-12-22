data "huaweicloud_dcs_instances" "test" {
  name = var.dcs_instance_name
}

locals {
  enterprise_project_id = try(data.huaweicloud_dcs_instances.test.instances[0].enterprise_project_id, null)
}

resource "huaweicloud_lts_group" "test" {
  group_name            = var.lts_group_name
  ttl_in_days           = 30
  enterprise_project_id = local.enterprise_project_id
}

resource "huaweicloud_lts_stream" "test" {
  group_id              = huaweicloud_lts_group.test.id
  stream_name           = var.lts_stream_name
  enterprise_project_id = local.enterprise_project_id
}

resource "huaweicloud_smn_topic" "test" {
  name                  = var.smn_topic_name
  enterprise_project_id = local.enterprise_project_id
}

resource "huaweicloud_smn_logtank" "test" {
  topic_urn     = huaweicloud_smn_topic.test.topic_urn
  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id
}

resource "huaweicloud_aom_alarm_action_rule" "test" {
  name                  = var.alarm_action_rule_name
  user_name             = var.alarm_action_rule_user_name
  type                  = var.alarm_action_rule_type
  notification_template = "aom.built-in.template.zh"

  smn_topics {
    topic_urn = huaweicloud_smn_topic.test.topic_urn
  }
}

data "huaweicloud_identity_projects" "test" {
  # ST.002 Disable
  name = var.region_name
  # ST.002 Enable
}

locals {
  exact_project_id = try([for v in data.huaweicloud_identity_projects.test.projects : v.id if v.name == var.region_name][0], null)
}

resource "huaweicloud_tms_resource_tags" "test" {
  project_id = local.exact_project_id

  resources {
    resource_type = "dcs"
    resource_id   = try(data.huaweicloud_dcs_instances.test.instances[0].id, null)
  }

  tags = var.alarm_rule_matric_dimension_tags
}

resource "huaweicloud_aom_prom_instance" "test" {
  depends_on = [huaweicloud_tms_resource_tags.test]

  prom_name             = var.prometheus_instance_name
  prom_type             = "CLOUD_SERVICE"
  enterprise_project_id = local.enterprise_project_id
}

resource "huaweicloud_aom_cloud_service_access" "test" {
  instance_id           = huaweicloud_aom_prom_instance.test.id
  service               = "DCS"
  tag_sync              = "auto"
  enterprise_project_id = local.enterprise_project_id

  provisioner "local-exec" {
    command = "sleep 240" # Waiting for the access center to complete the connection and generate indicators.
  }
}

resource "huaweicloud_aomv4_alarm_rule" "test" {
  depends_on = [huaweicloud_aom_cloud_service_access.test]

  name                  = var.alarm_rule_name
  type                  = "metric"
  enable                = true
  prom_instance_id      = huaweicloud_aom_prom_instance.test.id
  enterprise_project_id = local.enterprise_project_id

  alarm_notifications {
    notification_enable       = true
    notification_type         = "direct"
    bind_notification_rule_id = huaweicloud_aom_alarm_action_rule.test.id
    notify_resolved           = true
    notify_triggered          = true
    notify_frequency          = "0"
  }

  metric_alarm_spec {
    monitor_type = "all_metric"

    recovery_conditions {
      recovery_timeframe = 1
    }

    dynamic "trigger_conditions" {
      for_each = var.alarm_rule_trigger_conditions

      content {
        metric_query_mode       = "PROM"
        metric_name             = trigger_conditions.value.metric_name
        promql                  = trigger_conditions.value.promql
        promql_for              = trigger_conditions.value.promql_for
        aggregate_type          = trigger_conditions.value.aggregate_type
        aggregation_type        = trigger_conditions.value.aggregation_type
        aggregation_window      = trigger_conditions.value.aggregation_window
        metric_unit             = trigger_conditions.value.metric_unit
        metric_namespace        = trigger_conditions.value.metric_namespace
        operator                = trigger_conditions.value.operator
        metric_statistic_method = trigger_conditions.value.metric_statistic_method
        thresholds              = trigger_conditions.value.thresholds
        trigger_type            = trigger_conditions.value.trigger_type
        trigger_interval        = trigger_conditions.value.trigger_interval
        trigger_times           = trigger_conditions.value.trigger_times
        query_param             = trigger_conditions.value.query_param
        query_match             = trigger_conditions.value.query_match
      }
    }
  }

  lifecycle {
    ignore_changes = [
      metric_alarm_spec # If you want to update this configuration, please use a version higher than 1.82.3
    ]
  }
}

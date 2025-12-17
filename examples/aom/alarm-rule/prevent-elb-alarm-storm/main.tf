resource "huaweicloud_lts_group" "test" {
  group_name            = var.lts_group_name
  ttl_in_days           = 30
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

resource "huaweicloud_lts_stream" "test" {
  group_id              = huaweicloud_lts_group.test.id
  stream_name           = var.lts_stream_name
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

resource "huaweicloud_smn_topic" "test" {
  name                  = var.smn_topic_name
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
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

resource "huaweicloud_aom_alarm_group_rule" "test" {
  depends_on = [huaweicloud_aom_alarm_action_rule.test]

  name                  = var.alarm_group_rule_name
  group_by              = ["resource_provider"]
  group_interval        = var.alarm_group_rule_group_interval
  group_repeat_waiting  = var.alarm_group_rule_group_repeat_waiting
  group_wait            = var.alarm_group_rule_group_wait
  description           = var.alarm_group_rule_description != "" ? var.alarm_group_rule_description : null
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  detail {
    bind_notification_rule_ids = [huaweicloud_aom_alarm_action_rule.test.name]

    dynamic "match" {
      for_each = var.alarm_group_rule_condition_matching_rules

      content {
        key     = match.value.key
        operate = match.value.operate
        value   = match.value.value
      }
    }
  }
}

resource "huaweicloud_aomv4_alarm_rule" "test" {
  name             = var.alarm_rule_name
  type             = "metric"
  enable           = true
  prom_instance_id = var.prometheus_instance_id

  alarm_notifications {
    notification_enable = true
    notification_type   = "alarm_policy"
    route_group_enable  = true
    route_group_rule    = huaweicloud_aom_alarm_group_rule.test.name
    notify_resolved     = true
    notify_triggered    = true
    notify_frequency    = "-1"
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

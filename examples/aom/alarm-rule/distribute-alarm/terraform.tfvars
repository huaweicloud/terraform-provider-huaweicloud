dcs_instance_name                = "tf_test_aom_alarm_rule_distribute_alarm"
lts_group_name                   = "tf_test_aom_alarm_rule_distribute_alarm"
lts_stream_name                  = "tf_test_aom_alarm_rule_distribute_alarm"
smn_topic_name                   = "tf_test_aom_alarm_rule_distribute_alarm"
alarm_action_rule_name           = "tf_test_aom_alarm_rule_distribute_alarm_by_Ihn_tag"
alarm_action_rule_user_name      = "servicestage"
alarm_rule_matric_dimension_tags = {
  "Ihn" = "OPEN"
}

prometheus_instance_name      = "tf_test_aom_alarm_rule_distribute_alarm"
alarm_rule_name               = "tf_test_aom_alarm_rule_distribute_alarm_by_Ihn_tag"
alarm_rule_trigger_conditions = [
  {
    metric_name             = "huaweicloud_sys_dcs_cpu_usage"
    promql                  = "label_replace(avg_over_time(huaweicloud_sys_dcs_cpu_usage{Ihn=\"OPEN\"}[59999ms]),\"__name__\",\"huaweicloud_sys_dcs_cpu_usage\",\"\",\"\")"
    promql_for              = ""
    aggregate_type          = "by"
    aggregation_type        = "average"
    aggregation_window      = "1m"
    metric_unit             = "%"
    metric_query_mode       = "PROM"
    metric_namespace        = "SYS.DCS"
    operator                = ">"
    metric_statistic_method = "single"
    thresholds              = {
      "Critical" = 1
    }
    trigger_type            = "FIXED_RATE"
    trigger_interval        = "1m"
    trigger_times           = "3"
    query_param             = "{\"code\": \"a\", \"apmMetricReg\": []}"
    query_match             = "[{\"id\":\"first\",\"dimension\":\"Ihn\",\"conditionValue\":[{\"name\":\"OPEN\"}],\"conditionList\":[{\"name\":\"OPEN\"}],\"addMode\": \"first\",\"conditionCompare\":\"=\",\"regExpress\":null,\"dimensionSelected\":{\"label\":\"Ihn\",\"id\":\"Ihn\"}}]"
  }
]

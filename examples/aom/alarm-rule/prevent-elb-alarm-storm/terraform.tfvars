lts_group_name              = "tf_test_aom_prevent_elb_alarm_storm"
lts_stream_name             = "tf_test_aom_prevent_elb_alarm_storm"
smn_topic_name              = "tf_test_aom_prevent_elb_alarm_storm"
alarm_action_rule_user_name = "servicestage"
alarm_group_rule_name       = "tf_test_aom_prevent_elb_alarm_storm"
alarm_rule_name             = "tf_test_aom_prevent_elb_alarm_storm"

alarm_rule_trigger_conditions = [
  {
    metric_name             = "aom_metrics_total_per_hour"
    promql                  = "label_replace(avg_over_time(aom_metrics_total_per_hour{type=\"custom\"}[59999ms]),\"__name__\",\"aom_metrics_total_per_hour\",\"\",\"\")"
    promql_for              = "3m"
    aggregate_type          = "by"
    aggregation_type        = "average"
    aggregation_window      = "1m"
    metric_statistic_method = "single"
    thresholds              = {
      "Critical" = 1
    }
    trigger_type            = "FIXED_RATE"
    trigger_interval        = "1m"
    trigger_times           = "3"
    query_param             = "{\"code\": \"a\", \"apmMetricReg\": []}"
    query_match             = "{\"id\": \"first\", \"dimension\": \"type\", \"conditionValue\": [{\"name\": \"custom\"}], \"conditionList\": [{\"name\": \"custom\"}, {\"name\": \"basic\"}], \"addMode\": \"first\", \"conditionCompare\": \"=\", \"dimensionSelected\": {\"label\": \"type\", \"id\": \"type\"}}"
  }
]

# Prevent ELB Alarm Storm with AOM Alarm Group Rule

This example provides best practice code for using Terraform to prevent ELB alarm storm by configuring AOM alarm group
rules in HuaweiCloud. The example demonstrates how to use alarm grouping rules to reduce alarm noise and prevent alarm
storms when monitoring ELB business layer metrics.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* AOM service enabled in the target region
* LTS service enabled in the target region
* SMN service enabled in the target region

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `lts_group_name` - The name of the LTS group used to store SMN notification logs
* `lts_stream_name` - The name of the LTS stream used to store SMN notification logs
* `smn_topic_name` - The name of the SMN topic used to send notifications
* `alarm_action_rule_user_name` - The user name of the AOM alarm action rule
* `alarm_group_rule_name` - The name of the AOM alarm group rule
* `alarm_rule_name` - The name of the AOM alarm rule
* `prometheus_instance_id` - The ID of the Prometheus instance (default: "0", which represents the default
  Prometheus_AOM_Default instance)
* `alarm_rule_trigger_conditions` - The trigger conditions of the AOM alarm rule
  + `metric_name` - The name of the metric to monitor
  + `promql` - The PromQL query expression
  + `promql_for` - The duration for which the condition must be true
  + `aggregate_type` - The aggregation type (default: "by")
  + `aggregation_type` - The aggregation method (e.g., "average", "max", "min", "sum")
  + `aggregation_window` - The time window for aggregation
  + `metric_statistic_method` - The statistic method (e.g., "single")
  + `thresholds` - The alarm thresholds as a map (e.g., `{ "Critical" = 1 }`)
  + `trigger_type` - The trigger type (e.g., "FIXED_RATE")
  + `trigger_interval` - The interval between trigger checks
  + `trigger_times` - The number of consecutive times the condition must be met
  + `query_param` - The query parameters in JSON format
  + `query_match` - The query match conditions in JSON format

#### Optional Variables

* `alarm_action_rule_name` - The name of the AOM alarm action rule (default: "apm")
* `alarm_action_rule_type` - The type of the AOM alarm action rule (default: "1")
* `alarm_group_rule_group_interval` - The group interval of the alarm group rule in seconds (default: 60)
* `alarm_group_rule_group_repeat_waiting` - The group repeat waiting time in seconds (default: 3600)
* `alarm_group_rule_group_wait` - The group wait time in seconds (default: 15)
* `alarm_group_rule_description` - The description of the alarm group rule (default: "")
* `enterprise_project_id` - The ID of the enterprise project (default: "")
* `prometheus_instance_id` - The ID of the Prometheus instance (default: "0")
* `alarm_group_rule_condition_matching_rules` - The condition matching rules for the alarm group rule
  + `key` - The key of the matching condition (e.g., "event_severity", "resource_provider")
  + `operate` - The operation type (e.g., "EXIST", "EQUALS")
  + `value` - The list of values to match
  + Default: Filters for Critical and Major severity alarms from AOM

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  lts_group_name              = "tf_test_aom_prevent_elb_alarm_storm"
  lts_stream_name             = "tf_test_aom_prevent_elb_alarm_storm"
  smn_topic_name              = "tf_test_aom_prevent_elb_alarm_storm"
  alarm_action_rule_user_name = "your_ima_user_name"
  alarm_group_rule_name       = "tf_test_aom_prevent_elb_alarm_storm"
  alarm_rule_name             = "tf_test_aom_prevent_elb_alarm_storm"
  prometheus_instance_id      = "0"  # Optional, default is "0" (Prometheus_AOM_Default)

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
  ```

* Initialize Terraform:

  ```bash
  $ terraform init
  ```

* Review the Terraform plan:

  ```bash
  $ terraform plan
  ```

* Apply the configuration:

  ```bash
  $ terraform apply
  ```

* To clean up the resources:

  ```bash
  $ terraform destroy
  ```

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* The alarm group rule is dependent on the alarm action rule
* The alarm rule notification is automatically configured with alarm policy notification type and route group enabled
* The route group rule name in the alarm rule notification automatically matches the alarm group rule name
* The alarm group rule filters alarms by severity (Critical, Major) and source (AOM) by default
* Alarms are grouped by resource provider to merge similar alarms together
* The group wait time determines how long to wait before sending the first notification after creating a group
* The group interval determines how often to check for new alarms in a group
* The group repeat waiting time determines how long to wait before sending a repeat notification for the same group
* All resources will be created in the specified region
* The PromQL query must comply with Prometheus query syntax
* The thresholds is a map type with alarm levels as keys (Critical, Major, Minor, Info) and threshold values as values
* The default Prometheus instance ID is "0", which represents the Prometheus_AOM_Default instance
* The alarm rule uses `huaweicloud_aomv4_alarm_rule` resource with metric alarm type
* The `metric_alarm_spec` configuration is ignored in lifecycle to prevent unintended updates (requires provider version
  >= 1.82.3 for updates)

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.3.0 |
| huaweicloud | >= 1.80.4 |

# Distribute Alarms by HuaweiCloud Tags

This example provides best practice code for using Terraform to distribute alarms by HuaweiCloud tags (Tag) through
Prometheus monitoring and alarm management on HuaweiCloud. The example demonstrates how to monitor DCS instance CPU
utilization metrics and distribute alarms based on tags.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* DCS service enabled in the target region
* AOM service enabled in the target region
* LTS service enabled in the target region
* SMN service enabled in the target region
* TMS service enabled in the target region
* At least one DCS instance exists in the target region

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where resources will be created
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `lts_group_name` - The name of the LTS group used to store SMN notification logs
* `lts_stream_name` - The name of the LTS stream used to store SMN notification logs
* `smn_topic_name` - The name of the SMN topic used to send notifications
* `alarm_action_rule_name` - The name of the AOM alarm action rule used to send SMN notifications
* `alarm_action_rule_user_name` - The user name of the AOM alarm action rule
* `alarm_rule_matric_dimension_tags` - The custom tags to be added to the DCS instance for alarm distribution (map type,
  e.g. `{ "Ihn" = "OPEN" }`)
* `prometheus_instance_name` - The name of the Prometheus instance for cloud services
* `alarm_rule_name` - The name of the AOM alarm rule
* `alarm_rule_trigger_conditions` - The trigger conditions of the AOM alarm rule
  + `metric_name` - The name of the metric to monitor (e.g. "huaweicloud_sys_dcs_cpu_usage")
  + `promql` - The PromQL query expression (must include tag conditions,
    e.g. `huaweicloud_sys_dcs_cpu_usage{Ihn="OPEN"}`)
  + `aggregation_type` - The aggregation method (e.g., "average", "max", "min", "sum")
  + `aggregation_window` - The time window for aggregation (e.g., "1m")
  + `metric_unit` - The unit of the metric (e.g., "%")
  + `metric_query_mode` - The query mode (e.g., "PROM")
  + `metric_namespace` - The namespace of the metric (e.g., "SYS.DCS")
  + `operator` - The comparison operator (e.g., ">", "<", ">=", "<=", "=")
  + `metric_statistic_method` - The statistic method (e.g., "single")
  + `thresholds` - The alarm thresholds as a map (e.g., `{ "Critical" = 1 }`)
  + `trigger_type` - The trigger type (e.g., "FIXED_RATE")
  + `trigger_interval` - The interval between trigger checks (e.g., "1m")
  + `trigger_times` - The number of consecutive times the condition must be met (e.g., "3")
  + `query_param` - The query parameters in JSON format
  + `query_match` - The query match conditions in JSON format (must include tag matching conditions)

#### Optional Variables

* `dcs_instance_name` - The name of the existing DCS instance to be monitored (default: "", but required for the example
  to work properly)
* `alarm_action_rule_type` - The type of the AOM alarm action rule (default: "1" for notification)
* `promql_for` - The duration for which the condition must be true in trigger conditions (default: "")
* `aggregate_type` - The aggregation type in trigger conditions (default: "by")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  dcs_instance_name = "your_dcs_instance_name"
  lts_group_name    = "tf_test_aom_alarm_rule_distribute_alarm"
  lts_stream_name   = "tf_test_aom_alarm_rule_distribute_alarm"
  smn_topic_name    = "tf_test_aom_alarm_rule_distribute_alarm"
  
  alarm_action_rule_name      = "tf_test_aom_alarm_rule_distribute_alarm_by_Ihn_tag"
  alarm_action_rule_user_name = "your_iam_user_name"
  
  alarm_rule_matric_dimension_tags = {
    "Ihn" = "OPEN"
  }
  
  prometheus_instance_name = "tf_test_aom_alarm_rule_distribute_alarm"
  alarm_rule_name          = "tf_test_aom_alarm_rule_distribute_alarm_by_Ihn_tag"
  
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
  ```

* Initialize Terraform:

  ```bash
  terraform init
  ```

* Review the Terraform plan:

  ```bash
  terraform plan
  ```

* Apply the configuration:

  ```bash
  terraform apply
  ```

* To clean up the resources:

  ```bash
  terraform destroy
  ```

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* The DCS instance must exist before running this example
* Tags will be automatically added to the DCS instance through TMS (Tag Management Service)
* The Prometheus instance type is set to "CLOUD_SERVICE" to support cloud service monitoring
* The cloud service access configuration includes automatic tag synchronization (`tag_sync = "auto"`)
* After creating the cloud service access, the system waits 240 seconds for the access center to complete the connection
  and generate indicators
* The alarm rule uses direct notification type and binds to the alarm action rule
* The PromQL query expression must include tag conditions to filter metrics by tags
* The `query_match` parameter must include tag matching conditions in JSON format
* The `metric_alarm_spec` configuration is ignored in lifecycle to prevent unintended updates (requires provider version
  >= 1.82.3 for updates)
* All resources will be created in the specified region
* The enterprise project ID is automatically retrieved from the DCS instance if available
* The alarm rule will only trigger for DCS instances that match the specified tags

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.3.0 |
| huaweicloud | >= 1.80.4 |

---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_hpas"
description: |-
  Use this data source to get the list of CCI horizontal pod autoscaler within HuaweiCloud.
---

# huaweicloud_cciv2_hpas

Use this data source to get the list of CCI horizontal pod autoscaler within HuaweiCloud.

## Example Usage

```hcl
variable "namespace" {}

data "huaweicloud_cciv2_hpas" "test" {
  namespace = var.namespace
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `namespace` - (Required, String) Specifies the namespace.

* `name` - (Optional, String) Specifies the name of the CCI HPA.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `hpas` - The list of CCI horizontal pod autoscaler.
  The [hpas](#hpas) structure is documented below.

<a name="hpas"></a>
The `hpas` block supports:

* `behavior` - The behavior of the CCI HPA.
  The [behavior](#hpas_behavior) structure is documented below.

* `creation_timestamp` - The creation timestamp of the CCI HPA.

* `max_replicas` - The upper limit for the number of replicas
  to which the autoscaler can scale up.

* `metrics` - The metrics that can be used to calculate the desired replica count.
  The [metrics](#hpas_metrics) structure is documented below.

* `min_replicas` - The lower limit for the number of replicas
  to which the autoscaler can scale down.

* `name` - The name of the CCI HPA.

* `namespace` - The namespace.

* `resource_version` - The resource version of the CCI HPA.

* `scale_target_ref` - The scale target.
  The [scale_target_ref](#hpas_scale_target_ref) structure is documented below.

* `status` - The status.
  The [status](#hpas_status) structure is documented below.

* `uid` - The uid of the CCI HPA.

<a name="hpas_behavior"></a>
The `behavior` block supports:

* `scale_down` - The scale down of the behavior.
  The [scale_down](#hpas_behavior_scale) structure is documented below.

* `scale_up` - The scale up of the behavior.
  The [scale_up](#hpas_behavior_scale) structure is documented below.

<a name="hpas_behavior_scale"></a>
The `scale_down`, `scale_up` block supports:

* `policies` - The potential scaling policies which can be used during scaling.
  The [policies](#hpas_behavior_scale_policies) structure is documented below.

* `select_policy` - The select policy that should be used.

* `stabilization_window_seconds` - The seconds for which past recommendations should
  be considered while scaling up or scaling down.

<a name="hpas_behavior_scale_policies"></a>
The `policies` block supports:

* `period_seconds` - The window of time for which the policy should hold true.

* `type` - The type of the scaling policy.

* `value` - The value, it contains the amount of change which is permitted by the policy.

<a name="hpas_metrics"></a>
The `metrics` block supports:

* `container_resource` - The container resource metric source.
  The [container_resource](#hpas_metrics_container_resource) structure is documented below.

* `external` - The external metric resource.
  The [external](#hpas_metrics_external) structure is documented below.

* `object` - The object metric resource.
  The [object](#hpas_metrics_object) structure is documented below.

* `pods` - The pod metric resource.
  The [pods](#hpas_metrics_pods) structure is documented below.

* `resources` - The resource metric resource.
  The [resources](#hpas_metrics_resources) structure is documented below.

* `type` - The type.

<a name="hpas_metrics_container_resource"></a>
The `container_resource` block supports:

* `container` - The name of the container in the pods of the scaling target.

* `name` - The name of the resource in question.

* `target` - The target.
  The [target](#hpas_metrics_container_resource_target) structure is documented below.

<a name="hpas_metrics_container_resource_target"></a>
The `target` block supports:

* `average_utilization` - The target value of the resource metric across all elevant pods.

* `average_value` - The average value of the resource.

* `type` - The metric type, the value can be **Utilization**, **Value**, or **AverageValue**.

* `value` - The value of the resource.

<a name="hpas_metrics_external"></a>
The `external` block supports:

* `metric` - The metric.
  The [metric](#common_metric) structure is documented below.

* `target` - The target.
  The [target](#common_target) structure is documented below.

<a name="common_metric"></a>
The `metric` block supports:

* `name` - The name of the given metric.

* `selector` - The metric of external metric source.
  The [selector](#common_metric_selector) structure is documented below.

<a name="common_metric_selector"></a>
The `selector` block supports:

* `match_expressions` - The match expressions of the label selector requirements.
  The [match_expressions](#common_metric_selector_match_expressions) structure is documented below.

* `match_labels` - The match labels.

<a name="common_metric_selector_match_expressions"></a>
The `match_expressions` block supports:

* `key` - The label key.

* `operator` - The operator represents a key relationship to a set of values.

* `values` - The array of string values.

<a name="common_target"></a>
The `target` block supports:

* `average_utilization` - The target value of the resource metric across all elevant pods.

* `average_value` - The average value of the resource.

* `type` - The metric type, the value can be **Utilization**, **Value**, or **AverageValue**.

* `value` - The value of the resource.

<a name="hpas_metrics_object"></a>
The `object` block supports:

* `described_object` - The container resource metric source.
  The [described_object](#hpas_metrics_object_described_object) structure is documented below.

* `metric` - The metric.
  The [metric](#common_metric) structure is documented below.

* `target` - The target.
  The [target](#common_target) structure is documented below.

<a name="hpas_metrics_object_described_object"></a>
The `described_object` block supports:

* `api_version` - The API version of the referent.

* `kind` - The kind of the referent.

* `name` - The name of the referent.

<a name="hpas_metrics_pods"></a>
The `pods` block supports:

* `metric` - The metric.
  The [metric](#common_metric) structure is documented below.

* `target` - The target.
  The [target](#common_target) structure is documented below.

<a name="hpas_metrics_resources"></a>
The `resources` block supports:

* `name` - The name of the referent.

* `target` - The target of the referent.
  The [target](#common_target) structure is documented below.

<a name="hpas_scale_target_ref"></a>
The `scale_target_ref` block supports:

* `api_version` - The API version of the referent.

* `kind` - The kind of the referent.

* `name` - The name of the referent.

<a name="hpas_status"></a>
The `status` block supports:

* `conditions` - The conditions.
  The [conditions](#hpas_status_conditions) structure is documented below.

* `current_metrics` - The current metrics.
  The [current_metrics](#hpas_metrics) structure is documented below.

* `current_replicas` - The current replicas.

* `desired_replicas` - The desired replicas.

* `last_scale_time` - The last scale time.

* `observed_generation` - The observed generation.

<a name="hpas_status_conditions"></a>
The `conditions` block supports:

* `last_transition_time` - The last transition time of the conditions.

* `message` - The message of the conditions.

* `reason` - The reason of the conditions.

* `status` - Tthe status of the conditions.

* `type` - The type of the conditions.

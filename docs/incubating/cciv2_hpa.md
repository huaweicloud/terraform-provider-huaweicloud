---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_hpa"
description: |-
  Manages a CCI v2 horizontal pod autoscaler resource within HuaweiCloud.
---

# huaweicloud_cciv2_hpa

Manages a CCI v2 horizontal pod autoscaler resource within HuaweiCloud.

## Example Usage

```hcl
variable "namespace" {}
variable "name" {}

resource "huaweicloud_cciv2_hpa" "test" {
  name      = var.name
  namespace = var.namespace

  min_replicas = 1
  max_replicas = 5

  scale_target_ref {
    kind        = "Deployment"
    name        = "nginx"
    api_version = "cci/v2"
  }

  metrics {
    type = "Resource"

    resources {
      name = "memory"

      target {
        type                = "Utilization"
        average_utilization = 50
      }
    }
  }

  metrics {
    type = "Resource"

    resources {
      name = "cpu"

      target {
        type                = "Utilization"
        average_utilization = 50
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the name of the CCI HPA.

* `namespace` - (Required, String) Specifies the namespace.

* `behavior` - (Optional, List) Specifies the behavior of the CCI HPA.
  The [behavior](#behavior) structure is documented below.

* `max_replicas` - (Required, Int) Specifies the upper limit for the number of replicas
  to which the autoscaler can scale up.

* `metrics` - (Optional, List) Specifies the metrics that can be used to calculate the desired replica count.
  The [metrics](#metrics) structure is documented below.

* `min_replicas` - (Optional, Int) Specifies the lower limit for the number of replicas
  to which the autoscaler can scale down.

* `scale_target_ref` - (Required, List) Specifies the scale target.
  The [scale_target_ref](#scale_target_ref) structure is documented below.

<a name="behavior"></a>
The `behavior` block supports:

* `scale_down` - (Optional, List) Specifies the scale down of the behavior.
  The [scale_down](#behavior_scale) structure is documented below.

* `scale_up` - (Optional, List) Specifies the scale up of the behavior.
  The [scale_up](#behavior_scale) structure is documented below.

<a name="behavior_scale"></a>
The `scale_down`, `scale_up` block supports:

* `policies` - (Optional, List) Specifies the potential scaling policies which can be used during scaling.
  The [policies](#behavior_scale_policies) structure is documented below.

* `select_policy` - (Optional, String) Specifies select policy that should be used.

* `stabilization_window_seconds` - (Optional, Int) Specifies the seconds for which past recommendations should
  be considered while scaling up or scaling down.

<a name="behavior_scale_policies"></a>
The `policies` block supports:

* `period_seconds` - (Required, Int) Specifies the window of time for which the policy should hold true.

* `type` - (Required, String) Specifies the type of the scaling policy.

* `value` - (Required, Int) Specifies the value, it contains the amount of change which is permitted by the policy.

<a name="metrics"></a>
The `metrics`, `current_metrics` block supports:

* `container_resource` - (Optional, List) Specifies the container resource metric source.
  The [container_resource](#metrics_container_resource) structure is documented below.

* `external` - (Optional, List) Specifies the external metric resource.
  The [external](#metrics_external) structure is documented below.

* `object` - (Optional, List) Specifies the object metric resource.
  The [object](#metrics_object) structure is documented below.

* `pods` - (Optional, List) Specifies the pod metric resource.
  The [pods](#metrics_pods) structure is documented below.

* `resource` - (Optional, List) Specifies the resource metric resource.
  The [resource](#metrics_resource) structure is documented below.

* `type` - (Optional, String) Specifies the seconds for which past recommendations should
  be considered while scaling up or scaling down.

<a name="metrics_container_resource"></a>
The `container_resource` block supports:

* `container` - (Optional, String) Specifies the name of the container in the pods of the scaling target.

* `name` - (Optional, String) Specifies the name of the resource in question.

* `target` - (Optional, List) Specifies the target.
  The [target](#common_target) structure is documented below.

<a name="common_target"></a>
The `target` block supports:

* `average_utilization` - (Optional, Int) Specifies the target value of the resource metric across all elevant pods.

* `average_value` - (Optional, Map) Specifies the average value of the resource.

* `type` - (Optional, String) Specifies the metric type, the value can be **Utilization**, **Value**, or **AverageValue**.

* `value` - (Optional, Map) Specifies the value of the resource.

<a name="metrics_external"></a>
The `external` block supports:

* `metric` - (Required, List) Specifies the metric.
  The [metric](#common_metric) structure is documented below.

* `target` - (Required, List) Specifies the target.
  The [target](#common_target) structure is documented below.

<a name="common_metric"></a>
The `metric` block supports:

* `name` - (Required, String) Specifies the name of the given metric.

* `selector` - (Optional, List) Specifies the metric of external metric source.
  The [selector](#common_metric_selector) structure is documented below.

<a name="common_metric_selector"></a>
The `selector` block supports:

* `match_expressions` - (Optional, List) Specifies the match expressions of the label selector requirements.
  The [match_expressions](#common_metric_selector_match_expressions) structure is documented below.

* `match_labels` - (Optional, Map) Specifies the match labels.

<a name="common_metric_selector_match_expressions"></a>
The `match_expressions` block supports:

* `key` - (Optional, String) Specifies the label key that the selector applies to.

* `operator` - (Optional, String) Specifies the operator represents a key relationship to a set of values.

* `values` - (Optional, Map) Specifies the array of string values.

<a name="metrics_object"></a>
The `object` block supports:

* `described_object` - (Required, List) Specifies the container resource metric source.
  The [described_object](#metrics_object_described_object) structure is documented below.

* `metric` - (Required, List) Specifies the metric.
  The [metric](#common_metric) structure is documented below.

* `target` - (Required, List) Specifies the target.
  The [target](#common_target) structure is documented below.

<a name="metrics_object_described_object"></a>
The `described_object` block supports:

* `api_version` - (Optional, String) Specifies the API version of the referent.

* `kind` - (Optional, String) Specifies the kind of the referent.

* `name` - (Optional, String) Specifies the name of the referent.

<a name="metrics_pods"></a>
The `pods` block supports:

* `metric` - (Required, List) Specifies the container resource pod metric source.
  The [metric](#common_metric) structure is documented below.

* `target` - (Required, List) Specifies the target.
  The [target](#common_target) structure is documented below.

<a name="metrics_resource"></a>
The `resource` block supports:

* `name` - (Optional, String) Specifies the name of the resource in question.

* `target` - (Optional, List) Specifies the target.
  The [target](#common_target) structure is documented below.

<a name="scale_target_ref"></a>
The `scale_target_ref` block supports:

* `api_version` - (Optional, String) Specifies the API version of the referent.

* `kind` - (Optional, String) Specifies the kind of the referent.

* `name` - (Optional, String) Specifies the name of the referent.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `api_version` - The API version of the CCI HPA.

* `creation_timestamp` - The creation timestamp of the CCI HPA.

* `kind` - The kind of the CCI Image Snapshot.

* `resource_version` - The resource version of the CCI HPA.

* `status` - The status.
  The [status](#status) structure is documented below.

* `uid` - The uid of the CCI HPA.

<a name="status"></a>
The `status` block supports:

* `conditions` - The conditions.
  The [conditions](#status_conditions) structure is documented below.

* `current_metrics` - The current metrics.
  The [current_metrics](#metrics) structure is documented below.

* `current_replicas` - The current replicas.

* `desired_replicas` - The desired replicas.

* `last_scale_time` - The last scale time.

* `observed_generation` - The observed generation.

<a name="status_conditions"></a>
The `conditions` block supports:

* `last_transition_time` - The last transition time of the conditions.

* `message` - The message of the conditions.

* `reason` - The reason of the conditions.

* `status` - Tthe status of the conditions.

* `type` - The type of the conditions.

## Import

The CCI v2 horizontal pod autoscaler can be imported using `namespace` and `name`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_cciv2_hpa.test <namespace>/<name>
```

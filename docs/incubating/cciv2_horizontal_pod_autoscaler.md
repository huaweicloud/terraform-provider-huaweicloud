---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_horizontal_pod_autoscaler"
description: |-
  Manages a CCI v2 Horizontal Pod Autoscaler resource within HuaweiCloud.
---

# huaweicloud_cciv2_horizontal_pod_autoscaler

Manages a CCI v2 Horizontal Pod Autoscaler resource within HuaweiCloud.

## Example Usage

<!-- please add the usage of huaweicloud_cciv2_horizontal_pod_autoscaler -->
```hcl

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the name of the CCI Image Snapshot.

* `namespace` - (Required, String) Specifies the namespace.

* `annotations` - (Optional, Map) Specifies the annotations of the CCI Image Snapshot.

* `behavior` - (Optional, List) Specifies the behavior of the CCI horizontal pod autoscaler.
  The [behavior](#block--behavior) structure is documented below.

* `labels` - (Optional, Map) Specifies the annotations of the CCI Image Snapshot.

* `max_replicas` - (Optional, Int) Specifies the upper limit for the number of replicas
  to which the autoscaler can scale up.

* `metrics` - (Optional, List) Specifies the metrics that can be used to calculate the desired replica count.
  The [metrics](#block--metrics) structure is documented below.

* `min_replicas` - (Optional, Int) Specifies the lower limit for the number of replicas
  to which the autoscaler can scale down.

* `scale_target_ref` - (Optional, List) Specifies the scale target.
  The [scale_target_ref](#block--scale_target_ref) structure is documented below.

<a name="block--behavior"></a>
The `behavior` block supports:

* `scale_down` - (Optional, List) Specifies the scale down of the behavior.
  The [scale_down](#block--behavior--scale_down) structure is documented below.

* `scale_up` - (Optional, List) Specifies the scale up of the behavior.
  The [scale_up](#block--behavior--scale_up) structure is documented below.

<a name="block--behavior--scale_down"></a>
The `scale_down` block supports:

* `policies` - (Optional, List) Specifies the potential scaling policies which can be used during scaling.
  The [policies](#block--behavior--scale_down--policies) structure is documented below.

* `select_policy` - (Optional, String) Specifies select policy that should be used.

* `stabilization_window_seconds` - (Optional, Int) Specifies the seconds for which past recommendations should
  be considered while scaling up or scaling down.

<a name="block--behavior--scale_down--policies"></a>
The `policies` block supports:

* `period_seconds` - (Optional, Int) Specifies the window of time for which the policy should hold true.

* `type` - (Optional, String) Specifies the type of the scaling policy.

* `value` - (Optional, Int) Specifies the value, it contains the amount of change which is permitted by the policy.

<a name="block--behavior--scale_up"></a>
The `scale_up` block supports:

* `policies` - (Optional, List) Specifies the potential scaling policies which can be used during scaling.
  The [policies](#block--behavior--scale_up--policies) structure is documented below.

* `select_policy` - (Optional, String) Specifies select policy that should be used.

* `stabilization_window_seconds` - (Optional, Int) Specifies the seconds for which past recommendations should
  be considered while scaling up or scaling down.

<a name="block--behavior--scale_up--policies"></a>
The `policies` block supports:

* `period_seconds` - (Optional, Int) Specifies the window of time for which the policy should hold true.

* `type` - (Optional, String) Specifies the type of the scaling policy.

* `value` - (Optional, Int) Specifies the value, it contains the amount of change which is permitted by the policy.

<a name="block--metrics"></a>
The `metrics` block supports:

* `container_resource` - (Optional, List) Specifies the container resource metric source.
  The [container_resource](#block--metrics--container_resource) structure is documented below.

* `external` - (Optional, List) Specifies the external metric resource.
  The [external](#block--metrics--external) structure is documented below.

* `object` - (Optional, List) Specifies the object metric resource.
  The [object](#block--metrics--object) structure is documented below.

* `pods` - (Optional, List) Specifies the pod metric resource.
  The [pods](#block--metrics--pods) structure is documented below.

* `resource` - (Optional, List) Specifies the resource metric resource.
  The [resource](#block--metrics--resource) structure is documented below.

* `type` - (Optional, String) Specifies the seconds for which past recommendations should
  be considered while scaling up or scaling down.

<a name="block--metrics--container_resource"></a>
The `container_resource` block supports:

* `container` - (Optional, String) Specifies the name of the container in the pods of the scaling target.

* `name` - (Optional, String) Specifies the name of the resource in question.

* `target` - (Optional, List) Specifies the container resource metric source.
  The [target](#block--metrics--container_resource--target) structure is documented below.

<a name="block--metrics--container_resource--target"></a>
The `target` block supports:

* `average_utilization` - (Optional, Int) Specifies the target value of the resource metric across all elevant pods.

* `average_value` - (Optional, Map) Specifies the average value of the resource.

* `type` - (Optional, String) Specifies the metric type, the value can be **Utilization**, **Value**, or **AverageValue**.

* `value` - (Optional, Map) Specifies the value of the resource.

<a name="block--metrics--external"></a>
The `external` block supports:

* `metric` - (Optional, List) Specifies the metric of external metric source.
  The [metric](#block--metrics--external--metric) structure is documented below.

* `target` - (Optional, List) Specifies the target of external metric source.
  The [target](#block--metrics--external--target) structure is documented below.

<a name="block--metrics--external--metric"></a>
The `metric` block supports:

* `name` - (Optional, String) Specifies the name of the given metric.

* `selector` - (Optional, List) Specifies the metric of external metric source.
  The [selector](#block--metrics--external--metric--selector) structure is documented below.

<a name="block--metrics--external--metric--selector"></a>
The `selector` block supports:

* `match_expressions` - (Optional, List) Specifies the match expressions of the label selector requirements.
  The [match_expressions](#block--metrics--external--metric--selector--match_expressions) structure is documented below.

* `match_labels` - (Optional, Map) Specifies the match labels.

<a name="block--metrics--external--metric--selector--match_expressions"></a>
The `match_expressions` block supports:

* `key` - (Optional, String) Specifies the label key that the selector applies to.

* `operator` - (Optional, String) Specifies the operator represents a key relationship to a set of values.

* `values` - (Optional, Map) Specifies the array of string values.

<a name="block--metrics--external--target"></a>
The `target` block supports:

* `average_utilization` - (Optional, Int) Specifies the target value of the resource metric across all elevant pods.

* `average_value` - (Optional, Map) Specifies the average value of the resource.

* `type` - (Optional, String) Specifies the metric type, the value can be **Utilization**, **Value**, or **AverageValue**.

* `value` - (Optional, Map) Specifies the value of the resource.

<a name="block--metrics--object"></a>
The `object` block supports:

* `described_object` - (Optional, List) Specifies the container resource metric source.
  The [described_object](#block--metrics--object--described_object) structure is documented below.

* `metric` - (Optional, List) Specifies the container resource metric source.
  The [metric](#block--metrics--object--metric) structure is documented below.

* `target` - (Optional, List) Specifies the container resource metric source.
  The [target](#block--metrics--object--target) structure is documented below.

<a name="block--metrics--object--described_object"></a>
The `described_object` block supports:

* `api_version` - (Optional, String) Specifies the API version of the referent.

* `kind` - (Optional, String) Specifies the kind of the referent.

* `name` - (Optional, String) Specifies the name of the referent.

<a name="block--metrics--object--metric"></a>
The `metric` block supports:

* `name` - (Optional, String) Specifies the name of the given metric.

* `selector` - (Optional, String) Specifies the metric of external metric source.
  The [selector](#block--metrics--object--metric--selector) structure is documented below.

<a name="block--metrics--object--metric--selector"></a>
The `selector` block supports:

* `match_expressions` - (Optional, String) Specifies the match expressions of the label selector requirements.
  The [match_expressions](#block--metrics--object--metric--selector--match_expressions) structure is documented below.

* `match_labels` - (Optional, String) Specifies the match labels.

<a name="block--metrics--object--metric--selector--match_expressions"></a>
The `match_expressions` block supports:

* `key` - (Optional, String) Specifies the label key that the selector applies to.

* `operator` - (Optional, String) Specifies the operator represents a key relationship to a set of values.

* `values` - (Optional, String) Specifies the array of string values.

<a name="block--metrics--object--target"></a>
The `target` block supports:

* `average_utilization` - (Optional, Int) Specifies the target value of the resource metric across all elevant pods.

* `average_value` - (Optional, Map) Specifies the average value of the resource.

* `type` - (Optional, String) Specifies the metric type, the value can be **Utilization**, **Value**, or **AverageValue**.

* `value` - (Optional, Map) Specifies the value of the resource.

<a name="block--metrics--pods"></a>
The `pods` block supports:

* `metric` - (Optional, List) Specifies the container resource pod metric source.
  The [metric](#block--metrics--pods--metric) structure is documented below.

* `target` - (Optional, List) Specifies the container resource pod metric source.
  The [target](#block--metrics--pods--target) structure is documented below.

<a name="block--metrics--pods--metric"></a>
The `metric` block supports:

* `name` - (Optional, String) Specifies the name of the given metric.

* `selector` - (Optional, String) Specifies the metric of external metric source.
  The [selector](#block--metrics--pods--metric--selector) structure is documented below.

<a name="block--metrics--pods--metric--selector"></a>
The `selector` block supports:

* `match_expressions` - (Optional, String) Specifies the match expressions of the label selector requirements.
  The [match_expressions](#block--metrics--pods--metric--selector--match_expressions) structure is documented below.

* `match_labels` - (Optional, String) Specifies the match labels.

<a name="block--metrics--pods--metric--selector--match_expressions"></a>
The `match_expressions` block supports:

* `key` - (Optional, String) Specifies the label key that the selector applies to.

* `operator` - (Optional, String) Specifies the operator represents a key relationship to a set of values.

* `values` - (Optional, String) Specifies the array of string values.

<a name="block--metrics--pods--target"></a>
The `target` block supports:

* `average_utilization` - (Optional, Int) Specifies the target value of the resource metric across all elevant pods.

* `average_value` - (Optional, Map) Specifies the average value of the resource.

* `type` - (Optional, String) Specifies the metric type.
  The value can be **Utilization**, **Value**, or **AverageValue**.

* `value` - (Optional, Map) Specifies the value of the resource.

<a name="block--metrics--resource"></a>
The `resource` block supports:

* `name` - (Optional, String) Specifies the name of the resource in question.

* `target` - (Optional, List) Specifies the container resource pod metric source.
  The [target](#block--metrics--resource--target) structure is documented below.

<a name="block--metrics--resource--target"></a>
The `target` block supports:

* `average_utilization` - (Optional, Int) Specifies the target value of the resource metric across all elevant pods.

* `average_value` - (Optional, Map) Specifies the average value of the resource.

* `type` - (Optional, String) Specifies the metric type.
  The value can be **Utilization**, **Value**, or **AverageValue**.

* `value` - (Optional, Map) Specifies the value of the resource.

<a name="block--scale_target_ref"></a>
The `scale_target_ref` block supports:

* `api_version` - (Optional, String) Specifies the API version of the referent.

* `kind` - (Optional, String) Specifies the kind of the referent.

* `name` - (Optional, String) Specifies the name of the referent.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `api_version` - The API version of the CCI Image Snapshot.

* `creation_timestamp` - The creation timestamp of the CCI Image Snapshot.

* `kind` - The kind of the CCI Image Snapshot.

* `resource_version` - The resource version of the CCI Image Snapshot.

* `status` - The status.
  The [status](#attrblock--status) structure is documented below.

* `uid` - The uid of the CCI Image Snapshot.

<a name="attrblock--status"></a>
The `status` block supports:

* `conditions` - The status.
  The [conditions](#attrblock--status--conditions) structure is documented below.

* `current_metrics` - The status.
  The [current_metrics](#attrblock--status--current_metrics) structure is documented below.

* `current_replicas` - The current replicas.

* `desired_replicas` - The desired replicas.

* `last_scale_time` - The last scale time.

* `observed_generation` - The observed generation.

<a name="attrblock--status--conditions"></a>
The `conditions` block supports:

* `last_transition_time` - The last transition time of the conditions.

* `message` - The message of the conditions.

* `reason` - The reason of the conditions.

* `status` - Tthe status of the conditions.

* `type` - The type of the conditions.

<a name="attrblock--status--current_metrics"></a>
The `current_metrics` block supports:

## Import

The xxx can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_cciv2_horizontal_pod_autoscaler.test <id>
```

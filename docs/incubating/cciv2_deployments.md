---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_deployments"
description: |-
  Use this data source to get the list of CCI deployments within HuaweiCloud.
---

# huaweicloud_cciv2_deployments

Use this data source to get the list of CCI deployments within HuaweiCloud.

## Example Usage

```hcl
variable "namespace" {}

data "huaweicloud_cciv2_deployments" "test" {
  namespace = var.namespace
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `namespace` - (Required, String) Specifies the namespace.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `deployments` - The deployments.
  The [deployments](#deployments) structure is documented below.

<a name="deployments"></a>
The `deployments` block supports:

* `annotations` - The annotations of the CCI deployment.

* `creation_timestamp` - The creation time of the CCI deployment.

* `generation` - The generation of the CCI deployment.

* `min_ready_seconds` - The min ready seconds of the CCI deployment.

* `name` - The name of the CCI deployment.

* `namespace` - The namespace of the CCI deployment.

* `progress_deadline_seconds` - The progress deadline seconds of the CCI deployment.

* `replicas` - The replicas of the CCI deployment.

* `resource_version` - The resource version of the CCI deployment.

* `selector` - The selector of the CCI deployment.
  The [selector](#deployments_selector) structure is documented below.

* `status` - The status of the CCI deployment.
  The [status](#deployments_status) structure is documented below.

* `strategy` - The strategy of the CCI deployment.
  The [strategy](#deployments_strategy) structure is documented below.

* `template` - The template of the CCI deployment.
  The [template](#deployments_template) structure is documented below.

* `uid` - The uid of the CCI deployment.

<a name="deployments_selector"></a>
The `selector` block supports:

* `match_expressions` - The match expressions of the selector.
  The [match_expressions](#deployments_selector_match_expressions) structure is documented below.

* `match_labels` - The match labels of the selector.

<a name="deployments_selector_match_expressions"></a>
The `match_expressions` block supports:

* `key` - The key of the match expressions.

* `operator` - The operator of the match expressions.

* `values` - The values of the match expressions.

<a name="deployments_status"></a>
The `status` block supports:

* `conditions` - The conditions of the status.
  The [conditions](#deployments_status_conditions) structure is documented below.

* `observed_generation` - The observed generation of the status.

<a name="deployments_status_conditions"></a>
The `conditions` block supports:

* `last_transition_time` - The last transition time of the conditions.

* `last_update_time` - The last update time of the conditions.

* `message` - The message of the conditions.

* `reason` - The reason of the conditions.

* `status` - The status of the conditions.

* `type` - The type of the conditions.

<a name="deployments_strategy"></a>
The `strategy` block supports:

* `rolling_update` - The rolling update of the strategy.

* `type` - The type of the strategy.

<a name="deployments_template"></a>
The `template` block supports:

* `metadata` - The metadata of the template.
  The [metadata](#deployments_template_metadata) structure is documented below.

* `spec` - The spec of the template.
  The [spec](#deployments_template_spec) structure is documented below.

<a name="deployments_template_metadata"></a>
The `metadata` block supports:

* `annotations` - The annotations of the metadata.

* `labels` - The labels of the metadata.

<a name="deployments_template_spec"></a>
The `spec` block supports:

* `active_deadline_seconds` - The active deadline seconds of the spec.

* `affinity` - The affinity of the spec.
  The [affinity](#deployments_template_spec_affinity) structure is documented below.

* `containers` - The containers of the spec.
  The [containers](#deployments_template_spec_containers) structure is documented below.

* `dns_policy` - The DNS policy of the spec.

* `hostname` - The host name of the spec.

* `image_pull_secrets` - The image pull secrets of the spec.
  The [image_pull_secrets](#deployments_template_spec_image_pull_secrets) structure is documented below.

* `node_name` - The node name of the spec.

* `overhead` - The overhead of the spec.

* `restart_policy` - The restart policy of the spec.

* `scheduler_name` - The scheduler name of the spec.

* `set_hostname_as_pqdn` - The set host name as PQDN name of the spec.

* `share_process_namespace` - The share process namespace of the spec.

* `termination_grace_period_seconds` - The termination grace period seconds of the spec.

<a name="deployments_template_spec_affinity"></a>
The `termination_grace_period_seconds` block supports:

* `node_affinity` - The node affinity.
  The [node_affinity](#node_affinity) structure is documented below.

* `pod_anti_affinity` - The pod anti affinity.
  The [pod_anti_affinity](#pod_anti_affinity) structure is documented below.

<a name="node_affinity"></a>
The `node_affinity` block supports:

* `required_during_scheduling_ignored_during_execution` - The required during scheduling ignored during execution.
  The [required_during_scheduling_ignored_during_execution](#node_affinity_required) structure is documented below.

<a name="node_affinity_required"></a>
The `required_during_scheduling_ignored_during_execution` block supports:

* `node_selector_terms` - The node selector terms.
  The [node_selector_terms](#node_affinity_required_node_selector_terms) structure is documented below.

<a name="node_affinity_required_node_selector_terms"></a>
The `node_selector_terms` block supports:

* `match_expressions` - The match expressions of the node selector terms.
  The [match_expressions](#node_affinity_required_node_selector_terms_match_expressions) structure is documented below.

<a name="node_affinity_required_node_selector_terms_match_expressions"></a>
The `match_expressions` block supports:

* `key` - The key of the match expressions.

* `operator` - The operator of the match expressions.

* `values` - The values of the match expressions.

<a name="pod_anti_affinity"></a>
The `pod_anti_affinity` block supports:

* `preferred_during_scheduling_ignored_during_execution` - The preferred during scheduling ignored during execution.
  The [preferred_during_scheduling_ignored_during_execution](#pod_anti_affinity_preferred) structure is documented below.

* `required_during_scheduling_ignored_during_execution` - The required during scheduling ignored during execution.
  The [required_during_scheduling_ignored_during_execution](#pod_anti_affinity_required) structure is documented below.

<a name="pod_anti_affinity_preferred"></a>
The `required_during_scheduling_ignored_during_execution` block supports:

* `pod_affinity_term` - The pod affinity term.
  The [pod_affinity_term](#pod_anti_affinity_required_pod_affinity_term) structure is documented below.

* `weight` - The weight.

<a name="pod_anti_affinity_required_pod_affinity_term"></a>
The `weight` block supports:

<a name="pod_anti_affinity_required"></a>
The `required_during_scheduling_ignored_during_execution` block supports:

* `label_selector` - The label selector.
  The [label_selector](#pod_anti_affinity_required_label_selector) structure is documented below.

* `namespaces` - The namespaces.

* `topology_key` - The topology key.

<a name="pod_anti_affinity_required_label_selector"></a>
The `topology_key` block supports:

<a name="deployments_template_spec_containers"></a>
The `termination_grace_period_seconds` block supports:

* `env` - The env of the termination grace period seconds.
  The [env](#deployments_template_spec_termination_grace_period_seconds_env) structure is documented below.

* `image` - The image of the termination grace period seconds.

* `name` - The name of the termination grace period seconds.

* `resources` - The resources of the termination grace period seconds.
  The [resources](#deployments_template_spec_termination_grace_period_seconds_resources) structure is documented below.

<a name="deployments_template_spec_termination_grace_period_seconds_env"></a>
The `env` block supports:

* `name` - The name of the environment.

* `value` - The value of the environment.

<a name="deployments_template_spec_termination_grace_period_seconds_resources"></a>
The `resources` block supports:

* `limits` - The limits of the resources.

* `requests` - The requests of the resources.

<a name="deployments_template_spec_image_pull_secrets"></a>
The `termination_grace_period_seconds` block supports:

* `name` - The name of the termination grace period seconds.

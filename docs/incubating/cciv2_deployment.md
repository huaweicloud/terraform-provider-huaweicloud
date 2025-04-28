---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_deployment"
description: |-
  Manages a CCI v2 deployment resource within HuaweiCloud.
---
# huaweicloud_cciv2_deployment

Manages a CCI v2 deployment resource within HuaweiCloud.

## Example Usage

```hcl
variable "namespace" {}
variable "name" {}

resource "huaweicloud_cciv2_deployment" "test" {
  namespace = var.namespace
  name      = var.name

  selector {
    match_labels = {
      app = "template1"
    }
  }

  template {
    metadata {
      labels = {
        app = "template1"
      }

      annotations = {
        "resource.cci.io/instance-type" = "general-computing"
      }
    }

    spec {
      containers {
        name  = "c1"
        image = "alpine:latest"

        resources {
          limits = {
            cpu    = "1"
            memory = "2G"
          }

          requests = {
            cpu    = "1"
            memory = "2G"
          }
        }
      }

      image_pull_secrets {
        name = "imagepull-secret"
      }
    }
  }

  strategy {
    type = "RollingUpdate"

    rolling_update = {
      maxUnavailable = "25%"
      maxSurge       = "100%"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the name of the CCI deployment.

* `namespace` - (Required, String, NonUpdatable) Specifies the namespace.

* `selector` - (Required, List) Specifies the selector of the CCI deployment.
  The [selector](#selector) structure is documented below.

* `template` - (Required, List) Specifies the template of the CCI deployment.
  The [template](#template) structure is documented below.

* `min_ready_seconds` - (Optional, Int) Specifies the min ready seconds of the CCI deployment.

* `progress_deadline_seconds` - (Optional, Int) Specifies the progress deadline seconds of the CCI deployment.

* `replicas` - (Optional, Int) Specifies the replicas of the CCI deployment.

* `strategy` - (Optional, List) Specifies the strategy of the CCI deployment.
  The [strategy](#strategy) structure is documented below.

<a name="selector"></a>
The `selector` block supports:

* `match_expressions` - (Optional, List) Specifies the match expressions of the CCI deployment selector.
  The [match_expressions](#match_expressions) structure is documented below.

* `match_labels` - (Optional, Map) Specifies the match labels of the CCI deployment selector.

<a name="template"></a>
The `template` block supports:

* `metadata` - (Optional, List) Specifies the metadata of the CCI deployment template.
  The [metadata](#template_metadata) structure is documented below.

* `spec` - (Optional, List) Specifies the spec of the CCI deployment template.
  The [spec](#template_spec) structure is documented below.

<a name="template_metadata"></a>
The `metadata` block supports:

* `annotations` - (Optional, Map) Specifies the annotations.

* `labels` - (Optional, Map) Specifies the labels.

<a name="template_spec"></a>
The `spec` block supports:

* `containers` - (Required, List) Specifies the containers of the CCI deployment.
  The [containers](#template_spec_containers) structure is documented below.

* `active_deadline_seconds` - (Optional, Int) Specifies the active deadline seconds.

* `affinity` - (Optional, List) Specifies the affinity.
  The [affinity](#template_spec_affinity) structure is documented below.

* `dns_policy` - (Optional, String) Specifies the DNS policy.

* `hostname` - (Optional, String) Specifies the host name.

* `image_pull_secrets` - (Optional, List) Specifies the image pull secrets.
  The [image_pull_secrets](#image_pull_secrets) structure is documented below.

* `node_name` - (Optional, String) Specifies the node name.

* `overhead` - (Optional, Map) Specifies the overhead.

* `restart_policy` - (Optional, String) Specifies the restart policy.

* `scheduler_name` - (Optional, String) Specifies the scheduler name.

* `set_hostname_as_pqdn` - (Optional, Bool) Specifies whether to set hostname as PQDN.

* `share_process_namespace` - (Optional, Bool) Specifies whether to share process namespace.

* `termination_grace_period_seconds` - (Optional, Int) Specifies the period seconds of termination grace.

<a name="template_spec_containers"></a>
The `containers` block supports:

* `name` - (Required, String) Specifies the name of the container.

* `env` - (Optional, List) Specifies the environment of the container.
  The [env](#template_spec_containers_env) structure is documented below.

* `image` - (Optional, String) Specifies the image of the container.

* `resources` - (Optional, List) Specifies the resources of the container.
  The [resources](#template_spec_containers_resources) structure is documented below.

<a name="template_spec_containers_env"></a>
The `env` block supports:

* `name` - (Optional, String) Specifies the name of the environment.

* `value` - (Optional, String) Specifies the value of the environment.

<a name="template_spec_containers_resources"></a>
The `resources` block supports:

* `limits` - (Optional, Map) Specifies the limits of the resources.

* `requests` - (Optional, Map) Specifies the requests of the resources.

<a name="template_spec_affinity"></a>
The `affinity` block supports:

* `node_affinity` - (Optional, List) Specifies the node affinity.
  The [node_affinity](#affinity_node_affinity) structure is documented below.

* `pod_anti_affinity` - (Optional, List) Specifies the pod anti affinity.
  The [pod_anti_affinity](#affinity_pod_anti_affinity) structure is documented below.

<a name="affinity_node_affinity"></a>
The `pod_anti_affinity` block supports:

* `required_during_scheduling_ignored_during_execution` - (Optional, List) Specifies the required during scheduling
  the ignored during execution.
  The [required_during_scheduling_ignored_during_execution](#node_affinity_required) structure is documented below.

<a name="node_affinity_required"></a>
The `required_during_scheduling_ignored_during_execution` block supports:

* `node_selector_terms` - (Required, List) Specifies the node selector terms.
  The [node_selector_terms](#node_selector_terms) structure is documented below.

<a name="node_selector_terms"></a>
The `node_selector_terms` block supports:

* `match_expressions` - (Optional, List) Specifies the match expressions.
  The [match_expressions](#match_expressions) structure is documented below.

<a name="affinity_pod_anti_affinity"></a>
The `pod_anti_affinity` block supports:

* `preferred_during_scheduling_ignored_during_execution` - (Optional, List) Specifies the preferred during scheduling
  ignored during execution.
  The [preferred_during_scheduling_ignored_during_execution](#pod_anti_affinity_preferred) structure is documented below.

* `required_during_scheduling_ignored_during_execution` - (Optional, List) Specifies the required during schedulin
  ignored during execution.
  The [required_during_scheduling_ignored_during_execution](#pod_anti_affinity_required) structure is documented below.

<a name="pod_anti_affinity_preferred"></a>
The `preferred_during_scheduling_ignored_during_execution` block supports:

* `pod_affinity_term` - (Required, List) Specifies the pod affinity term.
  The [pod_affinity_term](#pod_affinity_term) structure is documented below.

* `weight` - (Required, Int) Specifies the weight.

<a name="pod_affinity_term"></a>
The `weight` block supports:

* `topology_key` - (Required, String) Specifies the topology key.

* `label_selector` - (Optional, List) Specifies the label selector.
  The [label_selector](#label_selector) structure is documented below.

* `namespaces` - (Optional, List) Specifies the namespaces.

<a name="pod_anti_affinity_required"></a>
The `required_during_scheduling_ignored_during_execution` block supports:

* `topology_key` - (Required, String) Specifies the topology key.

* `label_selector` - (Optional, List) Specifies the label selector.
  The [label_selector](#label_selector) structure is documented below.

* `namespaces` - (Optional, List) Specifies the namespaces.

<a name="label_selector"></a>
The `label_selector` block supports:

* `match_expressions` - (Optional, List) Specifies the match expressions of the CCI deployment selector.
  The [match_expressions](#match_expressions) structure is documented below.

* `match_labels` - (Optional, Map) Specifies the match labels of the CCI deployment selector.

<a name="match_expressions"></a>
The `match_expressions` block supports:

* `key` - (Required, String) Specifies the key of the match expressions.

* `operator` - (Required, String) Specifies the operator of the match expressions.

* `values` - (Optional, List) Specifies the values of the match expressions.

<a name="image_pull_secrets"></a>
The `image_pull_secrets` block supports:

* `name` - (Optional, String) Specifies the name of image pull secrets.

<a name="strategy"></a>
The `strategy` block supports:

* `rolling_update` - (Optional, Map) Specifies the rolling update config of the CCI deployment strategy.

* `type` - (Optional, String) Specifies the type of the CCI deployment strategy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `annotations` - The metadata annotations of the CCI deployment.

* `api_version` - The API version of the CCI deployment.

* `creation_timestamp` - The creation timestamp of the CCI deployment.

* `generation` - The generation of the CCI deployment.

* `kind` - The kind of the CCI deployment.

* `resource_version` - The resource version of the CCI deployment.

* `status` - The status of the CCI deployment.
  The [status](#status) structure is documented below.

* `uid` - The uid of the CCI deployment.

<a name="status"></a>
The `status` block supports:

* `conditions` - The conditions of the CCI deployment.
  The [conditions](#status_conditions) structure is documented below.

* `observed_generation` - The observed generation of the CCI deployment.

<a name="status_conditions"></a>
The `conditions` block supports:

* `last_transition_time` - The last transition time of the CCI deployment conditions.

* `last_update_time` - The last update time of the CCI deployment conditions.

* `message` - The message of the CCI deployment conditions.

* `reason` - The reason of the CCI deployment conditions.

* `status` - Tthe status of the CCI deployment conditions.

* `type` - The type of the CCI deployment conditions.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 10 minutes.

## Import

The CCI v2 deployment can be imported using `namespace` and `name`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_cciv2_deployment.test <namespace>/<name>
```

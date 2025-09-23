---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_replica_sets"
description: |-
  Use this data source to get the list of CCI replica sets within HuaweiCloud.
---

# huaweicloud_cciv2_replica_sets

Use this data source to get the list of CCI replica sets within HuaweiCloud.

## Example Usage

```hcl
variable "namespace" {}

data "huaweicloud_cciv2_replica_sets" "test" {
  namespace = var.namespace
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `namespace` - (Required, String) Specifies the namespace.

* `name` - (Optional, String) Specifies the name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `replica_sets` - The replica sets.
  The [replica_sets](#replica_sets) structure is documented below.

<a name="replica_sets"></a>
The `replica_sets` block supports:

* `annotations` - The annotations.

* `api_version` - The api version.

* `creation_timestamp` - The creation time.

* `kind` - The kind.

* `labels` - The labels.

* `min_ready_seconds` - The min ready seconds.

* `name` - The name.

* `replicas` - The replicas.

* `resource_version` - The resource version.

* `selector` - The selector.
  The [selector](#replica_sets_selector) structure is documented below.

* `status` - The status.
  The [status](#replica_sets_status) structure is documented below.

* `template` - The template.
  The [template](#replica_sets_template) structure is documented below.

* `uid` - The uid.

<a name="replica_sets_selector"></a>
The `selector` block supports:

* `match_expressions` - The match expressions.
  The [match_expressions](#replica_sets_selector_match_expressions) structure is documented below.

* `match_labels` - The match labels.

<a name="replica_sets_selector_match_expressions"></a>
The `match_expressions` block supports:

* `key` - The key.

* `operator` - The operator.

* `values` - The values.

<a name="replica_sets_status"></a>
The `status` block supports:

* `available_replicas` - The available replicas.

* `conditions` - The conditions.
  The [conditions](#replica_sets_status_conditions) structure is documented below.

* `fully_labeled_replicas` - The fully labeled replicas.

* `observed_generation` - The observed generation.

* `ready_replicas` - The ready replicas.

* `replicas` - The replicas.

<a name="replica_sets_status_conditions"></a>
The `conditions` block supports:

* `last_transition_time` - The last transition time.

* `message` - The message.

* `reason` - The reason.

* `status` - The status.

* `type` - The type.

<a name="replica_sets_template"></a>
The `template` block supports:

* `metadata` - The metadata.
  The [metadata](#replica_sets_template_metadata) structure is documented below.

* `spec` - The spec.
  The [spec](#replica_sets_template_spec) structure is documented below.

<a name="replica_sets_template_metadata"></a>
The `metadata` block supports:

* `annotations` - The annotations.

* `labels` - The labels.

<a name="replica_sets_template_spec"></a>
The `spec` block supports:

* `active_deadline_seconds` - The active deadline seconds.

* `affinity` - The affinity.
  The [affinity](#spec_affinity) structure is documented below.

* `containers` - The containers.
  The [containers](#spec_containers) structure is documented below.

* `dns_policy` - The DNS policy.

* `hostname` - The hostname.

* `image_pull_secrets` - The image pull secrets.
  The [image_pull_secrets](#spec_image_pull_secrets) structure is documented below.

* `node_name` - The node name.

* `overhead` - The overhead.

* `restart_policy` - The restart policy.

* `scheduler_name` - The scheduler name.

* `set_hostname_as_pqdn` - The set host name as PQDN name of the spec.

* `share_process_namespace` - The share process namespace of the spec.

* `termination_grace_period_seconds` - The termination grace period seconds of the spec.

<a name="spec_affinity"></a>
The `affinity` block supports:

* `node_affinity` - The node affinity.
  The [node_affinity](#spec_affinity_node_affinity) structure is documented below.

* `pod_anti_affinity` - The pod anti affinity.
  The [pod_anti_affinity](#spec_affinity_pod_anti_affinity) structure is documented below.

<a name="spec_affinity_node_affinity"></a>
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

<a name="spec_affinity_pod_anti_affinity"></a>
The `pod_anti_affinity` block supports:

* `preferred_during_scheduling_ignored_during_execution` - The preferred during scheduling ignored during execution.
  The [preferred_during_scheduling_ignored_during_execution](#pod_anti_affinity_preferred) structure is documented below.

* `required_during_scheduling_ignored_during_execution` - The required during scheduling ignored during execution.
  The [required_during_scheduling_ignored_during_execution](#pod_anti_affinity_preferred) structure is documented below.

<a name="pod_anti_affinity_preferred"></a>
The `required_during_scheduling_ignored_during_execution` block supports:

* `pod_affinity_term` - The pod affinity term.
  The [pod_affinity_term](#pod_anti_affinity_required_pod_affinity_term) structure is documented below.

* `weight` - The weight.

<a name="pod_anti_affinity_required_pod_affinity_term"></a>
The `required_during_scheduling_ignored_during_execution` block supports:

* `label_selector` - The label selector.
  The [label_selector](#pod_anti_affinity_required_label_selector) structure is documented below.

* `namespaces` - The namespaces.

* `topology_key` - The topology key.

<a name="pod_anti_affinity_required_label_selector"></a>
The `label_selector` block supports:

* `match_expressions` - The match expressions.
  The [match_expressions](#pod_anti_affinity_required_label_selector_match_expressions) structure is documented below.

* `match_labels` - The match labels.

<a name="pod_anti_affinity_required_label_selector_match_expressions"></a>
The `match_expressions` block supports:

* `key` - The key of the match expressions.

* `operator` - The operator of the match expressions.

* `values` - The values of the match expressions.

<a name="spec_containers"></a>
The `containers` block supports:

* `args` - The args of the containers.

* `command` - The command of the containers.

* `env` - The environment of the containers.
  The [env](#spec_containers_env) structure is documented below.

* `env_from` - The environment source of the containers.
  The [env_from](#spec_containers_env_from) structure is documented below.

* `image` -  The image of the containers.

* `lifecycle` - The lifecycle of the containers.
  The [lifecycle](#spec_containers_lifecycle) structure is documented below.

* `liveness_probe` - The liveness probe of the containers.
  The [liveness_probe](#spec_containers_probe) structure is documented below.

* `name` - The name of the containers.

* `ports` - The ports of the containers.
  The [ports](#spec_containers_ports) structure is documented below.

* `readiness_probe` - The readiness probe of the containers.
  The [readiness_probe](#spec_containers_probe) structure is documented below.

* `resources` - The resources of the containers.
  The [resources](#spec_containers_resources) structure is documented below.

* `security_context` - The security context of the containers.
  The [security_context](#spec_containers_security_context) structure is documented below.

* `startup_probe` - The startup probe of the containers.
  The [startup_probe](#spec_containers_probe) structure is documented below.

* `stdin` - The stdin of the containers.

* `stdin_once` - The stdin once of the containers.

* `termination_message_path` - The termination message path of the containers.

* `termination_message_policy` - The termination message policy of the containers.

* `tty` - The TTY of the containers.

* `working_dir` - The working dir of the containers.

<a name="spec_containers_env"></a>
The `env` block supports:

* `name` - The name of the environment.

* `value` - The value of the environment.

<a name="spec_containers_env_from"></a>
The `env_from` block supports:

* `config_map_ref` - The config map ref of the environment source.
  The [config_map_ref](#spec_containers_env_from_ref) structure is documented below.

* `prefix` - The prefix of the environment source.

* `secret_ref` - The secret ref of the environment source.
  The [secret_ref](#spec_containers_env_from_ref) structure is documented below.

<a name="spec_containers_env_from_ref"></a>
The `config_map_ref`, `secret_ref` block supports:

* `name` - The name.

* `optional` - The optional.

<a name="spec_containers_lifecycle"></a>
The `lifecycle` block supports:

* `post_start` - The post start of the lifecycle.
  The [post_start](#spec_containers_lifecycle_post_start_or_pre_stop) structure is documented below.

* `pre_stop` - The pre stop of the lifecycle.
  The [pre_stop](#spec_containers_lifecycle_post_start_or_pre_stop) structure is documented below.

<a name="spec_containers_lifecycle_post_start_or_pre_stop"></a>
The `post_start`, `pre_stop` block supports:

* `exec` - The exec.
  The [exec](#spec_containers_lifecycle_post_start_or_pre_stop_exec) structure is documented below.

* `http_get` - The http get.
  The [http_get](#spec_containers_lifecycle_post_start_or_pre_stop_http_get) structure is documented below.

<a name="spec_containers_lifecycle_post_start_or_pre_stop_exec"></a>
The `exec` block supports:

* `command` - The command of the exec.

<a name="spec_containers_lifecycle_post_start_or_pre_stop_http_get"></a>
The `http_get` block supports:

* `host` - The host of the http get.

* `http_headers` - The http headers.
  The [http_headers](#spec_containers_lifecycle_post_start_or_pre_stop_http_get_http_headers)
  structure is documented below.

* `path` - The path of the http get.

* `port` - The port of the http get.

* `scheme` - The scheme of the http get.

<a name="spec_containers_lifecycle_post_start_or_pre_stop_http_get_http_headers"></a>
The `http_headers` block supports:

* `name` - The name of the scheme.

* `value` - The value of the scheme.

<a name="spec_containers_probe"></a>
The `liveness_probe`, `readiness_probe`, `startup_probe` block supports:

* `exec` - The exec.
  The [exec](#spec_containers_probe_exec) structure is documented below.

* `failure_threshold` - The minimum consecutive failures for the probe to be
  considered failed after having succeeded.

* `http_get` - The http get.
  The [http_get](#spec_containers_probe_http_get) structure is documented below.

* `initial_delay_seconds` - The number of seconds after the container has started
  before liveness probes are initialed.

* `period_seconds` - How often to perform the probe.

* `success_threshold` - The success threshold.

* `termination_grace_period_seconds` - The termination grace period seconds.

<a name="spec_containers_probe_exec"></a>
The `exec` block supports:

* `command` - The command line to execute inside the container.

<a name="spec_containers_probe_http_get"></a>
The `http_get` block supports:

* `host` - The host name.

* `http_headers` - The custom headers to set in the request.
  The [http_headers](#spec_containers_probe_termination_grace_period_seconds_http_headers) structure is
  documented below.

* `path` - The path to access on the http server.

* `port` - The port to access on the http server.

* `scheme` - The scheme to use for connecting to the host.

<a name="spec_containers_probe_termination_grace_period_seconds_http_headers"></a>
The `http_headers` block supports:

* `name` - The name of the custom http headers.

* `value` - The value of the custom http headers.

<a name="spec_containers_ports"></a>
The `ports` block supports:

* `container_port` - The number of port to expose on the IP address of pod.

* `name` - The port name of the container.

* `protocol` - The protocol for container port.

<a name="spec_containers_resources"></a>
The `resources` block supports:

* `limits` - The limits of resource.

* `requests` - The requests of the resource.

<a name="spec_containers_security_context"></a>
The `security_context` block supports:

* `capabilities` - The capabilities of the security context.
  The [capabilities](#spec_containers_security_context_capabilities) structure is documented below.

* `proc_mount` - The proc mount to use for the containers.

* `read_only_root_file_system` - Whether this container has a read-only root file system.

* `run_as_group` - The GID TO run the entrypoint of the container process.

* `run_as_non_root` - The container must run as a non-root user.

* `run_as_user` - The UID to run the entrypoint of the container process.

<a name="spec_containers_security_context_capabilities"></a>
The `capabilities` block supports:

* `add` - The add of the capabilities.

* `drop` - The drop of the capabilities.

<a name="spec_image_pull_secrets"></a>
The `image_pull_secrets` block supports:

* `name` - The name of the image pull secrets.

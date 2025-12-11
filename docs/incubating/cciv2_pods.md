---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_pods"
description: |-
  Use this data source to get the list of CCI pods within HuaweiCloud.
---

# huaweicloud_cciv2_pods

Use this data source to get the list of CCI pods within HuaweiCloud.

## Example Usage

```hcl
variable "namespace" {}

data "huaweicloud_cciv2_pods" "test" {
  namespace = var.namespace
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `namespace` - (Required, String) Specifies the namespace of the CCI pods.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `pods` - The CCI pods.
  The [pods](#pods) structure is documented below.

<a name="pods"></a>
The `pods` block supports:

* `namespace` - The namespace of the pod.

* `name` - The name of the pod.

* `annotations` - The annotations of the pod.

* `labels` - The labels of the pod.

* `creation_timestamp` - The creation time of the pod.

* `resource_version` - The resource version of the pod.

* `uid` - The uid of the pod.

* `active_deadline_seconds` - The active deadline seconds the pod.

* `affinity` - The affinity of the pod.
  The [affinity](#pods_affinity) structure is documented below.

* `containers` - The container of the pod.
  The [containers](#pods_containers) structure is documented below.

* `dns_config` - The DNS config of the pod.
  The [dns_config](#pods_dns_config) structure is documented below.

* `dns_policy` - The DNS policy of the pod.

* `ephemeral_containers` - The ephemeral container of the pod.
  The [ephemeral_containers](#pods_containers) structure is documented below.

* `finalizers` - The finalizers of the pod.

* `host_aliases` - The host aliases of the pod.
  The [host_aliases](#pods_host_aliases) structure is documented below.

* `hostname` - The host name of the pod.

* `image_pull_secrets` - The image pull secrets of the pod.
  The [image_pull_secrets](#pods_image_pull_secrets) structure is documented below.

* `init_containers` - The init container of the pod.
  The [init_containers](#pods_containers) structure is documented below.

* `node_name` - The node name of the pod.

* `overhead` - The overhead of the pod.

* `readiness_gates` - The readiness gates of the pod.
  The [readiness_gates](#pods_readiness_gates) structure is documented below.

* `restart_policy` - The restart policy for all containers within the pod.

* `scheduler_name` - The scheduler name of the pod.

* `security_context` - The security context of the pod.
  The [security_context](#pods_security_context) structure is documented below.

* `set_hostname_as_fqdn` - Whether the pod hostname is configured as the pod FQDN.

* `share_process_namespace` - Whether to share a single process namespace between
  all of containers in a pod.

* `status` - The status of the pod.
  The [status](#pods_status) structure is documented below.

* `termination_grace_period_seconds` - The termination grace period seconds.

* `volumes` - The volumes of the pod.
  The [volumes](#pods_volumes) structure is documented below.

<a name="pods_affinity"></a>
The `affinity` block supports:

* `node_affinity` - The node affinity.
  The [node_affinity](#pods_affinity_node_affinity) structure is documented below.

* `pod_anti_affinity` - The pod anti affinity.
  The [pod_anti_affinity](#pods_affinity_pod_anti_affinity) structure is documented below.

<a name="pods_affinity_node_affinity"></a>
The `node_affinity` block supports:

* `required_during_scheduling_ignored_during_execution` - The required during scheduling
  the ignored during execution.
  The [required_during_scheduling_ignored_during_execution](#pods_affinity_node_affinity_required) structure is
  documented below.

<a name="pods_affinity_node_affinity_required"></a>
The `required_during_scheduling_ignored_during_execution` block supports:

* `node_selector_terms` - The node selector terms.
  The [node_selector_terms](#pods_affinity_node_affinity_required_node_selector_terms) structure is documented below.

<a name="pods_affinity_node_affinity_required_node_selector_terms"></a>
The `node_selector_terms` block supports:

* `match_expressions` - The match expressions.
  The [match_expressions](#pods_affinity_node_affinity_required_node_selector_terms_match_expressions) structure is
  documented below.

<a name="pods_affinity_node_affinity_required_node_selector_terms_match_expressions"></a>
The `match_expressions` block supports:

* `key` - The key of the match expressions.

* `operator` - The operator of the match expressions.

* `values` - The values of the match expressions.

<a name="pods_affinity_pod_anti_affinity"></a>
The `pod_anti_affinity` block supports:

* `preferred_during_scheduling_ignored_during_execution` - The preferred during scheduling ignored during execution.
  The [preferred_during_scheduling_ignored_during_execution](#pods_affinity_pod_anti_affinity_preferred) structure is
  documented below.

* `required_during_scheduling_ignored_during_execution` - The required during schedulin ignored during execution.
  The [required_during_scheduling_ignored_during_execution](#pods_affinity_pod_anti_affinity_required) structure is
  documented below.

<a name="pods_affinity_pod_anti_affinity_preferred"></a>
The `preferred_during_scheduling_ignored_during_execution` block supports:

* `pod_affinity_term` - The pod affinity term.
  The [pod_affinity_term](#pods_affinity_pod_anti_affinity_preferred_pod_affinity_term) structure is documented below.

* `weight` - The weight.

<a name="pods_affinity_pod_anti_affinity_preferred_pod_affinity_term"></a>
The `pod_affinity_term` block supports:

* `label_selector` - The label selector.
  The [label_selector](#pods_affinity_pod_anti_affinity_preferred_pod_affinity_term_label_selector) structure is
  documented below.

* `namespaces` - The namespaces.

* `topology_key` - The topology key.

<a name="pods_affinity_pod_anti_affinity_preferred_pod_affinity_term_label_selector"></a>
The `label_selector` block supports:

* `match_expressions` - The match expressions.
  The [match_expressions](#pods_affinity_pod_anti_affinity_preferred_pod_affinity_term_label_selector_match_expressions)
  structure is documented below.

* `match_labels` - The match labels.

<a name="pods_affinity_pod_anti_affinity_preferred_pod_affinity_term_label_selector_match_expressions"></a>
The `match_expressions` block supports:

* `key` - The key of the match expressions.

* `operator` - The operator of the match expressions.

* `values` - The values of the match expressions.

<a name="pods_affinity_pod_anti_affinity_required"></a>
The `required_during_scheduling_ignored_during_execution` block supports:

* `label_selector` - The label selector.
  The [label_selector](#pods_affinity_pod_anti_affinity_required_label_selector) structure is documented below.

* `namespaces` - The namespaces.

* `topology_key` - The topology key.

<a name="pods_affinity_pod_anti_affinity_required_label_selector"></a>
The `label_selector` block supports:

* `match_expressions` - The match expressions of the label selector.
  The [match_expressions](#pods_affinity_pod_anti_affinity_required_label_selector_match_expressions) structure is
  documented below.

* `match_labels` - The match labels of the label selector.

<a name="pods_affinity_pod_anti_affinity_required_label_selector_match_expressions"></a>
The `match_expressions` block supports:

* `key` - The key of the match labels.

* `operator` - The operator of the match labels.

* `values` - The values of the match labels.

<a name="pods_containers"></a>
The `containers` block supports:

* `args` - The args of the containers.

* `command` - The command of the containers.

* `env` - The environment of the containers.
  The [env](#pods_containers_env) structure is documented below.

* `env_from` - The environment source of the containers.
  The [env_from](#pods_containers_env_from) structure is documented below.

* `image` -  The image of the containers.

* `lifecycle` - The lifecycle of the containers.
  The [lifecycle](#pods_containers_lifecycle) structure is documented below.

* `liveness_probe` - The liveness probe of the containers.
  The [liveness_probe](#pods_containers_probe) structure is documented below.

* `name` - The name of the containers.

* `ports` - The ports of the containers.
  The [ports](#pods_containers_ports) structure is documented below.

* `readiness_probe` - The readiness probe of the containers.
  The [readiness_probe](#pods_containers_probe) structure is documented below.

* `resources` - The resources of the containers.
  The [resources](#pods_containers_resources) structure is documented below.

* `security_context` - The security context of the containers.
  The [security_context](#pods_containers_security_context) structure is documented below.

* `startup_probe` - The startup probe of the containers.
  The [startup_probe](#pods_containers_probe) structure is documented below.

* `stdin` - The stdin of the containers.

* `stdin_once` - The stdin once of the containers.

* `termination_message_path` - The termination message path of the containers.

* `termination_message_policy` - The termination message policy of the containers.

* `tty` - The TTY of the containers.

* `working_dir` - The working dir of the containers.

* `volume_mounts` - The volume mounts probe of the container.
  The [volume_mounts](#containers_volume_mounts) structure is documented below.

<a name="containers_volume_mounts"></a>
The `volume_mounts` block supports:

* `extend_path_mode` - The extend path mode of the volume mounts.

* `mount_path` - The mount path of the volume mounts.

* `name` - The name of the volume mounts.

* `read_only` - Whether to read only.

* `sub_path` - The sub path of the volume mounts.

* `sub_path_expr` - The sub path expression of the volume mounts.

<a name="pods_containers_env"></a>
The `env` block supports:

* `name` - The name of the environment.

* `value` - The value of the environment.

<a name="pods_containers_env_from"></a>
The `env_from` block supports:

* `config_map_ref` - The config map ref of the environment source.
  The [config_map_ref](#pods_containers_env_from_ref) structure is documented below.

* `prefix` - The prefix of the environment source.

* `secret_ref` - The secret ref of the environment source.
  The [secret_ref](#pods_containers_env_from_ref) structure is documented below.

<a name="pods_containers_env_from_ref"></a>
The `config_map_ref`, `secret_ref` block supports:

* `name` - The name.

* `optional` - The optional.

<a name="pods_containers_lifecycle"></a>
The `lifecycle` block supports:

* `post_start` - The post start of the lifecycle.
  The [post_start](#pods_containers_lifecycle_post_start_or_pre_stop) structure is documented below.

* `pre_stop` - The pre stop of the lifecycle.
  The [pre_stop](#pods_containers_lifecycle_post_start_or_pre_stop) structure is documented below.

<a name="pods_containers_lifecycle_post_start_or_pre_stop"></a>
The `post_start`, `pre_stop` block supports:

* `exec` - The exec.
  The [exec](#pods_containers_lifecycle_post_start_or_pre_stop_exec) structure is documented below.

* `http_get` - The http get.
  The [http_get](#pods_containers_lifecycle_post_start_or_pre_stop_http_get) structure is documented below.

<a name="pods_containers_lifecycle_post_start_or_pre_stop_exec"></a>
The `exec` block supports:

* `command` - The command of the exec.

<a name="pods_containers_lifecycle_post_start_or_pre_stop_http_get"></a>
The `http_get` block supports:

* `host` - The host of the http get.

* `http_headers` - The http headers.
  The [http_headers](#pods_containers_lifecycle_post_start_or_pre_stop_http_get_http_headers)
  structure is documented below.

* `path` - The path of the http get.

* `port` - The port of the http get.

* `scheme` - The scheme of the http get.

<a name="pods_containers_lifecycle_post_start_or_pre_stop_http_get_http_headers"></a>
The `http_headers` block supports:

* `name` - The name of the scheme.

* `value` - The value of the scheme.

<a name="pods_containers_probe"></a>
The `liveness_probe`, `readiness_probe`, `startup_probe` block supports:

* `exec` - The exec.
  The [exec](#pods_containers_probe_exec) structure is documented below.

* `failure_threshold` - The minimum consecutive failures for the probe to be
  considered failed after having succeeded.

* `http_get` - The http get.
  The [http_get](#pods_containers_probe_http_get) structure is documented below.

* `initial_delay_seconds` - The number of seconds after the container has started
  before liveness probes are initialed.

* `period_seconds` - How often to perform the probe.

* `success_threshold` - The success threshold.

* `termination_grace_period_seconds` - The termination grace period seconds.

<a name="pods_containers_probe_exec"></a>
The `exec` block supports:

* `command` - The command line to execute inside the container.

<a name="pods_containers_probe_http_get"></a>
The `http_get` block supports:

* `host` - The host name.

* `http_headers` - The custom headers to set in the request.
  The [http_headers](#pods_containers_probe_termination_grace_period_seconds_http_headers) structure is
  documented below.

* `path` - The path to access on the http server.

* `port` - The port to access on the http server.

* `scheme` - The scheme to use for connecting to the host.

<a name="pods_containers_probe_termination_grace_period_seconds_http_headers"></a>
The `http_headers` block supports:

* `name` - The name of the custom http headers.

* `value` - The value of the custom http headers.

<a name="pods_containers_ports"></a>
The `ports` block supports:

* `container_port` - The number of port to expose on the IP address of pod.

* `name` - The port name of the container.

* `protocol` - The protocol for container port.

<a name="pods_containers_resources"></a>
The `resources` block supports:

* `limits` - The limits of resource.

* `requests` - The requests of the resource.

<a name="pods_containers_security_context"></a>
The `security_context` block supports:

* `capabilities` - The capabilities of the security context.
  The [capabilities](#pods_containers_security_context_capabilities) structure is documented below.

* `proc_mount` - The proc mount to use for the containers.

* `read_only_root_file_system` - Whether this container has a read-only root file system.

* `run_as_group` - The GID TO run the entrypoint of the container process.

* `run_as_non_root` - The container must run as a non-root user.

* `run_as_user` - The UID to run the entrypoint of the container process.

<a name="pods_containers_security_context_capabilities"></a>
The `capabilities` block supports:

* `add` - The add of the capabilities.

* `drop` - The drop of the capabilities.

<a name="pods_dns_config"></a>
The `dns_config` block supports:

* `nameservers` - The name servers of the DNS config.

* `options` - The options of the DNS config.
  The [options](#pods_dns_config_options) structure is documented below.

* `searches` - The searches of the DNS config.

<a name="pods_dns_config_options"></a>
The `options` block supports:

* `name` - The name of the options.

* `value` - The value of the options.

<a name="pods_host_aliases"></a>
The `host_aliases` block supports:

* `hostnames` - The host names of the host aliases.

* `ip` - The IP of the host aliases.

<a name="pods_image_pull_secrets"></a>
The `image_pull_secrets` block supports:

* `name` - The name of the image pull secrets.

<a name="pods_readiness_gates"></a>
The `readiness_gates` block supports:

* `condition_type` - The condition type of the readiness gates.

<a name="pods_security_context"></a>
The `security_context` block supports:

* `fs_group` - The FS group of the security context.

* `fs_group_change_policy` - The fs group change policy of the security context.

* `run_as_group` - The GID TO run the entrypoint of the container process.

* `run_as_non_root` - The container must run as a non-root user.

* `run_as_user` - The UID to run the entrypoint of the container process.

* `supplemental_groups` - The supplemental groups.

* `sysctls` - The sysctls.
  The [sysctls](#pods_security_context_sysctls) structure is documented below.

<a name="pods_security_context_sysctls"></a>
The `sysctls` block supports:

* `name` - The name of the sysctls.

* `value` - The value of the sysctls.

<a name="pods_status"></a>
The `status` block supports:

* `conditions` - The conditions of the CCI pod.
  The [conditions](#pods_status_conditions) structure is documented below.

* `observed_generation` - The observed generation of the CCI pod.

<a name="pods_status_conditions"></a>
The `conditions` block supports:

* `last_transition_time` - The last transition time of the CCI pod conditions.

* `last_update_time` - The last update time of the CCI pod conditions.

* `message` - The message of the CCI pod conditions.

* `reason` - The reason of the CCI pod conditions.

* `status` - The status of the CCI pod conditions.

* `type` - The type of the CCI pod conditions.

<a name="pods_volumes"></a>
The `volumes` block supports:

* `config_map` - The config map of the volumes.
  The [config_map](#pods_volumes_config_map) structure is documented below.

* `name` - The name of the volumes.

* `nfs` - The nfs of the volumes.
  The [nfs](#pods_volumes_nfs) structure is documented below.

* `persistent_volume_claim` - The persistent volume claim of the volumes.
  The [persistent_volume_claim](#pods_volumes_persistent_volume_claim) structure is documented below.

* `projected` - The projected of the volumes.
  The [projected](#pods_volumes_projected) structure is documented below.

* `secret` - The secret of the volumes.
  The [secret](#pods_volumes_secret) structure is documented below.

<a name="pods_volumes_config_map"></a>
The `config_map` block supports:

* `default_mode` - The default mode of the config map.

* `items` - The items of the config map.
  The [items](#pods_volumes_config_map_items) structure is documented below.

* `name` - The name.

* `optional` - The optional.

<a name="pods_volumes_config_map_items"></a>
The `items` block supports:

* `key` - The key.

* `mode` - The mode.

* `path` - The path.

<a name="pods_volumes_nfs"></a>
The `nfs` block supports:

* `path` - The path.

* `read_only` - The read only.

* `server` - The server.

<a name="pods_volumes_persistent_volume_claim"></a>
The `persistent_volume_claim` block supports:

* `claim_name` - The claim name.

* `read_only` - The read only.
<a name="pods_volumes_projected"></a>
The `projected` block supports:

* `default_mode` - The default mode.

* `sources` - The sources.
  The [sources](#pods_volumes_projected_sources) structure is documented below.

<a name="pods_volumes_projected_sources"></a>
The `sources` block supports:

* `config_map` - The config map.
  The [config_map](#pods_volumes_projected_sources_config_map) structure is documented below.

* `downward_api` - The downward API.
  The [downward_api](#pods_volumes_projected_sources_downward_api) structure is documented below.

* `secret` - The secret.
  The [secret](#pods_volumes_projected_sources_secret) structure is documented below.

<a name="pods_volumes_projected_sources_config_map"></a>
The `config_map` block supports:

* `items` - The items of the config map.
  The [items](#pods_volumes_projected_sources_config_map_items) structure is documented below.

* `name` - The name of the config map.

* `optional` - The optional of the config map.

<a name="pods_volumes_projected_sources_config_map_items"></a>
The `items` block supports:

* `key` - The key.

* `mode` - The mode.

* `path` - The path.

<a name="pods_volumes_projected_sources_downward_api"></a>
The `downward_api` block supports:

* `items` - The items of the downward API.
  The [items](#pods_volumes_projected_sources_downward_api_items) structure is documented below.

<a name="pods_volumes_projected_sources_downward_api_items"></a>
The `items` block supports:

* `field_ref` - The field ref.
  The [field_ref](#pods_volumes_projected_sources_downward_api_items_field_ref) structure is documented below.

* `mode` - The mode.

* `path` - The path.

* `resource_file_ref` - The resource file ref.
  The [resource_file_ref](#pods_volumes_projected_sources_downward_api_items_resource_file_ref) structure is documented
  below.

<a name="pods_volumes_projected_sources_downward_api_items_field_ref"></a>
The `field_ref` block supports:

* `api_version` - The API version.

* `field_path` - The field path.

<a name="pods_volumes_projected_sources_downward_api_items_resource_file_ref"></a>
The `resource_file_ref` block supports:

* `container_name` - The container name.

* `resource` - The resource.

<a name="pods_volumes_projected_sources_secret"></a>
The `secret` block supports:

* `items` - The items.
  The [items](#pods_volumes_projected_sources_secret_items) structure is documented below.

* `name` - The name.

* `optional` - The optional.

<a name="pods_volumes_projected_sources_secret_items"></a>
The `items` block supports:

* `key` - The key.

* `mode` - The mode.

* `path` - The path.

<a name="pods_volumes_secret"></a>
The `secret` block supports:

* `default_mode` - The default mode.

* `items` - The items.
  The [items](#pods_volumes_secret_items) structure is documented below.

* `optional` - The optional.

* `secret_name` - The secret name.

<a name="pods_volumes_secret_items"></a>
The `items` block supports:

* `key` - The key.

* `mode` - The mode.

* `path` - The path.

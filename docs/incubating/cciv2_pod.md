---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_pod"
description: |-
  Manages a CCI pod resource within HuaweiCloud.
---

# huaweicloud_cciv2_pod

Manages a CCI pod resource within HuaweiCloud.

## Example Usage

```hcl
variable "namespace" {}
variable "name" {}

resource "huaweicloud_cciv2_pod" "test" {
  namespace = var.namespace
  name      = var.name

  annotations = {
    "description"                    = "test",
    "resource.cci.io/pod-size-specs" = "2.00_2.0",
    "resource.cci.io/instance-type"  = "general-computing",
  }

  containers {
    image = "nginx:stable-alpine-perl"
    name  = "c1"

    resources {
      limits = {
        cpu    = 2
        memory = "2G"
      }

      requests = {
        cpu    = 2
        memory = "2G"
      }
    }
  }

  image_pull_secrets {
    name = "imagepull-secret"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `containers` - (Required, List) Specifies the container of the CCI pod.
  The [containers](#containers) structure is documented below.

* `name` - (Required, String) Specifies the name of the CCI pod.

* `namespace` - (Required, String) Specifies the namespace of the CCI pod.

* `active_deadline_seconds` - (Optional, Int) The active deadline seconds the pod.

* `affinity` - (Optional, List) Specifies the affinity of the CCI pod.
  The [affinity](#affinity) structure is documented below.

* `annotations` - (Optional, Map) Specifies the annotations of the CCI pod.

* `dns_config` - (Optional, List) Specifies The DNS config of the pod.
  The [dns_config](#dns_config) structure is documented below.

* `dns_policy` - (Optional, String) Specifies the DNS policy of the pod.

* `ephemeral_containers` - (Optional, List) Specifies the ephemeral container of the CCI pod.
  The [ephemeral_containers](#containers) structure is documented below.

* `host_aliases` - (Optional, List) Specifies the host aliases of the CCI pod.
  The [host_aliases](#host_aliases) structure is documented below.

* `hostname` - (Optional, String) Specifies the host name of the pod.

* `image_pull_secrets` - (Optional, List) Specifies the image pull secrets of the pod.
  The [image_pull_secrets](#image_pull_secrets) structure is documented below.

* `init_containers` - (Optional, List) Specifies the init container of the CCI pod.
  The [init_containers](#containers) structure is documented below.

* `labels` - (Optional, Map) Specifies the labels of the CCI pod.

* `node_name` - (Optional, String) Specifies the node name of the CCI pod.

* `overhead` - (Optional, Map) Specifies the overhead of the CCI pod.

* `readiness_gates` - (Optional, List) Specifies the readiness gates of the CCI pod.
  The [readiness_gates](#readiness_gates) structure is documented below.

* `restart_policy` - (Optional, String) The restart policy for all containers within the pod.

* `scheduler_name` - (Optional, String) Specifies the scheduler name of the CCI pod.

* `security_context` - (Optional, List) Specifies the security context of the CCI pod.
  The [security_context](#security_context) structure is documented below.

* `set_hostname_as_fqdn` - (Optional, Bool) Specifies whether the pod hostname is configured as the pod FQDN.

* `share_process_namespace` - (Optional, Bool) Specifies whether to share a single process namespace between
  all of containers in a pod.

* `termination_grace_period_seconds` - (Optional, Int) Specifies the restart policy for all containers within the pod.

* `volumes` - (Optional, List) Specifies the volumes of the CCI pod.
  The [volumes](#volumes) structure is documented below.

<a name="containers"></a>
The `containers`, `ephemeral_containers`, `init_containers` block supports:

* `name` - (Required, String) Specifies the name of the container.

* `args` - (Optional, List) Specifies the arguments to the entrypoint of the container.

* `command` - (Optional, List) Specifies the command of the container.

* `env` - (Optional, List) Specifies the environment of the container.
  The [env](#containers_env) structure is documented below.

* `env_from` - (Optional, List) Specifies the sources to populate environment variables of the container.
  The [env_from](#containers_env_from) structure is documented below.

* `image` - (Optional, String) Specifies the image name of the container.

* `lifecycle` - (Optional, List) Specifies the lifecycle of the container.
  The [lifecycle](#containers_lifecycle) structure is documented below.

* `liveness_probe` - (Optional, List) Specifies the liveness probe of the container.
  The [liveness_probe](#containers_probe) structure is documented below.

* `ports` - (Optional, List) Specifies the ports of the container.
  The [ports](#containers_ports) structure is documented below.

* `readiness_probe` - (Optional, List) Specifies the readiness probe of the container.
  The [readiness_probe](#containers_probe) structure is documented below.

* `resources` - (Optional, List) Specifies the resources of the container.
  The [resources](#containers_resources) structure is documented below.

* `security_context` - (Optional, List) Specifies the security context of the container.
  The [security_context](#containers_security_context) structure is documented below.

* `startup_probe` - (Optional, List) Specifies the startup probe of the container.
  The [startup_probe](#containers_probe) structure is documented below.

* `stdin` - (Optional, Bool) Specifies whether this container should allocate a buffer for stdin in the container runtime.

* `stdin_once` - (Optional, Bool) Specifies whether this container runtime should close the stdin channel.

* `termination_message_path` - (Optional, String) Specifies the termination message path of the CCI pod container.

* `termination_message_policy` - (Optional, String) Specifies the termination message policy of the CCI pod container.

* `tty` - (Optional, Bool) Specifies whether this container should allocate a TTY for itself.

* `working_dir` - (Optional, String) Specifies the working directory of the CCI Pod container.

<a name="containers_env"></a>
The `env` block supports:

* `name` - (Optional, String) Specifies the name of the environment.

* `value` - (Optional, String) Specifies the value of the environment.

<a name="containers_env_from"></a>
The `env_from` block supports:

* `config_map_ref` - (Optional, List) Specifies the config map.
  The [config_map_ref](#containers_env_from_config_map_ref) structure is documented below.

* `prefix` - (Optional, String) Specifies the prefix.

* `secret_ref` - (Optional, List) Specifies the secret.
  The [secret_ref](#containers_env_from_secret_ref) structure is documented below.

<a name="containers_env_from_config_map_ref"></a>
The `config_map_ref` block supports:

* `name` - (Optional, String) Specifies the name.

* `optional` - (Optional, Bool) Specifies whether to be defined.

<a name="containers_env_from_secret_ref"></a>
The `secret_ref` block supports:

* `name` - (Optional, String) Specifies the name.

* `optional` - (Optional, Bool) Specifies whether to be defined.

<a name="containers_lifecycle"></a>
The `lifecycle` block supports:

* `post_start` - (Optional, List) Specifies the lifecycle post start of the CCI pod container.
  The [post_start](#containers_lifecycle_post_start) structure is documented below.

* `pre_stop` - (Optional, List) Specifies the lifecycle pre stop of the CCI pod container.
  The [pre_stop](#containers_lifecycle_pre_stop) structure is documented below.

<a name="containers_lifecycle_post_start"></a>
The `post_start` block supports:

* `exec` - (Optional, List) Specifies the lifecycle post start of the CCI pod container.
  The [exec](#containers_lifecycle_post_start_exec) structure is documented below.

* `http_get` - (Optional, List) Specifies the lifecycle pre stop of the CCI pod container.
  The [http_get](#containers_lifecycle_post_start_http_get) structure is documented below.

<a name="containers_lifecycle_post_start_exec"></a>
The `http_get` block supports:

* `command` - (Optional, List) Specifies the command line to execute inside the container.

<a name="containers_lifecycle_post_start_http_get"></a>
The `http_get` block supports:

* `host` - (Optional, String) Specifies the host name.

* `http_headers` - (Optional, List) Specifies the custom headers to set in the request.
  The [http_headers](#containers_lifecycle_post_start_http_get_http_headers) structure is documented below.

* `path` - (Optional, String) Specifies the path to access on the HTTP server.

* `port` - (Optional, String) Specifies the port to access on the HTTP server.

* `scheme` - (Optional, String) Specifies the scheme to use for connecting to the host.

<a name="containers_lifecycle_post_start_http_get_http_headers"></a>
The `http_headers` block supports:

* `name` - (Optional, String) Specifies the name of the custom HTTP headers.

* `value` - (Optional, String) Specifies the value of the custom HTTP headers.

<a name="containers_lifecycle_pre_stop"></a>
The `pre_stop` block supports:

* `exec` - (Optional, List) Specifies the lifecycle post start of the CCI Pod container.
  The [exec](#containers_lifecycle_pre_stop_exec) structure is documented below.

* `http_get` - (Optional, List) Specifies the lifecycle pre stop of the CCI Pod container.
  The [http_get](#containers_lifecycle_pre_stop_http_get) structure is documented below.

<a name="containers_lifecycle_pre_stop_exec"></a>
The `http_get` block supports:

* `command` - (Optional, List) Specifies the command line to execute inside the container.

<a name="containers_lifecycle_pre_stop_http_get"></a>
The `http_get` block supports:

* `host` - (Optional, String) Specifies the host name.

* `http_headers` - (Optional, List) Specifies the custom headers to set in the request.
  The [http_headers](#containers_lifecycle_pre_stop_http_get_http_headers) structure is documented below.

* `path` - (Optional, String) Specifies the path to access on the HTTP server.

* `port` - (Optional, String) Specifies the port to access on the HTTP server.

* `scheme` - (Optional, String) Specifies the scheme to use for connecting to the host.

<a name="containers_lifecycle_pre_stop_http_get_http_headers"></a>
The `http_headers` block supports:

* `name` - (Optional, String) Specifies the name of the custom HTTP headers.

* `value` - (Optional, String) Specifies the value of the custom HTTP headers.

<a name="containers_probe"></a>
The `liveness_probe` block supports:

* `exec` - (Optional, List) Specifies the exec.
  The [exec](#containers_probe_exec) structure is documented below.

* `failure_threshold` - (Optional, Int) Specifies the minimum consecutive failures for the probe to be
  considered failed after having succeeded.

* `http_get` - (Optional, List) Specifies the HTTP get.
  The [http_get](#containers_probe_http_get) structure is documented below.

* `initial_delay_seconds` - (Optional, Int) Specifies the number of seconds after the container has started
  before liveness probes are initialed.

* `period_seconds` - (Optional, Int) Specifies how often to perform the probe.

* `termination_grace_period_seconds` - (Optional, Int) Specifies the termination grace period seconds.

<a name="containers_probe_exec"></a>
The `exec` block supports:

* `command` - (Optional, List) Specifies the command line to execute inside the container.

<a name="containers_probe_http_get"></a>
The `http_get` block supports:

* `host` - (Optional, String) Specifies the host name.

* `http_headers` - (Optional, List) Specifies the custom headers to set in the request.
  The [http_headers](#containers_probe_http_get_http_headers) structure is documented below.

* `path` - (Optional, String) Specifies the path to access on the HTTP server.

* `port` - (Optional, String) Specifies the port to access on the HTTP server.

* `scheme` - (Optional, String) Specifies the scheme to use for connecting to the host.

<a name="containers_probe_http_get_http_headers"></a>
The `scheme` block supports:

* `name` - (Optional, String) Specifies the name of the custom HTTP headers.

* `value` - (Optional, String) Specifies the value of the custom HTTP headers.

<a name="containers_ports"></a>
The `ports` block supports:

* `container_port` - (Required, Int) Specifies the number of port to expose on the IP address of pod.

* `name` - (Optional, String) Specifies the port name of the container.

* `protocol` - (Optional, String) Specifies the protocol for container port.

<a name="containers_resources"></a>
The `resources` block supports:

* `limits` - (Optional, Map) Specifies the limits of resource.

* `requests` - (Optional, Map) Specifies the requests of the resource.

<a name="containers_security_context"></a>
The `security_context` block supports:

* `capabilities` - (Optional, List) Specifies the capabilities of the security context.
  The [capabilities](#containers_security_context_capabilities) structure is documented below.

* `proc_mount` - (Optional, String) Specifies the denotes the type of proc mount to use for the containers.

* `read_only_root_file_system` - (Optional, Bool) Specifies whether this container has a read-only root file system.

* `run_as_group` - (Optional, Int) Specifies the GID TO run the entrypoint of the container process.

* `run_as_non_root` - (Optional, Bool) Specifies the container must run as a non-root user.

* `run_as_user` - (Optional, Int) Specifies the UID to run the entrypoint of the container process.

<a name="containers_security_context_capabilities"></a>
The `capabilities` block supports:

* `add` - (Optional, List) Specifies the add of the capabilities.

* `drop` - (Optional, List) Specifies the drop of the capabilities.

<a name="affinity"></a>
The `affinity` block supports:

* `node_affinity` - (Optional, List) Specifies the node affinity.
  The [node_affinity](#affinity_node_affinity) structure is documented below.

* `pod_anti_affinity` - (Optional, List) Specifies the pod anti affinity.
  The [pod_anti_affinity](#affinity_pod_anti_affinity) structure is documented below.

<a name="affinity_node_affinity"></a>
The `node_affinity` block supports:

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

<a name="dns_config"></a>
The `dns_config` block supports:

* `nameservers` - (Optional, List) Specifies the name servers of the DNS config.

* `options` - (Optional, List) Specifies the options of the DNS config.
  The [options](#dns_config_options) structure is documented below.

* `searches` - (Optional, List) Specifies the searches of the DNS config.

<a name="dns_config_options"></a>
The `options` block supports:

* `name` - (Optional, String) Specifies the name of the options.

* `value` - (Optional, String) Specifies the value of the options.

<a name="host_aliases"></a>
The `host_aliases` block supports:

* `hostnames` - (Optional, List) Specifies the host names of the host aliases.

* `ip` - (Optional, String) Specifies the IP of the host aliases.

<a name="image_pull_secrets"></a>
The `image_pull_secrets` block supports:

* `name` - (Optional, String) Specifies the name of the image pull secrets.

<a name="readiness_gates"></a>
The `readiness_gates` block supports:

* `condition_type` - (Optional, String) Specifies the condition type of the readiness gates.

<a name="security_context"></a>
The `security_context` block supports:

* `fs_group` - (Optional, Int) Specifies the fs group of the security context.

* `fs_group_change_policy` - (Optional, String) Specifies the fs group change policy of the security context.

* `run_as_group` - (Optional, Int) Specifies the GID TO run the entrypoint of the container process.

* `run_as_non_root` - (Optional, Bool) Specifies the container must run as a non-root user.

* `run_as_user` - (Optional, Int) Specifies the UID to run the entrypoint of the container process.

* `supplemental_groups` - (Optional, List) Specifies the supplemental groups.

* `sysctls` - (Optional, List) Specifies the sysctls.
  The [sysctls](#security_context_sysctls) structure is documented below.

<a name="security_context_sysctls"></a>
The `sysctls` block supports:

* `name` - (Required, String) Specifies the name of the sysctls.

* `value` - (Required, String) Specifies the value of the sysctls.

<a name="volumes"></a>
The `volumes` block supports:

* `config_map` - (Optional, List) Specifies the config map of the volumes.
  The [config_map](#volumes_config_map) structure is documented below.

* `name` - (Optional, String) Specifies the name of the volumes.
* `nfs` - (Optional, List) Specifies the nfs of the volumes.
  The [nfs](#volumes_nfs) structure is documented below.

* `persistent_volume_claim` - (Optional, List) Specifies the persistent volume claim of the volumes.
  The [persistent_volume_claim](#volumes_persistent_volume_claim) structure is documented below.

* `projected` - (Optional, List) Specifies the projected of the volumes.
  The [projected](#volumes_projected) structure is documented below.

* `secret` - (Optional, List) Specifies the secret of the volumes.
  The [secret](#volumes_secret) structure is documented below.

<a name="volumes_config_map"></a>
The `config_map` block supports:

* `default_mode` - (Optional, Int) Specifies the default mode of the config map.

* `items` - (Optional, List) Specifies the items of the config map.
  The [items](#volumes_config_map_items) structure is documented below.

* `name` - (Optional, String) Specifies the name.

* `optional` - (Optional, Bool) Specifies the optional.

<a name="volumes_config_map_items"></a>
The `items` block supports:

* `key` - (Required, String) Specifies the key.

* `path` - (Required, String) Specifies the path.

* `mode` - (Optional, Int) Specifies the name.

<a name="volumes_nfs"></a>
The `nfs` block supports:

* `path` - (Required, String) Specifies the path.

* `server` - (Required, String) Specifies the server.

* `read_only` - (Optional, Bool) Specifies whether to read only.

<a name="volumes_persistent_volume_claim"></a>
The `persistent_volume_claim` block supports:

* `claim_name` - (Required, String) Specifies the claim name of the persistent volume claim.

* `read_only` - (Optional, Bool) Specifies whether to read only.

<a name="volumes_projected"></a>
The `projected` block supports:

* `default_mode` - (Optional, Int) Specifies the rolling update config of the CCI pod strategy.

* `sources` - (Optional, List) Specifies the type of the CCI pod strategy.
  The [sources](#volumes_projected_sources) structure is documented below.

<a name="volumes_projected_sources"></a>
The `sources` block supports:

* `config_map` - (Optional, List) Specifies the config map of the namevolumes projected sources.
  The [config_map](#volumes_projected_sources_config_map) structure is documented below.

* `downward_api` - (Optional, List) Specifies the downward API of the namevolumes projected sources.
  The [downward_api](#volumes_projected_sources_downward_api) structure is documented below.

* `secret` - (Optional, List) Specifies the secret of the namevolumes projected sources.
  The [secret](#volumes_projected_sources_secret) structure is documented below.

<a name="volumes_projected_sources_config_map"></a>
The `config_map` block supports:

* `items` - (Optional, List) Specifies the items of the config map.
  The [items](#volumes_projected_sources_secret_items) structure is documented below.

* `name` - (Optional, String) Specifies the name of the config map.

* `optional` - (Optional, Bool) Specifies the optional of the config map.

<a name="volumes_projected_sources_secret_items"></a>
The `items` block supports:

* `key` - (Required, String) Specifies the key of the items.

* `path` - (Required, String) Specifies the path of the items.

* `mode` - (Optional, Int) Specifies the mode of the items.

<a name="volumes_projected_sources_downward_api"></a>
The `downward_api` block supports:

* `items` - (Optional, List) Specifies the items of the downward API.
  The [items](#volumes_projected_sources_secret_items) structure is documented below.

<a name="volumes_projected_sources_secret_items"></a>
The `items` block supports:

* `field_ref` - (Optional, List) Specifies the field ref of the items.
  The [field_ref](#volumes_projected_sources_secret_items_field_ref) structure is documented below.

* `mode` - (Optional, Int) Specifies the mode of the items.

* `path` - (Optional, String) Specifies the path of the items.

* `resource_file_ref` - (Optional, List) Specifies the resource file ref of the items.
  The [resource_file_ref](#volumes_projected_sources_secret_items_resource_file_ref) structure is documented below.

<a name="volumes_projected_sources_secret_items_field_ref"></a>
The `field_ref` block supports:

* `field_path` - (Required, String) Specifies the field path of the file ref.

* `api_version` - (Optional, String) Specifies the API version of the file ref.

<a name="volumes_projected_sources_secret_items_resource_file_ref"></a>
The `resource_file_ref` block supports:

* `resource` - (Required, String) Specifies the resource of the resource file ref.

* `container_name` - (Optional, String) Specifies the container name of the resource file ref.

<a name="volumes_projected_sources_secret"></a>
The `secret` block supports:

* `items` - (Optional, List) Specifies the items of the secret.
  The [items](#volumes_projected_sources_secret_items) structure is documented below.

* `name` - (Optional, String) Specifies the name of the secret.

* `optional` - (Optional, Bool) Specifies the optional of the secret.

<a name="volumes_projected_sources_secret_items"></a>
The `items` block supports:

* `key` - (Required, String) Specifies the key of the items.

* `path` - (Required, String) Specifies the path of the items.

* `mode` - (Optional, Int) Specifies the mode of the items.

<a name="volumes_secret"></a>
The `secret` block supports:

* `default_mode` - (Optional, Int) Specifies the default mode of the secret.

* `items` - (Optional, List) Specifies the items of the secret.
  The [items](#volumes_secret_items) structure is documented below.

* `optional` - (Optional, Bool) Specifies the optional of the secret.

* `secret_name` - (Optional, String) Specifies the secret name of the secret.

<a name="volumes_secret_items"></a>
The `items` block supports:

* `key` - (Required, String) Specifies the key of the items.

* `path` - (Required, String) Specifies the path of the items.

* `mode` - (Optional, Int) Specifies the mode of the items.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `api_version` - The API version of the CCI pod.

* `creation_timestamp` - The creation timestamp of the CCI pod.

* `finalizers` - The finalizers of the namespace.

* `kind` - The kind of the CCI pod.

* `resource_version` - The resource version of the CCI pod.

* `status` - The status of the CCI pod.
  The [status](#attrstatus) structure is documented below.

* `uid` - The uid of the CCI pod.

<a name="attrstatus"></a>
The `status` block supports:

* `conditions` - The conditions of the CCI pod.
  The [conditions](#attrstatus_conditions) structure is documented below.

* `observed_generation` - The observed generation of the CCI pod.

<a name="attrstatus_conditions"></a>
The `conditions` block supports:

* `last_transition_time` - The last transition time of the CCI pod conditions.

* `last_update_time` - The last update time of the CCI pod conditions.

* `message` - The message of the CCI pod conditions.

* `reason` - The reason of the CCI pod conditions.

* `status` - The status of the CCI pod conditions.

* `type` - The type of the CCI pod conditions.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The xxx can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_cciv2_pod.test <id>
```

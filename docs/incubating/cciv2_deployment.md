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

* `delete_propagation_policy` - (Optional, String) Where and how garbage collection will be performed.
   The default policy is decided by the existing finalizer set in the metadata.finalizers
   and the resource-specific default policy.
   - `Orphan`: orphan the dependents.
   - `Background`: allow the garbage collector to delete the dependents in the background.
   - `Foreground`: a cascading policy that deletes all dependents in the foreground.

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

* `args` - (Optional, List) Specifies the arguments to the entrypoint of the container.

* `command` - (Optional, List) Specifies the command of the container.

* `env` - (Optional, List) Specifies the environment of the container.
  The [env](#template_spec_containers_env) structure is documented below.

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

* `volume_mounts` - (Optional, List) Specifies the volume mounts probe of the container.
  The [volume_mounts](#containers_volume_mounts) structure is documented below.

<a name="containers_volume_mounts"></a>
The `volume_mounts` block supports:

* `extend_path_mode` - (Optional, String) Specifies the extend path mode of the volume mounts.

* `mount_path` - (Required, String) Specifies the mount path of the volume mounts.

* `name` - (Required, String) Specifies the name of the volume mounts.

* `read_only` - (Optional, Bool) Specifies whether to read only.

* `sub_path` - (Optional, String) Specifies the sub path of the volume mounts.

* `sub_path_expr` - (Optional, String) Specifies the sub path expression of the volume mounts.

<a name="template_spec_containers_env"></a>
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

* `exec` - (Optional, List) Specifies the exec.
  The [exec](#containers_lifecycle_post_start_exec) structure is documented below.

* `http_get` - (Optional, List) Specifies the http get.
  The [http_get](#containers_lifecycle_post_start_http_get) structure is documented below.

<a name="containers_lifecycle_post_start_exec"></a>
The `exec` block supports:

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

* `exec` - (Optional, List) Specifies the exec.
  The [exec](#containers_lifecycle_pre_stop_exec) structure is documented below.

* `http_get` - (Optional, List) Specifies the http get.
  The [http_get](#containers_lifecycle_pre_stop_http_get) structure is documented below.

<a name="containers_lifecycle_pre_stop_exec"></a>
The `exec` block supports:

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
The `http_headers` block supports:

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

* `proc_mount` - (Optional, String) Specifies the proc mount to use for the containers.

* `read_only_root_file_system` - (Optional, Bool) Specifies whether this container has a read-only root file system.

* `run_as_group` - (Optional, Int) Specifies the GID TO run the entrypoint of the container process.

* `run_as_non_root` - (Optional, Bool) Specifies the container must run as a non-root user.

* `run_as_user` - (Optional, Int) Specifies the UID to run the entrypoint of the container process.

<a name="containers_security_context_capabilities"></a>
The `capabilities` block supports:

* `add` - (Optional, List) Specifies the add of the capabilities.

* `drop` - (Optional, List) Specifies the drop of the capabilities.

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

---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_training_jobs"
description: |-
  Use this data source to get the list of ModelArts training jobs within HuaweiCloud.
---

# huaweicloud_modelarts_training_jobs

Use this data source to get the list of ModelArts training jobs within HuaweiCloud.

## Example Usage

### Query all training jobs

```hcl
data "huaweicloud_modelarts_training_jobs" "test" {}
```

### Filter training jobs by its id

```hcl
variable "training_job_id" {}

data "huaweicloud_modelarts_training_jobs" "test" {
  filters {
    key      = "id"
    operator = "in"
    value    = [var.training_job_id]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the training jobs are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Optional, String) Specifies the workspace ID of the training jobs to be queried.  
  If omitted, all training jobs in the default workspace will be queried.

* `sort_by` - (Optional, String) Specifies the metric used to sort the training jobs.  
  The valid values are as follows:
  + **create_time**

  Defaults to **create_time**.

* `order` - (Optional, String) Specifies the sort order of the training jobs.  
  The valid values are as follows:
  + **asc**
  + **desc**

  Defaults to **desc**.

* `unified_jobs` - (Optional, Bool) Specifies whether to query custom jobs and fine-tuning jobs together.  
  Defaults to **false**.

* `train_type` - (Optional, String) Specifies the training job type to be queried.  
  The valid values are as follows:
  + **job**
  + **ftjob**

* `filters` - (Optional, List) Specifies the filter conditions used to query training jobs.  
  The [filters](#modelarts_training_jobs_filters_arg) structure is documented below.

<a name="modelarts_training_jobs_filters_arg"></a>
The `filters` block supports:

* `key` - (Required, String) Specifies the filter key.  
  The valid values are as follows:
  + **id**
  + **name**
  + **kind**
  + **phase**
  + **algorithm_id**
  + **algorithm_name**
  + **create_time**
  + **user_id**
  + **pool_id**
  + **training_experiment_id**
  + **runtime_type**
  + **priority**

* `operator` - (Optional, String) Specifies the filter operator.  
  The valid values are as follows:
  + **between**
  + **like**
  + **in**
  + **not**

* `value` - (Optional, List) Specifies the filter values.  
  The maximum length of the list is `10`.  
  When `key` is set to **create_time**, the value is a list of two elements, representing the start and end times
  of the creation time range, in RFC3339 format. The maximum value range of the creation time is `31` days.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `jobs` - The list of training jobs that match the filter parameters.  
  The [jobs](#modelarts_training_jobs) structure is documented below.

<a name="modelarts_training_jobs"></a>
The `jobs` block supports:

* `kind` - The type of the training job.
  + **job**
  + **federated_pool_job**
  + **edge_job**
  + **hetero_job**
  + **mrs_job**
  + **autosearch_job**
  + **diag_job**
  + **visualization_job**

* `metadata` - The metadata of the training job.  
  The [metadata](#modelarts_training_jobs_metadata) structure is documented below.

* `status` - The status of the training job.  
  The [status](#modelarts_training_jobs_status) structure is documented below.

* `algorithm` - The algorithm configuration of the training job.  
  The [algorithm](#modelarts_training_jobs_algorithm) structure is documented below.

* `spec` - The specification of the training job.  
  The [spec](#modelarts_training_jobs_spec) structure is documented below.

* `endpoints` - The remote access endpoints of the training job.  
  The [endpoints](#modelarts_training_jobs_endpoints) structure is documented below.

* `create_time` - The creation time of the training job, in RFC3339 format.

<a name="modelarts_training_jobs_metadata"></a>
The `metadata` block supports:

* `id` - The ID of the training job.

* `name` - The name of the training job.

* `workspace_id` - The workspace ID to which the training job belongs.

* `description` - The description of the training job.

* `user_name` - The user name that created the training job.

* `annotations` - The advanced feature configuration of the training job.

<a name="modelarts_training_jobs_status"></a>
The `status` block supports:

* `phase` - The primary status of the training job.
  + **Running**
  + **Failed**
  + **Completed**
  + **Terminated**
  + **Abnormal**

* `secondary_phase` - The secondary status of the training job.
  + **Running**
  + **Failed**
  + **Completed**
  + **Terminated**
  + **CreateFailed**
  + **TerminatedFailed**
  + **Unknown**
  + **Lost**: Abnormal.

* `duration` - The running duration of the training job, in milliseconds.

* `start_time` - The start time of the training job, in RFC3339 format.

* `tasks` - The subtask names of the training job.

* `node_count_metrics` - The node count metrics of the training job, in JSON format.

* `task_statuses` - The status of the first failed subtask.  
  The [task_statuses](#modelarts_training_jobs_status_task_statuses) structure is documented below.

* `running_records` - The running and fault recovery records of the training job.  
  The [running_records](#modelarts_training_jobs_status_running_records) structure is documented below.

<a name="modelarts_training_jobs_status_task_statuses"></a>
The `task_statuses` block supports:

* `task` - The subtask name.

* `exit_code` - The exit code of the subtask.

* `message` - The error message of the subtask.

<a name="modelarts_training_jobs_status_running_records"></a>
The `running_records` block supports:

* `start_at` - The start time of the run, in RFC3339 format.

* `end_at` - The end time of the run, in RFC3339 format.

* `xpu_start_at` - The accelerator start time of the run, in RFC3339 format.

* `start_type` - The start type of the run.

* `end_reason` - The end reason of the run.

* `end_related_task` - The task worker ID that caused the run to end.

* `end_recover` - The final fault tolerance strategy when the run ended abnormally.
  + **npu_proc_restart**
  + **proc_restart**
  + **npu_step_retry**
  + **pod_reschedule**
  + **job_reschedule**
  + **job_reschedule_with_taint**

* `end_recover_before_downgrade` - The fault tolerance strategy before downgrade when the run ended abnormally.

* `recover_records` - The fault tolerance strategy details when the run ended abnormally.  
  The [recover_records](#modelarts_training_jobs_status_running_records_recover_records) structure is documented below.

<a name="modelarts_training_jobs_status_running_records_recover_records"></a>
The `recover_records` block supports:

* `recover_start_at` - The start time of the fault tolerance strategy, in RFC3339 format.

* `recover_end_at` - The end time of the fault tolerance strategy, in RFC3339 format.

* `recover` - The fault tolerance strategy.
  + **npu_step_retry**
  + **npu_proc_restart**
  + **proc_restart**
  + **pod_reschedule**
  + **job_reschedule**
  + **job_reschedule_with_taint**

* `fault_scenario` - The fault scenario.
  + **chip_fault**
  + **node_fault**
  + **job_failed**
  + **job_hanged**
  + **job_subhealth**
  + **error_in_log**

* `reason` - The fault reason.

* `related_task` - The task worker ID that triggered the fault.

* `recover_result` - The execution result of the fault tolerance strategy.
  + **recovering**
  + **success**
  + **failed**
  + **downgrade**: Strategy downgrade.
  + **terminated**: Strategy terminated.
  + **quotaExceeded**: Strategy execution times exceeded.

<a name="modelarts_training_jobs_algorithm"></a>
The `algorithm` block supports:

* `id` - The algorithm ID of the training job.

* `name` - The algorithm name of the training job.

* `subscription_id` - The subscription ID of the subscribed algorithm.

* `item_version_id` - The version ID of the subscribed algorithm.

* `code_dir` - The code directory of the training job.

* `boot_file` - The boot file of the training job.

* `autosearch_config_path` - The YAML configuration path of the auto search job.

* `autosearch_framework_path` - The framework code directory of the auto search job.

* `command` - The startup command of the custom image training job.

* `local_code_dir` - The local code directory in the training container.

* `working_dir` - The working directory when running the algorithm.

* `parameters` - The runtime parameters of the training job.  
  The [parameters](#modelarts_training_jobs_algorithm_parameters) structure is documented below.

* `inputs` - The input channels of the training job.  
  The [inputs](#modelarts_training_jobs_algorithm_inputs) structure is documented below.

* `outputs` - The output channels of the training job.  
  The [outputs](#modelarts_training_jobs_algorithm_outputs) structure is documented below.

* `engine` - The engine configuration of the training job.  
  The [engine](#modelarts_training_jobs_algorithm_engine) structure is documented below.

* `environments` - The environment variables of the training job.

<a name="modelarts_training_jobs_algorithm_parameters"></a>
The `parameters` block supports:

* `name` - The parameter name.

* `value` - The parameter value.

* `description` - The parameter description.

* `constraint` - The parameter constraint.  
  The [constraint](#modelarts_training_jobs_algorithm_parameters_constraint) structure is documented below.

<a name="modelarts_training_jobs_algorithm_parameters_constraint"></a>
The `constraint` block supports:

* `type` - The parameter constraint type.
  + **Integer**
  + **Float**
  + **String**
  + **Boolean**

* `editable` - Whether the parameter is editable.

* `required` - Whether the parameter is required.

* `valid_type` - The parameter valid type.
  + **Choice**
  + **Range**
  + **None**

* `valid_range` - The parameter valid values.

<a name="modelarts_training_jobs_algorithm_engine"></a>
The `engine` block supports:

* `engine_id` - The engine specification ID of the training job.

* `engine_name` - The engine specification name of the training job.

* `engine_version` - The engine specification version of the training job.

* `image_url` - The custom image URL of the training job.

* `install_sys_packages` - Whether to install the moxing version specified by the training platform.

<a name="modelarts_training_jobs_algorithm_inputs"></a>
The `inputs` block supports:

* `name` - The name of the input channel.

* `description` - The description of the input channel.

* `local_dir` - The local directory mapped by the input channel.

* `access_method` - The delivery method of the input channel path.
  + **parameter**
  + **env**

* `remote` - The remote input information.  
  The [remote](#modelarts_training_jobs_algorithm_inputs_remote) structure is documented below.

* `remote_constraint` - The data input constraint.  
  The [remote_constraint](#modelarts_training_jobs_algorithm_inputs_remote_constraint) structure is documented below.

<a name="modelarts_training_jobs_algorithm_inputs_remote"></a>
The `remote` block supports:

* `dataset` - The dataset input information.  
  The [dataset](#modelarts_training_jobs_algorithm_inputs_remote_dataset) structure is documented below.

* `obs` - The OBS input information.  
  The [obs](#modelarts_training_jobs_algorithm_inputs_remote_obs) structure is documented below.

<a name="modelarts_training_jobs_algorithm_inputs_remote_dataset"></a>
The `dataset` block supports:

* `id` - The dataset ID.

* `version_id` - The dataset version ID.

* `obs_url` - The OBS URL of the dataset.

* `service_type` - The dataset service type.

* `name` - The dataset name.

<a name="modelarts_training_jobs_algorithm_inputs_remote_obs"></a>
The `obs` block supports:

* `obs_url` - The OBS URL of the input data.

<a name="modelarts_training_jobs_algorithm_inputs_remote_constraint"></a>
The `remote_constraint` block supports:

* `data_type` - The data type of the remote constraint.

* `attributes` - The attributes of the remote constraint.
  + **data_format**
  + **data_segmentation**
  + **dataset_type**

<a name="modelarts_training_jobs_algorithm_outputs"></a>
The `outputs` block supports:

* `name` - The name of the output channel.

* `description` - The description of the output channel.

* `local_dir` - The local directory mapped by the output channel.

* `access_method` - The delivery method of the output channel path.
  + **parameter**
  + **env**

* `remote` - The remote output information.  
  The [remote](#modelarts_training_jobs_algorithm_outputs_remote) structure is documented below.

<a name="modelarts_training_jobs_algorithm_outputs_remote"></a>
The `remote` block supports:

* `obs` - The OBS output information.  
  The [obs](#modelarts_training_jobs_algorithm_outputs_remote_obs) structure is documented below.

<a name="modelarts_training_jobs_algorithm_outputs_remote_obs"></a>
The `obs` block supports:

* `obs_url` - The OBS URL of the output data.

<a name="modelarts_training_jobs_spec"></a>
The `spec` block supports:

* `resource` - The resource specification of the training job.  
  The [resource](#modelarts_training_jobs_spec_resource) structure is documented below.

* `runtime_type` - The runtime type of the training job.

* `volumes` - The mounted volumes of the training job.  
  The [volumes](#modelarts_training_jobs_spec_volumes) structure is documented below.

* `log_export_path` - The log export path of the training job.  
  The [log_export_path](#modelarts_training_jobs_spec_log_export_path) structure is documented below.

* `schedule_policy` - The scheduling policy of the training job.  
  The [schedule_policy](#modelarts_training_jobs_spec_schedule_policy) structure is documented below.

* `custom_metrics` - The custom metrics configuration of the training job.  
  The [custom_metrics](#modelarts_training_jobs_spec_custom_metrics) structure is documented below.

<a name="modelarts_training_jobs_spec_resource"></a>
The `resource` block supports:

* `policy` - The resource specification mode of the training job.

* `flavor_id` - The flavor ID of the training job.

* `flavor_name` - The flavor name of the training job.

* `node_count` - The number of resource replicas selected by the training job.

* `pool_id` - The resource pool ID selected by the training job.

* `pool_group_id` - The federated resource pool ID selected by the training job.

* `main_container_allocated_resources` - The allocated resources of the main container.  
  The [main_container_allocated_resources](#modelarts_training_jobs_spec_resource_main_container_allocated_resources)
  structure is documented below.

* `main_container_customized_flavor` - The customized flavor of the main container.  
  The [main_container_customized_flavor](#modelarts_training_jobs_spec_resource_main_container_customized_flavor)
  structure is documented below.

<a name="modelarts_training_jobs_spec_resource_main_container_allocated_resources"></a>
The `main_container_allocated_resources` block supports:

* `cpu_arch` - The CPU architecture.

* `cpu_core_num` - The number of CPU cores.

* `mem_size` - The memory size.

* `accelerator_num` - The number of accelerator cards.

* `accelerator_type` - The accelerator type.

<a name="modelarts_training_jobs_spec_resource_main_container_customized_flavor"></a>
The `main_container_customized_flavor` block supports:

* `cpu_core_num` - The number of CPU cores.

* `mem_size` - The memory size.

* `accelerator_num` - The number of accelerator cards.

<a name="modelarts_training_jobs_spec_volumes"></a>
The `volumes` block supports:

* `nfs` - The NFS volume configuration.  
  The [nfs](#modelarts_training_jobs_spec_volumes_nfs) structure is documented below.

<a name="modelarts_training_jobs_spec_volumes_nfs"></a>
The `nfs` block supports:

* `nfs_server_path` - The NFS server path.

* `local_path` - The path for attaching volumes to the training container.

* `read_only` - Whether the disks attached in NFS mode are read-only.

<a name="modelarts_training_jobs_spec_log_export_path"></a>
The `log_export_path` block supports:

* `obs_url` - The OBS path where the training job logs are exported.

<a name="modelarts_training_jobs_spec_schedule_policy"></a>
The `schedule_policy` block supports:

* `priority` - The scheduling priority of the training job.

* `preemptible` - Whether the training job can be preempted.

* `required_affinity` - The required affinity policy of the training job.  
  The [required_affinity](#modelarts_training_jobs_spec_schedule_policy_required_affinity) structure is
  documented below.

* `preferred_affinity` - The preferred affinity configuration of the training job.  
  The [preferred_affinity](#modelarts_training_jobs_spec_schedule_policy_preferred_affinity) structure is
  documented below.

<a name="modelarts_training_jobs_spec_schedule_policy_required_affinity"></a>
The `required_affinity` block supports:

* `affinity_type` - The affinity scheduling policy.
  + **cabinet**: Strong cabinet scheduling.
  + **hyperinstance**: Hypernode affinity scheduling.

* `affinity_group_size` - The affinity group size.

* `node_affinity` - The node affinity configuration of the training job.  
  The [node_affinity](#modelarts_training_jobs_spec_schedule_policy_required_affinity_node_affinity) structure is
  documented below.

<a name="modelarts_training_jobs_spec_schedule_policy_required_affinity_node_affinity"></a>
The `node_affinity` block supports:

* `node_selector_terms` - The required node selector terms.  
  The [node_selector_terms](#modelarts_training_jobs_spec_schedule_policy_node_selector_terms) structure is
  documented below.

<a name="modelarts_training_jobs_spec_schedule_policy_preferred_affinity"></a>
The `preferred_affinity` block supports:

* `node_affinity` - The preferred node affinity terms.  
  The [node_affinity](#modelarts_training_jobs_spec_schedule_policy_preferred_affinity_node_affinity) structure
  is documented below.

<a name="modelarts_training_jobs_spec_schedule_policy_preferred_affinity_node_affinity"></a>
The `node_affinity` block supports:

* `weight` - The weight associated with the preferred node selector term.

* `preference` - The preferred node selector term.  
  The [preference](#modelarts_training_jobs_spec_schedule_policy_node_selector_terms) structure is documented below.

<a name="modelarts_training_jobs_spec_schedule_policy_node_selector_terms"></a>
The `node_selector_terms` and `preference` block supports:

* `match_expressions` - The node selector requirements based on node labels.  
  The [match_expressions](#modelarts_training_jobs_spec_schedule_policy_node_selector_requirement) structure is
  documented below.

* `match_fields` - The node selector requirements based on node fields.  
  The [match_fields](#modelarts_training_jobs_spec_schedule_policy_node_selector_requirement) structure is
  documented below.

<a name="modelarts_training_jobs_spec_schedule_policy_node_selector_requirement"></a>
The `match_expressions` and `match_fields` block supports:

* `key` - The label key used by the node selector requirement.

* `operator` - The operator used by the node selector requirement.

* `values` - The label values used by the node selector requirement.

<a name="modelarts_training_jobs_spec_custom_metrics"></a>
The `custom_metrics` block supports:

* `exec` - The command-based metrics collection configuration.  
  The [exec](#modelarts_training_jobs_spec_custom_metrics_exec) structure is documented below.

* `http_get` - The HTTP-based metrics collection configuration.  
  The [http_get](#modelarts_training_jobs_spec_custom_metrics_http_get) structure is documented below.

<a name="modelarts_training_jobs_spec_custom_metrics_exec"></a>
The `exec` block supports:

* `command` - The command for metrics collection.

<a name="modelarts_training_jobs_spec_custom_metrics_http_get"></a>
The `http_get` block supports:

* `path` - The URL path for HTTP metrics collection.

* `port` - The port for HTTP metrics collection.

<a name="modelarts_training_jobs_endpoints"></a>
The `endpoints` block supports:

* `ssh` - The SSH connection information.  
  The [ssh](#modelarts_training_jobs_endpoints_ssh) structure is documented below.

* `jupyter_lab` - The JupyterLab connection information.  
  The [jupyter_lab](#modelarts_training_jobs_endpoints_jupyter_lab) structure is documented below.

* `tensorboard` - The Tensorboard connection information.  
  The [tensorboard](#modelarts_training_jobs_endpoints_tensorboard) structure is documented below.

* `mindstudio_insight` - The MindStudio Insight connection information.  
  The [mindstudio_insight](#modelarts_training_jobs_endpoints_mindstudio_insight) structure is documented below.

<a name="modelarts_training_jobs_endpoints_ssh"></a>
The `ssh` block supports:

* `key_pair_names` - The SSH key pair names.

* `task_urls` - The SSH connection URLs.  
  The [task_urls](#modelarts_training_jobs_endpoints_ssh_task_urls) structure is documented below.

<a name="modelarts_training_jobs_endpoints_ssh_task_urls"></a>
The `task_urls` block supports:

* `task` - The task ID.

* `url` - The SSH connection URL.

<a name="modelarts_training_jobs_endpoints_jupyter_lab"></a>
The `jupyter_lab` block supports:

* `url` - The JupyterLab URL.

* `token` - The JupyterLab token.

<a name="modelarts_training_jobs_endpoints_tensorboard"></a>
The `tensorboard` block supports:

* `url` - The Tensorboard URL.

* `token` - The Tensorboard token.

<a name="modelarts_training_jobs_endpoints_mindstudio_insight"></a>
The `mindstudio_insight` block supports:

* `url` - The MindStudio Insight URL.

* `token` - The MindStudio Insight token.

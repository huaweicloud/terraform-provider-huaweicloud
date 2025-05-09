---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_cce_accesses"
description: |-
  Using this data source to query the list of CCE access configurations within HuaweiCloud.
---

# huaweicloud_lts_cce_accesses

Using this data source to query the list of CCE access configurations within HuaweiCloud.

## Example Usage

```hcl
variable "access_config_name" {}

data "huaweicloud_lts_cce_accesses" "test" {
  name = var.access_config_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query CCE access configurations.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the CCE access to be queried.

* `host_group_name` - (Optional, String) Specifies the name of the host group to be queried.

* `log_group_name` - (Optional, String) Specifies the name of the log group to which the access configurations and log
  streams belong.

* `log_stream_name` - (Optional, String) Specifies the name of the log stream to which the access configurations belong.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the CCE access.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `accesses` - The list of CCE access configurations.
  The [accesses](#data_accesses_attr) structure is documented below.

<a name="data_accesses_attr"></a>
The `accesses` block supports:

* `id` - The ID of the CCE access.

* `name` - The name of the CCE access.

* `log_group_id` - The ID of the log group.

* `log_group_name` - The name of the log group.

* `log_stream_id` - The ID of the log stream.

* `log_stream_name` - The name of the log stream.

* `access_config` - The configuration detail of the CCE access.  
  The [access_config](#data_accesses_elem_access_config) structure is documented below.

* `cluster_id` - The ID of the cluster corresponding to CCE access.

* `host_group_ids` - The ID list of the log access host groups.

* `tags` - The key/value pairs to associate with the CCE access.

* `binary_collect` - Whether collect in binary format.

* `log_split` - Whether to split log.

* `access_type` - The type of the log access.

* `processor_type` - The type of the ICAgent structuring parsing.
  + **SINGLE_LINE**
  + **MULTI_LINE**
  + **REGEX**
  + **MULTI_REGEX**
  + **SPLIT**
  + **JSON**
  + **NGINX**
  + **COMPOSE**

* `processors` - The list of the ICAgent structuring parsing rules.  
  The [processors](#data_cce_accesses_processors) structure is documented below.

* `demo_log` - The example log of the ICAgent structuring parsing.  
  This parameter is available when the `processor_type` parameter is specified.

* `demo_fields` - The list of the parsed fields of the example log.  
  The [demo_fields](#data_cce_accesses_demoFields) structure is documented below.

* `encoding_format` - The encoding format log file.
  + **UTF-8**
  + **GBK**

* `incremental_collect` - Whether to collect incrementally.

* `created_at` - The creation time of the CCE access, in RFC3339 format.

<a name="data_accesses_elem_access_config"></a>
The `access_config` block supports:

* `path_type` - The type of the CCE access.
  + **container_stdout**
  + **container_file**
  + **host_file**

* `paths` - The collection paths.

* `black_paths` - The collection path blacklist.

* `windows_log_info` - The configuration of Windows event logs.  
  The [windows_log_info](#data_access_config_elem_windows_log_info) structure is documented below.

* `single_log_format` - The configuration single-line logs.  
  The [single_log_format](#data_access_config_elem_single_log_format) structure is documented below.

* `multi_log_format` - The configuration multi-line logs.  
  The [multi_log_format](#data_access_config_elem_multi_log_format) structure is documented below.

* `stdout` - Whether output is standard.

* `stderr` - Whether error output is standard.

* `name_space_regex` - The regular expression matching of kubernetes namespaces.

* `pod_name_regex` - The regular expression matching of kubernetes pods.

* `container_name_regex` - The regular expression matching of kubernetes container names.

* `log_labels` - The container label log tag.

* `include_labels_logical` - The logical relationship between multiple container label whitelists.
  + **and**
  + **or**

* `include_labels` - The container label whitelist.

* `exclude_labels_logical` - The logical relationship between multiple container label blacklists.
  + **and**
  + **or**

* `exclude_labels` - The container label blacklist.

* `log_envs` - The environment variable tag.

* `include_envs_logical` - The logical relationship between multiple environment variable whitelists.
  + **and**
  + **or**

* `include_envs` - The environment variable whitelist.

* `exclude_envs_logical` - The logical relationship between multiple environment variable blacklist.
  + **and**
  + **or**

* `exclude_envs` - The environment variable blacklist.

* `log_k8s` - The kubernetes label log tag.

* `include_k8s_labels_logical` - The logical relationship between multiple kubernetes label whitelists.
  + **and**
  + **or**

* `include_k8s_labels` - The kubernetes label whitelist.

* `exclude_k8s_labels_logical` - The logical relationship between multiple kubernetes label blacklist.
  + **and**
  + **or**

* `exclude_k8s_labels` - The kubernetes label blacklist.

* `repeat_collect` - Whether to allow repeated file collection.

* `custom_key_value` - The custom key/value pairs of the CCE access.

* `system_fields` - The list of system built-in fields of the CCE access.

<a name="data_access_config_elem_windows_log_info"></a>
The `windows_log_info` block supports:

* `categorys` - The types of Windows event logs to be collected.
  + **Application**
  + **System**
  + **Security**
  + **Setup**

* `event_level` - The list of Windows event levels.  
  The element includes:
  + **information**
  + **warning**
  + **error**
  + **critical**
  + **verbose**

* `time_offset_unit` - The collection time offset unit.
  + **day**
  + **hour**
  + **sec**

* `time_offset` - The collection time offset.

<a name="data_access_config_elem_single_log_format"></a>
The `single_log_format` block supports:

* `mode` - The mode of single-line log format.
  + **system**: the system time.
  + **wildcard**: the time wildcard.

* `value` - The value of single-line log format.

<a name="data_access_config_elem_multi_log_format"></a>
The `multi_log_format` block supports:

* `mode` - The mode of multi-line log format.
  + **time**: the time wildcard.
  + **regular**: the regular expression.

* `value` - The value of multi-line log format.

<a name="data_cce_accesses_processors"></a>
The `processors` block supports:

* `type` - The type of the parser.
  + **processor_regex**
  + **processor_split_string**
  + **processor_json**
  + **processor_gotime**
  + **processor_filter_regex**
  + **processor_drop**
  + **processor_rename**

* `detail` - The configuration of the parser, in JSON format.

<a name="data_cce_accesses_demoFields"></a>
The `demo_fields` block supports:

* `field_name` - The name of the parsed field.

* `field_value` - The value of the parsed field.

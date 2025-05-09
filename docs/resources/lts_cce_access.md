---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_cce_access"
description: |-
  Manages an LTS CCE access resource within HuaweiCloud.
---

# huaweicloud_lts_cce_access

Manages an LTS CCE access resource within HuaweiCloud.

## Example Usage

### CCE Access With Container Stdout

```hcl
variable "name" {}
variable "log_group_id" {}
variable "log_stream_id" {}
variable "host_group_id" {}
variable "cluster_id" {}

resource "huaweicloud_lts_cce_access" "container_stdout" {
  name           = var.name
  log_group_id   = var.log_group_id
  log_stream_id  = var.log_stream_id
  host_group_ids = [var.host_group_id]
  cluster_id     = var.cluster_id

  access_config {
    path_type = "container_stdout"
    stdout    = true

    windows_log_info {
      categorys        = ["System", "Application"]
      event_level      = ["warning", "error"]
      time_offset_unit = "day"
      time_offset      = 7
    }

    single_log_format {
      mode = "system"
    }
  }
}
```

### CCE Access With Container File

```hcl
variable "name" {}
variable "log_group_id" {}
variable "log_stream_id" {}
variable "host_group_id" {}
variable "cluster_id" {}

resource "huaweicloud_lts_cce_access" "container_file" {
  name           = var.name
  log_group_id   = var.log_group_id
  log_stream_id  = var.log_stream_id
  host_group_ids = [var.host_group_id]
  cluster_id     = var.cluster_id

  access_config {
    path_type = "container_file"
    paths       = ["/var"]

    windows_log_info {
      categorys        = ["System", "Application"]
      event_level      = ["warning", "error"]
      time_offset_unit = "day"
      time_offset      = 7
    }

    single_log_format {
      mode = "system"
    }
  }
}
```

### CCE Access With Host File

```hcl
variable "name" {}
variable "log_group_id" {}
variable "log_stream_id" {}
variable "host_group_id" {}
variable "cluster_id" {}

resource "huaweicloud_lts_cce_access" "host_file" {
  name           = var.name
  log_group_id   = var.log_group_id
  log_stream_id  = var.log_stream_id
  host_group_ids = [var.host_group_id]
  cluster_id     = var.cluster_id

  access_config {
    path_type = "host_file"
    paths       = ["/var"]

    single_log_format {
      mode = "system"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create CCE access.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `name` - (Required, String) Specifies the name of the CCE access. The name consists of `1` to `64`
  characters. Only letters, digits, underscores (_), and periods (.) are allowed, and the period cannot be the first
  or last character.

* `log_group_id` - (Required, String, ForceNew) Specifies the log group ID. Changing this creates a new resource.

* `log_stream_id` - (Required, String, ForceNew) Specifies the log stream ID. Changing this creates a new resource.

* `access_config` - (Required, List) Specifies the configurations of CCE access.
  The [access_config](#block_access_config) structure is documented below.

* `cluster_id` - (Required, String, ForceNew) Specifies the CCE cluster ID. Changing this creates a new resource.

* `host_group_ids` - (Optional, List) Specifies the log access host group ID list.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the CCE access.

* `binary_collect` - (Optional, Bool) Specifies whether collect in binary format. Default is **false**.

* `log_split` - (Optional, Bool) Specifies whether to split log. Default is false.

* `processor_type` - (Optional, String) Specifies the type of the ICAgent structuring parsing.  
  This parameter must be set together with the `processors` parameter.  
  The valid values are as follows:
  + **SINGLE_LINE**
  + **MULTI_LINE**
  + **REGEX**
  + **MULTI_REGEX**
  + **SPLIT**
  + **JSON**
  + **NGINX**
  + **COMPOSE**

-> If `processor_type` is specified, it does not support modification to empty.

* `processors` - (Optional, List) Specifies the list of the ICAgent structuring parsing rules.  
  The [processors](#cce_access_processors) structure is documented below.  
  This parameter must be set together with the `processor_type` parameter.  
  Please refer to the [Setting ICAgent Structuring Parsing Rules](https://support.huaweicloud.com/intl/en-us/usermanual-lts/lts_07_0072.html).

 -> For the same log stream, If you have configured cloud structuring parsing, delete its configurations before configuring
    ICAgent structuring parsing.

* `demo_log` - (Optional, String) Specifies the example log of the ICAgent structuring parsing.  
  This parameter is available when the `processor_type` parameter is specified.

* `demo_fields` - (Optional, List) Specifies the list of the parsed fields of the example log.  
  The [demo_fields](#cce_access_demoFields) structure is documented below.  
  This parameter must be set together with the `demo_log` parameter.  
  This parameter is available when the `processor_type` parameter is specified.

* `encoding_format` - (Optional, String) Specifies the encoding format log file.  
  Defaults to **UTF-8**.  
  The valid values are as follows:
  + **UTF-8**
  + **GBK**

* `incremental_collect` - (Optional, Bool) Specifies whether to collect incrementally.  
  Defaults to **true**.  

<a name="block_access_config"></a>
The `access_config` block supports:

* `path_type` - (Required, String) Specifies the type of the CCE access. The options are as follows:
  + **container_stdout**
  + **container_file**
  + **host_file**

* `paths` - (Optional, List) Specifies the collection paths. Required when `path_type` is **container_file**
  or **host_file**.

* `black_paths` - (Optional, List) Specifies the collection path blacklist.

* `windows_log_info` - (Optional, List) Specifies the configuration of Windows event logs. Required when
  `path_type` is **container_file** or **container_stdout**.
  The [windows_log_info](#block_access_config_windows_log_info) structure is documented below.

* `single_log_format` - (Optional, List) Specifies the configuration single-line logs. Each log line is displayed
  as a single log event.
  The [single_log_format](#block_access_config_single_log_format) structure is documented below.

* `multi_log_format` - (Optional, List) Specifies the configuration multi-line logs. Multiple lines of exception log
  events can be displayed as a single log event. This is helpful when you check logs to locate problems.
  The [multi_log_format](#block_access_config_multi_log_format) structure is documented below.

 -> `single_log_format` or `multi_log_format` must be specified.

* `stdout` - (Optional, Bool) Specifies whether output is standard. Default is false.

* `stderr` - (Optional, Bool) Specifies whether error output is standard. Default is **false**.

 ->  If the value of `path_type` is **container_stdout**, `stdout` or `stderr` must be **true**.

* `name_space_regex` - (Optional, String) Specifies the regular expression matching of kubernetes namespaces.
  LTS will collect logs of the namespaces with names matching this expression. To collect logs of all namespaces,
  leave this field empty.

* `pod_name_regex` - (Optional, String) Specifies the regular expression matching of kubernetes pods.
  LTS will collect logs of the Pods with names matching this expression. To collect logs of all Pods,
  leave this field empty.

* `container_name_regex` - (Optional, String) Specifies the regular expression matching of kubernetes container names.
  LTS will collect logs of the containers with names matching this expression. To collect logs of all containers,
  leave this field empty.

* `log_labels` - (Optional, Map) Specifies the container label log tag. A maximum of `30` tags can be created.
  The key names must be unique. LTS adds the specified fields to the log when each label key has a corresponding
  label value. For example, if you enter `app` as the key and `app_alias` as the value, when the Container label
  contains `app=lts`, `{app_alias: lts}` will be added to the log.

* `include_labels_logical` - (Optional, String) Specifies the logical relationship between multiple container label
  whitelists.  
  Defaults to **or**.
  The valid values are as follows:
  + **and**
  + **or**

* `include_labels` - (Optional, Map) Specifies the container label whitelist. A maximum of `30` tags can be created.
  The key names must be unique. If labelValue is empty, all containers whose container label contains labelKey are
  matched. If labelValue is not empty, only containers whose container label contains `LabelKey=LabelValue` are
  matched. LabelKey must be fully matched, and labelValue supports regular expression matching. Multiple whitelists
  are in the OR relationship. That is, a container label can be matched as long as it meets any of the whitelists.

* `exclude_labels_logical` - (Optional, String) Specifies the logical relationship between multiple container label
  blacklists.  
  Defaults to **or**.
  The valid values are as follows:
  + **and**
  + **or**

* `exclude_labels` - (Optional, Map) Specifies the container label blacklist. A maximum of `30` tags can be created.
  The key names must be unique. If labelValue is empty, all containers whose container label contains labelKey are
  excluded. If labelValue is not empty, only containers whose container label contains `LabelKey=LabelValue` are
  excluded. LabelKey must be fully matched, and labelValue supports regular expression matching. Multiple blacklists
  are in the OR relationship. That is, a container label can be excluded as long as it meets any of the blacklists.

* `log_envs` - (Optional, Map) Specifies the environment variable tag. A maximum of `30` tags can be created.
  The key names must be unique. LTS adds the specified fields to the log when each environment variable key has a
  corresponding environment variable value. For example, if you enter `app` as the key and `app_alias` as the value,
  when the kubernetes environment variable contains `app=lts`, `{app_alias: lts}` will be added to the log.

* `include_envs_logical` - (Optional, String) Specifies the logical relationship between multiple environment variable
  whitelists.  
  Defaults to **or**.
  The valid values are as follows:
  + **and**
  + **or**

* `include_envs` - (Optional, Map) Specifies the environment variable whitelist. A maximum of `30` tags can be
  created. The key names must be unique. LTS will match all containers with environment variables containing either
  an environment variable key with an empty corresponding environment variable value, or an environment variable key
  with its corresponding environment variable value. LabelKey must be fully matched, and labelValue supports regular
  expression matching.

* `exclude_envs_logical` - (Optional, String) Specifies the logical relationship between multiple environment variable
  blacklists.  
  Defaults to **or**.
  The valid values are as follows:
  + **and**
  + **or**

* `exclude_envs` - (Optional, Map) Specifies the environment variable blacklist. A maximum of `30` tags can be
  created. The key names must be unique. LTS will exclude all containers with environment variables containing either
  an environment variable key with an empty corresponding environment variable value, or an environment variable key
  with its corresponding environment variable value. LabelKey must be fully matched, and labelValue supports regular
  expression matching.

* `log_k8s` - (Optional, Map) Specifies the kubernetes label log tag. A maximum of `30` tags can be created.
  The key names must be unique. LTS adds the specified fields to the log when each label key has a corresponding label
  value. For example, if you enter `app` as the key and `app_alias` as the value, when the K8s label contains
  `app=lts`, `{app_alias: lts}` will be added to the log.

* `include_k8s_labels_logical` - (Optional, String) Specifies the logical relationship between multiple kubernetes label
  whitelists.  
  Defaults to **or**.
  The valid values are as follows:
  + **and**
  + **or**

* `include_k8s_labels` - (Optional, Map) Specifies the kubernetes label whitelist. A maximum of `30` tags can be
  created. The key names must be unique. If labelValue is empty, all containers whose K8s label contains labelKey are
  matched. If labelValue is not empty, only containers whose K8s Label contains `LabelKey=LabelValue` are matched.
  LabelKey must be fully matched, and labelValue supports regular expression matching. Multiple whitelists are in the
  OR relationship. That is, a K8s label can be matched as long as it meets any of the whitelists.

* `exclude_k8s_labels_logical` - (Optional, String) Specifies the logical relationship between multiple kubernetes label
  blacklists.  
  Defaults to **or**.
  The valid values are as follows:
  + **and**
  + **or**

* `exclude_k8s_labels` - (Optional, Map) Specifies the kubernetes label blacklist. A maximum of `30` tags can be
  created. The key names must be unique. If labelValue is empty, all containers whose K8s label contains labelKey are
  excluded. If labelValue is not empty, only containers whose K8s label contains `LabelKey=LabelValue` are excluded.
  LabelKey must be fully matched, and labelValue supports regular expression matching. Multiple blacklists are in the
  OR relationship. That is, a K8s Label can be excluded as long as it meets any of the blacklists.

* `repeat_collect` - (Optional, Bool) Specifies whether to allow repeated file collection.  
  Defaults to **false**.
  + If this parameter is set to **true**, one host log file can be collected to multiple log streams.
    This function is available only to certain ICAgent versions, please refer to the [documentation](<https://support.huaweicloud.com/intl/en-us/usermanual-lts/lts_02_0014.html#lts_02_0014__section7761151916252>).
  + If this parameter is set to **false**, the same log file in the same host cannot be collected to different log streams.

* `custom_key_value` - (Optional, Map, ForceNew) Specifies the custom key/value pairs of the CCE access.  
  Changing this creates a new resource.

* `system_fields` - (Optional, List, ForceNew) Specifies the list of system built-in fields of the CCE access.  
  Changing this creates a new resource.  
  If `custom_key_value` is specified, the value of `system_fields` will be automatically assigned by
  the system as **pathfile**.  
  If `system_fields` is specified, **pathFile** must be included.  
  The valid values are as follows:
  + **pathFile**
  + **hostName**
  + **hostId**
  + **hostIP**
  + **hostIPv6**
  + **clusterId**
  + **podName**
  + **appName**
  + **containerName**
  + **nameSpace**
  + **category**
  + **serviceID**
  + **podIp**
  + **clusterName**
  + **workloadType**
  + **image_name**

 -> These parameters include `name_space_regex`, `pod_name_regex`, `container_name_regex`, `log_labels`,
    `include_labels_logical`, `include_labels`, `exclude_labels_logical`, `exclude_labels`, `log_envs`,
    `include_envs_logical`, `include_envs`, `exclude_envs_logical`, `exclude_envs`, `log_k8s`,
    `include_k8s_labels_logical`, `include_k8s_labels`, `exclude_k8s_labels_logical` and
    `exclude_k8s_labels` are available, only `path_type` is not **host_file**.

<a name="block_access_config_single_log_format"></a>
The `single_log_format` block supports:

* `mode` - (Required, String) Specifies mode of single-line log format. The options are as follows:
  + **system**: the system time.
  + **wildcard**: the time wildcard.

* `value` - (Optional, String) Specifies value of single-line log format.
  + If mode is **system**, the value is the current timestamp, the number of milliseconds elapsed
    since January 1, 1970 UTC.
  + If mode is **wildcard**, the value is **required** and is a time wildcard, which is used to look for
    the log printing time as the beginning of a log event. If the time format in a log event
    is `2019-01-01 23:59:59`, the time wildcard is **YYYY-MM-DD hh:mm:ss**. If the time format in
    a log event is `19-1-1 23:59:59`, the time wildcard is **YY-M-D hh:mm:ss**.

<a name="block_access_config_multi_log_format"></a>
The `multi_log_format` block supports:

* `mode` - (Required, String) Specifies mode of multi-line log format. The options are as follows:
  + **time**: the time wildcard.
  + **regular**: the regular expression.

* `value` - (Required, String) Specifies value of multi-line log format.
  + If mode is **regular**, the value is a regular expression.
  + If mode is **time**, the value is a time wildcard, which is used to look for the log printing time
    as the beginning of a log event. If the time format in a log event is `2019-01-01 23:59:59`, the time
    wildcard is **YYYY-MM-DD hh:mm:ss**. If the time format in a log event is `19-1-1 23:59:59`, the time
    wildcard is **YY-M-D hh:mm:ss**.

-> The time wildcard and regular expression will look for the specified pattern right from the beginning of each
   log line. If no match is found, the system time, which may be different from the time in the log event, is used.
   In general cases, you are advised to select **Single-line** for Log Format and **system** time for Log Time.

<a name="block_access_config_windows_log_info"></a>
The `windows_log_info` block supports:

* `categorys` - (Required, List) Specifies the types of Windows event logs to collect. The valid values are
  **Application**, **System**, **Security** and **Setup**.

* `event_level` - (Required, List) Specifies the Windows event severity. The valid values are **information**,
  **warning**, **error**, **critical** and **verbose**. Only Windows Vista or later is supported.

* `time_offset_unit` - (Required, String) Specifies the collection time offset unit. The valid values are
  **day**, **hour** and **sec**.

* `time_offset` - (Required, Int) Specifies the collection time offset. This time takes effect only for the first
  time to ensure that the logs are not collected repeatedly.

  + When `time_offset_unit` is set to **day**, the value ranges from `1` to `7` days.
  + When `time_offset_unit` is set to **hour**, the value ranges from `1` to `168` hours.
  + When `time_offset_unit` is set to **sec**, the value ranges from `1` to `604,800` seconds.

<a name="cce_access_processors"></a>
The `processors` block supports:

* `type` - (Optional, String) Specifies the type of the parser.  
  The valid values are as follows:
  + **processor_regex**
  + **processor_split_string**
  + **processor_json**
  + **processor_gotime**
  + **processor_filter_regex**
  + **processor_drop**
  + **processor_rename**

* `detail` - (Optional, String) Specifies the configuration of the parser, in JSON format.  
  For the keys, please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/api-lts/CreateAccessConfig.html#CreateAccessConfig__request_Detail).

<a name="cce_access_demoFields"></a>
The `demo_fields` block supports:

* `field_name` - (Required, String) Specifies the name of the parsed field.

* `field_value` - (Optional, String) Specifies the value of the parsed field.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `access_type` - The log access type.

* `created_at` - The creation time of the CCE access, in RFC3339 format.

* `log_group_name` - The log group name.

* `log_stream_name` - The log stream name.

## Import

The CCE access can be imported using `name`, e.g.

```bash
$ terraform import huaweicloud_lts_cce_access.test <name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `processors`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_lts_cce_access" "test" {
  ...

  lifecycle {
    ignore_changes = [
      processors,
    ]
  }
}
```

---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_host_accesses"
description: |-
  Use this data source to get the list of the host accesses with HuaweiCloud.
---

# huaweicloud_lts_host_accesses

Use this data source to get the list of the host accesses with HuaweiCloud.

## Example Usage

### Query all host accesses

```hcl
data "huaweicloud_lts_host_accesses" "test" {}
```

### Query the host accesses by the specified host access names

```hcl
variable "host_access_names" {
  type = list(string)
}

data "huaweicloud_lts_host_accesses" "test" {
  access_config_name_list = var.host_access_names
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `access_config_name_list` - (Optional, List) Specifies the list of the host access names.

* `host_group_name_list` - (Optional, List) Specifies the list host of the group names associated with the host access.

* `log_group_name_list` - (Optional, List) Specifies the list of log group names.

* `log_stream_name_list` - (Optional, List) Specifies the list of log stream names.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the host access.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `accesses` - All host accesses that match the filter parameters.

  The [accesses](#host_accesses_struct) structure is documented below.

<a name="host_accesses_struct"></a>
The `accesses` block supports:

* `id` - The ID of the host access.

* `name` - The name of the host access.

* `host_group_ids` - The ID list of the log access host groups.

* `log_group_id` - The ID of the log group to which the host access belongs.

* `log_group_name` - The name of the log group to which the host access belongs.

* `log_stream_id` - The ID of the log stream to which the host access belongs.

* `log_stream_name` - The name of the log stream to which the host access belongs.

* `access_config` - The configuration detail of the host access.

  The [access_config](#host_accesses_access_config_struct) structure is documented below.

* `access_type` - The type of the log access.
  + **AGENT**: ECS access.

* `tags` - The key/value pairs to associate with the host access.

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

  The [processors](#host_accesses_processors_struct) structure is documented below.

* `demo_log` - The example log of the ICAgent structuring parsing.

* `demo_fields` - The list of the parsed fields of the example log.

  The [demo_fields](#host_accesses_demo_fields_struct) structure is documented below.

* `binary_collect` - Whether to allow collection of binary log files.

* `encoding_format` - The encoding format log file.
  + **UTF-8**
  + **GBK**

* `incremental_collect` - Whether to collect logs incrementally.

* `log_split` - Whether log splitting is enabled.

* `created_at` - The creation time of the host access, in RFC3339 format.

<a name="host_accesses_access_config_struct"></a>
The `access_config` block supports:

* `paths` - The list of paths where collected logs are located.

* `black_paths` - The collection path blacklist.

* `single_log_format` - The configuration single-line logs.

  The [single_log_format](#host_format_single_log_format_struct) structure is documented below.

* `multi_log_format` - The configuration multi-line logs.

  The [multi_log_format](#host_format_multi_log_format_struct) structure is documented below.

* `windows_log_info` - The configuration of Windows event logs.

  The [windows_log_info](#host_access_config_windows_log_info_struct) structure is documented below.

* `custom_key_value` - The custom key/value pairs of the host access.

* `system_fields` - The list of system built-in fields of the host access.

* `repeat_collect` - Whether the file is allowed to be collected repeatedly.

<a name="host_format_single_log_format_struct"></a>
The `single_log_format` block supports:

* `mode` - The mode of single-line log format.
  + **system**: the system time.
  + **wildcard**: the time wildcard.

* `value` - The value of single-line log format.

<a name="host_format_multi_log_format_struct"></a>
The `multi_log_format` block supports:

* `mode` - The mode of multi-line log format.
  + **time**: the time wildcard.
  + **regular**: the regular expression.

* `value` - The value of multi-line log format.

<a name="host_access_config_windows_log_info_struct"></a>
The `windows_log_info` block supports:

* `categorys` - The types of Windows event logs to be collected.
  + **Application**
  + **System**
  + **Security**
  + **Setup**

* `event_level` - The list of Windows event levels.
  + **information**
  + **warning**
  + **error**
  + **critical**
  + **verbose**

* `time_offset` - The collection time offset.

* `time_offset_unit` - The collection time offset unit.
  + **day**
  + **hour**
  + **sec**

<a name="host_accesses_processors_struct"></a>
The `processors` block supports:

* `type` - The type of the parser.

* `detail` - The configuration of the parser, in JSON format.

  The [detail](#host_processors_detail_struct) structure is documented below.

<a name="host_processors_detail_struct"></a>
The `detail` block supports:

<a name="host_accesses_demo_fields_struct"></a>
The `demo_fields` block supports:

* `name` - The name of the parsed field.

* `value` - The value of the parsed field.

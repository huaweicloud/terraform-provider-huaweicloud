---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_host_access"
description: |-
  Manages an LTS host access resource within HuaweiCloud.
---

# huaweicloud_lts_host_access

Manages an LTS host access resource within HuaweiCloud.

## Example Usage

```hcl
variable "group_id" {}
variable "stream_id" {}
variable "host_group_id" {}

resource "huaweicloud_lts_host_access" "test" {
  name           = "access-demo"
  log_group_id   = var.group_id
  log_stream_id  = var.stream_id
  host_group_ids = [ var.host_group_id ]

  access_config {
    paths       = ["/var/log/*"]
    black_paths = ["/var/log/*/a.log"]

    single_log_format {
      mode = "system"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the host access name. The name consists of `1` to `64` characters.
  Only letters, digits, underscores (_), and periods (.) are allowed, and the period cannot be the first or last character.

* `log_group_id` - (Required, String, ForceNew) Specifies the log group ID.
  Changing this parameter will create a new resource.

* `log_stream_id` - (Required, String, ForceNew) Specifies the log stream ID.
  Changing this parameter will create a new resource.

* `access_config` - (Required, List) Specifies the configurations of host access.
  The [access_config](#HostAccessConfigDeatil) structure is documented below.

* `host_group_ids` - (Optional, List) Specifies the log access host group ID list.

* `tags` - (Optional, Map) Specifies the key/value to attach to the host access.

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

* `processors` - (Optional, List) Specifies the list of the ICAgent structuring parsing rules.  
  The [processors](#HostAccessProcessors) structure is documented below.  
  This parameter must be set together with the `processor_type` parameter.  
  Please refer to the [Setting ICAgent Structuring Parsing Rules](https://support.huaweicloud.com/intl/en-us/usermanual-lts/lts_07_0072.html).

 -> For the same log stream, If you have configured cloud structuring parsing, delete its configurations before configuring
    ICAgent structuring parsing.

* `demo_log` - (Optional, String) Specifies the example log of the ICAgent structuring parsing.  
  This parameter is available when the `processor_type` parameter is specified.

* `demo_fields` - (Optional, List) Specifies the list of the parsed fields of the example log.  
  The [demo_fields](#HostAccessDemoFields) structure is documented below.  
  This parameter must be set together with the `demo_log` parameter.  
  This parameter is available when the `processor_type` parameter is specified.

* `binary_collect` - (Optional, Bool, ForceNew) Specifies whether to allow collection of binary log files.  
  Defaults to **false**.  
  Changing this parameter will create a new resource.

* `encoding_format` - (Optional, String) Specifies the encoding format log file.  
  Defaults to **UTF-8**.  
  The valid values are as follows:
  + **UTF-8**
  + **GBK**

* `incremental_collect` - (Optional, Bool) Specifies whether to collect incrementally.  
  Defaults to **true**.  
  When incremental collection a new file, ICAgent reads the file from the end of the file.  
  When full collection a new file, ICAgent reads the file from the beginning of the file.

* `log_split` - (Optional, Bool) Specifies whether to enable log splitting.  
  Defaults to **false**.

<a name="HostAccessConfigDeatil"></a>
The `access_config` block supports:

* `paths` - (Required, List) Specifies the collection paths.

  + A path must start with `/` or `Letter:\`.
  + A path cannot contain only slashes (/). The following special characters are not allowed: <>'|"
  + A path cannot start with `/**` or `/*`.
  + Only one double asterisk (**) can be contained in a path.
  + Up to 10 paths can be specified.

* `black_paths` - (Optional, List) Specifies the collection path blacklist.

  + A path must start with `/` or `Letter:\`.
  + A path cannot contain only slashes (/). The following special characters are not allowed: <>'|"
  + A path cannot start with `/**` or `/*`.
  + Only one double asterisk (**) can be contained in a path.
  + Up to 10 paths can be specified.

  -> If you blacklist a file or directory that has been set as a collection path, the blacklist settings
    will be used and the file or files in the directory will be filtered out.

* `single_log_format` - (Optional, List) Specifies the configuration single-line logs. Each log line is displayed as a
  single log event. The [single_log_format](#HostAccessConfigSingleLogFormat) structure is documented below.

* `multi_log_format` - (Optional, List) Specifies the configuration multi-line logs. Multiple lines of exception log events
  can be displayed as a single log event. This is helpful when you check logs to locate problems.
  The [multi_log_format](#HostAccessConfigMultiLogFormat) structure is documented below.

* `windows_log_info` - (Optional, List) Specifies the configuration of Windows event logs.
  The [windows_log_info](#HostAccessConfigWindowsLogInfo) structure is documented below.

* `custom_key_value` - (Optional, Map, ForceNew) Specifies the custom key/value pairs of the host access.  
  Changing this parameter will create a new resource.

* `system_fields` - (Optional, List, ForceNew) Specifies the list of system built-in fields of the host access.  
  Changing this parameter will create a new resource.  
  If `custom_key_value` is specified, the value of `system_fields` will be automatically assigned by
  the system as **pathfile**.  
  If `system_fields` is specified, **pathFile** must be included.  
  The valid values are as follows:
  + **pathFile**
  + **hostName**
  + **hostId**
  + **hostIP**
  + **hostIPv6**

* `repeat_collect` - (Optional, Bool) Specifies whether to allow repeated flie collection.  
  Defaults to **false**.
  + If this parameter is set to **true**, one host log file can be collected to multiple log streams.
    This function is available only to certain ICAgent versions, please refer to the [documentation](<https://support.huaweicloud.com/intl/en-us/usermanual-lts/lts_02_0014.html#lts_02_0014__section7761151916252>).
  + If this parameter is set to **false**, the same log file in the same host cannot be collected to different log streams.

<a name="HostAccessConfigSingleLogFormat"></a>
The `single_log_format` blocks supports:

* `mode` - (Required, String) Specifies mode of single-line log format. The options are as follows:
  + **system**: the system time.
  + **wildcard**: the time wildcard.

* `value` - (Optional, String) Specifies value of single-line log format.
  + If mode is **system**, the value is the current timestamp, the number of milliseconds elapsed since January 1, 1970
    UTC.
  + If mode is **wildcard**, the value is **required** and is a time wildcard, which is used to look for the log
    printing time as the beginning of a log event. If the time format in a log event is `2019-01-01 23:59:59`,
    the time wildcard is **YYYY-MM-DD hh:mm:ss**. If the time format in a log event is `19-1-1 23:59:59`,
    the time wildcard is **YY-M-D hh:mm:ss**.

<a name="HostAccessConfigMultiLogFormat"></a>
The `multi_log_format` blocks supports:

* `mode` - (Required, String) Specifies mode of multi-line log format. The options are as follows:
  + **time**: the time wildcard.
  + **regular**: the regular expression.

* `value` - (Required, String) Specifies value of multi-line log format.
  + If mode is **regular**, the value is a regular expression.
  + If mode is **time**, the value is a time wildcard, which is used to look for the log printing time as the beginning
    of a log event. If the time format in a log event is `2019-01-01 23:59:59`,
    the time wildcard is **YYYY-MM-DD hh:mm:ss**.
    If the time format in a log event is `19-1-1 23:59:59`, the time wildcard is **YY-M-D hh:mm:ss**.

-> The time wildcard and regular expression will look for the specified pattern right from the beginning of each log
  line. If no match is found, the system time, which may be different from the time in the log event, is used.
  In general cases, you are advised to select **Single-line** for Log Format and **system** time for Log Time.

<a name="HostAccessConfigWindowsLogInfo"></a>
The `windows_log_info` block supports:

* `categorys` - (Required, List) Specifies the types of Windows event logs to collect. The valid values are
  **Application**, **System**, **Security** and **Setup**.

* `event_level` - (Required, List) Specifies the Windows event severity. The valid values are **information**,
  **warning**, **error**, **critical** and **verbose**.  Only Windows Vista or later is supported.

* `time_offset_unit` - (Required, String) Specifies the collection time offset unit. The valid values are
  **day**, **hour** and **sec**.

* `time_offset` - (Required, Int) Specifies the collection time offset. This time takes effect only for the first
  time to ensure that the logs are not collected repeatedly.

  + When `time_offset_unit` is set to **day**, the value ranges from `1` to `7` days.
  + When `time_offset_unit` is set to **hour**, the value ranges from `1` to `168` hours.
  + When `time_offset_unit` is set to **sec**, the value ranges from `1` to `604,800` seconds.

<a name="HostAccessProcessors"></a>
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

<a name="HostAccessDemoFields"></a>
The `demo_fields` block supports:

* `name` - (Optional, String) Specifies the name of the parsed field.

* `value` - (Optional, String) Specifies the value of the parsed field.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the host access.

* `access_type` - The log access type.

* `log_group_name` - The log group name.

* `log_stream_name` - The log stream name.

* `created_at` - The creation time of the host access, in RFC3339 format.

## Import

The host access can be imported using the `name`, e.g.

```bash
$ terraform import huaweicloud_lts_host_access.test <name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `processors`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```hcl
resource "huaweicloud_lts_host_access" "test" {
  ...

  lifecycle {
    ignore_changes = [
      processors,
    ]
  }
}
```

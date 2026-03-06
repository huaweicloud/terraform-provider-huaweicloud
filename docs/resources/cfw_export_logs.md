---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_export_logs"
description: |-
  Manages a resource to export firewall logs within HuaweiCloud.
---

# huaweicloud_cfw_export_logs

Manages a resource to export firewall logs within HuaweiCloud.

-> 1. This resource is a one-time action resource used to export firewall logs. Deleting this resource will
  not clear the corresponding request record, but will only remove the resource information from the tf state file.
  <br/>2. Executing this resource will generate a log file with the suffix **.csv** in the current working directory.

## Example Usage

```hcl
variable "fw_instance_id" {}
variable "start_time" {}
variable "end_time" {}
variable "log_type" {}
variable "type" {}

resource "huaweicloud_cfw_export_logs" "test" {
  fw_instance_id = var.fw_instance_id
  start_time     = var.start_time
  end_time       = var.end_time
  log_type       = var.log_type
  type           = var.type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `fw_instance_id` - (Required, String, NonUpdatable) Specifies the firewall instance ID.

* `start_time` - (Required, Int, NonUpdatable) Specifies the start time. Millisecond level timestamp.

* `end_time` - (Required, Int, NonUpdatable) Specifies the end time. Millisecond level timestamp.

* `log_type` - (Required, String, NonUpdatable) Specifies the log type.  
  The valid values are as follows:
  + **internet**: North-south (internet) logs.
  + **nat**: NAT scenario logs.
  + **vpc**: East-west (VPC) logs.
  + **vgw**: VGW scenario logs.

* `type` - (Required, String, NonUpdatable) Specifies the type.  
  The valid values are as follows:
  + **attack**: Attack logs.
  + **acl**: Access control logs.
  + **flow**: Flow logs.
  + **url**: URL logs.

* `filters` - (Optional, List, NonUpdatable) Specifies the filter conditions.

  The [filters](#filters_struct) structure is documented below.

* `time_zone` - (Optional, String, NonUpdatable) Specifies the time zone. The valid value is **GMT+08:00**.

* `export_file_name` - (Optional, String, NonUpdatable) Specifies the name of the exported firewall logs file.
  If omitted, the default file name `cfw-{log_type}-{type}-log-{start_time}.csv` will be used.
  If the file name does not end in **.csv**, **.csv** will be automatically appended.

<a name="filters_struct"></a>
The `filters` block supports:

* `field` - (Required, String, NonUpdatable) Specifies the log field, such as **src_ip**.

* `operator` - (Required, String, NonUpdatable) Specifies the operator.  
  The valid values are as follows:
  + **equal**: Equal.
  + **not_equal**: Not equal.
  + **contain**: Contain.
  + **starts_with**: Starts with.

* `values` - (Optional, List, NonUpdatable) Specifies the filter values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `fw_instance_id`.

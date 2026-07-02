---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_audit_logs"
description: |-
  Use this data source to get the list of audit logs of TaurusDB instance within HuaweiCloud.
---

# huaweicloud_taurusdb_audit_logs

Use this data source to get the list of audit logs of TaurusDB instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_taurusdb_audit_logs" "test" {
  instance_id = var.instance_id
  start_time  = var.start_time
  end_time    = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the TaurusDB instance.

* `start_time` - (Required, String) Specifies the start time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `end_time` - (Required, String) Specifies the end time in the **yyyy-mm-ddThh:mm:ssZ** format.
  The end time must be later than the start time and the time span cannot be longer than 30 days.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `audit_logs` - Indicates the list of audit logs.

  The [audit_logs](#audit_logs_struct) structure is documented below.

<a name="audit_logs_struct"></a>
The `audit_logs` block supports:

* `id` - Indicates the audit log ID.

* `name` - Indicates the audit log file name.

* `size` - Indicates the audit log size, in KB.

* `begin_time` - Indicates the start time of the audit log.

* `end_time` - Indicates the end time of the audit log.

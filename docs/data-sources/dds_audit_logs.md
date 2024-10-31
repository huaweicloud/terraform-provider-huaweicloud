---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_audit_logs"
description: |-
  Use this data source to get the list of DDS instance audit logs.
---

# huaweicloud_dds_audit_logs

Use this data source to get the list of DDS instance audit logs.

## Example Usage

```hcl
variable "instance_id" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_dds_audit_logs" "test" {
  instance_id = var.instance_id
  start_time  = var.start_time
  end_time    = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `start_time` - (Required, String) Specifies the start time. The format of the start time is **yyyy-MM-ddThh:mm:ssZ**.

* `end_time` - (Required, String) Specifies the end time. The format of the end time is **yyyy-mm-ddThh:mm:ssZ**.
  The end time must be later than the start time.
  The time span cannot be longer than 30 days.

* `node_id` - (Optional, String) Specifies the ID of the node whose audit logs are to be queried.
  The audit logs of cluster instances are distributed on mongos nodes.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `audit_logs` - Indicates the audit log details.

  The [audit_logs](#audit_logs_struct) structure is documented below.

<a name="audit_logs_struct"></a>
The `audit_logs` block supports:

* `id` - Indicates the audit log ID.

* `name` - Indicates the audit log file name.

* `node_id` - Indicates the node ID.

* `size` - Indicates the size of the audit log in byte.

* `start_time` - Indicates the start time of the audit log. The format is **yyyy-mm-ddThh:mm:ssZ**.

* `end_time` - Indicates the end time of the audit log. The format is **yyyy-mm-ddThh:mm:ssZ**.

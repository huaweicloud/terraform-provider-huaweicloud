---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_sql_audit_logs"
description: |-
  Use this data source to get the list of RDS SQL audit logs.
---

# huaweicloud_rds_sql_audit_logs

Use this data source to get the list of RDS SQL audit logs.

## Example Usage

```hcl
variable "instance_id" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_rds_sql_audit_logs" "test" {
  instance_id = var.instance_id
  start_time  = var.start_time
  end_time    = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `start_time` - (Required, String) Specifies the start time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `end_time` - (Required, String) Specifies the end time in the **yyyy-mm-ddThh:mm:ssZ** format.
  It must be later than the start time. The time span cannot be longer than 30 days.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `audit_logs` - Indicates the list of the SQL audit logs.

  The [audit_logs](#audit_logs_struct) structure is documented below.

<a name="audit_logs_struct"></a>
The `audit_logs` block supports:

* `id` - Indicates the ID of the audit log.

* `name` - Indicates the audit log file name.

* `size` - Indicates the size in KB of the audit log.

* `begin_time` - Indicates the start time of the audit log.

* `end_time` - Indicates the end time of the audit log.

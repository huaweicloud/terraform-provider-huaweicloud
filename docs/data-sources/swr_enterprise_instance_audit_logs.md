---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_instance_audit_logs"
description: |-
  Use this data source to get the list of SWR enterprise instance audit logs.
---

# huaweicloud_swr_enterprise_instance_audit_logs

Use this data source to get the list of SWR enterprise instance audit logs.

## Example Usage

```hcl
variable "instance_id" {}
variable "operation" {}

data "huaweicloud_swr_enterprise_instance_audit_logs" "test" {
  instance_id = var.instance_id
  operation   = var.operation
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the enterprise instance ID.

* `operation` - (Required, String) Specifies the operation type.
  Values can be **pull**, **delete**, **create**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `audit_logs` - Indicates the audit logs.

  The [audit_logs](#audit_logs_struct) structure is documented below.

* `total` - Indicates the total count.

<a name="audit_logs_struct"></a>
The `audit_logs` block supports:

* `id` - Indicates the audit log ID.

* `operation` - Indicates the operation type.

* `resource_type` - Indicates the resource type.

* `resource` - Indicates the resource name.

* `username` - Indicates the user name.

* `op_time` - Indicates the operation time.

---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_sql_audit_operations"
description: ""
---

# huaweicloud_rds_sql_audit_operations

Use this data source to get the list of RDS SQL audit operations.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_sql_audit_operations" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `operation_types` - (Optional, List) Specifies the list of the operation type. Value options: **DDL**, **DCL**,
  **DML**, **OTHER**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `operations` - Indicates the list of audit operations.

  The [operations](#operations_struct) structure is documented below.

<a name="operations_struct"></a>
The `operations` block supports:

* `type` - Indicates the type of the operation.

* `actions` - Indicates the list of the operation actions.

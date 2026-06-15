---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_dr_operation_records"
description: |-
  Use this data source to query the disaster recovery operation records of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_dr_operation_records

Use this data source to query the disaster recovery operation records of a GaussDB instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_dr_operation_records" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID to query the DR operation records.

* `entity_id` - (Required, String) Specifies the entity ID to filter the records.

* `entity_type` - (Required, String) Specifies the entity type to filter the records.
  The value can be **dr**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the instance ID.

* `records` - The list of DR operation records.
  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `id` - The record ID.

* `action` - The operation action type.

* `status` - The operation status.

* `message` - The operation message.

* `entity_id` - The entity ID.

* `entity_type` - The entity type.

* `job_id` - The job ID.

* `instance_id` - The instance ID.

* `created_at` - The creation time in the format of **yyyy-mm-ddThh:mm:ssZ**.

* `updated_at` - The update time in the format of **yyyy-mm-ddThh:mm:ssZ**.

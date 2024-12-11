---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_instance_parameter_modification_records"
description: |-
  Use this data source to get the list of DDS instance parameter modification records.
---

# huaweicloud_dds_instance_parameter_modification_records

Use this data source to get the list of DDS instance parameter modification records.

## Example Usage

```hcl
variable "instance_id"  {}

data "huaweicloud_dds_instance_parameter_modification_records" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `entity_id` - (Optional, String, ForceNew) Specifies the ID of a DDS instance entity.
  + If the DB instance type is cluster and the shard or config parameter template is to be changed, the value is the
  group ID. If the parameter template of the mongos node is to be changed, the value is the node ID.
  + If the DB instance to be changed is a replica set instance, the value should be empty.
  
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `histories` - Indicates the modification records.

  The [histories](#histories_struct) structure is documented below.

<a name="histories_struct"></a>
The `histories` block supports:

* `parameter_name` - Indicates the parameter name.

* `old_value` - Indicates the old value.

* `new_value` - Indicates the new value.

* `update_result` - Indicates the update result.

* `applied` - Indicates whether the parameter is applied.

* `updated_at` - Indicates the update time, in the **yyyy-mm-ddThh:mm:ssZ** format.

* `applied_at` - Indicates the apply time, in the **yyyy-mm-ddThh:mm:ssZ** format.

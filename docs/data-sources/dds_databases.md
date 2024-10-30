---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_databases"
description: |-
  Use this data source to get the list of DDS databases.
---

# huaweicloud_dds_databases

Use this data source to get the list of DDS databases.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dds_databases" "test"{
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `databases` - Indicates the databases list.

  The [databases](#databases_struct) structure is documented below.

<a name="databases_struct"></a>
The `databases` block supports:

* `name` - Indicates the database name.

* `data_size` - Indicates the data size, unit is GB.

* `storage_size` - Indicates the storage size, unit is GB.

* `collection_num` - Indicates the collection num.

---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_cluster_database_objects"
description: |-
  Use this data source to query database objects in a specified DWS cluster within HuaweiCloud.
---

# huaweicloud_dws_cluster_database_objects

Use this data source to query database objects in a specified DWS cluster within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_dws_cluster_database_objects" "test" {
  cluster_id = var.cluster_id
  type       = "DATABASE"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the database objects are located.  
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the ID of the cluster.

* `type` - (Required, String) Specifies the type of the database object.  
  The valid values are **DATABASE**, **SCHEMA**, **TABLE**, **VIEW**, **COLUMN**, **FUNCTION**,
  **SEQUENCE** and **NODEGROUP**.

* `name` - (Optional, String) Specifies the name of the database object.

* `database` - (Optional, String) Specifies the name of the database.

* `schema` - (Optional, String) Specifies the name of the schema.  
  Required if the `type` is set to **TABLE**, **VIEW**, **COLUMN**, **FUNCTION**
  or **SEQUENCE**.

* `table` - (Optional, String) Specifies the name of the table.  
  Required if the `type` is set to **COLUMN**.

* `is_fine_grained_disaster` - (Optional, String) Specifies whether fine-grained disaster recovery
  is enabled.  
  The valid values are **true** and **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `objects` - The list of the database objects that matched filter parameters.  
  The [objects](#dws_objects_struct) structure is documented below.

<a name="dws_objects_struct"></a>
The `objects` block supports:

* `obj_name` - The name of the database object.

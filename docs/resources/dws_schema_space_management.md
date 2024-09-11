---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_schema_space_management"
description: |-
  Use this resource to modify schema space limit under specified DWS cluster within HuaweiCloud.
---
# huaweicloud_dws_schema_space_management

Use this resource to modify schema space limit under specified DWS cluster within HuaweiCloud.

-> 1. This resource is supported only in `8.1.1` or later.
   <br>2. The space quota limit only common users but not database administrators.
   <br>3.This resource is only a one-time action resource for modifying schema space limit. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "dws_cluster_id" {}
variable "database_name" {}
variable "schema_name" {}

resource "huaweicloud_dws_schema_space_management" "test" {
  cluster_id    = var.dws_cluster_id
  database_name = var.database_name
  schema_name   = var.schema_name
  space_limit   = "1024"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `cluster_id` - (Required, String, ForceNew) Specifies the DWS cluster ID.
  Changing this creates a new resource.

* `database_name` - (Required, String, ForceNew) Specifies the database name to which the schema space management belongs.
  Changing this creates a new resource.

* `schema_name` - (Required, String, ForceNew) Specifies the name of the schema.
  Changing this creates a new resource.

* `space_limit` - (Required, String, ForceNew) Specifies space limit of the schema, in KB.  
  The valid value ranges from `-1` to `9,007,199,254,740,992`, `-1` and `0` means unlimited.

  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

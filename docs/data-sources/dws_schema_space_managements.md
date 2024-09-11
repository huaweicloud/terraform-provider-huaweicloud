---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_schema_space_managements"
description: |-
  Use this data source to get the list of schema space management information of the DWS Cluster within HuaweiCloud.
---

# huaweicloud_dws_schema_space_managements

Use this data source to get the list of schema space management information of the DWS Cluster within HuaweiCloud.

## Example Usage

```hcl
variable "dws_cluster_id" {}
variable "database_name" {}

data "huaweicloud_dws_schema_space_managements" "test" {
  cluster_id    = var.dws_cluster_id
  database_name = var.database_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the DWS cluster ID.

* `database_name` - (Required, String) Specifies the database name to which the schema space management belongs.

* `schema_name` - (Optional, String) Specifies the name of the schema. Fuzzy search is supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `schemas` - All schemas that match the filter parameters.

  The [schemas](#schemas_struct) structure is documented below.

<a name="schemas_struct"></a>
The `schemas` block supports:

* `database_name` - The database name corresponding to the schema.

* `schema_name` - The name of the schema.

* `used` - The number of schema spaces used, in bytes.

* `space_limit` - The number of available spaces, in bytes.

* `skew_percent` - The skew rate of the schema.

* `min_value` - The number of used spaces by the DN with the minimum usage, in bytes.

* `max_value` - The number of used spaces by the DN with the maximum usage, in bytes.

* `dn_num` - The number of DNs.

* `min_dn` - The DN that uses the least space.

* `max_dn` - The DN that uses the most space.

---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_cluster_database_schemas"
description: |-
  Use this data source to query the list of schemas in a specified database of a DWS cluster within HuaweiCloud.
---

# huaweicloud_dws_cluster_database_schemas

Use this data source to query the list of schemas in a specified database of a DWS cluster within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_dws_cluster_database_schemas" "test" {
  cluster_id    = var.cluster_id
  database_name = "gaussdb"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the database schemas are located.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the ID of the cluster to be queried.

* `database_name` - (Required, String) Specifies the name of the database to be queried.

* `sort_key` - (Optional, String) Specifies the field used to sort query results.
  The valid value is **schemaName**.

* `sort_dir` - (Optional, String) Specifies the direction of sorting query results.
  The valid values are **ASC** and **DESC**.

-> If the sort_key and sort_dir parameters are not specified, the results will be sorted according to
   the default system behavior.  
   Typically, data is sorted by creation time.

* `keywords` - (Optional, String) Specifies the keywords used for fuzzy query.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `schemas` - The list of the database schemas that matched filter parameters.  
  The [schemas](#dws_schemas_struct) structure is documented below.

<a name="dws_schemas_struct"></a>
The `schemas` block supports:

* `schema_name` - The name of the schema.

* `database_name` - The name of the database to which the schema belongs.

* `total_value` - The total used space value of the schema.

* `perm_space` - The space threshold of the schema.

* `skew_percent` - The skew percentage of the schema.

* `min_value` - The minimum value.

* `max_value` - The maximum value.

* `min_dn` - The minimum DN node.

* `max_dn` - The maximum DN node.

* `dn_num` - The number of DN nodes.

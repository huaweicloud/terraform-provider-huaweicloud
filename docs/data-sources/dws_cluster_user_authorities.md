---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_cluster_user_authorities"
description: |-
  Use this data source to query user or role authorities in a specified DWS cluster within HuaweiCloud.
---

# huaweicloud_dws_cluster_user_authorities

Use this data source to query user or role authorities in a specified DWS cluster within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}
variable "user_or_role_name" {}

data "huaweicloud_dws_cluster_user_authorities" "test" {
  cluster_id = var.cluster_id
  name       = var.user_or_role_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the cluster user authorities are located.  
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the ID of the DWS cluster to be queried.

* `name` - (Required, String) Specifies the user name or role name to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `authorities` - The list of user or role authorities.  
  The [authorities](#dws_authorities_struct) structure is documented below.

<a name="dws_authorities_struct"></a>
The `authorities` block supports:

* `type` - The authority type.

* `database` - The database name.

* `schema_name` - The schema name.

* `object_name` - The object name.

* `all_object` - Whether all objects are effective.

* `future` - Whether future objects are effective.

* `future_object_owners` - The owners of future objects.

* `column_names` - The list of column names.

* `privileges` - The privilege list under this authority record.  
  The [privileges](#authorities_privileges_struct) structure is documented below.

<a name="authorities_privileges_struct"></a>
The `privileges` block supports:

* `permission` - The privilege name.

* `grant_with` - Whether the grant option is included.

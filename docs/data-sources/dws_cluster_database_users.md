---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_cluster_database_users"
description: |-
  Use this data source to query database users or roles in a DWS cluster within HuaweiCloud.
---

# huaweicloud_dws_cluster_database_users

Use this data source to query database users or roles in a DWS cluster within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_dws_cluster_database_users" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the database users are located.  
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the ID of the cluster to be queried.

* `type` - (Optional, String) Specifies the object type to be queried.  
  The valid values are `USER` and `ROLE`.

* `user_type` - (Optional, String) Specifies the user type to be queried.  
  The valid values are `COMMON`, `IAM` and `OneAccess`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `users` - The list of database users or roles that matched filter parameters.  
  The [users](#dws_users_struct) structure is documented below.

<a name="dws_users_struct"></a>
The `users` block supports:

* `name` - The name of the database user or role.

* `login` - Whether the database user can login.

* `user_type` - The type of the database user.

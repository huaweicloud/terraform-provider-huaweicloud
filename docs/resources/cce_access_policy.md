---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_access_policy"
description: |-
  Manages a CCE access policy resource within huaweicloud.
---

# huaweicloud_cce_access_policy

Manages a CCE access policy resource within huaweicloud.

## Example Usage

### Basic Usage

```hcl
variable "name" {}
variable "user_id" {}

resource "huaweicloud_cce_access_policy" "test" {
  name     = var.name
  clusters = ["*"]

  access_scope {
    namespaces = ["default"]
  }

  policy_type = "CCEClusterAdminPolicy"

  principal {
    type = "user"
    ids  = [ var.user_id ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CCE access policy resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new access policy resource.

* `name` - (Required, String) Specifies the access policy name.

* `clusters` - (Required, List) Specifies the list of cluster IDs.
  Wildcards (*) are allowed, which indicate all clusters.

* `access_scope` - (Required, List) Specifies the access scope,
  which is used to specify the cluster and namespace to be authorized.
  The [structure](#cce_access_policy_access_scope) is documented below.

* `policy_type` - (Required, String) Specifies the permission type.
  The value can be:
  + **CCEClusterAdminPolicy**: Administrator permissions, including read and write permissions
    on all resources in all namespaces.
  + **CCEAdminPolicy**: O&M permissions, including read and write permissions on most resources
    in all namespaces and read-only permissions on nodes, storage volumes, namespaces, and quota management.
  + **CCEEditPolicy**: Developer permissions, including read and write permissions on most resources
    in all or selected namespaces. If this kind of permissions is configured for all namespaces,
    its capability is the same as the O&M permissions.
  + **CCEViewPolicy**: read-only permissions on most resources in all or selected namespaces.

* `principal` - (Required, List) Specifies the authorization object.
  The [structure](#cce_access_policy_principal) is documented below.

<a name="cce_access_policy_access_scope"></a>
The `access_scope` block supports:

* `namespaces` - (Required, List) Specifies the list of cluster namespaces.
  Wildcards (*) are allowed to indicate all namespaces. If different clusters are selected,
  the namespace list can be a collection of multiple clusters. When RBAC authorization is used,
  CCE automatically checks whether the namespaces exist in the clusters.

<a name="cce_access_policy_principal"></a>
The `principal` block supports:

* `type` - (Required, String) Specifies the type of the authorization object.
  The value can be: **user**, **group**, **agency**.

* `ids` - (Required, List) Specifies the list of IDs of authorized objects.
  Enter the IDs based on the object type, user, user group, and agency account.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the access policy resource.
  
* `created_at` - The time when the access policy was created.

* `updated_at` - The time when the access policy was updated.

## Import

The access policy can be imported using the access policy ID, e.g.

```bash
 $ terraform import huaweicloud_cce_access_policy.my_policy <policy_id>
```

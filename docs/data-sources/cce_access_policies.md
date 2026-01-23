---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_access_policies"
description: |-
  Use this data source to get the cluster access policies.
---

# huaweicloud_cce_access_policies

Use this data source to get the cluster access policies.

## Example Usage

```hcl
data "huaweicloud_cce_access_policies" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Optional, String) Specifies the ID of CCE cluster.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `access_policy_list` - The access policy data in cce cluster.

  The [access_policy_list](#access_policy_list_struct) structure is documented below.

<a name="access_policy_list_struct"></a>
The `access_policy_list` block supports:

* `kind` - The API type.

* `api_version` - The API version.

* `name` - The access policy name.

* `policy_id` - The access policy id.

* `clusters` - The list of cluster IDs.

* `access_scope` - The access scope, which is used to specify the cluster and namespace to be authorized.
  The [access_scope](#access_scope_struct) structure is documented below.

* `policy_type` - The access policy type.

* `principal` - The authorization object.
  The [principal](#principal_struct) structure is documented below.

* `create_time` - The access policy create time.

* `update_time` - The access policy update time.

<a name="access_scope_struct"></a>
The `access_scope` block supports:

* `namespaces` - The list of cluster namespaces.

<a name="principal_struct"></a>
The `principal` block supports:

* `type` - The type of the authorization object

* `ids` - The ist of IDs of authorized objects.

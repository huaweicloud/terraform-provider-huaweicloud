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

* `name` - The access policy name.

* `policy_id` - The access policy id.

* `clusters` - The list of clusters.

* `policy_type` - The access policy type.

* `create_time` - The access policy create time.

* `update_time` - The access policy update time.

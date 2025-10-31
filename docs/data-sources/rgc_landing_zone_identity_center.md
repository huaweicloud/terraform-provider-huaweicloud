---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_landing_zone_identity_center"
description: |
  Use this data source to get the landing zone identity center in Resource Governance Center.
---

# huaweicloud_rgc_landing_zone_identity_center

Use this data source to get the landing zone identity center in Resource Governance Center.

## Example Usage

```hcl
data "huaweicloud_rgc_landing_zone_identity_center" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `identity_store_id` - The ID of the identity store.

* `user_portal_url` - The URL of the user portal.

* `permission_sets` - Information about the permission sets of the landing zone identity center.

 The [permission_sets](#permission_sets) structure is documented below.

* `groups` - Information about the group of the landing zone identity center.

 The [groups](#groups) structure is documented below.

<a name="permission_sets"></a>
The `permission_sets` block supports:

* `permission_set_id` - The ID of the permission set.

* `permission_set_name` - The name of the permission set.

* `description` - The description of the permission set.

<a name="groups"></a>
The `groups` block supports:

* `group_id` - The ID of the group.

* `group_name` - The name of the group.
  
* `description` - The description of the group.

---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_group"
description: ""
---

# huaweicloud_identitycenter_group

Manages an Identity Center group resource within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_identitycenter_instance" "system" {}

resource "huaweicloud_identitycenter_group" "test"{
  identity_store_id = data.huaweicloud_identitycenter_instance.system.identity_store_id
  name              = "test_group"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `identity_store_id` - (Required, String, ForceNew) Specifies the ID of the identity store.

  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the group.

  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description of the group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - Indicates the creation time.

* `updated_at` - Indicates the last update time.

## Import

The Identity Center group can be imported using the `identity_store_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_identitycenter_group.test <identity_store_id>/<id>
```

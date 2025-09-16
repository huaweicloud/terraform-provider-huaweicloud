---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_long_term_credential"
description: |-
  Manages a SWR enterprise long term credential resource within HuaweiCloud.
---

# huaweicloud_swr_enterprise_long_term_credential

Manages a SWR enterprise long term credential resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "name" {}

resource "huaweicloud_swr_enterprise_long_term_credential" "test" {
  instance_id = var.instance_id
  name        = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the enterprise instance ID.

* `name` - (Required, String, NonUpdatable) Specifies the credential name.

* `enable` - (Optional, Bool) Specifies whether to enable the credential. Default to **true**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `auth_token` - Indicates the auth token.

* `created_at` - Indicates the creation time.

* `expire_date` - Indicates the expired time.

* `user_id` - Indicates the user ID.

* `user_profile` - Indicates the user profile.

## Import

The credential can be imported using `instance_id` and `id`, e.g.

```bash
$ terraform import huaweicloud_swr_enterprise_long_term_credential.test <instance_id>/<id>
```

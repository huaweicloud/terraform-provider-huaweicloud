---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_enterprise_temporary_credential"
description: |-
  Manages a SWR enterprise temporary credential resource within HuaweiCloud.
---

# huaweicloud_swr_enterprise_temporary_credential

Manages a SWR enterprise temporary credential resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_swr_enterprise_temporary_credential" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the enterprise instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `auth_token` - Indicates the auth token.

* `user_id` - Indicates the user ID.

* `created_at` - Indicates the creation time.

* `expire_date` - Indicates the expired time.

---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_password"
description: |-
  Use this data source to obtain the random password generated for user **Administrator** or the user configured in
  Cloudbase-Init when you use a Cloudbase-Init-enabled image to create a Windows ECS.
---

# huaweicloud_compute_password

Use this data source to obtain the random password generated for user **Administrator** or the user configured in
Cloudbase-Init when you use a Cloudbase-Init-enabled image to create a Windows ECS.

## Example Usage

```hcl
variable "server_id" {}

data "huaweicloud_compute_password" "test" {
  server_id = var.server_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `server_id` - (Required, String) Specifies the ECS ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `password` - Indicates the password in ciphertext.

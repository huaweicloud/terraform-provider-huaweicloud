---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_resource_tag"
description: |-
  Manages an IAM v5 resource tag within HuaweiCloud.
---

# huaweicloud_identityv5_resource_tag

Manages an IAM v5 resource tag within HuaweiCloud.

-> NOTE: You must have admin privileges to use this resource.

## Example Usage

```hcl
variable "resource_id" {}

resource "huaweicloud_identityv5_resource_tag" "test" {
  resource_type = "user"
  resource_id   = var.resource_id
  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `resource_type` - (Required, String, NonUpdatable) Specifies the resource type, which can be `trust agency` or `user`.

* `resource_id` - (Required, String, NonUpdatable) Specifies the resource id, a string of 1 to 64 characters containing
  only letters, numbers, and hyphens ("-").

* `tags` - (Optional, Map) Specifies the tags of a VM node, key/value pair format.

## Import

Resource tag can be imported using the `<resource_type>/<resource_id>`. For example,
if you have resource_type `user` and resource_id `xxxx`, you should use `user/xxxx` to import.

```bash
$ terraform import huaweicloud_identityv5_resource_tag.test <resource_type>/<resource_id>
```

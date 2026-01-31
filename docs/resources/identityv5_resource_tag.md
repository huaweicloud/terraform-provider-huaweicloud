---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_resource_tag"
description: |-
  Manages an IAM resource tag within HuaweiCloud.
---

# huaweicloud_identityv5_resource_tag

Manages an IAM resource tag within HuaweiCloud.

## Example Usage

```hcl
variable "resource_id" {}
variable "tags" {}

resource "huaweicloud_identityv5_resource_tag" "test" {
  resource_type = "user"
  resource_id   = var.resource_id
  tags          = var.tags
}
```

## Argument Reference

The following arguments are supported:

* `resource_type` - (Required, String, NonUpdatable) Specifies the resource type to be associated with the tags.  
  The valid values are as follows:
  + **agency**
  + **user**

* `resource_id` - (Required, String, NonUpdatable) Specifies the resource ID to be associated with the tags.

* `tags` - (Optional, Map) Specifies the key/value pairs to be associated with the resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The resource can be imported using `resource_type` and `resource_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_identityv5_resource_tag.test <resource_type>/<resource_id>
```

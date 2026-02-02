---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_resource_tags"
description: |-
  Use this data source to query resource tag list within HuaweiCloud.
---

# huaweicloud_identityv5_resource_tags

Use this data source to query resource tag list within HuaweiCloud.

## Example Usage

```hcl
variable "resource_id" {}

data "huaweicloud_identityv5_resource_tags" "test" {
  resource_id   = var.resource_id
  resource_type = "user"
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required, String) Specifies the resource ID to be queried.

* `resource_type` - (Required, String) Specifies the resource type to be queried.  
  The valid values are as follows:
  + **agency**
  + **user**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The key/value pairs associated with the resource.  

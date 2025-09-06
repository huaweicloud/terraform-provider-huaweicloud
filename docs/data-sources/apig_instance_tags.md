---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_instance_tags"
description: |-
  Use this data source to query the tags of APIG instance within HuaweiCloud.
---

# huaweicloud_apig_instance_tags

Use this data source to query the tags of APIG instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_apig_instance_tags" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the instance tags.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the tags belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The tags that belong to the dedicated instance.  
  The [tags](#apig_instance_tags_attr) structure is documented below.

<a name="apig_instance_tags_attr"></a>
The `tags` block supports:

* `key` - The key of the tag.

* `value` - The value of the tag.

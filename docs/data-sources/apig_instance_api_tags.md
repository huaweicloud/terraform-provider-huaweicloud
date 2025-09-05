---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_instance_api_tags"
description: |-
  Use this data source to get the tag list of all APIs under specified APIG instance within HuaweiCloud.
---

# huaweicloud_apig_instance_api_tags

Use this data source to get the tag list of all APIs under specified APIG instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_apig_instance_api_tags" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the APIG instance corresponding to the APIs to which the tags
  belong that to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of the tags.

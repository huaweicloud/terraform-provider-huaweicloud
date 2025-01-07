---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_resource_tags"
description: |-
  Use this data source to get the list of tags of a specified firewall instance.
---

# huaweicloud_cfw_resource_tags

Use this data source to get the list of tags of a specified firewall instance.

## Example Usage

```hcl
variable "fw_instance_id" {}

data "huaweicloud_cfw_resource_tags" "test" {
  fw_instance_id = var.fw_instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fw_instance_id` - (Required, String) Specifies the firewall ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The tag list.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - The tag key.

* `value` - The tag value.

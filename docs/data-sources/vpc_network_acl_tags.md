---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_network_acl_tags"
description: |-
  Use this data source to get the list of network ACL project tags.
---

# huaweicloud_vpc_network_acl_tags

Use this data source to get the list of network ACL project tags.

## Example Usage

```hcl
data "huaweicloud_network_acl_tags" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of tags.

  The [tags](#tags_struct) structure is documented below.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - The tag key.

* `values` - The tag values.

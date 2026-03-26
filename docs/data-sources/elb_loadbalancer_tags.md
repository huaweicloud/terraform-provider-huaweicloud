---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_loadbalancer_tags"
description: |-
  Use this data source to get the list of loadbalancer tags.
---

# huaweicloud_elb_loadbalancer_tags

Use this data source to get the list of loadbalancer tags.

## Example Usage

```hcl
data "huaweicloud_elb_loadbalancer_tags" "test" {}
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

* `values` - The list of tag values.

---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_tags"
description: |-
  Use this data source to get the list of DMS RocketMQ tags within HuaweiCloud.
---

# huaweicloud_dms_rocketmq_tags

Use this data source to get the list of DMS RocketMQ tags within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dms_rocketmq_tags" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of tags.  
  The [tags](#dms_rocketmq_tags_attr) structure is documented below.

<a name="dms_rocketmq_tags_attr"></a>
The `tags` block supports:

* `key` - The tag key.

* `values` - The list of tag values.

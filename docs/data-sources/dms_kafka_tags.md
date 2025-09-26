---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_tags"
description: |-
  Use this data source to get tag list for all Kafka instances under the specified project within HuaweiCloud.
---

# huaweicloud_dms_kafka_tags

Use this data source to get tag list for all Kafka instances under the specified project within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dms_kafka_tags" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region to which the resource tags belong.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - The list of tags for all Kafka instances.
  The [tags](#kafka_tags_attr) structure is documented below.

<a name="kafka_tags_attr"></a>
The `tags` block supports:

* `key` - The key of the tag.

* `values` - All values corresponding to the key.

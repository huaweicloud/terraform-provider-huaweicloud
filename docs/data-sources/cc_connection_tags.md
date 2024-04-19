---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_connection_tags"
description: |-
  Use this data source to get the list of cloud connection tags.
---

# huaweicloud_cc_connection_tags

Use this data source to get the list of cloud connection tags.

## Example Usage

```hcl
data "huaweicloud_cc_connection_tags" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tags` - All tags of the cloud connection.

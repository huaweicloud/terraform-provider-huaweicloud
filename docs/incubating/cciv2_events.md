---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_events"
description: |-
  Use this data source to get the list of CCI events within HuaweiCloud.
---

# huaweicloud_cciv2_events

Use this data source to get the list of CCI events within HuaweiCloud.

## Example Usage

```hcl

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `namespace` - (Required, String) <!-- please add the description of the argument -->

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `events` - <!-- please add the description of the attribute -->
  The [events](#attrblock--events) structure is documented below.

<a name="attrblock--events"></a>
The `events` block supports:

* `annotations` - <!-- please add the description of the attribute -->

* `creation_timestamp` - <!-- please add the description of the attribute -->

* `labels` - <!-- please add the description of the attribute -->

* `name` - <!-- please add the description of the attribute -->

* `namespace` - <!-- please add the description of the attribute -->

* `resource_version` - <!-- please add the description of the attribute -->

* `uid` - <!-- please add the description of the attribute -->

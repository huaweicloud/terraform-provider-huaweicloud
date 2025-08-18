---
subcategory: "EventGrid (EG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eg_eventrouter_availability_zones"
description: |-
  Use this data source to query EG event router availability zones within HuaweiCloud.
---

# huaweicloud_eg_eventrouter_availability_zones

Use this data source to query EG event router availability zones within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_eg_eventrouter_availability_zones" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the event router availability zones are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `names` - The list of availability zone names.

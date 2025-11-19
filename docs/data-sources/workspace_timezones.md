---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_timezones"
description: |-
  Use this data source to get the list of Workspace time zones within HuaweiCloud.
---

# huaweicloud_workspace_timezones

Use this data source to get the list of Workspace time zones within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_workspace_timezones" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the time zones are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `time_zones` - The list of time zones.  
  The [time_zones](#workspace_time_zones_attr) structure is documented below.

<a name="workspace_time_zones_attr"></a>
The `time_zones` block supports:

* `name` - The name of the time zone.

* `offset` - The offset of the time zone.

* `us_description` - The English description of the time zone.

* `cn_description` - The Chinese description of the time zone.

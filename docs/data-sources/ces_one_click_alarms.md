---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_one_click_alarms"
description: |-
  Use this data source to get the list of CES one-click alarms.
---

# huaweicloud_ces_one_click_alarms

Use this data source to get the list of CES one-click alarms.

## Example Usage

```hcl
data "huaweicloud_ces_one_click_alarms" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `one_click_alarms` - The one-click monitoring list.

  The [one_click_alarms](#one_click_alarms_struct) structure is documented below.

<a name="one_click_alarms_struct"></a>
The `one_click_alarms` block supports:

* `one_click_alarm_id` - The one-click monitoring ID for a service.

* `namespace` - The metric namespace.

* `description` - The supplementary information about one-click monitoring.

* `enabled` - Whether to enable one-click monitoring.
